package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/user-service/internal/api/requests"
)

func (h *UserHandler) Login(c fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err)
	}

	user, err := h.service.GetUserByEmail(c.Context(), req.Email)

	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
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
		"exp":     now.Add(refreshExpiration).Unix(),
		"type":    "refresh",
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
			"token":      accessTokenString,
			"expires_in": int(accessExpiration.Seconds()),
		},
		"refresh": fiber.Map{
			"token":      refreshTokenString,
			"expires_in": int(refreshExpiration.Seconds()),
		},
		"user": user.ToDTO(),
	})
}
