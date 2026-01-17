package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/user-service/internal/models"
	"github.com/tonysanin/brobar/user-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (h *UserService) CreateUser(ctx context.Context, user *models.User) error {
	user.RoleID = models.RoleUser

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.PasswordConfirm = ""

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	err = h.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return err
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserById(ctx, id)
}
