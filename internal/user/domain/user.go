package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the core user entity.
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Username  string         `gorm:"type:varchar(100);unique;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate hooks into GORM to auto-generate UUID IDs.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Find(ctx context.Context, limit, offset int) ([]*User, int64, error)
}

// UserUsecase defines the interface for user business logic.
type UserUsecase interface {
	Find(ctx context.Context, page, limit int) ([]*User, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}
