package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Branch represents the core branch entity.
type Branch struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate hooks into GORM to auto-generate UUID IDs.
func (b *Branch) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}

// BranchFilter represents criteria for branch data operations.
type BranchFilter struct {
	Search string
	request.PaginationParam
}

// BranchRepository defines the interface for branch data operations.
type BranchRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Branch, error)
	Find(ctx context.Context, filter BranchFilter) ([]*Branch, int64, error)
	Create(ctx context.Context, branch *Branch) error
	Update(ctx context.Context, branch *Branch) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// BranchUsecase defines the interface for branch business logic.
type BranchUsecase interface {
	Find(ctx context.Context, filter BranchFilter) ([]*Branch, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Branch, error)
	Create(ctx context.Context, branch *Branch) error
	Update(ctx context.Context, branch *Branch) error
	Delete(ctx context.Context, id uuid.UUID) error
}
