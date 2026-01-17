package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/product-service/internal/api/handlers"
	"github.com/tonysanin/brobar/product-service/internal/services"
)

type Server struct {
	app                   *fiber.App
	productService        *services.ProductService
	categoryService       *services.CategoryService
	variationService      *services.ProductVariationService
	variationGroupService *services.ProductVariationGroupService
	productHandler        *handlers.ProductHandler
	categoryHandler       *handlers.CategoryHandler
	variationHandler      *handlers.ProductVariationHandler
	variationGroupHandler *handlers.ProductVariationGroupHandler
	menuHandler           *handlers.MenuHandler
}

func NewServer(
	productService *services.ProductService,
	categoryService *services.CategoryService,
	variationService *services.ProductVariationService,
	variationGroupService *services.ProductVariationGroupService,
) *Server {
	s := &Server{
		app: fiber.New(fiber.Config{
			AppName: "Product Service",
		}),
		productService:        productService,
		categoryService:       categoryService,
		variationService:      variationService,
		variationGroupService: variationGroupService,
	}

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.productHandler = handlers.NewProductHandler(productService, categoryService)
	s.categoryHandler = handlers.NewCategoryHandler(categoryService, productService)
	s.variationHandler = handlers.NewProductVariationHandler(variationService, variationGroupService)
	s.variationGroupHandler = handlers.NewProductVariationGroupHandler(variationGroupService, productService)
	s.menuHandler = handlers.NewMenuHandler(categoryService)

	s.SetupRoutes()

	return s
}

func (s *Server) SetupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok",
		})
	})

	s.app.Get("/menu", s.menuHandler.GetMenu)

	productGroup := s.app.Group("/products")
	productGroup.Get("/", s.productHandler.GetProducts)
	productGroup.Get("/:id", s.productHandler.GetProduct)
	productGroup.Post("/", s.productHandler.CreateProduct)
	productGroup.Put("/:id", s.productHandler.UpdateProduct)
	productGroup.Delete("/:product_id/variation-groups", s.variationGroupHandler.DeleteGroupsByProduct)
	productGroup.Delete("/:id", s.productHandler.DeleteProduct)

	categoryGroup := s.app.Group("/categories")
	categoryGroup.Get("/", s.categoryHandler.GetCategories)
	categoryGroup.Get("/:id", s.categoryHandler.GetCategory)
	categoryGroup.Get("/:id/products", s.categoryHandler.GetProductsByCategory)
	categoryGroup.Post("/", s.categoryHandler.CreateCategory)
	categoryGroup.Put("/:id", s.categoryHandler.UpdateCategory)
	categoryGroup.Delete("/:id", s.categoryHandler.DeleteCategory)

	variationGroup := s.app.Group("/variations")
	variationGroup.Get("/", s.variationHandler.GetVariations)
	variationGroup.Get("/:id", s.variationHandler.GetVariation)
	variationGroup.Post("/", s.variationHandler.CreateVariation)
	variationGroup.Put("/:id", s.variationHandler.UpdateVariation)
	variationGroup.Delete("/:id", s.variationHandler.DeleteVariation)

	groupGroup := s.app.Group("/variation-groups")
	groupGroup.Get("/", s.variationGroupHandler.GetGroups)
	groupGroup.Get("/:id", s.variationGroupHandler.GetGroup)
	groupGroup.Post("/", s.variationGroupHandler.CreateGroup)
	groupGroup.Put("/:id", s.variationGroupHandler.UpdateGroup)
	groupGroup.Delete("/:id", s.variationGroupHandler.DeleteGroup)
}

func (s *Server) Listen(address string) error {
	return s.app.Listen(address)
}
