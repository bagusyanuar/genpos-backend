package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Material struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CategoryID   *uuid.UUID     `gorm:"type:uuid" json:"category_id"`
	SKU          string         `gorm:"type:varchar(50);not null" json:"sku"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Description  *string        `gorm:"type:text" json:"description"`
	MaterialType string         `gorm:"type:varchar(50)" json:"material_type"`
	ImageURL     *string        `gorm:"type:text" json:"image_url"`
	BaseCost     float64        `gorm:"type:decimal(15,4);not null;default:0" json:"base_cost"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UOMs []MaterialUOM `gorm:"foreignKey:MaterialID" json:"uoms"`
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
	Create(ctx context.Context, material *Material) error
	Update(ctx context.Context, material *Material) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDB() *gorm.DB
}

type MaterialUsecase interface {
	Find(ctx context.Context, filter MaterialFilter) ([]Material, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Material, error)
	Create(ctx context.Context, material *Material, uoms []MaterialUOM) error
	Update(ctx context.Context, material *Material) error
	UpdateImage(ctx context.Context, id uuid.UUID, imageURL string) error
	Delete(ctx context.Context, id uuid.UUID) error
	RecalibrateUOM(ctx context.Context, materialID uuid.UUID, targetUOMID uuid.UUID, userID uuid.UUID) error
}
