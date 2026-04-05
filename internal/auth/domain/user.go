package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the authentication user entity.
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

// AuthRepository defines the interface for authentication data operations.
type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// AuthUsecase defines the interface for authentication business logic.
type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (string, error)
}
