package rabbitmq

type QueueName string

const (
	QueueTelegram QueueName = "telegram_messages"
	QueueSyrve    QueueName = "syrve_orders"
)
