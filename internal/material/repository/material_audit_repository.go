package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type materialAuditRepository struct {
	db *gorm.DB
}

func NewMaterialAuditRepository(db *gorm.DB) domain.MaterialAuditRepository {
	return &materialAuditRepository{db: db}
}

func (r *materialAuditRepository) Create(ctx context.Context, audit *domain.MaterialAudit) error {
	if err := r.db.WithContext(ctx).Create(audit).Error; err != nil {
		return fmt.Errorf("material_audit_repo.Create: %w", err)
	}
	return nil
}

func (r *materialAuditRepository) FindByMaterialID(ctx context.Context, materialID uuid.UUID, filter request.PaginationParam) ([]domain.MaterialAudit, int64, error) {
	var audits []domain.MaterialAudit
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.MaterialAudit{}).Where("material_id = ?", materialID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("material_audit_repo.FindByMaterialID.Count: %w", err)
	}

	err := query.
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Order("created_at DESC").
		Find(&audits).Error

	if err != nil {
		return nil, 0, fmt.Errorf("material_audit_repo.FindByMaterialID.Data: %w", err)
	}

	return audits, total, nil
}
