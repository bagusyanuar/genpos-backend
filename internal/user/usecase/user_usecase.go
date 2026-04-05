package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (u *userUsecase) Find(ctx context.Context, page, limit int) ([]*domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, total, err := u.userRepo.Find(ctx, limit, offset)
	if err != nil {
		config.Log.Error("failed to find users",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("limit", limit),
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
