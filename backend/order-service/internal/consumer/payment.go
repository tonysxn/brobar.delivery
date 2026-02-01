package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/tonysanin/brobar/order-service/internal/services"
)

type PaymentConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	service *services.OrderService
	queue   string
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewPaymentConsumer(rabbitURL string, service *services.OrderService) (*PaymentConsumer, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare queue to ensure it exists
	_, err = ch.QueueDeclare(
		"payment_events", // name matching payment-service publisher
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &PaymentConsumer{
		conn:    conn,
		channel: ch,
		service: service,
		queue:   "payment_events",
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (c *PaymentConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false, // autoAck false
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
					log.Println("Payment consumer channel closed")
					return
				}

				var event services.PaymentSuccessEvent
				if err := json.Unmarshal(d.Body, &event); err != nil {
					log.Printf("Failed to unmarshal payment event: %v. Body: %s", err, d.Body)
					d.Ack(false)
					continue
				}

				if event.Status == "success" {
					if err := c.service.ProcessPaymentSuccess(event); err != nil {
						log.Printf("Failed to process payment success: %v", err)
						// Maybe Nack/Requeue? For now Ack to avoid loop
						d.Ack(false)
					} else {
						d.Ack(false)
					}
				} else {
					// Ignore other statuses
					d.Ack(false)
				}

			case <-c.ctx.Done():
				return
			}
		}
	}()

	log.Println("Payment consumer started")
	return nil
}

func (c *PaymentConsumer) Stop() {
	c.cancel()
	c.channel.Close()
	c.conn.Close()
}
