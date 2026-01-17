package services

import (
	"github.com/tonysanin/brobar/pkg/rabbitmq"
	"log"
)

func SendTelegramMessage(text string) {
	producer := rabbitmq.NewProducer()
	defer producer.Close()

	if err := producer.SendMessage(rabbitmq.QueueTelegram, text); err != nil {
		log.Printf("Failed to send message to RabbitMQ: %v", err)
		return
	}

	log.Println("Message successfully sent to telegram_messages")
}
