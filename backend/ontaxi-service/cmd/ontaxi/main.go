package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/ontaxi-service/internal/config"
	"github.com/tonysanin/brobar/ontaxi-service/internal/service"
	"github.com/tonysanin/brobar/ontaxi-service/internal/transport/rabbitmq" // Local package
	pkgRabbit "github.com/tonysanin/brobar/pkg/rabbitmq" // Shared package
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.NewConfig()

	// Services
	ontaxiSvc := service.NewOntaxiService(cfg)
	orderClient := service.NewOrderClient()

	// RabbitMQ Producer
	producer := pkgRabbit.NewProducer()
	defer producer.Close()

	// RabbitMQ Consumer
	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQURL, ontaxiSvc, orderClient, producer)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Ontaxi Service Started")

	// Graceful Shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down ontaxi-service...")
	// consumer.Stop() // If implemented
}
