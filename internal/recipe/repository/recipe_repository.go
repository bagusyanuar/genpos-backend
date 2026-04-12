package repository

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/recipe/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) domain.RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) FindByVariantID(ctx context.Context, variantID uuid.UUID) ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	err := r.db.WithContext(ctx).Find(&recipes, "product_variant_id = ?", variantID).Error
	if err != nil {
		return nil, fmt.Errorf("recipe_repo.FindByVariantID: %w", err)
	}
	return recipes, nil
}

func (r *recipeRepository) CreateBatch(ctx context.Context, recipes []domain.Recipe) error {
	err := r.db.WithContext(ctx).Create(&recipes).Error
	if err != nil {
		return fmt.Errorf("recipe_repo.CreateBatch: %w", err)
	}
	return nil
}

func (r *recipeRepository) ReplaceByVariantID(ctx context.Context, variantID uuid.UUID, recipes []domain.Recipe) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hard delete existing recipes for this variant
		if err := tx.Unscoped().Delete(&domain.Recipe{}, "product_variant_id = ?", variantID).Error; err != nil {
			return fmt.Errorf("recipe_repo.ReplaceByVariantID.Delete: %w", err)
		}

		// Insert new recipes
		if len(recipes) > 0 {
			if err := tx.Create(&recipes).Error; err != nil {
				return fmt.Errorf("recipe_repo.ReplaceByVariantID.Create: %w", err)
			}
		}

		return nil
	})
}

func (r *recipeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&domain.Recipe{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("recipe_repo.Delete: %w", err)
	}
	return nil
}
