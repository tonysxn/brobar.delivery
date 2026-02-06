package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/tonysanin/brobar/pkg/telegram"
)

type TelegramConsumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	client        *telegram.Client
	queue         string
	defaultChatID int64
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewTelegramConsumer(rabbitURL string, client *telegram.Client, defaultChatID int64) (*TelegramConsumer, error) {
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

	return &TelegramConsumer{
		conn:          conn,
		channel:       ch,
		client:        client,
		queue:         "telegram_messages",
		defaultChatID: defaultChatID,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func (c *TelegramConsumer) Start() error {
	_, err := c.channel.QueueDeclare(
		c.queue,
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
		c.queue,
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

	// Declare taxi_events queue
	_, err = c.channel.QueueDeclare(
		"taxi_events",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	taxiMsgs, err := c.channel.Consume(
		"taxi_events",
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
					log.Println("Message channel closed")
					return
				}
				c.handleTelegramMessage(d)
			case d, ok := <-taxiMsgs:
				if ok {
					c.handleTaxiEvent(d)
				}
			case <-c.ctx.Done():
				log.Println("Stopping consumer")
				return
			}
		}
	}()

	log.Println("Telegram consumer started")
	return nil
}

func (c *TelegramConsumer) Stop() {
	c.cancel()
	c.channel.Close()
	c.conn.Close()
}

func (c *TelegramConsumer) handleTelegramMessage(d amqp.Delivery) {
	message := string(d.Body)

	var payload struct {
		ChatID      int64           `json:"chat_id"`
		Text        string          `json:"text"`
		Phone       string          `json:"phone,omitempty"`
		Address     string          `json:"address,omitempty"`
		MapLink     string          `json:"map_link,omitempty"`
		ReplyMarkup json.RawMessage `json:"reply_markup,omitempty"`
	}

	if err := json.Unmarshal(d.Body, &payload); err != nil {
		log.Printf("Failed to unmarshal message: %v. Body: %s", err, message)
		d.Ack(false)
		return
	}

	if payload.ChatID == 0 {
		payload.ChatID = c.defaultChatID
	}

	if payload.ChatID == 0 {
		log.Printf("Skip message: ChatID 0. Body: %s", message)
		d.Ack(false)
		return
	}

	var err error
	log.Printf("Processing message for ChatID: %d", payload.ChatID)

	// Determine sending method
	// If ReplyMarkup is provided, we use standard SendMessage to respect the custom keyboard
	if len(payload.ReplyMarkup) > 0 {
		var markup interface{}
		var s string
		// Parse markup
		if err := json.Unmarshal(payload.ReplyMarkup, &s); err == nil {
			_ = json.Unmarshal([]byte(s), &markup)
		} else {
			_ = json.Unmarshal(payload.ReplyMarkup, &markup)
		}
		
		log.Printf("Sending message with custom markup: %+v", markup)
		err = c.client.SendMessage(payload.ChatID, payload.Text, markup)
	} else if payload.Phone != "" || payload.Address != "" {
		// Fallback to SendProfile if no custom markup but profile data exists
		err = c.client.SendProfile(payload.ChatID, payload.Text, payload.Phone, payload.Address, payload.MapLink)
	} else {
		// Just text
		err = c.client.SendMessage(payload.ChatID, payload.Text, nil)
	}

	if err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
	} else {
		d.Ack(false)
	}
}

func (c *TelegramConsumer) handleTaxiEvent(d amqp.Delivery) {
	log.Printf("Received taxi event: %s", string(d.Body))
	
	var base struct {
		ChatID int64 `json:"chat_id"`
	}
	if err := json.Unmarshal(d.Body, &base); err != nil {
		log.Printf("Failed to unmarshal into base: %v", err)
		d.Ack(false) 
		return
	}

	var est struct {
		OrderID   string  `json:"order_id"`
		Price     float64 `json:"price"`
		PayloadTo string  `json:"payload_to"`
	}
	// Basic check for price > 0 to identify estimate vs ordered result
	if err := json.Unmarshal(d.Body, &est); err == nil && est.Price > 0 {
		log.Printf("Processing estimate: OrderID=%s Price=%f PayloadTo=%s", est.OrderID, est.Price, est.PayloadTo)
		// Get short order ID (first 8 chars)
		shortID := est.OrderID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		text := fmt.Sprintf("🚕 Вартість таксі для #%s: %.0f ₴", shortID, est.Price)
		
		keyboard := map[string]interface{}{
			"inline_keyboard": [][]map[string]string{
				{
					{"text": "✅", "callback_data": fmt.Sprintf("confirm_taxi:%s", est.OrderID)},
					{"text": "❌", "callback_data": fmt.Sprintf("cancel_taxi:%s", est.OrderID)},
				},
			},
		}
		
		if err := c.client.SendMessage(base.ChatID, text, keyboard); err != nil {
			log.Printf("Failed to send estimate: %v", err)
		} else {
			d.Ack(false)
		}
		return
	}
	
	var ordered struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	if err := json.Unmarshal(d.Body, &ordered); err == nil && ordered.Status != "" {
		var text string
		if ordered.Status == "error" {
			text = "❌ " + ordered.Message
		} else {
			// Format: 🚕 Таксi для #shortID викликано
			shortID := ordered.OrderID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}
			text = fmt.Sprintf("🚕 Таксi для #%s викликано", shortID)
		}
		
		if err := c.client.SendMessage(base.ChatID, text, nil); err != nil {
			log.Printf("Failed to send ordered/status: %v", err)
		} else {
			d.Ack(false)
		}
		return
	}
	
	d.Ack(false)
}
