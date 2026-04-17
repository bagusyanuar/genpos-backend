package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/inventory/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type inventoryUsecase struct {
	repo domain.InventoryRepository
}

func NewInventoryUsecase(repo domain.InventoryRepository) domain.InventoryUsecase {
	return &inventoryUsecase{repo: repo}
}

func (u *inventoryUsecase) Find(ctx context.Context, filter domain.InventoryFilter) ([]domain.MaterialInventoryView, int64, error) {
	views, total, err := u.repo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find inventories",
			zap.Error(err),
			zap.String("branch_id", filter.BranchID.String()),
			zap.String("material_id", filter.MaterialID.String()),
		)
		return nil, 0, fmt.Errorf("inventory_uc.Find: %w", err)
	}

	return views, total, nil
}

func (u *inventoryUsecase) GetSummary(ctx context.Context, branchID uuid.UUID, filter domain.InventoryFilter) ([]domain.MaterialStockView, int64, error) {
	views, total, err := u.repo.GetSummary(ctx, branchID, filter)
	if err != nil {
		config.Log.Error("failed to get inventory summary",
			zap.Error(err),
			zap.String("branch_id", branchID.String()),
		)
		return nil, 0, fmt.Errorf("inventory_uc.GetSummary: %w", err)
	}

	return views, total, nil
}

func (u *inventoryUsecase) AdjustStock(ctx context.Context, move domain.StockMovement) error {
	// Simple validation: quantity must be positive
	if move.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if err := u.repo.UpdateStock(ctx, move); err != nil {
		config.Log.Error("failed to adjust stock",
			zap.Error(err),
			zap.String("branch_id", move.BranchID.String()),
			zap.String("material_id", move.MaterialID.String()),
		)
		return fmt.Errorf("inventory_uc.AdjustStock: %w", err)
	}

	return nil
}

func (u *inventoryUsecase) StockOpname(ctx context.Context, branchID uuid.UUID, materialID uuid.UUID, actualStock float64, note string) error {
	// To perform opname, we first find current system stock
	filter := domain.InventoryFilter{
		BranchID:   branchID,
		MaterialID: materialID,
	}

	views, _, err := u.repo.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to get current stock for opname: %w", err)
	}

	var currentStock float64
	if len(views) > 0 {
		currentStock = views[0].Stock
	}

	delta := actualStock - currentStock
	// If delta is 0, we don't need to do anything, but recording an adjustment of 0 is fine for audit.

	move := domain.StockMovement{
		BranchID:   branchID,
		MaterialID: materialID,
		Type:       domain.MovementAdjust,
		Quantity:   delta, // We send signed delta to repo for ADJUSTMENT type
		Note:       fmt.Sprintf("[Opname] Actual: %.2f, System: %.2f. Note: %s", actualStock, currentStock, note),
	}

	if err := u.repo.UpdateStock(ctx, move); err != nil {
		config.Log.Error("failed to record stock opname",
			zap.Error(err),
			zap.String("branch_id", branchID.String()),
			zap.String("material_id", materialID.String()),
		)
		return fmt.Errorf("inventory_uc.StockOpname: %w", err)
	}

	return nil
}

func (u *inventoryUsecase) GetStockMovements(ctx context.Context, filter domain.InventoryFilter) ([]domain.StockMovement, int64, error) {
	movements, total, err := u.repo.GetStockMovements(ctx, filter)
	if err != nil {
		config.Log.Error("failed to get stock movements",
			zap.Error(err),
			zap.String("branch_id", filter.BranchID.String()),
		)
		return nil, 0, fmt.Errorf("inventory_uc.GetStockMovements: %w", err)
	}

	return movements, total, nil
}
