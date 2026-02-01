package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/tonysanin/brobar/file-service/internal/api/handlers"
	"github.com/tonysanin/brobar/file-service/internal/services"
	"github.com/tonysanin/brobar/pkg/response"
)

type Server struct {
	app         *fiber.App
	fileService *services.FileService
	fileHandler *handlers.FileHandler
}

func NewServer(fileService *services.FileService) *Server {
	s := &Server{
		app:         fiber.New(fiber.Config{AppName: "File Service", BodyLimit: 10 * 1024 * 1024}),
		fileService: fileService,
	}

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.fileHandler = handlers.NewFileHandler(fileService)

	s.SetupRoutes()

	return s
}

func (s *Server) SetupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	files := s.app.Group("/files")
	files.Post("/upload", s.fileHandler.UploadFile)
	files.Get("/:filename", s.fileHandler.GetFile)
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
