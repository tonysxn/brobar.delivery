package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/payment-service/internal/services"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) InitPayment(c fiber.Ctx) error {
	var input services.InitPaymentInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.service.InitPayment(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

func (h *PaymentHandler) HandleMonobankWebhook(c fiber.Ctx) error {
	var payload services.WebhookPayload
	// Monobank sends JSON
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		// Just log and return ok to monobank to stop retrying if it's junk
		return c.SendStatus(fiber.StatusOK)
	}

	if err := h.service.HandleWebhook(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
