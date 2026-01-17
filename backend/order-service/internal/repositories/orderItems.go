package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tonysanin/brobar/order-service/internal/models"
)

type OrderItemRepository struct {
	db *sqlx.DB
}

func NewOrderItemRepository(db *sqlx.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) CreateOrderItem(ctx context.Context, item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (
			id, order_id, product_id, external_product_id, name, price, quantity, total_price, weight, total_weight
		) VALUES (
			:id, :order_id, :product_id, :external_product_id, :name, :price, :quantity, :total_price, :weight, :total_weight
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *OrderItemRepository) CreateOrderItems(ctx context.Context, items []models.OrderItem) error {
	for i := range items {
		if err := r.CreateOrderItem(ctx, &items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderItemRepository) DeleteOrderItemsByOrderID(ctx context.Context, orderID uuid.UUID) error {
	query := `DELETE FROM order_items WHERE order_id = $1`
	_, err := r.db.ExecContext(ctx, query, orderID)
	return err
}
