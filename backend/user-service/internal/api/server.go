package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/user-service/internal/api/handlers"
	"github.com/tonysanin/brobar/user-service/internal/services"
)

type Server struct {
	app         *fiber.App
	userService *services.UserService
	userHandler *handlers.UserHandler
}

func NewServer(
	userService *services.UserService,
) *Server {
	s := &Server{
		app: fiber.New(fiber.Config{
			AppName: "User Service",
		}),
		userService: userService,
	}

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	jwtSecret := []byte(helpers.GetEnv("JWT_SECRET", ""))

	s.userHandler = handlers.NewUserHandler(userService, jwtSecret)

	s.SetupRoutes()

	return s
}

func (s *Server) SetupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	auth := s.app.Group("/auth")

	auth.Post("/login", s.userHandler.Login)
	auth.Post("/register", s.userHandler.Register)
	auth.Post("/refresh", s.userHandler.Refresh)

	user := s.app.Group("/user")
	user.Get("/me", s.userHandler.GetUserMe)
	user.Get("/:id", s.userHandler.GetUserByID)
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
