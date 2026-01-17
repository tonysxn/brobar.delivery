package models

import (
	"github.com/google/uuid"
)

type ProductVariation struct {
	ID           uuid.UUID `json:"id" db:"id"`
	GroupID      uuid.UUID `json:"group_id" db:"group_id" validate:"required"`
	ExternalID   string    `json:"external_id" db:"external_id" validate:"max=100"`
	DefaultValue *int      `json:"default_value" db:"default_value"`
	Show         bool      `json:"show" db:"show"`
	Name         string    `json:"name" db:"name" validate:"required,min=1,max=255"`
}
