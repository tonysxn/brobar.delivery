package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/product-service/internal/services"
)

type MenuHandler struct {
	categoryService *services.CategoryService
}

func NewMenuHandler(categoryService *services.CategoryService) *MenuHandler {
	return &MenuHandler{
		categoryService: categoryService,
	}
}

func (h *MenuHandler) GetMenu(c fiber.Ctx) error {
	menu, err := h.categoryService.GetMenu(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("failed to get menu"))
	}

	return response.Success(c, menu)
}
