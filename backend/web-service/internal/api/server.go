package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/web-service/internal/api/handlers"
	"github.com/tonysanin/brobar/web-service/internal/config"
	"github.com/tonysanin/brobar/web-service/internal/services"
)

type Server struct {
	app            *fiber.App
	settingHandler *handlers.SettingHandler
	reviewHandler  *handlers.ReviewHandler
	timeHandler    *handlers.TimeHandler
}

func NewServer(cfg *config.Config, settingService *services.SettingService, reviewService *services.ReviewService) *Server {
	s := &Server{
		app: fiber.New(fiber.Config{
			AppName: "Web Service",
		}),
		settingHandler: handlers.NewSettingHandler(settingService),
		reviewHandler:  handlers.NewReviewHandler(reviewService),
		timeHandler:    handlers.NewTimeHandler(cfg.AppTimezone),
	}

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.SetupRoutes()

	return s
}

func (s *Server) SetupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	// Server time (Kiev timezone)
	s.app.Get("/time", s.timeHandler.GetServerTime)

	settings := s.app.Group("/settings")
	settings.Get("/", s.settingHandler.GetSettings)
	settings.Put("/:key", s.settingHandler.UpdateSetting)

	reviews := s.app.Group("/reviews")
	reviews.Post("/", s.reviewHandler.CreateReview)
	reviews.Get("/", s.reviewHandler.GetReviews)
	reviews.Get("/:id", s.reviewHandler.GetReview)
	reviews.Delete("/:id", s.reviewHandler.DeleteReview)
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
