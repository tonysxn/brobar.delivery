package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/user-service/internal/api/requests"
	customerrors "github.com/tonysanin/brobar/user-service/internal/errors"
	"github.com/tonysanin/brobar/user-service/internal/models"
)

func (h *UserHandler) Register(c fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err)
	}

	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		PromoCard: req.PromoCard,
		Address:   req.Address,
	}

	err := h.service.CreateUser(c.Context(), &user)
	if err != nil {
		if errors.Is(err, customerrors.UserAlreadyExists) {
			return response.Error(c, fiber.StatusConflict, err)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	accessExpiration := time.Hour * 24
	refreshExpiration := time.Hour * 24 * 7
	now := time.Now()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.RoleID,
		"exp":     now.Add(accessExpiration).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"type":    "refresh",
		"exp":     now.Add(refreshExpiration).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(h.jwtSecret)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	refreshTokenString, err := refreshToken.SignedString(h.jwtSecret)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.Map{
		"access": fiber.Map{
			"access_token": accessTokenString,
			"expires_in":   int(accessExpiration.Seconds()),
		},
		"refresh": fiber.Map{
			"refresh_token": refreshTokenString,
			"expires_in":    int(refreshExpiration.Seconds()),
		},
		"user": user.ToDTO(),
	})
}
