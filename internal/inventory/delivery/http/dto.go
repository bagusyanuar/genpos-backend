package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/inventory/domain"
	"github.com/google/uuid"
)

type InventoryResponse struct {
	ID         uuid.UUID `json:"id"`
	MaterialID uuid.UUID `json:"material_id"`
	BranchID   uuid.UUID `json:"branch_id"`
	Stock      float64   `json:"stock"`
	MinStock   float64   `json:"min_stock"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToInventoryResponse(i domain.Inventory) InventoryResponse {
	return InventoryResponse{
		ID:         i.ID,
		MaterialID: i.MaterialID,
		BranchID:   i.BranchID,
		Stock:      i.Stock,
		MinStock:   i.MinStock,
		CreatedAt:  i.CreatedAt,
		UpdatedAt:  i.UpdatedAt,
	}
}

func ToInventoryListResponse(inventories []domain.Inventory) []InventoryResponse {
	res := make([]InventoryResponse, 0)
	for _, i := range inventories {
		res = append(res, ToInventoryResponse(i))
	}
	return res
}

type StockAdjustmentRequest struct {
	MaterialID uuid.UUID `json:"material_id" validate:"required"`
	BranchID   uuid.UUID `json:"branch_id" validate:"required"`
	Type       string    `json:"type" validate:"required,oneof=STOCK_IN STOCK_OUT"`
	Quantity   float64   `json:"quantity" validate:"required,gt=0"`
	Note       string    `json:"note"`
}

type StockOpnameRequest struct {
	MaterialID  uuid.UUID `json:"material_id" validate:"required"`
	BranchID    uuid.UUID `json:"branch_id" validate:"required"`
	ActualStock float64   `json:"actual_stock" validate:"required,gte=0"`
	Note        string    `json:"note"`
}

type StockMovementResponse struct {
	ID          uuid.UUID  `json:"id"`
	MaterialID  uuid.UUID  `json:"material_id"`
	BranchID    uuid.UUID  `json:"branch_id"`
	Type        string     `json:"type"`
	Quantity    float64    `json:"quantity"`
	ReferenceID *uuid.UUID `json:"reference_id,omitempty"`
	Note        string     `json:"note"`
	CreatedAt   time.Time  `json:"created_at"`
}

func ToStockMovementResponse(m domain.StockMovement) StockMovementResponse {
	return StockMovementResponse{
		ID:          m.ID,
		MaterialID:  m.MaterialID,
		BranchID:    m.BranchID,
		Type:        m.Type,
		Quantity:    m.Quantity,
		ReferenceID: m.ReferenceID,
		Note:        m.Note,
		CreatedAt:   m.CreatedAt,
	}
}

func ToStockMovementListResponse(movements []domain.StockMovement) []StockMovementResponse {
	res := make([]StockMovementResponse, 0)
	for _, m := range movements {
		res = append(res, ToStockMovementResponse(m))
	}
	return res
}
