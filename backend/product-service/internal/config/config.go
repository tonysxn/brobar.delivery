package config

import (
	"fmt"

	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSSLMode  string
	RabbitMQURL string
}

func NewConfig() *Config {
	return &Config{
		Port:       helpers.GetEnv("PRODUCT_PORT", "1"),
		DBUser:     helpers.GetEnv("DB_USER", ""),
		DBPassword: helpers.GetEnv("DB_PASSWORD", ""),
		DBHost:     helpers.GetEnv("DB_HOST", ""),
		DBPort:     helpers.GetEnv("DB_PORT", ""),
		DBName:     helpers.GetEnv("DB_NAME", ""),
		DBSSLMode:  helpers.GetEnv("DB_SSLMODE", "disable"),
		RabbitMQURL: helpers.GetEnv("RABBITMQ_URL", "amqp://user:password@rabbitmq:5672/"),
	}
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}
