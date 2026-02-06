package config

import (
	"fmt"

	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	Port              string
	JWTSecret         []byte
	UserServiceURL    string
	ProductServiceURL string
	TelegramServiceURL string
	SyrveServiceURL   string
	FileServiceURL    string
	WebServiceURL     string
	PaymentServiceURL string
	OrderServiceURL   string
}

func NewConfig() *Config {
	return &Config{
		Port:      helpers.GetEnv("GATEWAY_PORT", "8000"),
		JWTSecret: []byte(helpers.GetEnv("JWT_SECRET", "")),
		
		UserServiceURL: buildServiceURL(
			helpers.GetEnv("USER_HOST", "http://user-service-dev"),
			helpers.GetEnv("USER_PORT", "3001"),
		),
		ProductServiceURL: buildServiceURL(
			helpers.GetEnv("PRODUCT_HOST", "http://product-service-dev"),
			helpers.GetEnv("PRODUCT_PORT", "3000"),
		),
		TelegramServiceURL: buildServiceURL(
			helpers.GetEnv("TELEGRAM_HOST", "http://telegram-service-dev"),
			helpers.GetEnv("TELEGRAM_PORT", "3010"),
		),
		SyrveServiceURL: buildServiceURL(
			helpers.GetEnv("SYRVE_HOST", "http://syrve-service-dev"),
			helpers.GetEnv("SYRVE_PORT", "3010"),
		),
		FileServiceURL: buildServiceURL(
			helpers.GetEnv("FILE_HOST", "http://file-service-dev"),
			helpers.GetEnv("FILE_PORT", "3001"),
		),
		WebServiceURL: buildServiceURL(
			helpers.GetEnv("WEB_HOST", "http://web-service-dev"),
			helpers.GetEnv("WEB_PORT", "3006"),
		),
		PaymentServiceURL: buildServiceURL(
			helpers.GetEnv("PAYMENT_SERVICE_HOST", "http://payment-service-dev"),
			helpers.GetEnv("PAYMENT_SERVICE_PORT", "8081"),
		),
		OrderServiceURL: buildServiceURL(
			helpers.GetEnv("ORDER_HOST", "http://order-service-dev"),
			helpers.GetEnv("ORDER_PORT", "3001"),
		),
	}
}

func buildServiceURL(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
