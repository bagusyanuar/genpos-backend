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
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
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
	GetDB() *gorm.DB
}

type MaterialUsecase interface {
	Find(ctx context.Context, filter MaterialFilter) ([]Material, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Material, error)
	Create(ctx context.Context, material *Material, uoms []MaterialUOM) error
}

type MaterialUOM struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MaterialID uuid.UUID `gorm:"type:uuid;not null;index" json:"material_id"`
	UnitID     uuid.UUID `gorm:"type:uuid;not null" json:"unit_id"`
	Multiplier float64   `gorm:"type:decimal(15,4);not null;default:1" json:"multiplier"`
	IsDefault  bool      `gorm:"not null;default:false" json:"is_default"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *MaterialUOM) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type MaterialUOMRepository interface {
	Find(ctx context.Context, materialID uuid.UUID) ([]MaterialUOM, error)
	CreateBatch(ctx context.Context, uoms []MaterialUOM) error
	ReplaceUOMs(ctx context.Context, materialID uuid.UUID, uoms []MaterialUOM) error
}

type MaterialUOMUsecase interface {
	Find(ctx context.Context, materialID uuid.UUID) ([]MaterialUOM, error)
	UpdateUOMs(ctx context.Context, materialID uuid.UUID, uoms []MaterialUOM) error
}
