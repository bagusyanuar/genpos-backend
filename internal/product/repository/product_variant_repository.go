package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/product/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) domain.ProductVariantRepository {
	return &productVariantRepository{db: db}
}

func (r *productVariantRepository) CreateBatch(ctx context.Context, variants []domain.ProductVariant) error {
	err := r.db.WithContext(ctx).Create(&variants).Error
	if err != nil {
		return fmt.Errorf("variant_repo.CreateBatch: %w", err)
	}
	return nil
}

func (r *productVariantRepository) UpdateBatch(ctx context.Context, variants []domain.ProductVariant) error {
	// GORM will use ID to update each variant
	err := r.db.WithContext(ctx).Save(&variants).Error
	if err != nil {
		return fmt.Errorf("variant_repo.UpdateBatch: %w", err)
	}
	return nil
}

func (r *productVariantRepository) DeleteByProductID(ctx context.Context, productID uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&domain.ProductVariant{}, "product_id = ?", productID).Error
	if err != nil {
		return fmt.Errorf("variant_repo.DeleteByProductID: %w", err)
	}
	return nil
}

func (r *productVariantRepository) FindByProductID(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error) {
	var variants []domain.ProductVariant
	err := r.db.WithContext(ctx).Find(&variants, "product_id = ?", productID).Error
	if err != nil {
		return nil, fmt.Errorf("variant_repo.FindByProductID: %w", err)
	}
	return variants, nil
}

func (r *productVariantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
	var variant domain.ProductVariant
	err := r.db.WithContext(ctx).First(&variant, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("variant_repo.FindByID: %w", err)
	}
	return &variant, nil
}
