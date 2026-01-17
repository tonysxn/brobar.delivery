package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/gateway-service/internal/api"
	"github.com/tonysanin/brobar/pkg/helpers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found or failed to load")
	}

	port := helpers.GetEnv("GATEWAY_PORT", "8000")

	server := api.NewServer(api.ServerConfig{
		UserServiceURL: helpers.GetEnv("USER_HOST", "http://user-service-dev") +
			":" +
			helpers.GetEnv("USER_PORT", "3001"),
		ProductServiceURL: helpers.GetEnv("PRODUCT_HOST", "http://product-service-dev") +
			":" +
			helpers.GetEnv("PRODUCT_PORT", "3000"),
		TelegramServiceURL: helpers.GetEnv("TELEGRAM_HOST", "http://telegram-service-dev") +
			":" +
			helpers.GetEnv("TELEGRAM_PORT", "3010"),
		SyrveServiceURL: helpers.GetEnv("SYRVE_HOST", "http://syrve-service-dev") +
			":" +
			helpers.GetEnv("SYRVE_PORT", "3010"),
		FileServiceURL: helpers.GetEnv("FILE_HOST", "http://file-service-dev") +
			":" +
			helpers.GetEnv("FILE_PORT", "3001"),
		WebServiceURL: helpers.GetEnv("WEB_HOST", "http://web-service-dev") +
			":" +
			helpers.GetEnv("WEB_PORT", "3006"),
		OrderServiceURL: helpers.GetEnv("ORDER_HOST", "http://order-service-dev") +
			":" +
			helpers.GetEnv("ORDER_PORT", "3001"),

		JWTSecret: []byte(helpers.GetEnv("JWT_SECRET", "")),
	})

	log.Printf("Starting server on :%s", port)
	if err := server.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
