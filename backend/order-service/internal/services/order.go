package services

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"math"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/tonysanin/brobar/order-service/internal/clients"
	"github.com/tonysanin/brobar/order-service/internal/models"
	"github.com/tonysanin/brobar/order-service/internal/repositories"
	"github.com/tonysanin/brobar/pkg/rabbitmq"
)

type OrderService struct {
	repository          *repositories.OrderRepository
	orderItemRepository *repositories.OrderItemRepository
	productClient       *clients.ProductClient
	validationService   *ValidationService
	producer            *rabbitmq.Producer
}

func NewOrderService(
	repository *repositories.OrderRepository,
	orderItemRepository *repositories.OrderItemRepository,
	productClient *clients.ProductClient,
	validationService *ValidationService,
	producer *rabbitmq.Producer,
) *OrderService {
	return &OrderService{
		repository:          repository,
		orderItemRepository: orderItemRepository,
		productClient:       productClient,
		validationService:   validationService,
		producer:            producer,
	}
}

// OrderItemInput represents minimal item data from frontend
type OrderItemInput struct {
	ProductID          uuid.UUID
	ProductVariationID *uuid.UUID
	Quantity           int
}

// CreateOrderInput represents minimal order data from frontend
type CreateOrderInput struct {
	Name           string
	Phone          string
	Email          string
	DeliveryTypeID string
	Address        string
	Zone           string
	Entrance       string
	DeliveryDoor   bool
	Coords         string
	Time           string
	PaymentMethod  string
	Cutlery        int
	PromoCode      string
	Wishes         string
	Items          []OrderItemInput
	ClientTotal    float64
}

func (s *OrderService) CreateOrderFromInput(ctx context.Context, input *CreateOrderInput) (*models.Order, error) {
	// 1. Validate time
	if err := s.validationService.ValidateOrderTime(input.Time, input.DeliveryTypeID); err != nil {
		return nil, err
	}

	// 2. Fetch products and build order items with actual prices
	var items []models.OrderItem
	var itemsTotal float64 = 0

	for _, itemInput := range input.Items {
		product, err := s.productClient.GetProduct(itemInput.ProductID)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrProductNotFound, itemInput.ProductID.String())
		}

		item := models.OrderItem{
			ID:                uuid.New(),
			ProductID:         itemInput.ProductID,
			ExternalProductID: product.ExternalID,
			Quantity:          itemInput.Quantity,
			// Use product price/weight (variations don't have separate prices in this system)
			Price:  product.Price,
			Weight: product.Weight,
		}

		// If variation is specified, fetch variation and group info
		if itemInput.ProductVariationID != nil {
			variation, err := s.productClient.GetVariation(*itemInput.ProductVariationID)
			if err != nil {
				return nil, fmt.Errorf("варіація не знайдена: %s", itemInput.ProductVariationID.String())
			}

			// Fetch variation group for the group name
			group, err := s.productClient.GetVariationGroup(variation.GroupID)
			if err == nil && group != nil {
				item.ProductVariationGroupID = &group.ID
				item.ProductVariationGroupName = &group.Name
			}

			item.ProductVariationID = itemInput.ProductVariationID
			item.ProductVariationExternalID = &variation.ExternalID
			item.ProductVariationName = &variation.Name
			item.Name = fmt.Sprintf("%s (%s)", product.Name, variation.Name)
		} else {
			item.Name = product.Name
		}

		item.TotalPrice = item.Price * float64(item.Quantity)
		item.TotalWeight = item.Weight * float64(item.Quantity)

		itemsTotal += item.TotalPrice
		items = append(items, item)
	}

	// 3. Calculate delivery cost by coordinates (now we have itemsTotal for free delivery check)
	var deliveryCost float64 = 0
	var deliveryDoorPrice float64 = 0
	var zoneName string

	if input.DeliveryTypeID == "delivery" && input.Coords != "" {
		cost, doorPrice, zone, err := s.validationService.GetDeliveryCost(input.Coords, input.DeliveryDoor, itemsTotal)
		if err != nil {
			return nil, err
		}
		deliveryCost = cost
		deliveryDoorPrice = doorPrice
		if zone != nil {
			zoneName = zone.Name
		}
	}

	// 4. Calculate server total
	serverTotal := itemsTotal + deliveryCost

	// 4. Compare totals (allow small difference for rounding)
	if math.Abs(serverTotal-input.ClientTotal) > 1.0 {
		return nil, fmt.Errorf("%w (очікувано: %.2f, отримано: %.2f)", ErrPriceMismatch, serverTotal, input.ClientTotal)
	}

	// 5. Parse time
	var orderTime time.Time
	if input.Time == "ASAP" {
		orderTime = time.Now()
	} else {
		parsed, err := time.Parse("2006-01-02 15:04", input.Time)
		if err != nil {
			return nil, fmt.Errorf("невірний формат часу")
		}
		orderTime = parsed
	}

	// 6. Create order
	order := &models.Order{
		ID:                uuid.New(),
		StatusID:          models.StatusPending,
		TotalPrice:        serverTotal,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Name:              input.Name,
		Phone:             input.Phone,
		Email:             input.Email,
		Address:           input.Address,
		Entrance:          input.Entrance,
		Zone:              &zoneName,
		Coords:            input.Coords,
		Time:              orderTime,
		PaymentMethod:     input.PaymentMethod,
		Cutlery:           input.Cutlery,
		Promo:             input.PromoCode,
		Wishes:            input.Wishes,
		DeliveryCost:      deliveryCost,
		DeliveryDoor:      input.DeliveryDoor,
		DeliveryDoorPrice: deliveryDoorPrice,
		DeliveryTypeID:    models.DeliveryType(input.DeliveryTypeID),
		Items:             items,
	}

	// 7. Save to database
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
	}

	if err := s.repository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	if err := s.orderItemRepository.CreateOrderItems(ctx, order.Items); err != nil {
		return nil, err
	}

	// 8. Send notification (async)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in notification sender:", r)
			}
		}()

		var itemsList string
		for _, item := range order.Items {
			// Escape item name for HTML
			itemsList += fmt.Sprintf("- %s x%d (%.0f ₴)\n", html.EscapeString(item.Name), item.Quantity, item.TotalPrice)
		}

		// Format address block (for copying) and escape it
		addressBlock := html.EscapeString(order.Address)
		if order.Entrance != "" {
			addressBlock += fmt.Sprintf(", п. %s", html.EscapeString(order.Entrance))
		}

		addInfo := fmt.Sprintf("Під'їзд/код: %s", html.EscapeString(order.Entrance))

		deliveryMethod := "Доставка"
		if order.DeliveryTypeID == "pickup" {
			deliveryMethod = "Самовивіз"
			addressBlock = "Самовивіз"
			addInfo = "Самовивіз"
		} else if order.Zone != nil {
			deliveryMethod += fmt.Sprintf(" (%s)", html.EscapeString(*order.Zone))
		}

		deliveryPriceDisplay := fmt.Sprintf("%.0f ₴", order.DeliveryCost)

		// Map link (ensure URL encoded)
		queryAddr := order.Address
		if order.Coords != "" {
			queryAddr = order.Coords
		}
		mapLink := fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%s", url.QueryEscape(queryAddr))

		msgText := fmt.Sprintf(
			"<a href=\"https://brobar.com.ua/admin/orders/%s\">Нове замовлення #%s</a>\n\n"+
				"👤 %s\n"+
				"📞 %s\n\n"+
				"<code>%s</code>\n"+
				"<code>%s</code>\n\n"+
				"🕐 На коли: %s\n"+
				"🍴 Кількість приборів: %d\n\n"+
				"💸 Загальна сума: %.0f ₴\n"+
				"🚚 %s: %s\n\n"+
				"Чек:\n"+
				"%s"+
				"💳 Оплата: <b>%s</b>",
			order.ID, order.ID.String()[:8],
			html.EscapeString(order.Name),
			html.EscapeString(order.Phone),
			addressBlock, // already escaped above
			addInfo,      // already escaped above
			order.Time.Format("15:04 02.01.2006"),
			order.Cutlery,
			order.TotalPrice,
			deliveryMethod, deliveryPriceDisplay,
			itemsList,
			html.EscapeString(order.PaymentMethod),
		)

		if order.Wishes != "" {
			msgText += fmt.Sprintf("\n\n💬 <b>Wishes:</b> %s", html.EscapeString(order.Wishes))
		}

		// Construct Inline Keyboard safely
		keyboard := map[string]interface{}{
			"inline_keyboard": [][]map[string]string{
				{
					{"text": "📍 Карта", "url": mapLink},
				},
			},
		}

		keyboardBytes, _ := json.Marshal(keyboard)

		payload := map[string]interface{}{
			"text":         msgText,
			"reply_markup": string(keyboardBytes),
			"phone":        order.Phone,
			"address":      order.Address,
			"map_link":     mapLink,
		}

		jsonBody, _ := json.Marshal(payload)
		_ = s.producer.SendMessage(rabbitmq.QueueTelegram, string(jsonBody))
	}()

	return order, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}

	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = now
	}

	if order.StatusID == "" {
		order.StatusID = models.StatusPending
	}

	var totalPrice float64

	for i := range order.Items {
		item := &order.Items[i]

		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}
		item.OrderID = order.ID

		item.TotalPrice = item.Price * float64(item.Quantity)
		item.TotalWeight = item.Weight * float64(item.Quantity)

		totalPrice += item.TotalPrice
	}

	order.TotalPrice = totalPrice

	err := s.repository.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	err = s.orderItemRepository.CreateOrderItems(ctx, order.Items)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrderById(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.repository.GetOrderById(ctx, id)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.repository.GetAllOrders(ctx)
}

func (s *OrderService) GetOrdersWithPagination(ctx context.Context, limit, offset int, orderBy, orderDir string) ([]models.Order, int, error) {
	return s.repository.GetOrdersWithPagination(ctx, limit, offset, orderBy, orderDir)
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	now := time.Now()
	order.UpdatedAt = now

	var totalPrice float64

	for i := range order.Items {
		item := &order.Items[i]

		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}
		item.OrderID = order.ID

		item.TotalPrice = item.Price * float64(item.Quantity)
		item.TotalWeight = item.Weight * float64(item.Quantity)

		totalPrice += item.TotalPrice
	}

	order.TotalPrice = totalPrice

	err := s.repository.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}

	err = s.orderItemRepository.DeleteOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return err
	}

	err = s.orderItemRepository.CreateOrderItems(ctx, order.Items)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteOrder(ctx, id)
}
