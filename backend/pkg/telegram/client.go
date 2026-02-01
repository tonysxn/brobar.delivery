package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const telegramAPIURL = "https://api.telegram.org/bot%s/%s"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type SendMessageRequest struct {
	ChatID      int64       `json:"chat_id"`
	Text        string      `json:"text"`
	ParseMode   string      `json:"parse_mode,omitempty"`
	ReplyMarkup interface{} `json:"reply_markup,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text     string                        `json:"text"`
	URL      string                        `json:"url,omitempty"`
	CopyText *InlineKeyboardButtonCopyText `json:"copy_text,omitempty"`
}

type InlineKeyboardButtonCopyText struct {
	Text string `json:"text"`
}

func (c *Client) SendMessage(chatID int64, text string, replyMarkup interface{}) error {
	req := SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: replyMarkup,
	}
	return c.send("sendMessage", req)
}

func (c *Client) SendProfile(chatID int64, text, phone, address, mapLink string) error {
	row := make([]InlineKeyboardButton, 0)

	if mapLink != "" {
		row = append(row, InlineKeyboardButton{
			Text: "📍",
			URL:  mapLink,
		})
	}

	if phone != "" {
		row = append(row, InlineKeyboardButton{
			Text: "📞",
			CopyText: &InlineKeyboardButtonCopyText{
				Text: phone,
			},
		})
	}

	if address != "" {
		row = append(row, InlineKeyboardButton{
			Text: "🏠",
			CopyText: &InlineKeyboardButtonCopyText{
				Text: address,
			},
		})
	}

	req := SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{row}},
	}

	return c.send("sendMessage", req)
}

func (c *Client) send(method string, payload interface{}) error {
	url := fmt.Sprintf(telegramAPIURL, c.token, method)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respBody bytes.Buffer
		_, _ = respBody.ReadFrom(resp.Body)
		return fmt.Errorf("telegram api returned status: %s, body: %s", resp.Status, respBody.String())
	}

	return nil
}

func (c *Client) SetWebhook(url string) error {
	payload := map[string]string{
		"url": url,
	}
	return c.send("setWebhook", payload)
}

func (c *Client) DeleteWebhook() error {
	return c.send("deleteWebhook", nil)
}
