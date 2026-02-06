package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/gateway-service/internal/api"
	"github.com/tonysanin/brobar/gateway-service/internal/config"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.NewConfig()

	server := api.NewServer(api.ServerConfig{
		UserServiceURL:    cfg.UserServiceURL,
		ProductServiceURL: cfg.ProductServiceURL,
		TelegramServiceURL: cfg.TelegramServiceURL,
		SyrveServiceURL:   cfg.SyrveServiceURL,
		FileServiceURL:    cfg.FileServiceURL,
		WebServiceURL:     cfg.WebServiceURL,
		PaymentServiceURL: cfg.PaymentServiceURL,
		OrderServiceURL:   cfg.OrderServiceURL,
		JWTSecret:         cfg.JWTSecret,
	})

	log.Printf("Starting server on :%s", cfg.Port)
	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
