package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	customerrors "github.com/tonysanin/brobar/product-service/internal/errors"
	"github.com/tonysanin/brobar/product-service/internal/models"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	const query = `SELECT * FROM categories`
	var categories []models.Category

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return nil, fmt.Errorf("database query timed out")
		}
		log.Printf("failed to get categories: %v", err)
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	if categories == nil {
		return []models.Category{}, nil
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoriesWithPagination(ctx context.Context, limit, offset int, orderBy, orderDir string) ([]models.Category, error) {
	query := fmt.Sprintf(`SELECT * FROM categories ORDER BY %s %s LIMIT $1 OFFSET $2`, orderBy, orderDir)

	var categories []models.Category

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &categories, query, limit, offset)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get paginated categories: %w", err)
	}

	if categories == nil {
		return []models.Category{}, nil
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoriesCount(ctx context.Context) (int, error) {
	const query = `SELECT COUNT(*) FROM categories`

	var count int

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return 0, fmt.Errorf("database query timed out")
		}
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return count, nil
}

func (r *CategoryRepository) GetCategoryById(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	const query = `SELECT * FROM categories WHERE id = $1`
	var category models.Category

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return nil, fmt.Errorf("database query timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.CategoryNotFound
		}
		log.Printf("failed to get category: %v", err)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	const query = `
		INSERT INTO categories (
			id, name, slug, icon, sort
		) VALUES (
			:id, :name, :slug, :icon, :sort
		)`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return fmt.Errorf("database query timed out")
		}
		log.Printf("failed to create category: %v", err)
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	const query = `
		UPDATE categories SET
			name = :name,
			slug = :slug,
			icon = :icon,
			sort = :sort
		WHERE id = :id
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return fmt.Errorf("database query timed out")
		}
		log.Printf("failed to update category: %v", err)
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get affected rows: %v", err)
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.CategoryNotFound
	}

	return nil
}

func (r *CategoryRepository) GetMenuData(ctx context.Context) ([]models.MenuCategory, error) {
	// Query to get all categories, products, variation groups, and variations
	// We'll use multiple queries for clarity and then assemble in memory

	// Get all categories ordered by sort
	categoriesQuery := `SELECT * FROM categories ORDER BY sort, name`
	var categories []models.MenuCategory

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.SelectContext(ctx, &categories, categoriesQuery)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get categories for menu: %w", err)
	}

	if len(categories) == 0 {
		return []models.MenuCategory{}, nil
	}

	// Get all products ordered by category and sort
	productsQuery := `SELECT * FROM products ORDER BY category_id, sort, name`
	var products []models.MenuProduct

	err = r.db.SelectContext(ctx, &products, productsQuery)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get products for menu: %w", err)
	}

	// Get all variation groups where show = true
	groupsQuery := `SELECT * FROM product_variation_groups WHERE show = true ORDER BY product_id, name`
	var groups []models.MenuVariationGroup

	err = r.db.SelectContext(ctx, &groups, groupsQuery)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get variation groups for menu: %w", err)
	}

	// Get all variations where show = true
	variationsQuery := `SELECT * FROM product_variations WHERE show = true ORDER BY group_id, name`
	var variations []models.MenuVariation

	err = r.db.SelectContext(ctx, &variations, variationsQuery)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out")
		}
		return nil, fmt.Errorf("failed to get variations for menu: %w", err)
	}

	// Build the tree structure
	// Map variations to their groups
	groupVariationsMap := make(map[uuid.UUID][]models.MenuVariation)
	for _, variation := range variations {
		groupVariationsMap[variation.GroupID] = append(groupVariationsMap[variation.GroupID], variation)
	}

	// Assign variations to groups
	for i := range groups {
		groups[i].Variations = groupVariationsMap[groups[i].ID]
		if groups[i].Variations == nil {
			groups[i].Variations = []models.MenuVariation{}
		}
	}

	// Map groups to their products
	productGroupsMap := make(map[uuid.UUID][]models.MenuVariationGroup)
	for _, group := range groups {
		productGroupsMap[group.ProductID] = append(productGroupsMap[group.ProductID], group)
	}

	// Assign groups to products
	for i := range products {
		products[i].VariationGroups = productGroupsMap[products[i].ID]
		if products[i].VariationGroups == nil {
			products[i].VariationGroups = []models.MenuVariationGroup{}
		}
	}

	// Map products to their categories
	categoryProductsMap := make(map[uuid.UUID][]models.MenuProduct)
	for _, product := range products {
		categoryProductsMap[product.CategoryID] = append(categoryProductsMap[product.CategoryID], product)
	}

	// Assign products to categories and filter
	filteredCategories := make([]models.MenuCategory, 0)
	for i := range categories {
		categories[i].Products = categoryProductsMap[categories[i].ID]
		if categories[i].Products == nil {
			categories[i].Products = []models.MenuProduct{}
		}

		if len(categories[i].Products) > 0 {
			filteredCategories = append(filteredCategories, categories[i])
		}
	}

	return filteredCategories, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {

	const query = `DELETE FROM categories WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return fmt.Errorf("database query timed out")
		}
		log.Printf("failed to delete category: %v", err)
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get affected rows: %v", err)
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return customerrors.CategoryNotFound
	}

	return nil
}
