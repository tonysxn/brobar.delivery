package provider

import (
	"github.com/tonysanin/brobar/pkg/monobank"
)

type PaymentProvider interface {
	CreateInvoice(amount int, orderID string, redirectURL string, webhookURL string, basket []monobank.BasketOrder) (*monobank.InvoiceData, error)
}

type MonobankProvider struct {
	client *monobank.Acquiring
}

func NewMonobankProvider(client *monobank.Acquiring) *MonobankProvider {
	return &MonobankProvider{
		client: client,
	}
}

func (m *MonobankProvider) CreateInvoice(amount int, orderID string, redirectURL string, webhookURL string, basket []monobank.BasketOrder) (*monobank.InvoiceData, error) {
	return m.client.CreateInvoice(&monobank.Invoice{
		Amount: amount,
		Ccy:    monobank.UAH,
		MerchantPaymInfo: monobank.MerchantPaymInfo{
			Reference:   orderID,
			Destination: "Оплата замовлення",
			BasketOrder: basket,
		},
		RedirectURL: redirectURL,
		WebHookURL:  webhookURL,
	})
}
