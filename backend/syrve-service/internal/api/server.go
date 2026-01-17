package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/syrve-service/internal/api/handlers"
	"github.com/tonysanin/brobar/syrve-service/internal/services/syrve"
)

type Server struct {
	app          *fiber.App
	syrveClient  *syrve.Client
	syrveHandler *handlers.SyrveHandler
}

func NewServer(
	syrveClient *syrve.Client,
) *Server {
	s := &Server{
		app: fiber.New(fiber.Config{
			AppName: "Syrve Service",
		}),
		syrveClient: syrveClient,
	}

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.syrveHandler = handlers.NewSyrveHandler(syrveClient)

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
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
