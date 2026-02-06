package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/pkg/rabbitmq"
	"github.com/tonysanin/brobar/pkg/syrve"
	"github.com/tonysanin/brobar/syrve-service/internal/api"
	"github.com/tonysanin/brobar/syrve-service/internal/config"
	"github.com/tonysanin/brobar/syrve-service/internal/consumer"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.NewConfig()

	client := syrve.NewClient(cfg.APILogin, cfg.OrganizationID).WithTimeout(10)

	// Initialize Producer
	producer := rabbitmq.NewProducer()
	defer producer.Close()

	server := api.NewServer(client, producer, cfg.OrderServiceURL)

	// Start Consumer
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}
	
	cons := consumer.NewConsumer(client, producer)
	go cons.Start(rabbitMQURL)

	log.Printf("Starting server on :%s", cfg.Port)
	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
