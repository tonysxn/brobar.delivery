package services

import (
	"context"
	"errors"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/repositories"
)

type ProductVariationGroupService struct {
	repo *repositories.ProductVariationGroupRepository
}

func NewProductVariationGroupService(repo *repositories.ProductVariationGroupRepository) *ProductVariationGroupService {
	return &ProductVariationGroupService{repo: repo}
}

func (s *ProductVariationGroupService) GetAllByProductID(ctx context.Context, productID string) ([]models.ProductVariationGroup, error) {
	parsedID, err := uuid.Parse(productID)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}
	return s.repo.GetAllByProductID(ctx, parsedID)
}

func (s *ProductVariationGroupService) GetByID(ctx context.Context, id uuid.UUID) (*models.ProductVariationGroup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductVariationGroupService) Create(ctx context.Context, g *models.ProductVariationGroup) error {
	err := validation.ValidateStruct(g,
		validation.Field(&g.ProductID, validation.Required),
		validation.Field(&g.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&g.ExternalID, validation.Length(0, 100)),
	)
	if err != nil {
		return err
	}

	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}

	return s.repo.Create(ctx, g)
}

func (s *ProductVariationGroupService) Update(ctx context.Context, id string, updated *models.ProductVariationGroup) (*models.ProductVariationGroup, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid group ID format")
	}

	existing, err := s.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	if updated.Name != "" {
		existing.Name = updated.Name
	}
	if updated.ExternalID != "" {
		existing.ExternalID = updated.ExternalID
	}
	existing.DefaultValue = updated.DefaultValue
	existing.Show = updated.Show
	existing.Required = updated.Required
	if updated.ProductID != uuid.Nil {
		existing.ProductID = updated.ProductID
	}

	err = validation.ValidateStruct(existing,
		validation.Field(&existing.ProductID, validation.Required),
		validation.Field(&existing.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&existing.ExternalID, validation.Length(0, 100)),
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *ProductVariationGroupService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid group ID format")
	}
	return s.repo.Delete(ctx, parsedID)
}

func (s *ProductVariationGroupService) DeleteByProductID(ctx context.Context, productID uuid.UUID) error {
	return s.repo.DeleteByProductID(ctx, productID)
}
