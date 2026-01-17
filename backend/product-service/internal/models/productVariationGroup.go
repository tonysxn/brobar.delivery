package models

import (
	"github.com/google/uuid"
)

type ProductVariationGroup struct {
	ID           uuid.UUID          `json:"id" db:"id"`
	ProductID    uuid.UUID          `json:"product_id" db:"product_id" validate:"required"`
	Name         string             `json:"name" db:"name" validate:"required,min=2,max=255"`
	ExternalID   string             `json:"external_id" db:"external_id" validate:"max=100"`
	DefaultValue *int               `json:"default_value" db:"default_value"`
	Show         bool               `json:"show" db:"show"`
	Required     bool               `json:"required" db:"required"`
	Variations   []ProductVariation `json:"variations,omitempty"`
}
