package models

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
	Slug string    `json:"slug" db:"slug"`
	Icon string    `json:"icon" db:"icon"`
	Sort int       `json:"sort" db:"sort,omitempty"`
}
