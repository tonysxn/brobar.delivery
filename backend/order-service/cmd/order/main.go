package main

import (
	"context"
	"log"
	"time"

	"github.com/tonysanin/brobar/order-service/internal/api"
	"github.com/tonysanin/brobar/order-service/internal/clients"
	"github.com/tonysanin/brobar/order-service/internal/consumer"
	"github.com/tonysanin/brobar/order-service/internal/repositories"
	"github.com/tonysanin/brobar/order-service/internal/services"
	"github.com/tonysanin/brobar/pkg/clients/payment"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/pkg/rabbitmq"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load(".env")

	port := helpers.GetEnv("ORDER_PORT", "3001")

	db, err := InitDatabase()
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
	paymentClient := payment.NewClient(helpers.GetEnv("PAYMENT_SERVICE_URL", "http://payment-service:8081"))

	// Initialize Message Broker
	producer := rabbitmq.NewProducer()
	defer producer.Close()

	// Initialize repositories
	orderRepository := repositories.NewOrderRepository(db)
	orderItemsRepository := repositories.NewOrderItemRepository(db)

	// Initialize services
	validationService := services.NewValidationService(productClient, webClient)
	orderService := services.NewOrderService(orderRepository, orderItemsRepository, productClient, paymentClient, validationService, producer)

	// Initialize Consumer
	rabbitURL := helpers.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	paymentConsumer, err := consumer.NewPaymentConsumer(rabbitURL, orderService)
	if err != nil {
		log.Fatalf("Failed to initialize payment consumer: %v", err)
	}

	if err := paymentConsumer.Start(); err != nil {
		log.Fatalf("Failed to start payment consumer: %v", err)
	}
	defer paymentConsumer.Stop()

	server := api.NewServer(orderService)

	log.Printf("Starting order service on :%s", port)
	if err := server.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func InitDatabase() (*sqlx.DB, error) {
	connStr := "postgres://" +
		helpers.GetEnv("DB_USER", "") + ":" +
		helpers.GetEnv("DB_PASSWORD", "") + "@" +
		helpers.GetEnv("DB_HOST", "") + ":" +
		helpers.GetEnv("DB_PORT", "") + "/" +
		helpers.GetEnv("DB_NAME", "") + "?sslmode=" +
		helpers.GetEnv("DB_SSLMODE", "disable")

	return sqlx.Connect("postgres", connStr)
}
