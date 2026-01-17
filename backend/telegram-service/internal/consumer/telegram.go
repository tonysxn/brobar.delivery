package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/tonysanin/brobar/telegram-service/internal/tdlib"
)

type TelegramConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	client  *tdlib.Client
	queue   string
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewTelegramConsumer(rabbitURL string, client *tdlib.Client) (*TelegramConsumer, error) {
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
		conn:    conn,
		channel: ch,
		client:  client,
		queue:   "telegram_messages",
		ctx:     ctx,
		cancel:  cancel,
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
		false, // autoAck false, чтобы подтверждать вручную
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

				message := string(d.Body)

				// Try to unmarshal as JSON first
				var payload struct {
					ChatID  int64  `json:"chat_id"`
					Text    string `json:"text"`
					Phone   string `json:"phone,omitempty"`
					Address string `json:"address,omitempty"`
					MapLink string `json:"map_link,omitempty"`
				}

				// Basic validation and parsing
				if err := json.Unmarshal(d.Body, &payload); err != nil {
					log.Printf("Failed to unmarshal message: %v. Body: %s", err, message)
					// Reject invalid messages
					d.Ack(false)
					continue
				}

				var err error
				log.Printf("Processing message for ChatID: %d. Phone: %s, Address: %s, MapLink: %s", payload.ChatID, payload.Phone, payload.Address, payload.MapLink)

				// Check if this is a profile message (has phone or address)
				if payload.Phone != "" || payload.Address != "" {
					_, err = c.client.SendProfile(payload.ChatID, payload.Text, payload.Phone, payload.Address, payload.MapLink)
				} else {
					// Standard text message
					_, err = c.client.SendMessage(payload.ChatID, payload.Text, nil)
				}

				if err != nil {
					log.Printf("Failed to send Telegram message: %v", err)
					// Retry logic could be handled here, but for now we just log
				} else {
					d.Ack(false)
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
