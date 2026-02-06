package config

import (
	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	RabbitMQURL       string
	OntaxiToken       string
	OntaxiBusinessID  string
	OntaxiClientID    string
	OntaxiPaymentMethodID string
	OntaxiLocalCoords string
	OntaxiLocalPlace  string
	OntaxiBaseURL     string
}

func NewConfig() *Config {
	return &Config{
		RabbitMQURL:       helpers.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		OntaxiToken:       helpers.GetEnv("ONTAXI_TOKEN", ""),
		OntaxiBusinessID:  helpers.GetEnv("ONTAXI_BUSINESS_ID", ""),
		OntaxiClientID:    helpers.GetEnv("ONTAXI_CLIENT_ID", ""),
		OntaxiPaymentMethodID: helpers.GetEnv("ONTAXI_PAYMENT_METHOD_ID", ""),
		OntaxiLocalCoords: helpers.GetEnv("ONTAXI_LOCAL_COORDS", ""),
		OntaxiLocalPlace:  helpers.GetEnv("ONTAXI_LOCAL_PLACE", ""),
		OntaxiBaseURL:     helpers.GetEnv("ONTAXI_BASE_URL", "https://business.ontaxi.com.ua/api/"),
	}
}
