package api

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/tonysanin/brobar/gateway-service/internal/middleware"
	"github.com/tonysanin/brobar/pkg/response"
)

type Server struct {
	app             *fiber.App
	userService     string
	productService  string
	telegramService string
	syrveService    string
	fileService     string
	webService      string
	orderService    string
	jwtSecret       []byte
}

type ServerConfig struct {
	UserServiceURL     string
	ProductServiceURL  string
	TelegramServiceURL string
	SyrveServiceURL    string
	FileServiceURL     string
	WebServiceURL      string
	OrderServiceURL    string
	JWTSecret          []byte
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		app:             fiber.New(fiber.Config{AppName: "Gateway Service"}),
		userService:     cfg.UserServiceURL,
		productService:  cfg.ProductServiceURL,
		telegramService: cfg.TelegramServiceURL,
		syrveService:    cfg.SyrveServiceURL,
		fileService:     cfg.FileServiceURL,
		webService:      cfg.WebServiceURL,
		orderService:    cfg.OrderServiceURL,
		jwtSecret:       []byte(cfg.JWTSecret),
	}

	s.app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	s.app.Use(recover.New())
	s.app.Use(logger.New())
	s.app.Use(requestid.New())
	s.app.Use(func(c fiber.Ctx) error {
		// Debug logging for CORS
		origin := c.Get("Origin")
		if origin != "" {
			log.Printf("Request Origin: %s", origin)
		}
		return c.Next()
	})
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3005", "http://localhost:4000", "http://127.0.0.1:3000", "http://127.0.0.1:3005", "http://127.0.0.1:4000"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}))
	s.app.Use(etag.New())

	s.SetupRoutes()

	return s
}

// proxyWithCORS wraps proxy.Do to preserve CORS headers
func (s *Server) proxyWithCORS(c fiber.Ctx, url string) error {
	// Store origin before proxying
	origin := c.Get("Origin")

	// Perform the proxy request
	if err := proxy.Do(c, url); err != nil {
		return err
	}

	// Re-apply CORS headers after proxying (they get stripped by proxy.Do)
	if origin != "" {
		// Check if origin is allowed
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:3005", "http://localhost:4000", "http://127.0.0.1:3000", "http://127.0.0.1:3005", "http://127.0.0.1:4000"}
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				c.Set("Access-Control-Allow-Origin", origin)
				c.Set("Access-Control-Allow-Credentials", "true")
				c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
				c.Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS")
				break
			}
		}
	}

	return nil
}

func (s *Server) ProxyToUserService(c fiber.Ctx) error {
	return s.proxyWithCORS(c, s.userService+c.OriginalURL())
}

func (s *Server) ProxyToFileService(c fiber.Ctx) error {
	return s.proxyWithCORS(c, s.fileService+c.OriginalURL())
}

func (s *Server) ProxyToProductService(c fiber.Ctx) error {
	return s.proxyWithCORS(c, s.productService+c.OriginalURL())
}

func (s *Server) ProxyToTelegramService(c fiber.Ctx) error {
	return s.proxyWithCORS(c, s.telegramService+c.OriginalURL())
}

func (s *Server) ProxyToSyrveService(c fiber.Ctx) error {
	return s.proxyWithCORS(c, s.syrveService+c.OriginalURL())
}

func (s *Server) ProxyToWebService(c fiber.Ctx) error {
	return s.proxyWithCORS(c, s.webService+c.OriginalURL())
}

func (s *Server) ProxyToOrderService(c fiber.Ctx) error {
	return s.proxyWithCORS(c, s.orderService+c.OriginalURL())
}

func (s *Server) SetupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	// Menu
	s.app.Get("/menu", s.ProxyToProductService)

	jwtMiddleware := middleware.NewJWTMiddleware(middleware.JWTConfig{
		Secret: s.jwtSecret,
	})

	// Telegram webhook
	telegramGroupAuthorized := s.app.Group("webhooks/telegram")
	telegramGroupAuthorized.Use(jwtMiddleware)
	telegramGroupAuthorized.Post("/", s.ProxyToTelegramService, middleware.AdminOnly)

	// File
	fileGroup := s.app.Group("/files")
	fileGroup.Get("/:id", s.ProxyToFileService)
	fileGroup.Use(jwtMiddleware)
	fileGroup.Post("/", s.ProxyToFileService, middleware.AdminOnly)

	// Server Time
	s.app.Get("/time", s.ProxyToWebService)

	// Settings
	settingsGroup := s.app.Group("/settings")
	settingsGroup.Get("/", s.ProxyToWebService)
	settingsGroup.Use(jwtMiddleware)
	settingsGroup.Put("/:key", s.ProxyToWebService, middleware.AdminOnly)

	// Reviews
	reviewsGroup := s.app.Group("/reviews")
	reviewsGroup.Post("/", s.ProxyToWebService)
	reviewsGroup.Get("/", s.ProxyToWebService)
	reviewsGroup.Get("/:id", s.ProxyToWebService)
	reviewsGroup.Use(jwtMiddleware)
	reviewsGroup.Delete("/:id", s.ProxyToWebService, middleware.AdminOnly)

	// Orders (public POST, admin GET/PUT/DELETE)
	ordersGroup := s.app.Group("/orders")
	ordersGroup.Post("/", s.ProxyToOrderService)
	ordersGroup.Use(jwtMiddleware)
	ordersGroup.Get("/", s.ProxyToOrderService, middleware.AdminOnly)
	ordersGroup.Get("/:id", s.ProxyToOrderService, middleware.AdminOnly)
	ordersGroup.Put("/:id", s.ProxyToOrderService, middleware.AdminOnly)
	ordersGroup.Delete("/:id", s.ProxyToOrderService, middleware.AdminOnly)

	// Syrve
	syrveGroupAuthorized := s.app.Group("/syrve")
	syrveGroupAuthorized.Use(jwtMiddleware)
	syrveGroupAuthorized.Get("/products", s.ProxyToSyrveService, middleware.AdminOnly)

	// Auth
	authGroup := s.app.Group("/auth")
	authGroup.Post("/login", s.ProxyToUserService)
	authGroup.Post("/register", s.ProxyToUserService)
	authGroup.Post("/refresh", s.ProxyToUserService)

	// User
	userGroup := s.app.Group("/user")
	userGroup.Use(jwtMiddleware)
	userGroup.Get("/me", s.ProxyToUserService)
	userGroup.Get("/:id", s.ProxyToUserService, middleware.AdminOnly)

	// Products
	categoryGroup := s.app.Group("/categories")
	categoryGroup.Get("/", s.ProxyToProductService)
	categoryGroup.Get("/:id", s.ProxyToProductService)
	categoryGroup.Get("/:id/products", s.ProxyToProductService)

	categoryGroupAuthorized := s.app.Group("/categories")
	categoryGroupAuthorized.Use(jwtMiddleware)
	categoryGroupAuthorized.Post("/", s.ProxyToProductService, middleware.AdminOnly)
	categoryGroupAuthorized.Put("/:id", s.ProxyToProductService, middleware.AdminOnly)
	categoryGroupAuthorized.Delete("/:id", s.ProxyToProductService, middleware.AdminOnly)

	productGroup := s.app.Group("/products")
	productGroup.Get("/", s.ProxyToProductService)
	productGroup.Get("/:id", s.ProxyToProductService)

	productGroupAuthorized := s.app.Group("/products")
	productGroupAuthorized.Use(jwtMiddleware)
	productGroupAuthorized.Post("/", s.ProxyToProductService, middleware.AdminOnly)
	productGroupAuthorized.Put("/:id", s.ProxyToProductService, middleware.AdminOnly)
	productGroupAuthorized.Delete("/:product_id/variation-groups", s.ProxyToProductService, middleware.AdminOnly)
	productGroupAuthorized.Delete("/:id", s.ProxyToProductService, middleware.AdminOnly)

	variationGroup := s.app.Group("/variations")
	variationGroup.Get("/", s.ProxyToProductService)
	variationGroup.Get("/:id", s.ProxyToProductService)

	variationGroupAuth := variationGroup.Group("", jwtMiddleware)
	variationGroupAuth.Post("/", s.ProxyToProductService, middleware.AdminOnly)
	variationGroupAuth.Put("/:id", s.ProxyToProductService, middleware.AdminOnly)
	variationGroupAuth.Delete("/:id", s.ProxyToProductService, middleware.AdminOnly)

	groupGroup := s.app.Group("/variation-groups")
	groupGroup.Get("/", s.ProxyToProductService)
	groupGroup.Get("/:id", s.ProxyToProductService)

	groupGroupAuth := groupGroup.Group("", jwtMiddleware)
	groupGroupAuth.Post("/", s.ProxyToProductService, middleware.AdminOnly)
	groupGroupAuth.Put("/:id", s.ProxyToProductService, middleware.AdminOnly)
	groupGroupAuth.Delete("/:id", s.ProxyToProductService, middleware.AdminOnly)
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
