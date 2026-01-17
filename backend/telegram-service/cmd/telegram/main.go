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
	"github.com/tonysanin/brobar/telegram-service/internal/consumer"
	"github.com/tonysanin/brobar/telegram-service/internal/tdlib"
	"github.com/zelenin/go-tdlib/client"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found or failed to load")
	}

	rabbitURL := helpers.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	botToken := helpers.GetEnv("TELEGRAM_BOT_TOKEN", "")
	apiIDStr := helpers.GetEnv("TELEGRAM_API_ID", "0")
	apiHash := helpers.GetEnv("TELEGRAM_API_HASH", "placeholder")

	if rabbitURL == "" || botToken == "" {
		log.Fatal("One or more required env variables are missing: RABBITMQ_URL, TELEGRAM_BOT_TOKEN")
	}

	apiID, _ := strconv.Atoi(apiIDStr)
	// Initialize TDLib client
	tdClient, err := tdlib.NewClient(int32(apiID), apiHash, botToken)
	if err != nil {
		log.Fatalf("Failed to create TDLib client: %v", err)
	}

	// Initialize RabbitMQ Consumer
	telegramConsumer, err := consumer.NewTelegramConsumer(rabbitURL, tdClient)
	if err != nil {
		log.Fatalf("Failed to create Telegram consumer: %v", err)
	}

	if err := telegramConsumer.Start(); err != nil {
		log.Fatalf("Failed to start Telegram consumer: %v", err)
	}

	// Handle Telegram Updates
	updateCh := make(chan *client.UpdateNewMessage, 100)
	tdClient.AddUpdateListener(updateCh)

	go func() {
		for update := range updateCh {
			msg := update.Message
			chatID := msg.ChatId

			switch content := msg.Content.(type) {
			case *client.MessageText:
				text := content.Text.Text
				if text == "/start" {
					_, err := tdClient.SendMessage(chatID, "Привет!", nil)
					if err != nil {
						log.Printf("Failed to send welcome message: %v", err)
					}
				}
			}
		}
	}()

	// Start Healthcheck Server
	app := fiber.New()
	app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{"status": "ok"})
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
