package usecase

import (
	"context"
	"fmt"
	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type materialUOMUsecase struct {
	repo domain.MaterialUOMRepository
}

func NewMaterialUOMUsecase(repo domain.MaterialUOMRepository) domain.MaterialUOMUsecase {
	return &materialUOMUsecase{repo: repo}
}

func (u *materialUOMUsecase) Find(ctx context.Context, materialID uuid.UUID) ([]domain.MaterialUOM, error) {
	uoms, err := u.repo.Find(ctx, materialID)
	if err != nil {
		config.Log.Error("failed to find material UOMs",
			zap.Error(err),
			zap.String("material_id", materialID.String()),
		)
		return nil, fmt.Errorf("material_uom_uc.Find: %w", err)
	}

	return uoms, nil
}

func (u *materialUOMUsecase) UpdateUOMs(ctx context.Context, materialID uuid.UUID, uoms []domain.MaterialUOM) error {
	// 1. Validation: Must have exactly 1 default (Base Unit)
	hasDefault := false
	for _, uom := range uoms {
		if uom.IsDefault {
			if hasDefault {
				return fmt.Errorf("multiple default UOMs provided")
			}
			hasDefault = true

			// Base Unit multiplier must be 1
			if uom.Multiplier != 1 {
				return fmt.Errorf("base unit (default) multiplier must be 1")
			}
		}
	}

	if !hasDefault {
		return fmt.Errorf("default UOM (Base Unit) must be provided")
	}

	// 2. Execute Replace (Transactional Sync)
	if err := u.repo.ReplaceUOMs(ctx, materialID, uoms); err != nil {
		config.Log.Error("failed to sync material UOMs",
			zap.Error(err),
			zap.String("material_id", materialID.String()),
		)
		return fmt.Errorf("material_uom_uc.UpdateUOMs: %w", err)
	}

	return nil
}
