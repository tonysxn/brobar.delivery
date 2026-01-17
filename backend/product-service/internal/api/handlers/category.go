package handlers

import (
	"errors"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/response"
	customerrors "github.com/tonysanin/brobar/product-service/internal/errors"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/services"
)

type CategoryHandler struct {
	service        *services.CategoryService
	productService *services.ProductService
}

func NewCategoryHandler(s *services.CategoryService, ps *services.ProductService) *CategoryHandler {
	return &CategoryHandler{service: s, productService: ps}
}

func (h *CategoryHandler) GetCategories(c fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")
	orderBy := c.Query("order_by", "name")
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
		"name": true,
		"slug": true,
		"sort": true,
	}
	if !allowedOrderFields[orderBy] {
		orderBy = "sort"
	}

	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "asc"
	}

	offset := (page - 1) * limit

	categories, totalCount, err := h.service.GetCategoriesWithPagination(c.Context(), limit, offset, orderBy, orderDir)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	resp := response.PaginatedResponse[models.Category]{
		Data: categories,
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

func (h *CategoryHandler) GetCategory(c fiber.Ctx) error {
	id := c.Params("id")
	category, err := h.service.GetCategory(c.Context(), id)
	if err != nil {
		if errors.Is(err, customerrors.CategoryNotFound) {
			return response.NotFound(c)
		}
		return response.BadRequest(c, err)
	}
	return response.Success(c, category)
}

func (h *CategoryHandler) CreateCategory(c fiber.Ctx) error {
	var category models.Category
	if err := c.Bind().Body(&category); err != nil {
		return response.BadRequest(c, err)
	}

	err := validation.ValidateStruct(&category,
		validation.Field(&category.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&category.Icon, validation.Required, validation.Length(2, 50)),
		validation.Field(&category.Sort, validation.Min(-100), validation.Max(100)),
	)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err)
	}

	err = h.service.CreateCategory(c.Context(), &category)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, category)
}

func (h *CategoryHandler) UpdateCategory(c fiber.Ctx) error {
	id := c.Params("id")

	var updatedCategory models.Category
	if err := c.Bind().Body(&updatedCategory); err != nil {
		return response.BadRequest(c, err)
	}

	err := validation.ValidateStruct(&updatedCategory,
		validation.Field(&updatedCategory.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&updatedCategory.Icon, validation.Required, validation.Length(2, 50)),
		validation.Field(&updatedCategory.Sort, validation.Min(-100), validation.Max(100)),
	)

	updated, err := h.service.UpdateCategory(c.Context(), id, &updatedCategory)
	if err != nil {
		if errors.Is(err, customerrors.CategoryNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, updated)
}

func (h *CategoryHandler) DeleteCategory(c fiber.Ctx) error {
	id := c.Params("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid product id"))
	}

	err = h.service.DeleteCategory(c.Context(), categoryID)
	if err != nil {
		if errors.Is(err, customerrors.CategoryNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, nil)
}

func (h *CategoryHandler) GetProductsByCategory(c fiber.Ctx) error {
	id := c.Params("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid product id"))
	}
	products, err := h.productService.GetProductsByCategory(c.Context(), productID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, products)
}
