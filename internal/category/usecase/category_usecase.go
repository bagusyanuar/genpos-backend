package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/category/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUsecase) Find(ctx context.Context, filter domain.CategoryFilter) ([]*domain.Category, int64, error) {
	categories, total, err := u.categoryRepo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find categories",
			zap.Error(err),
			zap.Any("filter", filter),
		)
		return nil, 0, fmt.Errorf("category_uc.Find: %w", err)
	}

	return categories, total, nil
}

func (u *categoryUsecase) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		config.Log.Error("failed to find category by id",
			zap.Error(err),
			zap.String("category_id", id.String()),
		)
		return nil, fmt.Errorf("category_uc.FindByID: %w", err)
	}

	return category, nil
}

func (u *categoryUsecase) Create(ctx context.Context, category *domain.Category) error {
	if category.ParentID != nil {
		parent, err := u.categoryRepo.FindByID(ctx, *category.ParentID)
		if err != nil {
			return fmt.Errorf("category_uc.Create (parent check): %w", err)
		}
		if parent.Level >= 2 {
			return fmt.Errorf("category_uc.Create: maximum depth reached (level 3)")
		}
		category.Level = parent.Level + 1
	} else {
		category.Level = 0
	}

	if err := u.categoryRepo.Create(ctx, category); err != nil {
		return fmt.Errorf("category_uc.Create: %w", err)
	}
	return nil
}

func (u *categoryUsecase) Update(ctx context.Context, category *domain.Category) error {
	oldCategory, err := u.categoryRepo.FindByID(ctx, category.ID)
	if err != nil {
		return fmt.Errorf("category_uc.Update (fetch old): %w", err)
	}

	parentChanged := false
	if (oldCategory.ParentID == nil && category.ParentID != nil) ||
		(oldCategory.ParentID != nil && category.ParentID == nil) ||
		(oldCategory.ParentID != nil && category.ParentID != nil && *oldCategory.ParentID != *category.ParentID) {
		parentChanged = true
	}

	if parentChanged {
		if err := u.validateMove(ctx, category.ID, category.ParentID); err != nil {
			return fmt.Errorf("category_uc.Update (validate move): %w", err)
		}

		newLevel := 0
		if category.ParentID != nil {
			parent, _ := u.categoryRepo.FindByID(ctx, *category.ParentID)
			newLevel = parent.Level + 1
		}
		category.Level = newLevel
	} else {
		category.Level = oldCategory.Level
	}

	tx := u.categoryRepo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.WithContext(ctx).Save(category).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("category_uc.Update (save): %w", err)
	}

	if parentChanged {
		levelDiff := category.Level - oldCategory.Level
		if levelDiff != 0 {
			// Batched SQL Update for children (Level 1)
			if err := tx.WithContext(ctx).Model(&domain.Category{}).
				Where("parent_id = ?", category.ID).
				Update("level", gorm.Expr("level + ?", levelDiff)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("category_uc.Update (batch child): %w", err)
			}

			// Batched SQL Update for grandchildren (Level 2)
			// Manual query to handle nested level shift
			if err := tx.WithContext(ctx).Exec(
				"UPDATE categories SET level = level + ? WHERE parent_id IN (SELECT id FROM categories WHERE parent_id = ?)",
				levelDiff, category.ID,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("category_uc.Update (batch grandchild): %w", err)
			}
		}
	}

	return tx.Commit().Error
}

func (u *categoryUsecase) validateMove(ctx context.Context, categoryID uuid.UUID, newParentID *uuid.UUID) error {
	if newParentID != nil {
		if *newParentID == categoryID {
			return fmt.Errorf("category cannot be its own parent")
		}

		parent, err := u.categoryRepo.FindByID(ctx, *newParentID)
		if err != nil {
			return fmt.Errorf("new parent not found")
		}

		// Use Recursive CTE to find current sub-tree depth
		var currentMaxLevel int
		u.categoryRepo.GetDB().WithContext(ctx).Raw(`
			WITH RECURSIVE sub_tree AS (
				SELECT id, level FROM categories WHERE id = ?
				UNION ALL
				SELECT c.id, c.level FROM categories c JOIN sub_tree st ON c.parent_id = st.id
			)
			SELECT MAX(level) FROM sub_tree
		`, categoryID).Scan(&currentMaxLevel)

		// Get current level of moving category to find relative height
		var currentLevel int
		u.categoryRepo.GetDB().WithContext(ctx).Model(&domain.Category{}).Select("level").Where("id = ?", categoryID).Scan(&currentLevel)

		subtreeHeight := currentMaxLevel - currentLevel
		if parent.Level+1+subtreeHeight > 2 {
			return fmt.Errorf("maximum depth reached (level 3)")
		}

		// Circular reference: ensure new parent is NOT a descendant
		var count int64
		u.categoryRepo.GetDB().WithContext(ctx).Raw(`
			WITH RECURSIVE sub_tree AS (
				SELECT id FROM categories WHERE id = ?
				UNION ALL
				SELECT c.id FROM categories c JOIN sub_tree st ON c.parent_id = st.id
			)
			SELECT COUNT(*) FROM sub_tree WHERE id = ?
		`, categoryID, *newParentID).Scan(&count)

		if count > 0 {
			return fmt.Errorf("circular reference: new parent is a descendant")
		}
	}
	return nil
}

func (u *categoryUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.categoryRepo.Delete(ctx, id); err != nil {
		config.Log.Error("failed to delete category", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("category_uc.Delete: %w", err)
	}
	return nil
}
