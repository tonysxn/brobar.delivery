package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/tonysanin/brobar/web-service/internal/models"
)

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (food_rating, service_rating, comment, phone, email, name)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRowx(query,
		review.FoodRating,
		review.ServiceRating,
		review.Comment,
		review.Phone,
		review.Email,
		review.Name,
	).StructScan(review)
}

func (r *ReviewRepository) GetAll() ([]models.Review, error) {
	var reviews []models.Review
	query := `
		SELECT id, food_rating, service_rating, comment, phone, email, name, created_at 
		FROM reviews 
		ORDER BY created_at DESC
	`
	err := r.db.Select(&reviews, query)
	return reviews, err
}

func (r *ReviewRepository) GetByID(id string) (*models.Review, error) {
	var review models.Review
	query := `
		SELECT id, food_rating, service_rating, comment, phone, email, name, created_at 
		FROM reviews 
		WHERE id = $1
	`
	err := r.db.Get(&review, query, id)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) Delete(id string) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
