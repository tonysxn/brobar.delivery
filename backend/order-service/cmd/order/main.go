package main

import (
	"context"
	"log"
	"time"

	"github.com/tonysanin/brobar/order-service/internal/api"
	"github.com/tonysanin/brobar/order-service/internal/clients"
	"github.com/tonysanin/brobar/order-service/internal/repositories"
	"github.com/tonysanin/brobar/order-service/internal/services"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/pkg/rabbitmq"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found or failed to load")
	}

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

	// Initialize Message Broker
	producer := rabbitmq.NewProducer()
	defer producer.Close()

	// Initialize repositories
	orderRepository := repositories.NewOrderRepository(db)
	orderItemsRepository := repositories.NewOrderItemRepository(db)

	// Initialize services
	validationService := services.NewValidationService(productClient, webClient)
	orderService := services.NewOrderService(orderRepository, orderItemsRepository, productClient, validationService, producer)

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
