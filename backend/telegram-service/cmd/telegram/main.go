package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/pkg/telegram"
	"github.com/tonysanin/brobar/telegram-service/internal/consumer"
)

func main() {
	_ = godotenv.Load(".env")

	rabbitURL := helpers.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	botToken := helpers.GetEnv("TELEGRAM_BOT_TOKEN", "")

	if rabbitURL == "" || botToken == "" {
		log.Fatal("One or more required env variables are missing: RABBITMQ_URL, TELEGRAM_BOT_TOKEN")
	}

	// Initialize Telegram HTTP Client
	tgClient := telegram.NewClient(botToken)

	defaultChatIDStr := helpers.GetEnv("TELEGRAM_CHAT_ID", "0")
	defaultChatID, _ := strconv.ParseInt(defaultChatIDStr, 10, 64)

	// Initialize RabbitMQ Consumer
	telegramConsumer, err := consumer.NewTelegramConsumer(rabbitURL, tgClient, defaultChatID)
	if err != nil {
		log.Fatalf("Failed to create Telegram consumer: %v", err)
	}

	if err := telegramConsumer.Start(); err != nil {
		log.Fatalf("Failed to start Telegram consumer: %v", err)
	}

	// Start Healthcheck Server
	app := fiber.New()
	app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{"status": "ok"})
	})

	app.Post("/", func(c fiber.Ctx) error {
		log.Printf("Received Telegram Update: %s", string(c.Body()))
		return c.SendStatus(200)
	})

	port := helpers.GetEnv("TELEGRAM_PORT", "8080")
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Fiber server failed: %v", err)
		}
	}()

	log.Printf("Telegram service started on port %s", port)

	// Graceful Shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down telegram-service...")
	telegramConsumer.Stop()
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber server: %v", err)
	}
}
