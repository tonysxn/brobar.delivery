package services

import (
	"context"
	"errors"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/product-service/internal/models"
	"github.com/tonysanin/brobar/product-service/internal/repositories"
)

type ProductVariationService struct {
	repo *repositories.ProductVariationRepository
}

func NewProductVariationService(repo *repositories.ProductVariationRepository) *ProductVariationService {
	return &ProductVariationService{repo: repo}
}

func (s *ProductVariationService) GetAllByGroupID(ctx context.Context, groupID string) ([]models.ProductVariation, error) {
	parsedID, err := uuid.Parse(groupID)
	if err != nil {
		return nil, errors.New("invalid group ID format")
	}
	return s.repo.GetAllByGroupID(ctx, parsedID)
}

func (s *ProductVariationService) GetByID(ctx context.Context, id string) (*models.ProductVariation, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid variation ID format")
	}
	return s.repo.GetByID(ctx, parsedID)
}

func (s *ProductVariationService) Create(ctx context.Context, v *models.ProductVariation) error {
	err := validation.ValidateStruct(v,
		validation.Field(&v.GroupID, validation.Required),
		validation.Field(&v.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&v.ExternalID, validation.Length(0, 100)),
	)
	if err != nil {
		return err
	}

	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}

	return s.repo.Create(ctx, v)
}

func (s *ProductVariationService) Update(ctx context.Context, id string, updated *models.ProductVariation) (*models.ProductVariation, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid variation ID format")
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
	if updated.GroupID != uuid.Nil {
		existing.GroupID = updated.GroupID
	}

	err = validation.ValidateStruct(existing,
		validation.Field(&existing.GroupID, validation.Required),
		validation.Field(&existing.Name, validation.Required, validation.Length(1, 255)),
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

func (s *ProductVariationService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid variation ID format")
	}
	return s.repo.Delete(ctx, parsedID)
}
