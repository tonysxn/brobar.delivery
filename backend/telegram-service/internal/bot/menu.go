package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID         string   `json:"id"`
	ExternalID string   `json:"external_id"`
	Name       string   `json:"name"`
	Stock      *float64 `json:"stock"`
}

type ProductResponse struct {
	Success bool      `json:"success"`
	Data    []Product `json:"data"`
}

type SyrveStopListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		TerminalGroupStopLists []struct {
			OrganizationID string `json:"organizationId"`
			Items          []struct {
				TerminalGroupID string `json:"terminalGroupId"`
				Items           []struct {
					ProductID string  `json:"productId"`
					Balance   float64 `json:"balance"`
				} `json:"items"`
			} `json:"items"`
		} `json:"terminalGroupStopLists"`
	} `json:"data"`
}

func (h *Handler) getMenuState(chatID int64) (string, interface{}, error) {
	// Fetch current status
	settings, _ := h.fetchSettings()
	salesPaused := "false"
	if settings != nil {
		if v, ok := settings["sales_paused"]; ok {
			salesPaused = v.(string)
		}
	}

	stateIcon := "🟢"
	statusText := "ВІДКРИТО"
	salesBtnText := "⏸ Зупинити продажі"

	if salesPaused == "true" {
		stateIcon = "🔴"
		statusText = "ЗАКРИТО"
		salesBtnText = "▶️ Відновити продажі"
	}

	text := fmt.Sprintf("<b>Меню адміністратора</b>\nСтатус: %s %s\n\nОберіть дію:", stateIcon, statusText)

	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{
					"text":          "📊 Звірити Стоп-Лист",
					"callback_data": "compare_stop_list",
				},
				{
					"text":          "📦 Перевірка складу",
					"callback_data": "show_stock",
				},
			},
			{
				{
					"text":          salesBtnText,
					"callback_data": "toggle_sales_paused",
				},
			},
		},
	}
	
	return text, keyboard, nil
}

func (h *Handler) handleMenu(chatID int64) error {
	text, keyboard, err := h.getMenuState(chatID)
	if err != nil {
		return err
	}
	return h.client.SendMessage(chatID, text, keyboard)
}

func (h *Handler) handleToggleSalesPaused(chatID int64, messageID int) error {
	settings, err := h.fetchSettings()
	if err != nil {
		h.client.SendMessage(chatID, fmt.Sprintf("❌ Помилка отримання налаштувань: %v", err), nil)
		return nil
	}

	currentVal := "false"
	if v, ok := settings["sales_paused"]; ok {
		currentVal = v.(string)
	}

	newVal := "true"
	if currentVal == "true" {
		newVal = "false"
	}

	// Update
	payload := map[string]string{
		"value": newVal,
		"type":  "boolean",
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/settings/sales_paused", h.webURL), bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.client.SendMessage(chatID, fmt.Sprintf("❌ Помилка оновлення: %v", err), nil)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.client.SendMessage(chatID, fmt.Sprintf("❌ Помилка API: статус %d", resp.StatusCode), nil)
		return nil
	}
	
	// Re-render menu by editing the message
	text, keyboard, _ := h.getMenuState(chatID)
	
	// If opening shop (newVal == "false"), trigger sync via RabbitMQ
	if newVal == "false" {
		 if err := h.producer.SendMessage("syrve.sync.start", fmt.Sprintf(`{"chat_id": %d, "initiator": "auto_open"}`, chatID)); err != nil {
			h.client.SendMessage(chatID, fmt.Sprintf("❌ Помилка запуску синхронізації: %v", err), nil)
		 } else {
			h.client.SendMessage(chatID, "🔄 Запущено процес звірки та синхронізації стоп-листів...", nil)
		 }
	}

	return h.client.EditMessageText(chatID, messageID, text, keyboard)
}

func (h *Handler) handleCompareStopList(chatID int64) error {
	payload := fmt.Sprintf(`{"chat_id": %d, "initiator": "manual"}`, chatID)
	if err := h.producer.SendMessage("syrve.sync.start", payload); err != nil {
		return h.client.SendMessage(chatID, fmt.Sprintf("❌ Помилка: %v", err), nil)
	}
	return h.client.SendMessage(chatID, "🔄 Запит на звірку надіслано...", nil)
}

func (h *Handler) handleShowStock(chatID int64) error {
	payload := fmt.Sprintf(`{"chat_id": %d}`, chatID)
	if err := h.producer.SendMessage("product.report.stock", payload); err != nil {
		return h.client.SendMessage(chatID, fmt.Sprintf("❌ Помилка: %v", err), nil)
	}
	return h.client.SendMessage(chatID, "🔄 Формуємо звіт по наявності...", nil)
}

func (h *Handler) fetchSettings() (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/settings", h.webURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var wrapper struct {
		Success bool `json:"success"`
		Data    []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	if !wrapper.Success {
		return nil, fmt.Errorf("api returned success=false")
	}

	result := make(map[string]interface{})
	for _, s := range wrapper.Data {
		result[s.Key] = s.Value
	}

	return result, nil
}
