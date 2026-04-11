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
