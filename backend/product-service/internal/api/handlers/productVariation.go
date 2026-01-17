package handlers

import (
	"errors"
	"github.com/google/uuid"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/pkg/response"
	customerrors "github.com/tonysanin/brobar/product-service/internal/errors"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/services"
)

type ProductVariationHandler struct {
	service                      *services.ProductVariationService
	productVariationGroupService *services.ProductVariationGroupService
}

func NewProductVariationHandler(s *services.ProductVariationService, pvgs *services.ProductVariationGroupService) *ProductVariationHandler {
	return &ProductVariationHandler{service: s, productVariationGroupService: pvgs}
}

func (h *ProductVariationHandler) GetVariations(c fiber.Ctx) error {
	groupID := c.Query("group_id")
	if groupID == "" {
		return response.BadRequest(c, errors.New("group_id query parameter is required"))
	}

	variations, err := h.service.GetAllByGroupID(c.Context(), groupID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, variations)
}

func (h *ProductVariationHandler) GetVariation(c fiber.Ctx) error {
	id := c.Params("id")
	variation, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, customerrors.ProductVariationNotFound) {
			return response.NotFound(c)
		}
		return response.BadRequest(c, err)
	}
	return response.Success(c, variation)
}

func (h *ProductVariationHandler) CreateVariation(c fiber.Ctx) error {
	var variation models.ProductVariation
	if err := c.Bind().Body(&variation); err != nil {
		return response.BadRequest(c, err)
	}

	err := validation.ValidateStruct(&variation,
		validation.Field(&variation.GroupID,
			validation.Required,
			validation.By(func(value interface{}) error {
				id, ok := value.(uuid.UUID)
				if !ok {
					return errors.New("invalid UUID format")
				}
				if id == uuid.Nil {
					return errors.New("group ID cannot be empty")
				}
				group, err := h.productVariationGroupService.GetByID(c.Context(), id)
				if err != nil || group == nil {
					return errors.New("group does not exist")
				}
				return nil
			}),
		),
		validation.Field(&variation.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&variation.ExternalID, validation.Required, validation.Length(0, 100)),
	)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err)
	}

	err = h.service.Create(c.Context(), &variation)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, variation)
}

func (h *ProductVariationHandler) UpdateVariation(c fiber.Ctx) error {
	id := c.Params("id")
	var variation models.ProductVariation
	if err := c.Bind().Body(&variation); err != nil {
		return response.BadRequest(c, err)
	}

	updatedVariation, err := h.service.Update(c.Context(), id, &variation)
	if err != nil {
		if errors.Is(err, customerrors.ProductVariationNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, updatedVariation)
}

func (h *ProductVariationHandler) DeleteVariation(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, customerrors.ProductVariationNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, nil)
}
