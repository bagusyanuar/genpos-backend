package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Inventory struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	MaterialID uuid.UUID      `gorm:"type:uuid;not null;index" json:"material_id"`
	BranchID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"branch_id"`
	Stock      float64        `gorm:"type:decimal(15,2);not null;default:0" json:"stock"`
	MinStock   float64        `gorm:"type:decimal(15,2);not null;default:0" json:"min_stock"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}

type InventoryFilter struct {
	BranchID   uuid.UUID `json:"branch_id"`
	MaterialID uuid.UUID `json:"material_id"` // Optional filter
	Search     string    `json:"search"`      // Search by Material Name or SKU
	request.PaginationParam
}

type InventoryRepository interface {
	Find(ctx context.Context, filter InventoryFilter) ([]Inventory, int64, error)
	GetSummary(ctx context.Context, branchID uuid.UUID, filter InventoryFilter) ([]MaterialStockView, int64, error)
}

type InventoryUsecase interface {
	Find(ctx context.Context, filter InventoryFilter) ([]Inventory, int64, error)
	GetSummary(ctx context.Context, branchID uuid.UUID, filter InventoryFilter) ([]MaterialStockView, int64, error)
}

// MaterialStockView represents a unified view of Material master data and its total stock from inventories.
type MaterialStockView struct {
	ID         uuid.UUID `json:"id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	TotalStock float64   `json:"total_stock"`
}
