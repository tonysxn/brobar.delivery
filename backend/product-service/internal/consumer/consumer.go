package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/tonysanin/brobar/product-service/internal/services"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	service *services.ProductService
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewConsumer(rabbitURL string, service *services.ProductService) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Consumer{
		conn:    conn,
		channel: ch,
		service: service,
		ctx:     ctx,
		cancel:  cancel,
	}

	return c, nil
}

func (c *Consumer) Start() error {
	if err := c.setupStopListConsumer(); err != nil {
		return err
	}
	if err := c.setupStockReportConsumer(); err != nil {
		return err
	}
	
	log.Println("Product Service Consumer started")
	return nil
}

func (c *Consumer) Stop() {
	c.cancel()
	c.channel.Close()
	c.conn.Close()
}

type StopListEventItem struct {
	ProductID string  `json:"product_id"`
	Balance   float64 `json:"balance"`
}

func (c *Consumer) setupStopListConsumer() error {
	qName := "syrve.stop_list.updated"

	// Declare queue
	_, err := c.channel.QueueDeclare(
		qName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}
	
	msgs, err := c.channel.Consume(
		qName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					return
				}
				c.handleStopListUpdate(d)
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Consumer) setupStockReportConsumer() error {
	qName := "product.report.stock"
	_, err := c.channel.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil { return err }
	
	msgs, err := c.channel.Consume(qName, "", false, false, false, false, nil)
	if err != nil { return err }

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok { return }
				c.handleStockReport(d)
			case <-c.ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (c *Consumer) handleStockReport(d amqp.Delivery) {
	var payload struct {
		ChatID int64 `json:"chat_id"`
	}
	_ = json.Unmarshal(d.Body, &payload)
	
	products, _ := c.service.GetProducts(context.Background())
	
	report := "📦 <b>Звіт по складу</b>\n\n"
	
	// Sort or categorize? Just list for now.
	for _, p := range products {
		stockVal := "∞"
		if p.Stock != nil {
			stockVal = fmt.Sprintf("%.0f", *p.Stock)
		}
		hidden := ""
		if p.Hidden { hidden = " [HIDDEN]" }
		
		report += fmt.Sprintf("• %s: <b>%s</b>%s\n", p.Name, stockVal, hidden)
	}
	
	c.sendTelegramMessage(payload.ChatID, report)
	d.Ack(false)
}

func (c *Consumer) handleStopListUpdate(d amqp.Delivery) {
	// Parse Wrapper
	var wrapper struct {
		Items  []StopListEventItem `json:"items"`
		ChatID int64               `json:"chat_id"`
	}
	
	// Fallback for legacy array format (if any legacy events are in queue)
	if err := json.Unmarshal(d.Body, &wrapper); err != nil {
		// Try array
		var items []StopListEventItem
		if err2 := json.Unmarshal(d.Body, &items); err2 == nil {
			wrapper.Items = items
		} else {
			log.Printf("Failed to unmarshal stop list update: %v", err)
			d.Ack(false)
			return
		}
	}
	
	// 1. Get Current State for Comparison
	ctx := context.Background()
	products, err := c.service.GetProducts(ctx)
	if err != nil {
		log.Printf("Failed to fetch products: %v", err)
		d.Ack(false)
		return
	}
	
	// Map Products by ExternalID
	prodMap := make(map[string]float64) // ID -> Stock
	prodNames := make(map[string]string)
	
	for _, p := range products {
		if p.ExternalID != "" {
			val := -1.0 // "Infinite" marker
			if p.Stock != nil {
				val = *p.Stock
			}
			prodMap[p.ExternalID] = val
			prodNames[p.ExternalID] = p.Name
		}
	}
	
	updates := 0
	mismatches := 0
	report := "🔄 <b>Результат Синхронізації</b>\n\n"

	// 2. Process Updates
	incomingMap := make(map[string]float64)
	
	for _, item := range wrapper.Items {
		incomingMap[item.ProductID] = item.Balance
		
		currentStock, exists := prodMap[item.ProductID]
		if !exists { continue }
		
		// If currentStock != newBalance
		if currentStock != item.Balance {
			// Update
			if err := c.service.UpdateStock(ctx, item.ProductID, item.Balance); err != nil {
				log.Printf("Failed to update stock for %s: %v", item.ProductID, err)
			} else {
				updates++
				mismatches++
				// Add to report
				oldVal := "∞"
				if currentStock != -1.0 {
					oldVal = fmt.Sprintf("%.0f", currentStock)
				}
				name := prodNames[item.ProductID]
				report += fmt.Sprintf("✏️ <b>%s</b>: %s ➝ %.0f\n", name, oldVal, item.Balance)
			}
		}
	}
	
	// 3. Mark unlimited if not in incoming list?
	// Syrve logic: If item is NOT in stop list -> it is available (unlimited?)
	// Our logic: If we have stock 0 but it's not in incoming list -> Set to nil (Unlimited).
	// 3. (Optional) Handled unlimited stock logic if needed - currently trusted Syrve updates.

	if mismatches == 0 {
		report += "✅ Розбіжностей не знайдено! Все ідеально."
	} else {
		report += fmt.Sprintf("\n♻️ Оновлено позицій: %d", updates)
	}

	if wrapper.ChatID != 0 {
		c.sendTelegramMessage(wrapper.ChatID, report)
	}

	d.Ack(false)
}

func (c *Consumer) sendTelegramMessage(chatID int64, text string) {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	body, _ := json.Marshal(payload)
	
	_ = c.channel.Publish(
		"",
		"telegram_messages",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
