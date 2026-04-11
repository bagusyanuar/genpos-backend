package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type materialUsecase struct {
	materialRepo domain.MaterialRepository
	uomRepo      domain.MaterialUOMRepository
}

func NewMaterialUsecase(materialRepo domain.MaterialRepository, uomRepo domain.MaterialUOMRepository) domain.MaterialUsecase {
	return &materialUsecase{
		materialRepo: materialRepo,
		uomRepo:      uomRepo,
	}
}

func (u *materialUsecase) Create(ctx context.Context, material *domain.Material, uoms []domain.MaterialUOM) error {
	// 1. Validation for MaterialType
	if material.MaterialType != "RAW" && material.MaterialType != "SEMI_FINISHED" {
		return fmt.Errorf("invalid material type: %s", material.MaterialType)
	}

	// 2. Validation for UOMs: Must have exactly 1 default
	hasDefault := false
	for _, uom := range uoms {
		if uom.IsDefault {
			if hasDefault {
				return fmt.Errorf("multiple default UOMs provided")
			}
			hasDefault = true
		}
	}
	if !hasDefault {
		return fmt.Errorf("default UOM (Base Unit) must be provided")
	}

	// 3. Start Transaction
	tx := u.materialRepo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. Save Material in transaction
	if err := tx.WithContext(ctx).Create(material).Error; err != nil {
		tx.Rollback()
		config.Log.Error("failed to create material in transaction", zap.Error(err))
		return fmt.Errorf("material_uc.Create.Material: %w", err)
	}

	// 5. Prepare UOMs with MaterialID and Save
	for i := range uoms {
		uoms[i].MaterialID = material.ID
	}

	if err := tx.WithContext(ctx).Create(&uoms).Error; err != nil {
		tx.Rollback()
		config.Log.Error("failed to create material uoms in transaction", zap.Error(err))
		return fmt.Errorf("material_uc.Create.UOMs: %w", err)
	}

	// 6. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("material_uc.Create.Commit: %w", err)
	}

	return nil
}

func (u *materialUsecase) FindByID(ctx context.Context, id uuid.UUID) (*domain.Material, error) {
	material, err := u.materialRepo.FindByID(ctx, id)
	if err != nil {
		config.Log.Error("failed to find material by id",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, fmt.Errorf("material_uc.FindByID: %w", err)
	}
	return material, nil
}

func (u *materialUsecase) Find(ctx context.Context, filter domain.MaterialFilter) ([]domain.Material, int64, error) {
	materials, total, err := u.materialRepo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find materials", 
			zap.Error(err), 
			zap.String("search", filter.Search),
		)
		return nil, 0, fmt.Errorf("material_uc.Find: %w", err)
	}

	return materials, total, nil
}
