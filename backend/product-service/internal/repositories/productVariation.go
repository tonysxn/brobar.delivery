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

type ProductVariationRepository struct {
	db dbExecutor
}

func NewProductVariationRepository(db *sqlx.DB) *ProductVariationRepository {
	return &ProductVariationRepository{db: db}
}

func (r *ProductVariationRepository) WithTx(tx *sqlx.Tx) *ProductVariationRepository {
	return &ProductVariationRepository{db: tx}
}

func (r *ProductVariationRepository) GetAllByGroupID(ctx context.Context, groupID uuid.UUID) ([]models.ProductVariation, error) {
	const query = `SELECT * FROM product_variations WHERE group_id = $1`

	var variations []models.ProductVariation

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &variations, query, groupID)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get product variations: %w", err)
	}

	if variations == nil {
		return []models.ProductVariation{}, nil
	}

	return variations, nil
}

func (r *ProductVariationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ProductVariation, error) {
	const query = `SELECT * FROM product_variations WHERE id = $1`

	var variation models.ProductVariation

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &variation, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ProductVariationNotFound
		}
		return nil, fmt.Errorf("failed to get product variation: %w", err)
	}

	return &variation, nil
}

func (r *ProductVariationRepository) Create(ctx context.Context, variation *models.ProductVariation) error {
	const query = `
		INSERT INTO product_variations (
			id, group_id, external_id, default_value, show, name
		) VALUES (
			:id, :group_id, :external_id, :default_value, :show, :name
		) ON CONFLICT (group_id, external_id) DO UPDATE SET
			default_value = EXCLUDED.default_value,
			show = EXCLUDED.show,
			name = EXCLUDED.name
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(ctx, query, variation)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to create product variation: %w", err)
	}

	return nil
}

func (r *ProductVariationRepository) Update(ctx context.Context, variation *models.ProductVariation) error {
	const query = `
		UPDATE product_variations SET
			group_id = :group_id,
			external_id = :external_id,
			default_value = :default_value,
			show = :show,
			name = :name
		WHERE id = :id
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.NamedExecContext(ctx, query, variation)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to update product variation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.ProductVariationNotFound
	}

	return nil
}

func (r *ProductVariationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM product_variations WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to delete product variation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.ProductVariationNotFound
	}

	return nil
}
