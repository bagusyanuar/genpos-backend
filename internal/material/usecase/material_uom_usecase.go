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
