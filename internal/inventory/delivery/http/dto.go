package http

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/inventory/domain"
	"github.com/google/uuid"
)

type MaterialUOMResponse struct {
	ID         uuid.UUID `json:"id"`
	UnitID     uuid.UUID `json:"unit_id"`
	UnitName   string    `json:"unit_name"`
	Multiplier float64   `json:"multiplier"`
	IsDefault  bool      `json:"is_default"`
	Stock      float64   `json:"stock"`
}

type InventoryMaterialResponse struct {
	ID   uuid.UUID `json:"id"`
	SKU  string    `json:"sku"`
	Name string    `json:"name"`
}

type InventoryResponse struct {
	ID             *uuid.UUID                `json:"inventory_id"`
	Material       InventoryMaterialResponse `json:"material"`
	Stock          float64                   `json:"stock"`
	FormattedStock []string                  `json:"formatted_stock"`
	MinStock       float64                   `json:"min_stock"`
	UpdatedAt      *time.Time                `json:"updated_at"`
	UOMs           []MaterialUOMResponse     `json:"uoms"`
}

func formatStock(baseStock float64, uoms []domain.MaterialUOMView) []string {
	if len(uoms) == 0 {
		return []string{fmt.Sprintf("%g", math.Round(baseStock*10000)/10000)}
	}

	sortedUOMs := make([]domain.MaterialUOMView, len(uoms))
	copy(sortedUOMs, uoms)
	// Sort by Multiplier DESCENDING (Largest unit first)
	sort.Slice(sortedUOMs, func(i, j int) bool {
		return sortedUOMs[i].Multiplier > sortedUOMs[j].Multiplier
	})

	isNegative := baseStock < 0
	currentBaseVal := math.Abs(baseStock)
	var parts []string

	for i, uom := range sortedUOMs {
		if uom.Multiplier <= 0 {
			continue // Safety check
		}

		unitVal := currentBaseVal / uom.Multiplier

		if i == len(sortedUOMs)-1 {
			if unitVal > 0 || len(parts) == 0 {
				roundedVal := math.Round(unitVal*10000) / 10000
				parts = append(parts, fmt.Sprintf("%g %s", roundedVal, uom.UnitName))
			}
			break
		}

		intPart := math.Trunc(unitVal)
		if intPart > 0 {
			parts = append(parts, fmt.Sprintf("%g %s", intPart, uom.UnitName))
		}

		remUnit := unitVal - intPart
		currentBaseVal = remUnit * uom.Multiplier

		if currentBaseVal <= 0.0001 {
			break
		}
	}

	if isNegative && len(parts) > 0 {
		parts[0] = "- " + parts[0]
	}
	return parts
}

func ToInventoryResponse(v domain.MaterialInventoryView) InventoryResponse {
	uoms := make([]MaterialUOMResponse, 0)
	for _, uom := range v.UOMs {
		uomStock := float64(0)
		if uom.Multiplier > 0 {
			uomStock = math.Round((v.Stock/uom.Multiplier)*10000) / 10000
		}

		uoms = append(uoms, MaterialUOMResponse{
			ID:         uom.ID,
			UnitID:     uom.UnitID,
			UnitName:   uom.UnitName,
			Multiplier: uom.Multiplier,
			IsDefault:  uom.IsDefault,
			Stock:      uomStock,
		})
	}

	return InventoryResponse{
		ID: v.ID,
		Material: InventoryMaterialResponse{
			ID:   v.MaterialID,
			SKU:  v.MaterialSKU,
			Name: v.MaterialName,
		},
		Stock:          v.Stock,
		FormattedStock: formatStock(v.Stock, v.UOMs),
		MinStock:       v.MinStock,
		UpdatedAt:      v.UpdatedAt,
		UOMs:           uoms,
	}
}

func ToInventoryListResponse(views []domain.MaterialInventoryView) []InventoryResponse {
	res := make([]InventoryResponse, 0)
	for _, v := range views {
		res = append(res, ToInventoryResponse(v))
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
type InventorySummaryResponse struct {
	Material   InventoryMaterialResponse `json:"material"`
	TotalStock float64                   `json:"total_stock"`
}

func ToInventorySummaryResponse(v domain.MaterialStockView) InventorySummaryResponse {
	return InventorySummaryResponse{
		Material: InventoryMaterialResponse{
			ID:   v.ID,
			SKU:  v.SKU,
			Name: v.Name,
		},
		TotalStock: v.TotalStock,
	}
}

func ToInventorySummaryListResponse(views []domain.MaterialStockView) []InventorySummaryResponse {
	res := make([]InventorySummaryResponse, 0)
	for _, v := range views {
		res = append(res, ToInventorySummaryResponse(v))
	}
	return res
}
