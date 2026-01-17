package tdlib

import (
	"fmt"
	"path/filepath"

	"github.com/zelenin/go-tdlib/client"
)

// BotAuthorizer handles the authentication process for a bot
type BotAuthorizer struct {
	apiID    int32
	apiHash  string
	botToken string
}

func (a *BotAuthorizer) Handle(tdClient *client.Client, state client.AuthorizationState) error {
	switch state.AuthorizationStateType() {
	case client.TypeAuthorizationStateWaitTdlibParameters:
		param := &client.SetTdlibParametersRequest{
			UseTestDc:           false,
			DatabaseDirectory:   filepath.Join(".tdlib", "database"),
			FilesDirectory:      filepath.Join(".tdlib", "files"),
			UseFileDatabase:     true,
			UseChatInfoDatabase: true,
			UseMessageDatabase:  true,
			UseSecretChats:      false,
			ApiId:               a.apiID,
			ApiHash:             a.apiHash,
			SystemLanguageCode:  "en",
			DeviceModel:         "Server",
			SystemVersion:       "1.0",
			ApplicationVersion:  "1.0",
		}
		_, err := tdClient.SetTdlibParameters(param)
		if err != nil {
			fmt.Printf("SetTdlibParameters error: %v\n", err)
			return err
		}

	case client.TypeAuthorizationStateWaitPhoneNumber:
		_, err := tdClient.CheckAuthenticationBotToken(&client.CheckAuthenticationBotTokenRequest{
			Token: a.botToken,
		})
		if err != nil {
			fmt.Printf("CheckAuthenticationBotToken error: %v\n", err)
			return err
		}
	}
	return nil
}

func (a *BotAuthorizer) Close() {}

type Client struct {
	tdlibClient *client.Client
	listeners   []chan *client.UpdateNewMessage
	stopCh      chan struct{}
}

func NewClient(apiID int32, apiHash, botToken string) (*Client, error) {
	authorizer := &BotAuthorizer{
		apiID:    apiID,
		apiHash:  apiHash,
		botToken: botToken,
	}

	// client.NewClient takes the authorizer and handles the auth loop
	tdClient, err := client.NewClient(authorizer)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	c := &Client{
		tdlibClient: tdClient,
		listeners:   make([]chan *client.UpdateNewMessage, 0),
		stopCh:      make(chan struct{}),
	}

	go c.runUpdateLoop()

	return c, nil
}

func (c *Client) runUpdateLoop() {
	listener := c.tdlibClient.GetListener()
	defer listener.Close()

	for {
		select {
		case update, ok := <-listener.Updates:
			if !ok {
				return
			}
			if update.GetClass() == client.ClassUpdate {
				switch u := update.(type) {
				case *client.UpdateNewMessage:
					for _, l := range c.listeners {
						select {
						case l <- u:
						default:
						}
					}
				}
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c *Client) AddUpdateListener(listener chan *client.UpdateNewMessage) {
	c.listeners = append(c.listeners, listener)
}

func (c *Client) SendMessage(chatID int64, text string, replyMarkup client.ReplyMarkup) (*client.Message, error) {
	formattedText := &client.FormattedText{
		Text:     text,
		Entities: []*client.TextEntity{},
	}

	inputMessageContent := &client.InputMessageText{
		Text:       formattedText,
		ClearDraft: true,
	}

	req := &client.SendMessageRequest{
		ChatId:              chatID,
		InputMessageContent: inputMessageContent,
		ReplyMarkup:         replyMarkup,
	}

	// Set WithDefaults is likely not needed or different in zelenin
	// req.WithOptions(...) ? No, simple struct.

	return c.tdlibClient.SendMessage(req)
}

func (c *Client) SendProfile(chatID int64, text, phone, address, mapLink string) (*client.Message, error) {
	formattedText := &client.FormattedText{
		Text:     text,
		Entities: []*client.TextEntity{},
	}

	rows := make([][]*client.InlineKeyboardButton, 0)

	if mapLink != "" {
		rows = append(rows, []*client.InlineKeyboardButton{
			{
				Text: "📍 Карта",
				Type: &client.InlineKeyboardButtonTypeUrl{
					Url: mapLink,
				},
			},
		})
	}

	if phone != "" {
		rows = append(rows, []*client.InlineKeyboardButton{
			{
				Text: "📱 Скопировать телефон",
				Type: &client.InlineKeyboardButtonTypeCopyText{
					Text: phone,
				},
			},
		})
	}

	if address != "" {
		rows = append(rows, []*client.InlineKeyboardButton{
			{
				Text: "📍 Скопировать адрес",
				Type: &client.InlineKeyboardButtonTypeCopyText{
					Text: address,
				},
			},
		})
	}

	replyMarkup := &client.ReplyMarkupInlineKeyboard{
		Rows: rows,
	}

	inputMessageContent := &client.InputMessageText{
		Text:       formattedText,
		ClearDraft: true,
	}

	req := &client.SendMessageRequest{
		ChatId:              chatID,
		InputMessageContent: inputMessageContent,
		ReplyMarkup:         replyMarkup,
	}

	return c.tdlibClient.SendMessage(req)
}

func (c *Client) Stop() {
	close(c.stopCh)
	c.tdlibClient.Close()
}
