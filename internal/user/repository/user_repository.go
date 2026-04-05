package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Find(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.User{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		db = db.Where("email LIKE ? OR username LIKE ?", search, search)
	}

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := db.Limit(filter.GetLimit()).Offset(filter.GetOffset()).Order(filter.GetSort()).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}
