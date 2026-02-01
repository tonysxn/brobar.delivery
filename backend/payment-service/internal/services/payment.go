package services

import (
	"log"

	"github.com/tonysanin/brobar/payment-service/internal/provider"
	"github.com/tonysanin/brobar/payment-service/internal/rabbitmq"
	"github.com/tonysanin/brobar/pkg/monobank"
)

type PaymentService struct {
	provider  provider.PaymentProvider
	publisher *rabbitmq.Publisher
}

func NewPaymentService(provider provider.PaymentProvider, publisher *rabbitmq.Publisher) *PaymentService {
	return &PaymentService{
		provider:  provider,
		publisher: publisher,
	}
}

type InitPaymentInput struct {
	Amount      int                    `json:"amount"`
	OrderID     string                 `json:"order_id"`
	RedirectURL string                 `json:"redirect_url"`
	WebhookURL  string                 `json:"webhook_url"`
	Basket      []monobank.BasketOrder `json:"basket"`
}

type InitPaymentOutput struct {
	PaymentURL string `json:"payment_url"`
	InvoiceID  string `json:"invoice_id"`
}

func (s *PaymentService) InitPayment(input InitPaymentInput) (*InitPaymentOutput, error) {
	invoice, err := s.provider.CreateInvoice(
		input.Amount,
		input.OrderID,
		input.RedirectURL,
		input.WebhookURL,
		input.Basket,
	)
	if err != nil {
		return nil, err
	}

	return &InitPaymentOutput{
		PaymentURL: invoice.PageUrl,
		InvoiceID:  invoice.InvoiceId,
	}, nil
}

type WebhookPayload struct {
	InvoiceID string `json:"invoiceId"`
	Status    string `json:"status"`
	Amount    int    `json:"amount"`
	Ccy       int    `json:"ccy"`
	// Add other fields if necessary
}

// HandleWebhook processes the webhook from payment provider
func (s *PaymentService) HandleWebhook(payload WebhookPayload) error {
	log.Printf("Received webhook for invoice %s with status %s", payload.InvoiceID, payload.Status)

	if payload.Status == "success" {
		// Publish event to RabbitMQ
		// We publish a generic PaymentSuccess event that order-service can consume
		event := map[string]interface{}{
			"type":       "payment_success",
			"invoice_id": payload.InvoiceID,
			"amount":     payload.Amount,
			"status":     payload.Status,
		}

		// For now, let's assume we publish to a 'payment_events' queue or similar
		// But based on current setup, we might need to send to 'telegram_messages' if we wanted direct notification
		// However, correct architecture is: Payment Service -> (event) -> Order Service -> (msg) -> Telegram Service
		// So we publish to 'payment_events' which Order Service will listen to.
		// Wait, I need to check if OrderService listens to RabbitMQ.
		// The requirement said "Integrate ... in order-service ... we must via rabbitmq from webhook send message".

		// Publish to payment_events queue
		err := s.publisher.PublishPaymentEvent(event)

		return err
	}

	return nil
}
