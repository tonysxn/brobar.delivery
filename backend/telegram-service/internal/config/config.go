package config

import (
	"strconv"

	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	Port        string
	RabbitMQURL string
	BotToken    string
	ChatID      int64
	SyrveServiceURL   string
	ProductServiceURL string
	WebServiceURL     string
}

func NewConfig() *Config {
	chatIDStr := helpers.GetEnv("TELEGRAM_CHAT_ID", "0")
	chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

	return &Config{
		Port:        helpers.GetEnv("TELEGRAM_PORT", "8080"),
		RabbitMQURL: helpers.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		BotToken:    helpers.GetEnv("TELEGRAM_BOT_TOKEN", ""),
		ChatID:      chatID,
		SyrveServiceURL:   helpers.GetEnv("SYRVE_SERVICE_URL", "http://syrve-service:3004"),
		ProductServiceURL: helpers.GetEnv("PRODUCT_SERVICE_URL", "http://product-service:3000"),
		WebServiceURL:     helpers.GetEnv("WEB_SERVICE_URL", "http://web-service:3006"),
	}
}
