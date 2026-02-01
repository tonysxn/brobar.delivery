package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/payment-service/internal/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	server := app.NewServer()
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Payment service starting on port %s", port)
	if err := server.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
