package requests

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/tonysanin/brobar/pkg/validator"
)

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	PromoCard       string `json:"promo_card,omitempty"`
	Address         string `json:"address,omitempty"`
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&r.Email, validation.Required, validation.Match(emailRegex).Error("invalid email format")),
		validation.Field(&r.Phone, validation.Required, validator.IsPhone),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 100)),
		validation.Field(&r.PromoCard, validation.Length(1, 20)),
		validation.Field(&r.Address, validation.Length(1, 256)),
		validation.Field(&r.PasswordConfirm, validation.Required, validation.By(func(value interface{}) error {
			if str, ok := value.(string); ok {
				if str != r.Password {
					return errors.New("passwords do not match")
				}
			}
			return nil
		})),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, validation.Match(emailRegex).Error("invalid email format")),
		validation.Field(&r.Password, validation.Required),
	)
}
