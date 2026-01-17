package services

import (
	"fmt"

	pkgServices "github.com/tonysanin/brobar/pkg/services"
	"github.com/tonysanin/brobar/web-service/internal/models"
	"github.com/tonysanin/brobar/web-service/internal/repositories"
)

type ReviewService struct {
	repo *repositories.ReviewRepository
}

func NewReviewService(repo *repositories.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) CreateReview(review *models.Review) error {
	if err := s.repo.Create(review); err != nil {
		return err
	}

	// Send notification to Telegram
	go s.sendTelegramNotification(review)

	return nil
}

func (s *ReviewService) GetAllReviews() ([]models.Review, error) {
	return s.repo.GetAll()
}

func (s *ReviewService) GetReview(id string) (*models.Review, error) {
	return s.repo.GetByID(id)
}

func (s *ReviewService) DeleteReview(id string) error {
	return s.repo.Delete(id)
}

func (s *ReviewService) sendTelegramNotification(review *models.Review) {
	stars := func(n int) string {
		result := ""
		for i := 0; i < n; i++ {
			result += "⭐"
		}
		return result
	}

	message := fmt.Sprintf(`🆕 Новий відгук!

🍽 Страви: %s
🛎 Сервіс: %s`, stars(review.FoodRating), stars(review.ServiceRating))

	if review.Comment != "" {
		message += fmt.Sprintf("\n\n💬 Коментар:\n%s", review.Comment)
	}

	if review.Name != nil && *review.Name != "" {
		message += fmt.Sprintf("\n\n👤 Ім'я: %s", *review.Name)
	}
	if review.Phone != nil && *review.Phone != "" {
		message += fmt.Sprintf("\n📱 Телефон: %s", *review.Phone)
	}
	if review.Email != nil && *review.Email != "" {
		message += fmt.Sprintf("\n📧 Email: %s", *review.Email)
	}

	pkgServices.SendTelegramMessage(message)
}
