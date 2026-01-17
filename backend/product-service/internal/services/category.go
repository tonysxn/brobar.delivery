package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAllCategories(ctx)
}

func (s *CategoryService) GetCategoriesWithPagination(ctx context.Context, limit, offset int, orderBy, orderDir string) ([]models.Category, int, error) {
	categories, err := s.repo.GetCategoriesWithPagination(ctx, limit, offset, orderBy, orderDir)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := s.repo.GetCategoriesCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return categories, totalCount, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, id string) (*models.Category, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}
	return s.repo.GetCategoryById(ctx, parsedID)
}

func (s *CategoryService) GetCategoryById(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	return s.repo.GetCategoryById(ctx, id)
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	if category.Slug == "" {
		category.Slug = helpers.GenerateSlug(category.Name)
	}

	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	return s.repo.CreateCategory(ctx, category)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id string, updatedCategory *models.Category) (*models.Category, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}

	existingCategory, err := s.repo.GetCategoryById(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	if updatedCategory.Name != "" {
		existingCategory.Name = updatedCategory.Name
		if updatedCategory.Slug == "" {
			existingCategory.Slug = helpers.GenerateSlug(updatedCategory.Name)
		} else {
			existingCategory.Slug = updatedCategory.Slug
		}
	}
	if updatedCategory.Icon != "" {
		existingCategory.Icon = updatedCategory.Icon
	}
	if updatedCategory.Sort != 0 {
		existingCategory.Sort = updatedCategory.Sort
	}

	err = s.repo.UpdateCategory(ctx, existingCategory)
	if err != nil {
		return nil, err
	}

	return existingCategory, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *CategoryService) GetMenu(ctx context.Context) ([]models.MenuCategory, error) {
	return s.repo.GetMenuData(ctx)
}
