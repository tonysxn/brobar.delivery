package config

import (
	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	Port          string
	RabbitMQURL   string
	MonobankToken string
	PublicDomain  string
}

func NewConfig() *Config {
	return &Config{
		Port:          helpers.GetEnv("PAYMENT_SERVICE_PORT", "8081"),
		RabbitMQURL:   helpers.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		MonobankToken: helpers.GetEnv("MONOBANK_TOKEN", ""),
		PublicDomain:  helpers.GetEnv("PUBLIC_DOMAIN", ""),
	}
}
