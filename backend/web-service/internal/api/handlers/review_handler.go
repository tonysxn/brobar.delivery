package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/pkg/response"
	"github.com/tonysanin/brobar/web-service/internal/models"
	"github.com/tonysanin/brobar/web-service/internal/services"
)

type ReviewHandler struct {
	service *services.ReviewService
}

func NewReviewHandler(service *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

type CreateReviewRequest struct {
	FoodRating    int     `json:"food_rating"`
	ServiceRating int     `json:"service_rating"`
	Comment       string  `json:"comment"`
	Phone         *string `json:"phone,omitempty"`
	Email         *string `json:"email,omitempty"`
	Name          *string `json:"name,omitempty"`
}

func (h *ReviewHandler) CreateReview(c fiber.Ctx) error {
	var req CreateReviewRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid request body"))
	}

	if req.FoodRating < 1 || req.FoodRating > 5 {
		return response.Error(c, fiber.StatusBadRequest, errors.New("food_rating must be between 1 and 5"))
	}
	if req.ServiceRating < 1 || req.ServiceRating > 5 {
		return response.Error(c, fiber.StatusBadRequest, errors.New("service_rating must be between 1 and 5"))
	}

	review := &models.Review{
		FoodRating:    req.FoodRating,
		ServiceRating: req.ServiceRating,
		Comment:       req.Comment,
		Phone:         req.Phone,
		Email:         req.Email,
		Name:          req.Name,
	}

	if err := h.service.CreateReview(review); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("failed to create review"))
	}

	return response.Success(c, review)
}

func (h *ReviewHandler) GetReviews(c fiber.Ctx) error {
	reviews, err := h.service.GetAllReviews()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("failed to fetch reviews"))
	}
	return response.Success(c, reviews)
}

func (h *ReviewHandler) GetReview(c fiber.Ctx) error {
	id := c.Params("id")
	review, err := h.service.GetReview(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, errors.New("review not found"))
	}
	return response.Success(c, review)
}

func (h *ReviewHandler) DeleteReview(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteReview(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("failed to delete review"))
	}
	return response.Success(c, fiber.Map{"status": "deleted"})
}
