package rabbitmq

import (
	"fmt"
	"github.com/tonysanin/brobar/pkg/helpers"
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer() *Producer {
	user := helpers.GetEnv("RABBITMQ_USER", "guest")
	pass := helpers.GetEnv("RABBITMQ_PASS", "guest")
	host := helpers.GetEnv("RABBITMQ_HOST", "localhost")
	port := helpers.GetEnv("RABBITMQ_PORT", "5672")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	return &Producer{
		conn:    conn,
		channel: ch,
	}
}

func (p *Producer) SendMessage(queueName QueueName, body string) error {
	_, err := p.channel.QueueDeclare(
		string(queueName),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		string(queueName),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
}

func (p *Producer) Close() {
	_ = p.channel.Close()
	_ = p.conn.Close()
}
