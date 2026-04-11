package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Material struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SKU       string         `gorm:"type:varchar(50);not null" json:"sku"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (m *Material) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

type MaterialFilter struct {
	Search string `json:"search"`
	request.PaginationParam
}

type MaterialRepository interface {
	Find(ctx context.Context, filter MaterialFilter) ([]Material, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Material, error)
}

type MaterialUsecase interface {
	Find(ctx context.Context, filter MaterialFilter) ([]Material, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Material, error)
}

type MaterialUOM struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MaterialID uuid.UUID `gorm:"type:uuid;not null;index" json:"material_id"`
	UnitID     uuid.UUID `gorm:"type:uuid;not null" json:"unit_id"`
	Multiplier float64   `gorm:"type:decimal(15,4);not null;default:1" json:"multiplier"`
	IsDefault  bool      `gorm:"not null;default:false" json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MaterialUOMRepository interface {
	Find(ctx context.Context, materialID uuid.UUID) ([]MaterialUOM, error)
}

type MaterialUOMUsecase interface {
	Find(ctx context.Context, materialID uuid.UUID) ([]MaterialUOM, error)
}
