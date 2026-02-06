package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/tonysanin/brobar/order-service/internal/api/handlers"
	"github.com/tonysanin/brobar/order-service/internal/services"
	"github.com/tonysanin/brobar/pkg/response"
)

type Server struct {
	app          *fiber.App
	orderService *services.OrderService
	orderHandler *handlers.OrderHandler
}

func NewServer(
	orderService *services.OrderService,
) *Server {
	s := &Server{
		app: fiber.New(fiber.Config{
			AppName: "Order Service",
		}),
		orderService: orderService,
	}

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.orderHandler = handlers.NewOrderHandler(orderService)

	s.SetupRoutes()

	return s
}

func (s *Server) SetupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	orderGroup := s.app.Group("/orders")
	orderGroup.Get("/", s.orderHandler.GetOrders)
	orderGroup.Get("/:id", s.orderHandler.GetOrder)
	orderGroup.Post("/", s.orderHandler.CreateOrder)
	orderGroup.Put("/:id", s.orderHandler.UpdateOrder)
	orderGroup.Delete("/:id", s.orderHandler.DeleteOrder)
	orderGroup.Post("/:id/syrve-notified", s.orderHandler.MarkSyrveNotified)
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
