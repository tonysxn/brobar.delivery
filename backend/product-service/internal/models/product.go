package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID              uuid.UUID               `json:"id" db:"id"`
	ExternalID      string                  `json:"external_id" db:"external_id" validate:"max=100"`
	Name            string                  `json:"name" db:"name" validate:"required,min=2,255"`
	Slug            string                  `json:"slug" db:"slug" validate:"required,min=2,max=255"`
	Description     *string                 `json:"description" db:"description"`
	Price           float64                 `json:"price" db:"price" validate:"required,min=0"`
	Weight          *float64                `json:"weight" db:"weight" validate:"omitempty,min=0"`
	CategoryID      uuid.UUID               `json:"category_id" db:"category_id" validate:"required"`
	Sort            int                     `json:"sort" db:"sort" validate:"min=0"`
	Hidden          bool                    `json:"hidden" db:"hidden"`
	Alcohol         bool                    `json:"alcohol" db:"alcohol"`
	Sold            bool                    `json:"sold" db:"sold"`
	Image           string                  `json:"image" db:"image"`
	Stock           *float64                `json:"stock" db:"stock"`
	VariationGroups []ProductVariationGroup `json:"variation_groups,omitempty" db:"-"`
}
