package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/inventory/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) domain.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Find(ctx context.Context, filter domain.InventoryFilter) ([]domain.Inventory, int64, error) {
	var inventories []domain.Inventory
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Inventory{})

	// Strict isolated filter for multi-tenancy
	query = query.Where("branch_id = ?", filter.BranchID)

	// Material specific filter if provided
	if filter.MaterialID != uuid.Nil {
		query = query.Where("material_id = ?", filter.MaterialID)
	}

	if filter.Search != "" {
		// Needs to join material table to search by name/SKU
		query = query.Joins("LEFT JOIN materials ON inventories.material_id = materials.id").
			Where("materials.name ILIKE ? OR materials.sku ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("inventory_repo.Find.Count: %w", err)
	}

	err := query.
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Order(filter.GetSort()).
		Find(&inventories).Error

	if err != nil {
		return nil, 0, fmt.Errorf("inventory_repo.Find.Data: %w", err)
	}

	return inventories, total, nil
}

func (r *inventoryRepository) GetSummary(ctx context.Context, branchID uuid.UUID, filter domain.InventoryFilter) ([]domain.MaterialStockView, int64, error) {
	var views []domain.MaterialStockView
	var total int64

	// Start from materials table to ensure we get all materials even if stock is 0
	baseQuery := r.db.WithContext(ctx).Table("materials").Where("materials.deleted_at IS NULL")

	if filter.MaterialID != uuid.Nil {
		baseQuery = baseQuery.Where("materials.id = ?", filter.MaterialID)
	}

	if filter.Search != "" {
		baseQuery = baseQuery.Where("materials.name ILIKE ? OR materials.sku ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// FAST COUNT: Count purely from materials without heavy JOINs
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("inventory_repo.GetSummary.Count: %w", err)
	}

	// For the actual data, clone the baseQuery and attach JOINs
	query := baseQuery.Session(&gorm.Session{})

	// Left join with inventories, filtering by branch_id if provided
	if branchID != uuid.Nil {
		query = query.Joins("LEFT JOIN inventories ON materials.id = inventories.material_id AND inventories.branch_id = ? AND inventories.deleted_at IS NULL", branchID)
	} else {
		query = query.Joins("LEFT JOIN inventories ON materials.id = inventories.material_id AND inventories.deleted_at IS NULL")
	}

	// Calculate stock and aggregate
	err := query.
		Select("materials.id, materials.sku, materials.name, COALESCE(SUM(inventories.stock), 0) as total_stock").
		Group("materials.id, materials.sku, materials.name").
		Order("materials.name ASC"). // Default sorting, can be extended to use filter.GetSort()
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Scan(&views).Error

	if err != nil {
		return nil, 0, fmt.Errorf("inventory_repo.GetSummary.Data: %w", err)
	}

	return views, total, nil
}

func (r *inventoryRepository) UpdateStock(ctx context.Context, move domain.StockMovement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		delta := move.Quantity
		if move.Type == domain.MovementOut || move.Type == domain.MovementDeduction {
			delta = -move.Quantity
		} else if move.Type == domain.MovementAdjust {
			// Assessment: Adjust uses the signed delta calculated in UC
			delta = move.Quantity
		}

		// Refactored to Atomic Upsert (PostgreSQL ON CONFLICT)
		// This is faster and prevents race conditions in stock calculation
		inv := domain.Inventory{
			BranchID:   move.BranchID,
			MaterialID: move.MaterialID,
			Stock:      delta, // Initial stock if record doesn't exist
		}

		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "branch_id"}, {Name: "material_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"stock":      gorm.Expr("inventories.stock + ?", delta),
				"updated_at": time.Now(),
			}),
		}).Create(&inv).Error

		if err != nil {
			return fmt.Errorf("failed to upsert inventory: %w", err)
		}

		// Record movement as audit log
		if err := tx.Create(&move).Error; err != nil {
			return fmt.Errorf("failed to create stock movement: %w", err)
		}

		return nil
	})
}

func (r *inventoryRepository) GetStockMovements(ctx context.Context, filter domain.InventoryFilter) ([]domain.StockMovement, int64, error) {
	var movements []domain.StockMovement
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.StockMovement{}).
		Where("branch_id = ?", filter.BranchID)

	if filter.MaterialID != uuid.Nil {
		query = query.Where("material_id = ?", filter.MaterialID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("inventory_repo.GetStockMovements.Count: %w", err)
	}

	err := query.
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Order("created_at DESC").
		Find(&movements).Error

	if err != nil {
		return nil, 0, fmt.Errorf("inventory_repo.GetStockMovements.Data: %w", err)
	}

	return movements, total, nil
}

func (r *inventoryRepository) RecalibrateStock(ctx context.Context, tx *gorm.DB, materialID uuid.UUID, cf float64) error {
	if tx == nil {
		tx = r.db
	}

	// Update stock and min_stock for all records matching material_id
	err := tx.WithContext(ctx).Model(&domain.Inventory{}).
		Where("material_id = ?", materialID).
		Updates(map[string]interface{}{
			"stock":      gorm.Expr("stock / ?", cf),
			"min_stock":  gorm.Expr("min_stock / ?", cf),
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("inventory_repo.RecalibrateStock: %w", err)
	}

	return nil
}
