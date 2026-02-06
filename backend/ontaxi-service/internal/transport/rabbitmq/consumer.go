package rabbitmq

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
	"github.com/tonysanin/brobar/ontaxi-service/internal/service"
	"github.com/tonysanin/brobar/pkg/rabbitmq"
)

const (
	QueueTaxiRequests = "taxi_requests"
	QueueTaxiConfirms = "taxi_confirms"
	
	// Destination queues (handled by telegram service)
	QueueTelegramMessages = "telegram_messages" // Existing queue
	QueueTaxiEvents       = "taxi_events"       // New queue for specific taxi events if needed, but telegram consumer reads telegram_messages
)

// We likely want to send formatted messages directly to telegram_messages queue
// OR we send an event to 'taxi_events' and have TelegramService format it.
// The plan said "Publish taxi_estimate_ready". Telegram Service consumes it.
// So let's use "taxi_events" queue for responses to Telegram Service.

type Consumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	ontaxiService *service.OntaxiService
	orderClient   *service.OrderClient
	producer      *rabbitmq.Producer
	rabbitURL     string
}

func NewConsumer(rabbitURL string, ontaxiService *service.OntaxiService, orderClient *service.OrderClient, producer *rabbitmq.Producer) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:          conn,
		channel:       ch,
		ontaxiService: ontaxiService,
		orderClient:   orderClient,
		producer:      producer,
		rabbitURL:     rabbitURL,
	}, nil
}

func (c *Consumer) Start() error {
	// 1. Consumer for Requests
	if err := c.consume(QueueTaxiRequests, c.handleTaxiRequest); err != nil {
		return err
	}
	
	// 2. Consumer for Confirms
	if err := c.consume(QueueTaxiConfirms, c.handleTaxiConfirm); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) consume(queueName string, handler func([]byte) error) error {
	_, err := c.channel.QueueDeclare(
		queueName,
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queueName, "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message from %s: %v", queueName, err)
				// Determine if we should ack/nack. Simple ack for now to avoid loops
			}
			d.Ack(false)
		}
	}()
	return nil
}

func (c *Consumer) handleTaxiRequest(body []byte) error {
	var req TaxiRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return err
	}

	log.Printf("Received taxi request for order %s", req.OrderID)

	// 1. Get Order
	order, err := c.orderClient.GetOrder(req.OrderID)
	if err != nil {
		log.Printf("Failed to get order details: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Не вдалося отримати дані замовлення")
		return nil
	}

	// 2. Resolve destination payload using reverse geocoding (like PHP getPayloadByCoords)
	if order.Coords == "" {
		log.Printf("Coords are missing for order %s", req.OrderID)
		c.sendError(req.ChatID, req.OrderID, "Не вдалося визначити координати для замовлення")
		return nil
	}

	// Support both comma and semicolon as delimiters
	coordsStr := strings.ReplaceAll(order.Coords, ";", ",")
	parts := strings.Split(coordsStr, ",")
	if len(parts) != 2 {
		log.Printf("Invalid coords format: %s", order.Coords)
		c.sendError(req.ChatID, req.OrderID, "Невірний формат координат")
		return nil
	}
	
	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Printf("Failed to parse lat: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Невірний формат координат")
		return nil
	}
	
	lng, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Printf("Failed to parse lng: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Невірний формат координат")
		return nil
	}
	
	payloadTo, err := c.ontaxiService.GetPayloadByCoords(lat, lng)
	if err != nil {
		log.Printf("Failed to resolve destination: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Не вдалося визначити адресу доставки в OnTaxi")
		return nil
	}

	// 3. Get Estimate
	price, err := c.ontaxiService.GetDeliveryEstimate(lat, lng)
     if err != nil {
         log.Printf("Failed to get estimate: %v", err)
         c.sendError(req.ChatID, req.OrderID, "Не вдалося розрахувати вартість таксі: "+err.Error())
         return nil
     }
     
     // 4. Send Estimate Ready Event
     evt := TaxiEstimateReady{
         ChatID:    req.ChatID,
         OrderID:   req.OrderID,
         Price:     price,
         PayloadTo: payloadTo,
     }
     
     log.Printf("Sending estimate ready event: %+v", evt)
     msg, _ := json.Marshal(evt)
     return c.producer.SendMessage("taxi_events", string(msg))
}

func (c *Consumer) handleTaxiConfirm(body []byte) error {
	var req TaxiConfirm
	if err := json.Unmarshal(body, &req); err != nil {
		return err
	}

	log.Printf("Received taxi confirm for order %s", req.OrderID)
	
	order, err := c.orderClient.GetOrder(req.OrderID)
	if err != nil {
        c.sendError(req.ChatID, req.OrderID, "Err fetching order")
		return nil
	}

	// 2. Resolve destination payload using reverse geocoding (like PHP getPayloadByCoords)
	if order.Coords == "" {
		log.Printf("Coords are missing for order %s (confirm)", req.OrderID)
		c.sendError(req.ChatID, req.OrderID, "Не вдалося визначити координати для замовлення")
		return nil
	}

	// Support both comma and semicolon as delimiters
	coordsStr := strings.ReplaceAll(order.Coords, ";", ",")
	parts := strings.Split(coordsStr, ",")
	if len(parts) != 2 {
		log.Printf("Invalid coords format for confirm: %s", order.Coords)
		c.sendError(req.ChatID, req.OrderID, "Невірний формат координат")
		return nil
	}
	
	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Printf("Failed to parse lat for confirm: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Невірний формат координат")
		return nil
	}
	
	lng, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Printf("Failed to parse lng for confirm: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Невірний формат координат")
		return nil
	}
	
	payloadTo, err := c.ontaxiService.GetPayloadByCoords(lat, lng)
	if err != nil {
		log.Printf("Failed to resolve destination for confirm: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Не вдалося визначити адресу доставки")
		return nil
	}

	// Build comment string matching PHP implementation
	comment := ""
	
	if order.Entrance != "" {
		comment += "Під'їзд: " + order.Entrance + ". "
	}
	
	if order.DeliveryDoor {
		comment += "Доставка до дверей! "
		
		if order.Floor != "" {
			comment += "Поверх: " + order.Floor + ". "
		}
		
		if order.Flat != "" {
			comment += "Квартира: " + order.Flat + ". "
		}
	}
	
	comment += "Після приїзду набрати ПАСАЖИРА. "
	
	if order.AddressWishes != "" {
		comment += "Коментарiй клієнта за адресою: " + order.AddressWishes
	}

	resp, err := c.ontaxiService.CreateOrder(payloadTo, order.Phone, order.Name, comment, order.Entrance, order.DeliveryDoor)
	if err != nil {
		log.Printf("Failed to create ontaxi order: %v", err)
		c.sendError(req.ChatID, req.OrderID, "Помилка створення замовлення OnTaxi: "+err.Error())
		return nil
	}
	
	log.Printf("Ontaxi order created response: %s", resp)

	evt := TaxiOrdered{
		ChatID:  req.ChatID,
		OrderID: req.OrderID,
		Status:  "success",
		Message: "Таксі успішно викликано!",
	}

	msg, _ := json.Marshal(evt)
	return c.producer.SendMessage("taxi_events", string(msg))
}

func (c *Consumer) sendError(chatID int64, orderID, message string) {
    evt := TaxiOrdered{
        ChatID: chatID,
        OrderID: orderID,
        Status: "error",
        Message: message,
    }
    msg, _ := json.Marshal(evt)
    _ = c.producer.SendMessage("taxi_events", string(msg))
}
