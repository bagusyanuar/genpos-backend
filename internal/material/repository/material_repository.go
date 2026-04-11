package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type materialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) domain.MaterialRepository {
	return &materialRepository{db: db}
}

func (r *materialRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Material, error) {
	var material domain.Material
	err := r.db.WithContext(ctx).First(&material, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("material_repo.FindByID: %w", err)
	}
	return &material, nil
}

func (r *materialRepository) Find(ctx context.Context, filter domain.MaterialFilter) ([]domain.Material, int64, error) {
	var materials []domain.Material
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Material{})

	// Filter Search (SKU or Name)
	if filter.Search != "" {
		searchText := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(sku) LIKE ? OR LOWER(name) LIKE ?", searchText, searchText)
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("material_repo.Find.Count: %w", err)
	}

	// Pagination & Sorting
	err := query.
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Order(filter.GetSort()).
		Find(&materials).Error

	if err != nil {
		return nil, 0, fmt.Errorf("material_repo.Find.Find: %w", err)
	}

	return materials, total, nil
}
