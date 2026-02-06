package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/payment-service/internal/app"
	"github.com/tonysanin/brobar/payment-service/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.NewConfig()
	server := app.NewServer(cfg)

	log.Printf("Payment service starting on port %s", cfg.Port)
	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
