package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/product/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type branchProductRepository struct {
	db *gorm.DB
}

func NewBranchProductRepository(db *gorm.DB) domain.BranchProductRepository {
	return &branchProductRepository{db: db}
}

func (r *branchProductRepository) Assign(ctx context.Context, branchID uuid.UUID, productIDs []uuid.UUID) error {
	var relations []domain.BranchProduct
	for _, pid := range productIDs {
		relations = append(relations, domain.BranchProduct{
			BranchID:  branchID,
			ProductID: pid,
			IsActive:  true,
		})
	}

	// upsert logic if needed, or just create
	err := r.db.WithContext(ctx).Save(&relations).Error
	if err != nil {
		return fmt.Errorf("branch_product_repo.Assign: %w", err)
	}
	return nil
}

func (r *branchProductRepository) Unassign(ctx context.Context, branchID uuid.UUID, productIDs []uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("branch_id = ? AND product_id IN ?", branchID, productIDs).
		Delete(&domain.BranchProduct{}).Error

	if err != nil {
		return fmt.Errorf("branch_product_repo.Unassign: %w", err)
	}
	return nil
}

func (r *branchProductRepository) FindByBranch(ctx context.Context, branchID uuid.UUID, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	// Logic is similar to ProductRepo.Find with filter.BranchID set
	// This can be delegated to ProductRepo or implemented here.
	return nil, 0, nil
}
