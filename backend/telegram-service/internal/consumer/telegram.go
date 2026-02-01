package consumer

import (
	"context"
	"encoding/json"
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
					ChatID      int64           `json:"chat_id"`
					Text        string          `json:"text"`
					Phone       string          `json:"phone,omitempty"`
					Address     string          `json:"address,omitempty"`
					MapLink     string          `json:"map_link,omitempty"`
					ReplyMarkup json.RawMessage `json:"reply_markup,omitempty"`
				}

				// Basic validation and parsing
				if err := json.Unmarshal(d.Body, &payload); err != nil {
					log.Printf("Failed to unmarshal message: %v. Body: %s", err, message)
					// Reject invalid messages
					d.Ack(false)
					continue
				}

				if payload.ChatID == 0 {
					payload.ChatID = c.defaultChatID
				}

				if payload.ChatID == 0 {
					log.Printf("Skip message: ChatID is still 0 after applying default. Body: %s", message)
					d.Ack(false)
					continue
				}

				var err error
				log.Printf("Processing message for ChatID: %d. Phone: %s, Address: %s, MapLink: %s", payload.ChatID, payload.Phone, payload.Address, payload.MapLink)

				// Check if this is a profile message (has phone or address)
				if payload.Phone != "" || payload.Address != "" {
					err = c.client.SendProfile(payload.ChatID, payload.Text, payload.Phone, payload.Address, payload.MapLink)
				} else {
					// Standard text message with optional markup
					var markup interface{}
					if len(payload.ReplyMarkup) > 0 {
						// Decode it back to interface or use RawMessage directly if client supports it
						// The Client struct uses interface{} for ReplyMarkup, so RawMessage is fine as it's a []byte
						// but Telegram API expects an object. We need to unmarshal the string/byte if it was a stringified JSON.

						// In order-service, keyboard is marshaled to string then put in map.
						// So RawMessage will be a JSON string like "\"{\\\"inline_keyboard\\\":...}\"".
						// We need to unquote it then unmarshal.
						var s string
						if err := json.Unmarshal(payload.ReplyMarkup, &s); err == nil {
							_ = json.Unmarshal([]byte(s), &markup)
						} else {
							_ = json.Unmarshal(payload.ReplyMarkup, &markup)
						}
					}
					err = c.client.SendMessage(payload.ChatID, payload.Text, markup)
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
