package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/repositories"
)

type ProductService struct {
	db                 *sqlx.DB
	repo               *repositories.ProductRepository
	variationGroupRepo *repositories.ProductVariationGroupRepository
	variationRepo      *repositories.ProductVariationRepository
}

func NewProductService(
	db *sqlx.DB,
	repo *repositories.ProductRepository,
	variationGroupRepo *repositories.ProductVariationGroupRepository,
	variationRepo *repositories.ProductVariationRepository,
) *ProductService {
	return &ProductService{
		db:                 db,
		repo:               repo,
		variationGroupRepo: variationGroupRepo,
		variationRepo:      variationRepo,
	}
}

func (s *ProductService) GetProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAllProducts(ctx)
}

func (s *ProductService) GetProductsWithPagination(ctx context.Context, limit, offset int, orderBy, orderDir string) ([]models.Product, int, error) {
	products, err := s.repo.GetProductsWithPagination(ctx, limit, offset, orderBy, orderDir)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := s.repo.GetProductsCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}
	return s.repo.GetProductByID(ctx, parsedID)
}

func (s *ProductService) GetProductById(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product, fileHeader *multipart.FileHeader) error {
	return s.CreateProductWithNested(ctx, product, fileHeader)
}

func (s *ProductService) CreateProductWithNested(ctx context.Context, product *models.Product, fileHeader *multipart.FileHeader) error {
	if product.Slug == "" {
		product.Slug = helpers.GenerateSlug(product.Name)
	}

	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	if fileHeader != nil {
		filename, err := s.uploadFileToFileService(fileHeader)
		if err != nil {
			return err
		}
		product.Image = filename
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Use transactional repositories
	txRepo := s.repo.WithTx(tx)
	txGroupRepo := s.variationGroupRepo.WithTx(tx)
	txVarRepo := s.variationRepo.WithTx(tx)

	if err := txRepo.CreateProduct(ctx, product); err != nil {
		return err
	}

	for i := range product.VariationGroups {
		group := &product.VariationGroups[i]
		if group.ID == uuid.Nil {
			group.ID = uuid.New()
		}
		group.ProductID = product.ID

		if err := txGroupRepo.Create(ctx, group); err != nil {
			return err
		}

		for j := range group.Variations {
			variation := &group.Variations[j]
			if variation.ID == uuid.Nil {
				variation.ID = uuid.New()
			}
			variation.GroupID = group.ID

			if err := txVarRepo.Create(ctx, variation); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, updatedProduct *models.Product, fileHeader *multipart.FileHeader) (*models.Product, error) {
	existingProduct, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updatedProduct.Name != "" {
		existingProduct.Name = updatedProduct.Name
		if updatedProduct.Slug == "" {
			existingProduct.Slug = helpers.GenerateSlug(updatedProduct.Name)
		} else {
			existingProduct.Slug = updatedProduct.Slug
		}
	}
	if updatedProduct.Description != nil {
		existingProduct.Description = updatedProduct.Description
	}
	if updatedProduct.Price != 0 {
		existingProduct.Price = updatedProduct.Price
	}
	existingProduct.Weight = updatedProduct.Weight
	if updatedProduct.ExternalID != "" {
		existingProduct.ExternalID = updatedProduct.ExternalID
	}
	existingProduct.Hidden = updatedProduct.Hidden
	existingProduct.Alcohol = updatedProduct.Alcohol
	existingProduct.Sold = updatedProduct.Sold
	if updatedProduct.CategoryID != uuid.Nil {
		existingProduct.CategoryID = updatedProduct.CategoryID
	}

	if fileHeader != nil {
		resultChan := make(chan struct {
			filename string
			err      error
		})

		go func() {
			filename, err := s.uploadFileToFileService(fileHeader)
			resultChan <- struct {
				filename string
				err      error
			}{filename, err}
		}()

		res := <-resultChan
		if res.err != nil {
			return nil, res.err
		}

		existingProduct.Image = res.filename
	}

	err = s.repo.UpdateProduct(ctx, existingProduct)
	if err != nil {
		return nil, err
	}

	return existingProduct, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) GetProductsByCategory(ctx context.Context, id uuid.UUID) ([]models.Product, error) {
	return s.repo.GetProductsByCategoryID(ctx, id)
}
