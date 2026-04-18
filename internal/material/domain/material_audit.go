package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaterialAudit struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MaterialID uuid.UUID `gorm:"type:uuid;not null;index" json:"material_id"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`
	Note       string    `gorm:"type:text" json:"note"`
	CreatedBy  uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *MaterialAudit) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}

type MaterialAuditRepository interface {
	Create(ctx context.Context, audit *MaterialAudit) error
	FindByMaterialID(ctx context.Context, materialID uuid.UUID, filter request.PaginationParam) ([]MaterialAudit, int64, error)
}
