package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/branch/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type branchUsecase struct {
	branchRepo domain.BranchRepository
}

func NewBranchUsecase(branchRepo domain.BranchRepository) domain.BranchUsecase {
	return &branchUsecase{
		branchRepo: branchRepo,
	}
}

func (u *branchUsecase) Find(ctx context.Context, filter domain.BranchFilter) ([]*domain.Branch, int64, error) {
	branches, total, err := u.branchRepo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find branches",
			zap.Error(err),
			zap.Any("filter", filter),
		)
		return nil, 0, fmt.Errorf("branch_usecase.Find: %w", err)
	}

	return branches, total, nil
}

func (u *branchUsecase) FindByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	branch, err := u.branchRepo.FindByID(ctx, id)
	if err != nil {
		config.Log.Error("failed to find branch by id",
			zap.Error(err),
			zap.String("branch_id", id.String()),
		)
		return nil, fmt.Errorf("branch_usecase.FindByID: %w", err)
	}

	return branch, nil
}

func (u *branchUsecase) Create(ctx context.Context, branch *domain.Branch) error {
	if err := u.branchRepo.Create(ctx, branch); err != nil {
		config.Log.Error("failed to create branch",
			zap.Error(err),
			zap.String("name", branch.Name),
		)
		return fmt.Errorf("branch_usecase.Create: %w", err)
	}

	return nil
}

func (u *branchUsecase) Update(ctx context.Context, branch *domain.Branch) error {
	if err := u.branchRepo.Update(ctx, branch); err != nil {
		config.Log.Error("failed to update branch",
			zap.Error(err),
			zap.String("branch_id", branch.ID.String()),
		)
		return fmt.Errorf("branch_usecase.Update: %w", err)
	}

	return nil
}

func (u *branchUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.branchRepo.Delete(ctx, id); err != nil {
		config.Log.Error("failed to delete branch",
			zap.Error(err),
			zap.String("branch_id", id.String()),
		)
		return fmt.Errorf("branch_usecase.Delete: %w", err)
	}

	return nil
}
