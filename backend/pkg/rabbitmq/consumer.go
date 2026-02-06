package rabbitmq

import "github.com/streadway/amqp"

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *Consumer) Consume(queueName QueueName) (<-chan amqp.Delivery, error) {
	q, err := c.channel.QueueDeclare(
		string(queueName),
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return c.channel.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (c *Consumer) Close() {
	_ = c.channel.Close()
	_ = c.conn.Close()
}
