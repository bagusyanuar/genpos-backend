package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/category/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Find(ctx context.Context, filter domain.CategoryFilter) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.Category{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		db = db.Where("name ILIKE ?", search)
	}

	if filter.ParentID != nil {
		db = db.Where("parent_id = ?", filter.ParentID)
	}

	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
	}

	if filter.IsActive != nil {
		db = db.Where("is_active = ?", *filter.IsActive)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(filter.GetLimit()).Offset(filter.GetOffset()).Order(filter.GetSort()).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("category_repo.Create: %w", err)
	}
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return fmt.Errorf("category_repo.Update: %w", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Cascading soft-delete using a single query for efficiency
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// PostgreSQL specific recursive DELETE for soft-delete (using GORM's deleted_at)
		return tx.Exec(`
			UPDATE categories 
			SET deleted_at = NOW() 
			WHERE id IN (
				WITH RECURSIVE sub_tree AS (
					SELECT id FROM categories WHERE id = ?
					UNION ALL
					SELECT c.id FROM categories c JOIN sub_tree st ON c.parent_id = st.id
				)
				SELECT id FROM sub_tree
			) AND deleted_at IS NULL
		`, id).Error
	})

	if err != nil {
		return fmt.Errorf("category_repo.Delete: %w", err)
	}
	return nil
}

func (r *categoryRepository) GetDB() *gorm.DB {
	return r.db
}
