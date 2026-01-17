package handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/user-service/internal/services"
	"strings"
)

type UserHandler struct {
	service   *services.UserService
	jwtSecret []byte
}

func NewUserHandler(s *services.UserService, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		service:   s,
		jwtSecret: jwtSecret,
	}
}

func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid user ID"))
	}

	user, err := h.service.GetUserById(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err)
	}

	return response.Success(c, user.ToDTO())
}

func (h *UserHandler) GetUserMe(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return h.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	user, err := h.service.GetUserById(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err)
	}

	return response.Success(c, user.ToDTO())
}
