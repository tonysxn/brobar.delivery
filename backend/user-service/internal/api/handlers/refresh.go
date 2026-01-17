package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/response"
	"time"
)

func (h *UserHandler) Refresh(c fiber.Ctx) error {
	type RefreshRequest struct {
		Token string `json:"refresh_token"`
	}

	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, err)
	}

	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	userID, err := uuid.Parse(claims["user_id"].(string))

	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	user, err := h.service.GetUserById(c.Context(), userID)
	if err != nil {
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
