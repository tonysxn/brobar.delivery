package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/product-service/internal/api/requests"
	customerrors "github.com/tonysanin/brobar/product-service/internal/errors"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/services"
)

type ProductHandler struct {
	service         *services.ProductService
	categoryService *services.CategoryService
}

func NewProductHandler(s *services.ProductService, cs *services.CategoryService) *ProductHandler {
	return &ProductHandler{service: s, categoryService: cs}
}

func (h *ProductHandler) GetProducts(c fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")
	orderBy := c.Query("order_by", "sort")
	orderDir := c.Query("order_dir", "asc")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	allowedOrderFields := map[string]bool{
		"price": true,
		"name":  true,
		"sort":  true,
		"slug":  true,
	}
	if !allowedOrderFields[orderBy] {
		orderBy = "name"
	}

	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "desc"
	}

	offset := (page - 1) * limit

	products, totalCount, err := h.service.GetProductsWithPagination(c.Context(), limit, offset, orderBy, orderDir)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	resp := response.PaginatedResponse[models.Product]{
		Data: products,
		Pagination: response.Pagination{
			TotalCount: totalCount,
			Page:       page,
			Limit:      limit,
			OrderBy:    orderBy,
			OrderDir:   orderDir,
		},
	}

	return response.Success(c, resp)
}

func (h *ProductHandler) GetProduct(c fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.service.GetProduct(c.Context(), id)
	if err != nil {
		if errors.Is(err, customerrors.ProductNotFound) {
			return response.NotFound(c)
		}
		return response.BadRequest(c, err)
	}
	return response.Success(c, product)
}

func (h *ProductHandler) CreateProduct(c fiber.Ctx) error {
	var req requests.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, err)
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, errors.New("image is required"))
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err)
	}

	product := req.ToModel()

	err = h.service.CreateProduct(c.Context(), product, fileHeader)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, product)
}

func (h *ProductHandler) UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid product id"))
	}

	var req requests.UpdateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, err)
	}

	fileHeader, _ := c.FormFile("image")

	if err := req.Validate(); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err)
	}

	// Validate Category Existence
	category, err := h.categoryService.GetCategoryById(c.Context(), req.CategoryID)
	if err != nil || category == nil {
		return response.BadRequest(c, errors.New("category does not exist"))
	}

	product := models.Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Price:       req.Price,
		Weight:      req.Weight,
		ExternalID:  req.ExternalID,
		Hidden:      req.Hidden,
		Alcohol:     req.Alcohol,
		Sold:        req.Sold,
		CategoryID:  req.CategoryID,
	}

	updatedProduct, err := h.service.UpdateProduct(c.Context(), productID, &product, fileHeader)
	if err != nil {
		if errors.Is(err, customerrors.ProductNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, updatedProduct)
}

func (h *ProductHandler) DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid product id"))
	}

	err = h.service.DeleteProduct(c.Context(), productID)
	if err != nil {
		if errors.Is(err, customerrors.ProductNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, nil)
}
