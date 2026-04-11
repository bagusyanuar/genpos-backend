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

func (u *inventoryUsecase) Find(ctx context.Context, filter domain.InventoryFilter) ([]domain.Inventory, int64, error) {
	inventories, total, err := u.repo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find inventories",
			zap.Error(err),
			zap.String("branch_id", filter.BranchID.String()),
			zap.String("material_id", filter.MaterialID.String()),
		)
		return nil, 0, fmt.Errorf("inventory_uc.Find: %w", err)
	}

	return inventories, total, nil
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
