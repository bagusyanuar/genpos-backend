package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bagusyanuar/genpos-backend/internal/product/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	err := r.db.WithContext(ctx).Create(product).Error
	if err != nil {
		return fmt.Errorf("product_repo.Create: %w", err)
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	err := r.db.WithContext(ctx).Model(product).
		Select("category_id", "name", "description", "image_url", "is_active", "updated_at").
		Updates(product).Error

	if err != nil {
		return fmt.Errorf("product_repo.Update: %w", err)
	}
	return nil
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Preload("Variants").First(&product, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("product_repo.FindByID: %w", err)
	}
	return &product, nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("product_repo.Delete: %w", err)
	}
	return nil
}

func (r *productRepository) Find(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{})

	// Filter Search (Name)
	if filter.Search != "" {
		searchText := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchText)
	}

	// Filter Category
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	// Filter Branch (Many-to-Many join)
	if filter.BranchID != nil {
		query = query.Joins("JOIN branch_products ON branch_products.product_id = products.id").
			Where("branch_products.branch_id = ?", *filter.BranchID).
			Where("branch_products.is_active = ?", true)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("product_repo.Find.Count: %w", err)
	}

	// Pagination
	err := query.
		Preload("Variants").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Order(filter.GetSort()).
		Find(&products).Error

	if err != nil {
		return nil, 0, fmt.Errorf("product_repo.Find.Find: %w", err)
	}

	return products, total, nil
}
