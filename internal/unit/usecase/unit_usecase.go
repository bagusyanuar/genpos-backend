package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/unit/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type unitUsecase struct {
	unitRepo domain.UnitRepository
}

func NewUnitUsecase(unitRepo domain.UnitRepository) domain.UnitUsecase {
	return &unitUsecase{
		unitRepo: unitRepo,
	}
}

func (u *unitUsecase) Find(ctx context.Context, filter domain.UnitFilter) ([]*domain.Unit, int64, error) {
	units, total, err := u.unitRepo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find units",
			zap.Error(err),
			zap.Any("filter", filter),
		)
		return nil, 0, fmt.Errorf("unit_uc.Find: %w", err)
	}

	return units, total, nil
}

func (u *unitUsecase) FindByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error) {
	unit, err := u.unitRepo.FindByID(ctx, id)
	if err != nil {
		config.Log.Error("failed to find unit by id",
			zap.Error(err),
			zap.String("unit_id", id.String()),
		)
		return nil, fmt.Errorf("unit_uc.FindByID: %w", err)
	}

	return unit, nil
}

func (u *unitUsecase) Create(ctx context.Context, unit *domain.Unit) error {
	if err := u.unitRepo.Create(ctx, unit); err != nil {
		config.Log.Error("failed to create unit",
			zap.Error(err),
			zap.String("name", unit.Name),
		)
		return fmt.Errorf("unit_uc.Create: %w", err)
	}

	return nil
}

func (u *unitUsecase) Update(ctx context.Context, unit *domain.Unit) error {
	if err := u.unitRepo.Update(ctx, unit); err != nil {
		config.Log.Error("failed to update unit",
			zap.Error(err),
			zap.String("unit_id", unit.ID.String()),
		)
		return fmt.Errorf("unit_uc.Update: %w", err)
	}

	return nil
}

func (u *unitUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.unitRepo.Delete(ctx, id); err != nil {
		config.Log.Error("failed to delete unit",
			zap.Error(err),
			zap.String("unit_id", id.String()),
		)
		return fmt.Errorf("unit_uc.Delete: %w", err)
	}

	return nil
}
