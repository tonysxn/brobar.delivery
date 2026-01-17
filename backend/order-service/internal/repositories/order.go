package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	customerrors "github.com/tonysanin/brobar/order-service/internal/errors"
	"github.com/tonysanin/brobar/order-service/internal/models"
	"log"
	"time"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	const query = `SELECT * FROM orders`
	var orders []models.Order

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return nil, fmt.Errorf("database query timed out")
		}
		log.Printf("failed to get orders: %v", err)
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	if orders == nil {
		return []models.Order{}, nil
	}

	return orders, nil
}

func (r *OrderRepository) GetOrdersWithPagination(ctx context.Context, limit, offset int, orderBy, orderDir string) ([]models.Order, int, error) {
	queryOrders := fmt.Sprintf(`
		SELECT * FROM orders
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, orderBy, orderDir)

	const queryCount = `SELECT COUNT(*) FROM orders`

	var orders []models.Order

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &orders, queryOrders, limit, offset)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, 0, fmt.Errorf("database query timed out")
		}
		return nil, 0, fmt.Errorf("failed to get orders with pagination: %w", err)
	}

	var totalCount int
	err = r.db.GetContext(ctx, &totalCount, queryCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders total count: %w", err)
	}

	if orders == nil {
		orders = []models.Order{}
	}

	return orders, totalCount, nil
}

func (r *OrderRepository) GetOrderById(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	const query = `
		SELECT
			o.id as "order.id",
			o.user_id as "order.user_id",
			o.status_id as "order.status_id",
			o.total_price as "order.total_price",
			o.created_at as "order.created_at",
			o.updated_at as "order.updated_at",
			o.address as "order.address",
			o.entrance as "order.entrance",
			o.floor as "order.floor",
			o.flat as "order.flat",
			o.address_wishes as "order.address_wishes",
			o.name as "order.name",
			o.phone as "order.phone",
			o.time as "order.time",
			o.email as "order.email",
			o.wishes as "order.wishes",
			o.promo as "order.promo",
			o.coords as "order.coords",
			o.cutlery as "order.cutlery",
			o.delivery_cost as "order.delivery_cost",
			o.delivery_door as "order.delivery_door",
			o.delivery_type_id as "order.delivery_type_id",

			oi.id as "items.id",
			oi.order_id as "items.order_id",
			oi.product_id as "items.product_id",
			oi.external_product_id as "items.external_product_id",
			oi.name as "items.name",
			oi.price as "items.price",
			oi.quantity as "items.quantity",
			oi.total_price as "items.total_price",
			oi.weight as "items.weight",
			oi.total_weight as "items.total_weight",

			oi.product_variation_group_id as "items.product_variation_group_id",
			oi.product_variation_group_name as "items.product_variation_group_name",
			oi.product_variation_id as "items.product_variation_id",
			oi.product_variation_external_id as "items.product_variation_external_id",
			oi.product_variation_name as "items.product_variation_name"

		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	rows, err := r.db.QueryxContext(ctx, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.OrderNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var order models.Order
	itemsMap := make(map[uuid.UUID]models.OrderItem)

	first := true
	for rows.Next() {
		var (
			o  models.Order
			oi models.OrderItem
		)
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.StatusID,
			&o.TotalPrice,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Address,
			&o.Entrance,
			&o.Floor,
			&o.Flat,
			&o.AddressWishes,
			&o.Name,
			&o.Phone,
			&o.Time,
			&o.Email,
			&o.Wishes,
			&o.Promo,
			&o.Coords,
			&o.Cutlery,
			&o.DeliveryCost,
			&o.DeliveryDoor,
			&o.DeliveryTypeID,

			&oi.ID,
			&oi.OrderID,
			&oi.ProductID,
			&oi.ExternalProductID,
			&oi.Name,
			&oi.Price,
			&oi.Quantity,
			&oi.TotalPrice,
			&oi.Weight,
			&oi.TotalWeight,

			&oi.ProductVariationGroupID,
			&oi.ProductVariationGroupName,
			&oi.ProductVariationID,
			&oi.ProductVariationExternalID,
			&oi.ProductVariationName,
		)
		if err != nil {
			return nil, err
		}

		if first {
			order = o
			first = false
		}

		if oi.ID != uuid.Nil {
			itemsMap[oi.ID] = oi
		}
	}

	order.Items = make([]models.OrderItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		order.Items = append(order.Items, item)
	}

	if first {
		return nil, customerrors.OrderNotFound
	}

	return &order, nil
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	const query = `
		INSERT INTO orders (
			id, user_id, status_id, total_price, created_at, updated_at,
			address, entrance, floor, flat, address_wishes,
			name, phone, time, email, wishes, promo, coords,
			cutlery, delivery_cost, delivery_door, delivery_type_id
		) VALUES (
			:id, :user_id, :status_id, :total_price, :created_at, :updated_at,
			:address, :entrance, :floor, :flat, :address_wishes,
			:name, :phone, :time, :email, :wishes, :promo, :coords,
			:cutlery, :delivery_cost, :delivery_door, :delivery_type_id
		)
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	// Если created_at и updated_at не заполнены, ставим now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = time.Now()
	}

	_, err := r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return fmt.Errorf("database query timed out")
		}
		log.Printf("failed to create order: %v", err)
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	const query = `
		UPDATE orders SET
			user_id = :user_id,
			status_id = :status_id,
			total_price = :total_price,
			updated_at = :updated_at,
			address = :address,
			entrance = :entrance,
			floor = :floor,
			flat = :flat,
			address_wishes = :address_wishes,
			name = :name,
			phone = :phone,
			time = :time,
			email = :email,
			wishes = :wishes,
			promo = :promo,
			coords = :coords,
			cutlery = :cutlery,
			delivery_cost = :delivery_cost,
			delivery_door = :delivery_door,
			delivery_type_id = :delivery_type_id
		WHERE id = :id
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = time.Now()
	}

	result, err := r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return fmt.Errorf("database query timed out")
		}
		log.Printf("failed to update order: %v", err)
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get affected rows: %v", err)
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.OrderNotFound
	}

	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM orders WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return fmt.Errorf("database query timed out")
		}
		log.Printf("failed to delete order: %v", err)
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get affected rows: %v", err)
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.OrderNotFound
	}

	return nil
}
