package models

import (
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/user-service/internal/dto"
)

type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	RoleID          Role      `json:"role_id" db:"role_id"`
	Email           string    `json:"email" db:"email" validate:"required,email"`
	Password        string    `json:"password" db:"password"`
	PasswordConfirm string    `json:"password_confirm,omitempty" db:"-"`
	Name            string    `json:"name" db:"name"`
	Address         string    `json:"address" db:"address"`
	AddressCoords   string    `json:"address_coords" db:"address_coords"`
	Phone           string    `json:"phone" db:"phone"`
	PromoCard       string    `json:"promo_card" db:"promo_card"`
}

func (u *User) ToDTO() *dto.UserDTO {
	return &dto.UserDTO{
		ID:            u.ID,
		RoleID:        string(u.RoleID),
		Email:         u.Email,
		Name:          u.Name,
		Address:       u.Address,
		AddressCoords: u.AddressCoords,
		Phone:         u.Phone,
		PromoCard:     u.PromoCard,
	}
}
