package requests

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/validator"
	"github.com/tonysanin/brobar/product-service/internal/models"
)

type NestedVariationRequest struct {
	Name         string `json:"name" form:"name"`
	ExternalID   string `json:"external_id" form:"external_id"`
	DefaultValue *int   `json:"default_value" form:"default_value"`
	Show         bool   `json:"show" form:"show"`
}

type NestedVariationGroupRequest struct {
	Name         string                   `json:"name" form:"name"`
	ExternalID   string                   `json:"external_id" form:"external_id"`
	DefaultValue *int                     `json:"default_value" form:"default_value"`
	Show         bool                     `json:"show" form:"show"`
	Required     bool                     `json:"required" form:"required"`
	Variations   []NestedVariationRequest `json:"variations" form:"variations"`
}

type CreateProductRequest struct {
	Name            string                        `json:"name" form:"name"`
	Slug            string                        `json:"slug" form:"slug"`
	Description     *string                       `json:"description" form:"description"`
	Price           float64                       `json:"price" form:"price"`
	Weight          *float64                      `json:"weight" form:"weight"`
	ExternalID      string                        `json:"external_id" form:"external_id"`
	Hidden          bool                          `json:"hidden" form:"hidden"`
	Alcohol         bool                          `json:"alcohol" form:"alcohol"`
	Sold            bool                          `json:"sold" form:"sold"`
	CategoryID      uuid.UUID                     `json:"category_id" form:"category_id"`
	VariationGroups []NestedVariationGroupRequest `json:"variation_groups" form:"variation_groups"`
}

func (r CreateProductRequest) ToModel() *models.Product {
	p := &models.Product{
		Name:            r.Name,
		Slug:            r.Slug,
		Description:     r.Description,
		Price:           r.Price,
		Weight:          r.Weight,
		ExternalID:      r.ExternalID,
		Hidden:          r.Hidden,
		Alcohol:         r.Alcohol,
		Sold:            r.Sold,
		CategoryID:      r.CategoryID,
		VariationGroups: make([]models.ProductVariationGroup, len(r.VariationGroups)),
	}

	for i, g := range r.VariationGroups {
		vg := models.ProductVariationGroup{
			Name:         g.Name,
			ExternalID:   g.ExternalID,
			DefaultValue: g.DefaultValue,
			Show:         g.Show,
			Required:     g.Required,
			Variations:   make([]models.ProductVariation, len(g.Variations)),
		}
		for j, v := range g.Variations {
			vg.Variations[j] = models.ProductVariation{
				Name:         v.Name,
				ExternalID:   v.ExternalID,
				DefaultValue: v.DefaultValue,
				Show:         v.Show,
			}
		}
		p.VariationGroups[i] = vg
	}
	return p
}

func (r CreateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&r.Description, validation.Length(0, 2048)),
		validation.Field(&r.Slug, validation.By(func(value interface{}) error {
			str, _ := value.(string)
			if str == "" {
				return nil
			}
			length := len(str)
			if length < 2 || length > 255 {
				return errors.New("slug length must be between 2 and 255 characters")
			}
			return nil
		})),
		validation.Field(&r.Price, validation.Required, validator.IsNonNegative),
		validation.Field(&r.Weight, validator.IsNonNegative),
		validation.Field(&r.CategoryID, validation.Required, validator.IsUUID),
		validation.Field(&r.ExternalID, validation.Required, validation.Length(0, 100)),
	)
}

type UpdateProductRequest struct {
	Name        string    `json:"name" form:"name"`
	Slug        string    `json:"slug" form:"slug"`
	Description *string   `json:"description" form:"description"`
	Price       float64   `json:"price" form:"price"`
	Weight      *float64  `json:"weight" form:"weight"`
	ExternalID  string    `json:"external_id" form:"external_id"`
	Hidden      bool      `json:"hidden" form:"hidden"`
	Alcohol     bool      `json:"alcohol" form:"alcohol"`
	Sold        bool      `json:"sold" form:"sold"`
	CategoryID  uuid.UUID `json:"category_id" form:"category_id"`
}

func (r UpdateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&r.Description, validation.Length(0, 2048)),
		validation.Field(&r.Slug, validation.By(func(value interface{}) error {
			str, _ := value.(string)
			if str == "" {
				return nil
			}
			length := len(str)
			if length < 2 || length > 255 {
				return errors.New("slug length must be between 2 and 255 characters")
			}
			return nil
		})),
		validation.Field(&r.Price, validation.Required, validator.IsNonNegative),
		validation.Field(&r.Weight, validator.IsNonNegative),
		validation.Field(&r.CategoryID, validation.Required, validator.IsUUID),
		validation.Field(&r.ExternalID, validation.Required, validation.Length(0, 100)),
	)
}
