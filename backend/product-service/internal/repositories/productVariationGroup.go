package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	customerrors "github.com/tonysanin/brobar/product-service/internal/errors"
	"github.com/tonysanin/brobar/product-service/internal/models"
)

type ProductVariationGroupRepository struct {
	db dbExecutor
}

func NewProductVariationGroupRepository(db *sqlx.DB) *ProductVariationGroupRepository {
	return &ProductVariationGroupRepository{db: db}
}

func (r *ProductVariationGroupRepository) WithTx(tx *sqlx.Tx) *ProductVariationGroupRepository {
	return &ProductVariationGroupRepository{db: tx}
}

func (r *ProductVariationGroupRepository) GetAllByProductID(ctx context.Context, productID uuid.UUID) ([]models.ProductVariationGroup, error) {
	const query = `SELECT * FROM product_variation_groups WHERE product_id = $1`

	var groups []models.ProductVariationGroup

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &groups, query, productID)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get product variation groups: %w", err)
	}

	if groups == nil {
		return []models.ProductVariationGroup{}, nil
	}

	return groups, nil
}

func (r *ProductVariationGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ProductVariationGroup, error) {
	const query = `SELECT * FROM product_variation_groups WHERE id = $1`

	var group models.ProductVariationGroup

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &group, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ProductVariationGroupNotFound
		}
		return nil, fmt.Errorf("failed to get product variation group: %w", err)
	}

	return &group, nil
}

func (r *ProductVariationGroupRepository) Create(ctx context.Context, group *models.ProductVariationGroup) error {
	const query = `
		INSERT INTO product_variation_groups (
			id, product_id, name, external_id, default_value, show, required
		) VALUES (
			:id, :product_id, :name, :external_id, :default_value, :show, :required
		) ON CONFLICT (product_id, external_id) DO UPDATE SET
			name = EXCLUDED.name,
			default_value = EXCLUDED.default_value,
			show = EXCLUDED.show,
			required = EXCLUDED.required
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(ctx, query, group)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to create product variation group: %w", err)
	}

	return nil
}

func (r *ProductVariationGroupRepository) Update(ctx context.Context, group *models.ProductVariationGroup) error {
	const query = `
		UPDATE product_variation_groups SET
			product_id = :product_id,
			name = :name,
			external_id = :external_id,
			default_value = :default_value,
			show = :show,
			required = :required
		WHERE id = :id
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.NamedExecContext(ctx, query, group)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to update product variation group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.ProductVariationGroupNotFound
	}

	return nil
}

func (r *ProductVariationGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM product_variation_groups WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to delete product variation group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.ProductVariationGroupNotFound
	}

	return nil
}

func (r *ProductVariationGroupRepository) DeleteByProductID(ctx context.Context, productID uuid.UUID) error {
	const query = `DELETE FROM product_variation_groups WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to delete product variation groups by product ID: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	return nil
}
