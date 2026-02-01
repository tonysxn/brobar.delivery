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
	Address         *string   `json:"address,omitempty" db:"address"`
	AddressCoords   *string   `json:"address_coords,omitempty" db:"address_coords"`
	Phone           *string   `json:"phone,omitempty" db:"phone"`
	PromoCard       *string   `json:"promo_card,omitempty" db:"promo_card"`
}

func (u *User) ToDTO() *dto.UserDTO {
	var address, addressCoords, phone, promoCard string
	if u.Address != nil {
		address = *u.Address
	}
	if u.AddressCoords != nil {
		addressCoords = *u.AddressCoords
	}
	if u.Phone != nil {
		phone = *u.Phone
	}
	if u.PromoCard != nil {
		promoCard = *u.PromoCard
	}
	return &dto.UserDTO{
		ID:            u.ID,
		RoleID:        string(u.RoleID),
		Email:         u.Email,
		Name:          u.Name,
		Address:       address,
		AddressCoords: addressCoords,
		Phone:         phone,
		PromoCard:     promoCard,
	}
}
