package repository

import (
	"context"
	"fmt"
	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type materialUOMRepository struct {
	db *gorm.DB
}

func NewMaterialUOMRepository(db *gorm.DB) domain.MaterialUOMRepository {
	return &materialUOMRepository{db: db}
}

func (r *materialUOMRepository) Find(ctx context.Context, materialID uuid.UUID) ([]domain.MaterialUOM, error) {
	var uoms []domain.MaterialUOM

	err := r.db.WithContext(ctx).
		Where("material_id = ?", materialID).
		Find(&uoms).Error

	if err != nil {
		return nil, fmt.Errorf("material_uom_repo.Find: %w", err)
	}

	return uoms, nil
}
