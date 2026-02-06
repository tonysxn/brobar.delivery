package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/tonysanin/brobar/pkg/rabbitmq"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/pkg/syrve"
	"github.com/tonysanin/brobar/syrve-service/internal/api/handlers"
)

type Server struct {
	app             *fiber.App
	syrveClient     *syrve.Client
	producer        *rabbitmq.Producer
	orderServiceURL string
	syrveHandler    *handlers.SyrveHandler
}

func NewServer(
	syrveClient *syrve.Client,
	producer *rabbitmq.Producer,
	orderServiceURL string,
) *Server {
	s := &Server{
		app: fiber.New(fiber.Config{
			AppName: "Syrve Service",
		}),
		syrveClient:     syrveClient,
		producer:        producer,
		orderServiceURL: orderServiceURL,
	}

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.syrveHandler = handlers.NewSyrveHandler(syrveClient, producer, orderServiceURL)

	s.SetupRoutes()

	return s
}

func (s *Server) SetupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	syrveGroup := s.app.Group("/syrve")
	syrveGroup.Get("/products", s.syrveHandler.GetProducts)
	syrveGroup.Get("/stop-lists", s.syrveHandler.GetStopLists)
	syrveGroup.Post("/stop-lists/sync", s.syrveHandler.SyncStopLists)

	// Webhook endpoint
	// Gateway proxies /webhooks/syrve -> /webhooks/syrve if configured simply
	// But in existing logic:
	// gateway: telegramGroup.Post("/", s.ProxyToTelegramService) where ProxyToTelegramService -> service + url
	// So if gateway calls telegram-service/webhooks/telegram, then route must match.
	// But here I'm adding `syrveGroup` which is /syrve
	// I'll add a webhook specific group or route.
	
	webhookGroup := s.app.Group("/webhooks")
	webhookGroup.Post("/syrve", s.syrveHandler.HandleWebhook)
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
