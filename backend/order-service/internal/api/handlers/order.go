package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/order-service/internal/api/requests"
	customerrors "github.com/tonysanin/brobar/order-service/internal/errors"
	"github.com/tonysanin/brobar/order-service/internal/models"
	"github.com/tonysanin/brobar/order-service/internal/services"
	"github.com/tonysanin/brobar/pkg/response"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrder(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid order id"))
	}

	order, err := h.service.GetOrderById(c.Context(), id)
	if err != nil {
		if errors.Is(err, customerrors.OrderNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, order)
}

func (h *OrderHandler) GetOrders(c fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")
	orderBy := c.Query("order_by", "created_at")
	orderDir := c.Query("order_dir", "desc")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	allowedOrderFields := map[string]bool{
		"created_at":  true,
		"total_price": true,
		"name":        true,
	}
	if !allowedOrderFields[orderBy] {
		orderBy = "created_at"
	}

	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "desc"
	}

	offset := (page - 1) * limit

	orders, totalCount, err := h.service.GetOrdersWithPagination(c.Context(), limit, offset, orderBy, orderDir)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	resp := response.PaginatedResponse[models.Order]{
		Data: orders,
		Pagination: response.Pagination{
			TotalCount: totalCount,
			Page:       page,
			Limit:      limit,
			OrderBy:    orderBy,
			OrderDir:   orderDir,
		},
	}

	return response.Success(c, resp)
}

// CreateOrder handles public order creation from frontend
func (h *OrderHandler) CreateOrder(c fiber.Ctx) error {
	var req requests.CreateOrderRequest

	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err)
	}

	// Map to input for service
	input := &services.CreateOrderInput{
		Name:           req.Name,
		Phone:          req.Phone,
		Email:          req.Email,
		DeliveryTypeID: req.DeliveryTypeID,
		Address:        req.Address,
		Zone:           req.Zone,
		Entrance:       req.Entrance,
		DeliveryDoor:   req.DeliveryDoor,
		Coords:         req.Coords,
		Time:           req.Time,
		PaymentMethod:  req.PaymentMethod,
		Cutlery:        req.Cutlery,
		PromoCode:      req.PromoCode,
		Wishes:         req.Wishes,
		ClientTotal:    req.ClientTotal,
		Items:          make([]services.OrderItemInput, len(req.Items)),
	}

	for i, itemReq := range req.Items {
		if err := itemReq.Validate(); err != nil {
			return response.BadRequest(c, err)
		}

		input.Items[i] = services.OrderItemInput{
			ProductID:          itemReq.ProductID,
			ProductVariationID: itemReq.ProductVariationID,
			Quantity:           itemReq.Quantity,
		}
	}

	order, err := h.service.CreateOrderFromInput(c.Context(), input)
	if err != nil {
		// Return validation errors as bad request
		if errors.Is(err, services.ErrTimeNotAvailable) ||
			errors.Is(err, services.ErrPriceMismatch) ||
			errors.Is(err, services.ErrProductNotFound) {
			return response.BadRequest(c, err)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, order)
}

func (h *OrderHandler) UpdateOrder(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid order id"))
	}

	var req requests.UpdateOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, err)
	}
	req.ID = id

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err)
	}

	order := models.Order{
		ID:             req.ID,
		Name:           req.Name,
		Address:        req.Address,
		Phone:          req.Phone,
		StatusID:       models.Status(req.StatusID),
		DeliveryTypeID: models.DeliveryType(req.DeliveryTypeID),
		PaymentMethod:  req.PaymentMethod,
		Time:           req.Time,
		Items:          make([]models.OrderItem, len(req.Items)),
	}

	for i, itemReq := range req.Items {
		if err := itemReq.Validate(); err != nil {
			return response.BadRequest(c, err)
		}

		order.Items[i] = models.OrderItem{
			ProductID:                  itemReq.ProductID,
			Quantity:                   itemReq.Quantity,
			Price:                      itemReq.Price,
			Weight:                     itemReq.Weight,
			Name:                       itemReq.Name,
			ExternalProductID:          itemReq.ExternalProductID,
			ProductVariationGroupID:    itemReq.ProductVariationGroupID,
			ProductVariationGroupName:  itemReq.ProductVariationGroupName,
			ProductVariationID:         itemReq.ProductVariationID,
			ProductVariationExternalID: itemReq.ProductVariationExternalID,
			ProductVariationName:       itemReq.ProductVariationName,
		}

		if order.Items[i].ProductVariationID != nil && *order.Items[i].ProductVariationID == uuid.Nil {
			order.Items[i].ProductVariationID = nil
		}
		if order.Items[i].ProductVariationGroupID != nil && *order.Items[i].ProductVariationGroupID == uuid.Nil {
			order.Items[i].ProductVariationGroupID = nil
		}
	}

	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = time.Now()
	}

	err = h.service.UpdateOrder(c.Context(), &order)
	if err != nil {
		if errors.Is(err, customerrors.OrderNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, order)
}

func (h *OrderHandler) DeleteOrder(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, errors.New("invalid order id"))
	}

	err = h.service.DeleteOrder(c.Context(), id)
	if err != nil {
		if errors.Is(err, customerrors.OrderNotFound) {
			return response.NotFound(c)
		}
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, nil)
}
