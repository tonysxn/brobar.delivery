package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/syrve-service/internal/services/syrve"
)

type SyrveHandler struct {
	client *syrve.Client
}

func NewSyrveHandler(c *syrve.Client) *SyrveHandler {
	return &SyrveHandler{client: c}
}

func (h *SyrveHandler) GetProducts(c fiber.Ctx) error {
	tokenResp, err := h.client.GetAccessToken(c.Context())
	if err != nil {
		log.Fatal("failed to get access token:", err)
	}

	products, err := h.client.GetProducts(c.Context(), tokenResp.Token, h.client.OrganizationID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, products)
}
