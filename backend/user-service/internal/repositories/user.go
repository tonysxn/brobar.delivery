package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	customerrors "github.com/tonysanin/brobar/user-service/internal/errors"
	"github.com/tonysanin/brobar/user-service/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `SELECT * FROM users WHERE email = $1`
	var user models.User

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return nil, fmt.Errorf("database query timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.UserNotFound
		}
		log.Printf("failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const query = `SELECT * FROM users WHERE id = $1`
	var user models.User

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return nil, fmt.Errorf("database query timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.UserNotFound
		}
		log.Printf("failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, category *models.User) error {
	const query = `
		INSERT INTO users (
			id, role_id, email, password, name, address, address_coords, phone, promo_card
		) VALUES (
			:id, :role_id, :email, :password, :name, :address, :address_coords, :phone, :promo_card
		)`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	_, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("database query timed out")
			return fmt.Errorf("database query timed out")
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return customerrors.UserAlreadyExists
			}
		}

		log.Printf("failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
