package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tonysanin/brobar/gateway-service/internal/models"
	"github.com/tonysanin/brobar/pkg/response"
	"log"
)

func AdminOnly(c fiber.Ctx) error {
	claimsRaw := c.Locals("user_claims")
	log.Println(claimsRaw)
	if claimsRaw == nil {
		return response.Error(c, fiber.StatusForbidden, fiber.ErrForbidden)
	}

	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return response.Error(c, fiber.StatusForbidden, fiber.ErrForbidden)
	}

	roleVal, ok := claims["role"]

	if !ok {
		return response.Error(c, fiber.StatusForbidden, fiber.ErrForbidden)
	}

	roleStr, ok := roleVal.(string)
	if !ok {
		return fiber.ErrForbidden
	}

	role := models.Role(roleStr)
	if role != models.RoleAdmin {
		return response.Error(c, fiber.StatusForbidden, fiber.ErrForbidden)
	}

	return c.Next()
}
