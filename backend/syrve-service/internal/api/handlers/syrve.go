package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/pkg/rabbitmq"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/pkg/syrve"
)

type SyrveHandler struct {
	client          *syrve.Client
	producer        *rabbitmq.Producer
	orderServiceURL string
}

func NewSyrveHandler(c *syrve.Client, p *rabbitmq.Producer, orderServiceURL string) *SyrveHandler {
	return &SyrveHandler{
		client:          c,
		producer:        p,
		orderServiceURL: orderServiceURL,
	}
}


func (h *SyrveHandler) GetProducts(c fiber.Ctx) error {
	tokenResp, err := h.client.GetAccessToken(c.Context())
	if err != nil {
		log.Printf("failed to get access token: %v", err)
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	products, err := h.client.GetProducts(c.Context(), tokenResp.Token, h.client.OrganizationID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, products)
}

func (h *SyrveHandler) GetStopLists(c fiber.Ctx) error {
	tokenResp, err := h.client.GetAccessToken(c.Context())
	if err != nil {
		log.Printf("failed to get access token: %v", err)
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	stopLists, err := h.client.GetStopLists(c.Context(), tokenResp.Token, h.client.OrganizationID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, stopLists)
}

func (h *SyrveHandler) HandleWebhook(c fiber.Ctx) error {
	body := c.Body()
	log.Printf("Received Syrve Webhook Body: %s", string(body))

	// Syrve webhooks are often sent as an array of events
	var events []struct {
		EventType      string `json:"eventType"`
		EventTime      string `json:"eventTime"`
		OrganizationID string `json:"organizationId"`
		EventInfo      struct {
			ID             string `json:"id"`
			CreationStatus string `json:"creationStatus"`
			ErrorInfo      struct {
				Message     string `json:"message"`
				Description string `json:"description"`
			} `json:"errorInfo"`
		} `json:"eventInfo"`
	}

	if err := json.Unmarshal(body, &events); err != nil {
		// Try unmarshalling as single object if array fails
		log.Printf("Webhook is not an array, skipping structured parsing for now")
		return c.SendStatus(200)
	}

	for _, event := range events {
		if event.EventType == "DeliveryOrderUpdate" || event.EventType == "TableOrderUpdate" ||
			event.EventType == "DeliveryOrderError" || event.EventType == "TableOrderError" {

			shortID := event.EventInfo.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			statusText := event.EventInfo.CreationStatus
			if statusText == "" {
				statusText = "Updated"
			}

			// Noise Filtering: Only Success or Error
			if statusText != "Success" && statusText != "Error" {
				log.Printf("Ignoring non-final status update: %s for order %s", statusText, shortID)
				continue
			}

			// Persistent Idempotency check: Call order-service
			markUrl := fmt.Sprintf("%s/orders/%s/syrve-notified", h.orderServiceURL, event.EventInfo.ID)
			resp, err := http.Post(markUrl, "application/json", nil)
			if err != nil {
				log.Printf("Failed to call order-service for idempotency: %v", err)
				// Log error but maybe continue? If we can't check, we might double-notify.
				// But usually order-service should be up.
			} else {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					var result struct {
						Status string `json:"status"`
						Data   struct {
							Notified bool `json:"notified"`
						} `json:"data"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
						if !result.Data.Notified {
							log.Printf("Already notified about order %s (from DB), skipping", event.EventInfo.ID)
							continue
						}
					}
				} else {
					log.Printf("Order service returned non-OK status: %d", resp.StatusCode)
				}
			}

			statusUA := statusText
			switch statusText {
			case "Success":
				statusUA = "Успішно"
			case "Error":
				statusUA = "Помилка"
			case "Updated":
				statusUA = "Оновлено"
			}

			statusEmoji := "ℹ️"
			if statusText == "Error" || event.EventType == "DeliveryOrderError" || event.EventType == "TableOrderError" {
				statusEmoji = "❌"
			} else if statusText == "Success" {
				statusEmoji = "✅"
			}

			text := fmt.Sprintf("%s <b>Результат замовлення Syrve</b>\n\n", statusEmoji)
			text += fmt.Sprintf("ID: #%s\n", shortID)
			text += fmt.Sprintf("Статус: %s\n", statusUA)

			if event.EventInfo.ErrorInfo.Description != "" {
				text += fmt.Sprintf("\n⚠️ <i>%s</i>", event.EventInfo.ErrorInfo.Description)
			} else if event.EventInfo.ErrorInfo.Message != "" {
				text += fmt.Sprintf("\n⚠️ <i>%s</i>", event.EventInfo.ErrorInfo.Message)
			}

			notification := map[string]interface{}{
				"text": text,
			}

			noteBytes, _ := json.Marshal(notification)
			if err := h.producer.SendMessage(rabbitmq.QueueTelegram, string(noteBytes)); err != nil {
				log.Printf("Failed to send Telegram notification: %v", err)
			} else {
				log.Printf("Telegram notification sent for order %s (Status: %s)", shortID, statusText)
			}
		} else if event.EventType == "StopListUpdate" {
			log.Printf("Received StopListUpdate event: %s", event.EventInfo.ID)

			tokenResp, err := h.client.GetAccessToken(c.Context())
			if err != nil {
				log.Printf("Failed to get access token for stop list update: %v", err)
				continue
			}

			stopLists, err := h.client.GetStopLists(c.Context(), tokenResp.Token, h.client.OrganizationID)
			if err != nil {
				log.Printf("Failed to get stop lists: %v", err)
				continue
			}

			// Publish to RabbitMQ
			if err := h.publishStopListUpdate(c.Context(), stopLists); err != nil {
				log.Printf("Failed to publish stop list update: %v", err)
			}
		}
	}

	return c.SendStatus(200)
}

func (h *SyrveHandler) SyncStopLists(c fiber.Ctx) error {
	tokenResp, err := h.client.GetAccessToken(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	stopLists, err := h.client.GetStopLists(c.Context(), tokenResp.Token, h.client.OrganizationID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	if err := h.publishStopListUpdate(c.Context(), stopLists); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, stopLists)
}

func (h *SyrveHandler) publishStopListUpdate(ctx interface{}, stopLists *syrve.StopListResponse) error {
	// We need to flatten the structure for easier consumption
	type StopListEventItem struct {
		ProductID string  `json:"product_id"`
		Balance   float64 `json:"balance"`
	}
	var items []StopListEventItem

	for _, orgList := range stopLists.TerminalGroupStopLists {
		for _, group := range orgList.Items {
			for _, item := range group.Items {
				items = append(items, StopListEventItem{
					ProductID: item.ProductID,
					Balance:   item.Balance,
				})
			}
		}
	}
	
	// Publish to RabbitMQ
	wrapper := map[string]interface{}{
		"items":   items,
		"chat_id": 0, // Webhook has no chat context
	}
	eventBytes, _ := json.Marshal(wrapper)
	if err := h.producer.SendMessage("syrve.stop_list.updated", string(eventBytes)); err != nil {
		return err
	}
	
	log.Printf("Published stop list update with %d items", len(items))
	return nil
}
