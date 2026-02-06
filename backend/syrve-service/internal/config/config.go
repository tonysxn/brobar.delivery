package config

import (
	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	Port           string
	APILogin       string
	OrganizationID string
	OrderServiceURL string
}

func NewConfig() *Config {
	return &Config{
		Port:            helpers.GetEnv("SYRVE_PORT", "3011"),
		APILogin:        helpers.GetEnv("SYRVE_TOKEN", ""),
		OrganizationID:  helpers.GetEnv("SYRVE_ORGANIZATION", ""),
		OrderServiceURL: helpers.GetEnv("ORDER_SERVICE_URL", "http://order-service-dev:3003"),
	}
}
