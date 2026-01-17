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

type ProductRepository struct {
	db dbExecutor
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) WithTx(tx *sqlx.Tx) *ProductRepository {
	return &ProductRepository{db: tx}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	const query = `SELECT * FROM products`
	var products []models.Product

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &products, query)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	if products == nil {
		return []models.Product{}, nil
	}

	return products, nil
}

func (r *ProductRepository) GetProductsWithPagination(ctx context.Context, limit, offset int, orderBy, orderDir string) ([]models.Product, error) {
	query := fmt.Sprintf(`SELECT * FROM products ORDER BY %s %s LIMIT $1 OFFSET $2`, orderBy, orderDir)

	var products []models.Product

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &products, query, limit, offset)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get paginated products: %w", err)
	}

	if products == nil {
		return []models.Product{}, nil
	}

	return products, nil
}

func (r *ProductRepository) GetProductsCount(ctx context.Context) (int, error) {
	const query = `SELECT COUNT(*) FROM products`

	var count int

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return 0, fmt.Errorf("database query timed out")
		}
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	const query = `SELECT * FROM products WHERE id = $1`
	var product models.Product

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	const query = `
		INSERT INTO products (
			id, name, description, slug, price, weight, category_id, external_id,
			hidden, alcohol, sold, image
		) VALUES (
			:id, :name, :description, :slug, :price, :weight, :category_id, :external_id,
			:hidden, :alcohol, :sold, :image
		) ON CONFLICT (external_id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			slug = EXCLUDED.slug,
			price = EXCLUDED.price,
			weight = EXCLUDED.weight,
			category_id = EXCLUDED.category_id,
			hidden = EXCLUDED.hidden,
			alcohol = EXCLUDED.alcohol,
			sold = EXCLUDED.sold,
			image = EXCLUDED.image
		`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	const query = `
		UPDATE products SET
			name = :name,
			description = :description,
			slug = :slug,
			price = :price,
			weight = :weight,
			category_id = :category_id,
			external_id = :external_id,
			hidden = :hidden,
			alcohol = :alcohol,
			sold = :sold,
			image = :image
		WHERE id = :id
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.ProductNotFound
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM products WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("database query timed out")
		}
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.ProductNotFound
	}

	return nil
}

func (r *ProductRepository) GetProductsByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]models.Product, error) {
	const query = `SELECT * FROM products WHERE category_id = $1`
	var products []models.Product

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &products, query, categoryID)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	if products == nil {
		return []models.Product{}, nil
	}

	return products, nil
}
