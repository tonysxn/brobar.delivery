package config

import (
	"fmt"

	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	Port              string
	RabbitMQURL       string
	PaymentServiceURL string
	DBUser            string
	DBPassword        string
	DBHost            string
	DBPort            string
	DBName            string
	DBSSLMode         string
	AppTimezone       string
}

func NewConfig() *Config {
	return &Config{
		Port:              helpers.GetEnv("ORDER_PORT", "3001"),
		RabbitMQURL:       helpers.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		PaymentServiceURL: helpers.GetEnv("PAYMENT_SERVICE_URL", "http://payment-service:8081"),
		DBUser:            helpers.GetEnv("DB_USER", ""),
		DBPassword:        helpers.GetEnv("DB_PASSWORD", ""),
		DBHost:            helpers.GetEnv("DB_HOST", ""),
		DBPort:            helpers.GetEnv("DB_PORT", ""),
		DBName:            helpers.GetEnv("DB_NAME", ""),
		DBSSLMode:         helpers.GetEnv("DB_SSLMODE", "disable"),
		AppTimezone:       helpers.GetEnv("APP_TIMEZONE", "Europe/Kyiv"),
	}
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}
