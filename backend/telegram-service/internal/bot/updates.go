package bot

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/tonysanin/brobar/pkg/rabbitmq"
	"github.com/tonysanin/brobar/pkg/telegram"
)

type Update struct {
	UpdateID      int            `json:"update_id"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
	Message       *Message       `json:"message"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Data    string   `json:"data"`
	Message *Message `json:"message"`
}

type Message struct {
	MessageID int   `json:"message_id"`
	Chat      struct {
		ID int64 `json:"id"`
	} `json:"chat"`
	Text string `json:"text"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Handler struct {
	client        *telegram.Client
	producer      *rabbitmq.Producer
	allowedChatID int64
	syrveURL      string
	productURL    string
	webURL        string
}

func NewHandler(client *telegram.Client, producer *rabbitmq.Producer, allowedChatID int64, syrveURL, productURL, webURL string) *Handler {
	return &Handler{
		client:        client,
		producer:      producer,
		allowedChatID: allowedChatID,
		syrveURL:      syrveURL,
		productURL:    productURL,
		webURL:        webURL,
	}
}

func (h *Handler) HandleUpdate(body []byte) error {
	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		return err
	}

	if update.CallbackQuery != nil {
		return h.handleCallbackQuery(update.CallbackQuery)
	}

	if update.Message != nil {
		if strings.EqualFold(update.Message.Text, "menu") || 
		   strings.EqualFold(update.Message.Text, "/menu") || 
		   update.Message.Text == "Меню" {
			// Security check
			if update.Message.Chat.ID == h.allowedChatID {
				return h.handleMenu(update.Message.Chat.ID)
			}
		}
	}

	return nil
}

func (h *Handler) handleCallbackQuery(cq *CallbackQuery) error {
	log.Printf("Received callback: %s from %d", cq.Data, cq.From.ID)

	// Security: Only allow callbacks from the configured chat
	if cq.Message == nil || cq.Message.Chat.ID != h.allowedChatID {
		log.Printf("Rejected callback from unauthorized chat: %d (allowed: %d)", cq.Message.Chat.ID, h.allowedChatID)
		return nil
	}

	// Format: action:order_id(:payload_to)
	parts := strings.Split(cq.Data, ":")
	action := parts[0]
	
	if action == "call_taxi" {
		if len(parts) < 2 {
			return nil
		}
		orderID := parts[1]
		
		// Publish Taxi Request
		h.publishTaxiRequest(cq.Message.Chat.ID, orderID)
		
		h.client.AnswerCallbackQuery(cq.ID, "Викликаємо таксі... чекайте розрахунок")
	} else if action == "confirm_taxi" {
		if len(parts) < 2 {
			return nil
		}
		orderID := parts[1]
		
		// Delete the estimate message
		if cq.Message != nil {
			h.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		}
		
		// PayloadTo will be re-derived by ontaxi-service from order address
		h.publishTaxiConfirm(cq.Message.Chat.ID, orderID)
		
		h.client.AnswerCallbackQuery(cq.ID, "Підтверджуємо замовлення...")
	} else if action == "cancel_taxi" {
		// Delete the estimate message
		if cq.Message != nil {
			h.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		}
		h.client.AnswerCallbackQuery(cq.ID, "Скасовано")
	} else if action == "compare_stop_list" {
		// New action for comparing stop lists
		h.handleCompareStopList(cq.Message.Chat.ID)
		h.client.AnswerCallbackQuery(cq.ID, "Звіряємо...")
	} else if action == "toggle_sales_paused" {
		if cq.Message != nil {
			h.handleToggleSalesPaused(cq.Message.Chat.ID, cq.Message.MessageID)
		}
		h.client.AnswerCallbackQuery(cq.ID, "Оновлюємо...")
	} else if action == "show_stock" {
		h.handleShowStock(cq.Message.Chat.ID)
		h.client.AnswerCallbackQuery(cq.ID, "Формуємо...")
	} else {
		// Unknown callback
		h.client.AnswerCallbackQuery(cq.ID, "")
	}

	return nil
}

func (h *Handler) publishTaxiRequest(chatID int64, orderID string) {
	payload := map[string]interface{}{
		"chat_id":  chatID,
		"order_id": orderID,
	}
	msg, _ := json.Marshal(payload)
	_ = h.producer.SendMessage("taxi_requests", string(msg))
}

func (h *Handler) publishTaxiConfirm(chatID int64, orderID string) {
	payload := map[string]interface{}{
		"chat_id":  chatID,
		"order_id": orderID,
	}
	msg, _ := json.Marshal(payload)
	_ = h.producer.SendMessage("taxi_confirms", string(msg))
}
