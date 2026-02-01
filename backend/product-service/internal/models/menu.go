package models

import "github.com/google/uuid"

// MenuVariation represents a variation in the menu tree
type MenuVariation struct {
	ID           uuid.UUID `json:"id" db:"id"`
	GroupID      uuid.UUID `json:"-" db:"group_id"`
	Name         string    `json:"name" db:"name"`
	ExternalID   string    `json:"external_id,omitempty" db:"external_id"`
	DefaultValue *int      `json:"default_value,omitempty" db:"default_value"`
	Show         bool      `json:"show" db:"show"`
}

// MenuVariationGroup represents a variation group with its variations in the menu tree
type MenuVariationGroup struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	ProductID    uuid.UUID       `json:"-" db:"product_id"`
	Name         string          `json:"name" db:"name"`
	ExternalID   string          `json:"external_id,omitempty" db:"external_id"`
	DefaultValue *int            `json:"default_value,omitempty" db:"default_value"`
	Show         bool            `json:"show" db:"show"`
	Required     bool            `json:"required" db:"required"`
	Variations   []MenuVariation `json:"variations"`
}

// MenuProduct represents a product with its variation groups in the menu tree
type MenuProduct struct {
	ID              uuid.UUID            `json:"id" db:"id"`
	ExternalID      string               `json:"external_id,omitempty" db:"external_id"`
	Name            string               `json:"name" db:"name"`
	Slug            string               `json:"slug" db:"slug"`
	Description     *string              `json:"description,omitempty" db:"description"`
	Price           float64              `json:"price" db:"price"`
	Weight          *float64             `json:"weight" db:"weight"`
	CategoryID      uuid.UUID            `json:"-" db:"category_id"`
	Sort            int                  `json:"sort" db:"sort"`
	Hidden          bool                 `json:"hidden" db:"hidden"`
	Alcohol         bool                 `json:"alcohol" db:"alcohol"`
	Sold            bool                 `json:"sold" db:"sold"`
	Image           string               `json:"image" db:"image"`
	VariationGroups []MenuVariationGroup `json:"variation_groups"`
}

// MenuCategory represents a category with its products in the menu tree
type MenuCategory struct {
	ID       uuid.UUID     `json:"id" db:"id"`
	Name     string        `json:"name" db:"name"`
	Slug     string        `json:"slug" db:"slug"`
	Icon     string        `json:"icon,omitempty" db:"icon"`
	Sort     int           `json:"sort" db:"sort"`
	Products []MenuProduct `json:"products"`
}
