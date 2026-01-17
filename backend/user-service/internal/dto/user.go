package dto

import "github.com/google/uuid"

type UserDTO struct {
	ID            uuid.UUID `json:"id"`
	RoleID        string    `json:"role_id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	AddressCoords string    `json:"address_coords"`
	Phone         string    `json:"phone"`
	PromoCard     string    `json:"promo_card"`
}
