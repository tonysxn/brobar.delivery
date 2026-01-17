package handlers

import (
	"errors"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/response"
	customerrors "github.com/tonysanin/brobar/product-service/internal/errors"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/services"
)

type ProductVariationGroupHandler struct {
	service        *services.ProductVariationGroupService
	productService *services.ProductService
}

func NewProductVariationGroupHandler(s *services.ProductVariationGroupService, ps *services.ProductService) *ProductVariationGroupHandler {
	return &ProductVariationGroupHandler{service: s, productService: ps}
}

func (h *ProductVariationGroupHandler) GetGroups(c fiber.Ctx) error {
	productID := c.Query("product_id")
	if productID == "" {
		return response.BadRequest(c, errors.New("product_id query parameter is required"))
	}

	groups, err := h.service.GetAllByProductID(c.Context(), productID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, groups)
}

func (h *ProductVariationGroupHandler) GetGroup(c fiber.Ctx) error {
	id := c.Params("id")
	groupId, err := uuid.Parse(id)
	group, err := h.service.GetByID(c.Context(), groupId)
	if err != nil {
		if errors.Is(err, customerrors.ProductVariationGroupNotFound) {
			return response.NotFound(c)
		}
		return response.BadRequest(c, err)
	}
	return response.Success(c, group)
}

func (h *ProductVariationGroupHandler) CreateGroup(c fiber.Ctx) error {
	var group models.ProductVariationGroup
	if err := c.Bind().Body(&group); err != nil {
		return response.BadRequest(c, err)
	}

	err := validation.ValidateStruct(&group,
		validation.Field(&group.ProductID, validation.Required),
		validation.Field(&group.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&group.ExternalID, validation.Required, validation.Length(2, 100)),
		validation.Field(&group.DefaultValue, validation.Min(0), validation.Max(1024)),
		validation.Field(&group.ProductID,
			validation.Required,
			validation.By(func(value interface{}) error {
				id, ok := value.(uuid.UUID)
				if !ok {
					return errors.New("invalid UUID format")
				}
				if id == uuid.Nil {
					return errors.New("product ID cannot be empty")
				}
				product, err := h.productService.GetProductById(c.Context(), id)
				if err != nil || product == nil {
					return errors.New("product does not exist")
				}
				return nil
			}),
		),
	)

	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err)
	}

	err = h.service.Create(c.Context(), &group)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, group)
}

func (h *ProductVariationGroupHandler) UpdateGroup(c fiber.Ctx) error {
	id := c.Params("id")
	var group models.ProductVariationGroup
	if err := c.Bind().Body(&group); err != nil {
		return response.BadRequest(c, err)
	}

	err := validation.ValidateStruct(&group,
		validation.Field(&group.ProductID, validation.Required),
		validation.Field(&group.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&group.ExternalID, validation.Required, validation.Length(2, 100)),
		validation.Field(&group.DefaultValue, validation.Min(0), validation.Max(1024)),
	)

	updatedGroup, err := h.service.Update(c.Context(), id, &group)
	if err != nil {
		if errors.Is(err, customerrors.ProductVariationGroupNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, updatedGroup)
}

func (h *ProductVariationGroupHandler) DeleteGroup(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, customerrors.ProductVariationGroupNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}
	return response.Success(c, nil)
}

func (h *ProductVariationGroupHandler) DeleteGroupsByProduct(c fiber.Ctx) error {
	productIDStr := c.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid product id"))
	}

	err = h.service.DeleteByProductID(c.Context(), productID)
	if err != nil {
		if errors.Is(err, customerrors.ProductNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, nil)
}
