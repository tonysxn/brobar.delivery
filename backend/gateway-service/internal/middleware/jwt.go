package middleware

import (
	"errors"
	"github.com/tonysanin/brobar/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret []byte
}

func NewJWTMiddleware(cfg JWTConfig) fiber.Handler {
	return fiber.Handler(func(c fiber.Ctx) error {
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
				return nil, errors.New("unexpected signing method")
			}
			return cfg.Secret, nil
		})

		if err != nil || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
		}

		c.Locals("user_claims", claims)
		return c.Next()
	})
}
