package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
	conf     *config.Config
}

func NewUserUsecase(userRepo domain.UserRepository, conf *config.Config) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		conf:     conf,
	}
}

func (u *userUsecase) Find(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int64, error) {
	users, total, err := u.userRepo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find users",
			zap.Error(err),
			zap.Any("filter", filter),
		)
		return nil, 0, fmt.Errorf("user_usecase.Find: %w", err)
	}

	return users, total, nil
}

func (u *userUsecase) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		config.Log.Error("failed to find user by id",
			zap.Error(err),
			zap.String("user_id", id.String()),
		)
		return nil, fmt.Errorf("user_usecase.FindByID: %w", err)
	}

	return user, nil
}

func (u *userUsecase) Create(ctx context.Context, user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		config.Log.Error("failed to hash password", zap.Error(err))
		return fmt.Errorf("user_usecase.Create (hash): %w", err)
	}
	user.Password = string(hashedPassword)

	if err := u.userRepo.Create(ctx, user); err != nil {
		config.Log.Error("failed to create user in repository",
			zap.Error(err),
			zap.String("email", user.Email),
		)
		return fmt.Errorf("user_usecase.Create: %w", err)
	}

	return nil
}
