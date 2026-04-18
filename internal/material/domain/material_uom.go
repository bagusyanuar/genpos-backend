package domain

import (
	"context"
	"time"

	unitDomain "github.com/bagusyanuar/genpos-backend/internal/unit/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaterialUOM struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	MaterialID uuid.UUID      `gorm:"type:uuid;not null;index" json:"material_id"`
	UnitID     uuid.UUID      `gorm:"type:uuid;not null" json:"unit_id"`
	Unit       unitDomain.Unit `gorm:"foreignKey:UnitID" json:"unit"`
	Multiplier float64        `gorm:"type:decimal(15,4);not null;default:1" json:"multiplier"`
	IsDefault  bool           `gorm:"not null;default:false" json:"is_default"`
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
	RecalibrateUOMs(ctx context.Context, tx *gorm.DB, materialID uuid.UUID, cf float64, targetUOMID uuid.UUID) error
}

type MaterialUOMUsecase interface {
	Find(ctx context.Context, materialID uuid.UUID) ([]MaterialUOM, error)
	UpdateUOMs(ctx context.Context, materialID uuid.UUID, uoms []MaterialUOM) error
}
