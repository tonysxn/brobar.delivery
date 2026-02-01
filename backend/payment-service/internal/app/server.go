package app

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/tonysanin/brobar/payment-service/internal/handlers"
	"github.com/tonysanin/brobar/payment-service/internal/provider"
	"github.com/tonysanin/brobar/payment-service/internal/rabbitmq"
	"github.com/tonysanin/brobar/payment-service/internal/services"
	"github.com/tonysanin/brobar/pkg/monobank"
)

type Server struct {
	app *fiber.App
}

func NewServer() *Server {
	s := &Server{
		app: fiber.New(fiber.Config{
			AppName: "Payment Service",
		}),
	}

	s.app.Use(logger.New())
	s.app.Use(cors.New())
	s.app.Use(compress.New())

	// Initialize dependencies here
	// 1. RabbitMQ Publisher
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	rabbitPublisher, err := rabbitmq.NewPublisher(rabbitURL)
	if err != nil {
		// Log error but don't crash, maybe retry or fail depending on strictness
		// For now simple log
		println("Failed to connect to RabbitMQ:", err.Error())
	}

	// 2. Monobank Client
	monoToken := os.Getenv("MONOBANK_TOKEN")
	monoClient := monobank.NewAcquiring(monoToken, os.Getenv("PUBLIC_DOMAIN"))

	// 3. Payment Provider (Strategy)
	// For now we only have Monobank, but interface allows more
	monoProvider := provider.NewMonobankProvider(monoClient)

	// 4. Service
	paymentService := services.NewPaymentService(monoProvider, rabbitPublisher)

	// 5. Handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Health Check
	s.app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Routes
	api := s.app.Group("/api/v1")
	payments := api.Group("/payments")
	payments.Post("/init", paymentHandler.InitPayment)

	webhooks := s.app.Group("/webhooks")
	webhooks.Post("/monobank", paymentHandler.HandleMonobankWebhook)

	return s
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}
