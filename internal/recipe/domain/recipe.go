package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recipe struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ProductVariantID uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_variant_id"`
	MaterialID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"material_id"`
	UomID            uuid.UUID      `gorm:"type:uuid;not null" json:"uom_id"`
	Quantity         float64        `gorm:"type:decimal(15,4);not null;default:0" json:"quantity"`
	SubtotalCost     float64        `gorm:"type:decimal(15,4);not null;default:0" json:"subtotal_cost"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (r *Recipe) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

type RecipeFilter struct {
	ProductVariantID *uuid.UUID `json:"product_variant_id"`
	MaterialID       *uuid.UUID `json:"material_id"`
}

type RecipeRepository interface {
	FindByVariantID(ctx context.Context, variantID uuid.UUID) ([]Recipe, error)
	CreateBatch(ctx context.Context, recipes []Recipe) error
	ReplaceByVariantID(ctx context.Context, variantID uuid.UUID, recipes []Recipe) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RecipeUsecase interface {
	GetByVariantID(ctx context.Context, variantID uuid.UUID) ([]Recipe, error)
	SyncRecipe(ctx context.Context, variantID uuid.UUID, recipes []Recipe) error
	CalculateEstimatedCOGS(ctx context.Context, variantID uuid.UUID) (float64, error)
}
