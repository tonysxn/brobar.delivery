package services

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/pkg/rabbitmq"
)

func SendTelegramMessage(text string) {
	producer := rabbitmq.NewProducer()
	defer producer.Close()

	chatIDStr := helpers.GetEnv("TELEGRAM_CHAT_ID", "0")
	chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	jsonBody, _ := json.Marshal(payload)

	if err := producer.SendMessage(rabbitmq.QueueTelegram, string(jsonBody)); err != nil {
		log.Printf("Failed to send message to RabbitMQ: %v", err)
		return
	}

	log.Println("Message successfully sent to telegram_messages")
}
