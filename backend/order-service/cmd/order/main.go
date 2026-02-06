package main

import (
	"context"
	"log"
	"time"

	"github.com/tonysanin/brobar/order-service/internal/api"
	"github.com/tonysanin/brobar/order-service/internal/clients"
	"github.com/tonysanin/brobar/order-service/internal/config"
	"github.com/tonysanin/brobar/order-service/internal/consumer"
	"github.com/tonysanin/brobar/order-service/internal/repositories"
	"github.com/tonysanin/brobar/order-service/internal/services"
	"github.com/tonysanin/brobar/pkg/clients/payment"
	"github.com/tonysanin/brobar/pkg/rabbitmq"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.NewConfig()

	db, err := InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Close database failed: %v", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// Initialize clients
	productClient := clients.NewProductClient()
	webClient := clients.NewWebClient()
	paymentClient := payment.NewClient(cfg.PaymentServiceURL)

	// Initialize Message Broker
	producer := rabbitmq.NewProducer()
	defer producer.Close()

	// Initialize repositories
	orderRepository := repositories.NewOrderRepository(db)
	orderItemsRepository := repositories.NewOrderItemRepository(db)

	// Initialize services
	validationService := services.NewValidationService(productClient, webClient)
	orderService := services.NewOrderService(orderRepository, orderItemsRepository, productClient, paymentClient, validationService, producer, cfg.AppTimezone)

	// Initialize Consumer
	paymentConsumer, err := consumer.NewPaymentConsumer(cfg.RabbitMQURL, orderService)
	if err != nil {
		log.Fatalf("Failed to initialize payment consumer: %v", err)
	}

	if err := paymentConsumer.Start(); err != nil {
		log.Fatalf("Failed to start payment consumer: %v", err)
	}
	defer paymentConsumer.Stop()

	server := api.NewServer(orderService)

	log.Printf("Starting order service on :%s", cfg.Port)
	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", cfg.GetDatabaseURL())
}
