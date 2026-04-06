package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/branch/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) domain.BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	var branch domain.Branch
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&branch).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) Find(ctx context.Context, filter domain.BranchFilter) ([]*domain.Branch, int64, error) {
	var branches []*domain.Branch
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.Branch{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		db = db.Where("name LIKE ?", search)
	}

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := db.Limit(filter.GetLimit()).Offset(filter.GetOffset()).Order(filter.GetSort()).Find(&branches).Error; err != nil {
		return nil, 0, err
	}

	return branches, total, nil
}

func (r *branchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	if err := r.db.WithContext(ctx).Create(branch).Error; err != nil {
		return err
	}
	return nil
}

func (r *branchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	if err := r.db.WithContext(ctx).Save(branch).Error; err != nil {
		return err
	}
	return nil
}

func (r *branchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Branch{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
