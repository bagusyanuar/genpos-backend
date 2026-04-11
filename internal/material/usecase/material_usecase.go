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
}

func NewMaterialUsecase(materialRepo domain.MaterialRepository) domain.MaterialUsecase {
	return &materialUsecase{
		materialRepo: materialRepo,
	}
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
