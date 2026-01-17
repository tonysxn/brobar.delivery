package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/web-service/internal/services"
)

type SettingHandler struct {
	service *services.SettingService
}

func NewSettingHandler(service *services.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

func (h *SettingHandler) GetSettings(c fiber.Ctx) error {
	settings, err := h.service.GetAllSettings()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("failed to fetch settings"))
	}
	return response.Success(c, settings)
}

func (h *SettingHandler) UpdateSetting(c fiber.Ctx) error {
	key := c.Params("key")

	type UpdateRequest struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	}

	var req UpdateRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid request body"))
	}

	if err := h.service.UpdateSetting(key, req.Value, req.Type); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("failed to update setting"))
	}

	return response.Success(c, fiber.Map{"status": "updated"})
}
