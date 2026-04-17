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

const (
	MovementIn        = "STOCK_IN"
	MovementOut       = "STOCK_OUT"
	MovementAdjust    = "ADJUSTMENT"
	MovementDeduction = "DEDUCTION"
)

type StockMovement struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BranchID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"branch_id"`
	MaterialID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"material_id"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"`
	Quantity    float64        `gorm:"type:decimal(15,2);not null" json:"quantity"`
	ReferenceID *uuid.UUID     `gorm:"type:uuid" json:"reference_id,omitempty"`
	Note        string         `gorm:"type:text" json:"note"`
	CreatedAt   time.Time      `json:"created_at"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (s *StockMovement) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
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
	// UpdateStock performs atomic inventory update and creates a stock movement record.
	// quantity can be positive (increment) or negative (decrement).
	UpdateStock(ctx context.Context, move StockMovement) error
	GetStockMovements(ctx context.Context, filter InventoryFilter) ([]StockMovement, int64, error)
	RecalibrateStock(ctx context.Context, tx *gorm.DB, materialID uuid.UUID, cf float64) error
}

type InventoryUsecase interface {
	Find(ctx context.Context, filter InventoryFilter) ([]Inventory, int64, error)
	GetSummary(ctx context.Context, branchID uuid.UUID, filter InventoryFilter) ([]MaterialStockView, int64, error)
	AdjustStock(ctx context.Context, move StockMovement) error
	StockOpname(ctx context.Context, branchID uuid.UUID, materialID uuid.UUID, actualStock float64, note string) error
	GetStockMovements(ctx context.Context, filter InventoryFilter) ([]StockMovement, int64, error)
}

// MaterialStockView represents a unified view of Material master data and its total stock from inventories.
type MaterialStockView struct {
	ID         uuid.UUID `json:"id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	TotalStock float64   `json:"total_stock"`
}
