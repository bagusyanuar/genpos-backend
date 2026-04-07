package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Unit represents the core unit entity.
type Unit struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate hooks into GORM to auto-generate UUID IDs.
func (u *Unit) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// UnitFilter represents criteria for unit data operations.
type UnitFilter struct {
	Search string
	request.PaginationParam
}

// UnitRepository defines the interface for unit data operations.
type UnitRepository interface {
	Find(ctx context.Context, filter UnitFilter) ([]*Unit, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Unit, error)
	Create(ctx context.Context, unit *Unit) error
	Update(ctx context.Context, unit *Unit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UnitUsecase defines the interface for unit business logic.
type UnitUsecase interface {
	Find(ctx context.Context, filter UnitFilter) ([]*Unit, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Unit, error)
	Create(ctx context.Context, unit *Unit) error
	Update(ctx context.Context, unit *Unit) error
	Delete(ctx context.Context, id uuid.UUID) error
}
