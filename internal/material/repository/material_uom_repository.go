package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *materialUOMRepository) CreateBatch(ctx context.Context, uoms []domain.MaterialUOM) error {
	err := r.db.WithContext(ctx).Create(&uoms).Error
	if err != nil {
		return fmt.Errorf("material_uom_repo.CreateBatch: %w", err)
	}
	return nil
}

func (r *materialUOMRepository) ReplaceUOMs(ctx context.Context, materialID uuid.UUID, uoms []domain.MaterialUOM) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Collect incoming IDs to keep (Pre-allocate capacity)
		incomingIDs := make([]uuid.UUID, 0, len(uoms))
		for _, u := range uoms {
			if u.ID != uuid.Nil {
				incomingIDs = append(incomingIDs, u.ID)
			}
		}

		// 2. Soft delete UOMs not in the incoming list
		query := tx.Where("material_id = ?", materialID)
		if len(incomingIDs) > 0 {
			query = query.Where("id NOT IN ?", incomingIDs)
		}
		if err := query.Delete(&domain.MaterialUOM{}).Error; err != nil {
			return fmt.Errorf("failed to delete removed UOMs: %w", err)
		}

		// 3. Prepare and Save (Upsert) incoming UOMs
		for i := range uoms {
			uoms[i].MaterialID = materialID
		}

		if len(uoms) > 0 {
			if err := tx.Save(&uoms).Error; err != nil {
				return fmt.Errorf("failed to save UOMs: %w", err)
			}
		}

		return nil
	})
}

func (r *materialUOMRepository) RecalibrateUOMs(ctx context.Context, tx *gorm.DB, materialID uuid.UUID, cf float64, targetUOMID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	// Bulk update multipliers and is_default
	err := tx.WithContext(ctx).Model(&domain.MaterialUOM{}).
		Where("material_id = ?", materialID).
		Updates(map[string]interface{}{
			"multiplier": gorm.Expr("multiplier / ?", cf),
			"is_default": gorm.Expr("id = ?", targetUOMID),
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("material_uom_repo.RecalibrateUOMs: %w", err)
	}

	return nil
}
