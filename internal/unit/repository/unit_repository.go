package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/unit/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) domain.UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) Find(ctx context.Context, filter domain.UnitFilter) ([]*domain.Unit, int64, error) {
	var units []*domain.Unit
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.Unit{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		db = db.Where("name ILIKE ?", search)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(filter.GetLimit()).Offset(filter.GetOffset()).Order(filter.GetSort()).Find(&units).Error; err != nil {
		return nil, 0, err
	}

	return units, total, nil
}

func (r *unitRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error) {
	var unit domain.Unit
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&unit).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) Create(ctx context.Context, unit *domain.Unit) error {
	if err := r.db.WithContext(ctx).Create(unit).Error; err != nil {
		return err
	}
	return nil
}

func (r *unitRepository) Update(ctx context.Context, unit *domain.Unit) error {
	if err := r.db.WithContext(ctx).Save(unit).Error; err != nil {
		return err
	}
	return nil
}

func (r *unitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Unit{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
