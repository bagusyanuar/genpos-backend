package usecase

import (
	"context"
	"fmt"

	matDomain "github.com/bagusyanuar/genpos-backend/internal/material/domain"
	prodDomain "github.com/bagusyanuar/genpos-backend/internal/product/domain"
	"github.com/bagusyanuar/genpos-backend/internal/recipe/domain"
	"github.com/google/uuid"
)

type recipeUsecase struct {
	recipeRepo         domain.RecipeRepository
	materialRepo       matDomain.MaterialRepository
	productVariantRepo prodDomain.ProductVariantRepository
}

func NewRecipeUsecase(
	recipeRepo domain.RecipeRepository,
	materialRepo matDomain.MaterialRepository,
	productVariantRepo prodDomain.ProductVariantRepository,
) domain.RecipeUsecase {
	return &recipeUsecase{
		recipeRepo:         recipeRepo,
		materialRepo:       materialRepo,
		productVariantRepo: productVariantRepo,
	}
}

func (u *recipeUsecase) GetByVariantID(ctx context.Context, variantID uuid.UUID) ([]domain.Recipe, error) {
	return u.recipeRepo.FindByVariantID(ctx, variantID)
}

func (u *recipeUsecase) SyncRecipe(ctx context.Context, variantID uuid.UUID, recipes []domain.Recipe) error {
	// Set the variant ID for each recipe to ensure consistency
	for i := range recipes {
		recipes[i].ProductVariantID = variantID
	}

	return u.recipeRepo.ReplaceByVariantID(ctx, variantID, recipes)
}

func (u *recipeUsecase) CalculateEstimatedCOGS(ctx context.Context, variantID uuid.UUID) (float64, error) {
	recipes, err := u.recipeRepo.FindByVariantID(ctx, variantID)
	if err != nil {
		return 0, fmt.Errorf("recipe_uc.CalculateCOGS.FindRecipes: %w", err)
	}

	if len(recipes) == 0 {
		// Even without recipes, a product might have overhead costs
		return u.getOverheadCost(ctx, variantID)
	}

	// Fetch Product Variant for OverheadCost
	overhead, err := u.getOverheadCost(ctx, variantID)
	if err != nil {
		return 0, err
	}

	var totalCOGS float64 = overhead

	for _, r := range recipes {
		itemCost := r.SubtotalCost

		// Middle Ground Logic: If SubtotalCost is 0, calculate dynamically
		if itemCost <= 0 {
			material, err := u.materialRepo.FindByID(ctx, r.MaterialID)
			if err != nil {
				return 0, fmt.Errorf("recipe_uc.CalculateCOGS.GetMaterial %s: %w", r.MaterialID, err)
			}
			
			// COGS = Quantity * (BaseCost * Multiplier)
			// BaseCost is cost per base unit. Multiplier is how many base units in this UOM.
			itemCost = r.Quantity * (material.BaseCost * r.MaterialUOM.Multiplier)
		}

		totalCOGS += itemCost
	}

	return totalCOGS, nil
}

func (u *recipeUsecase) getOverheadCost(ctx context.Context, variantID uuid.UUID) (float64, error) {
	variant, err := u.productVariantRepo.FindByID(ctx, variantID)
	if err != nil {
		return 0, fmt.Errorf("recipe_uc.getOverheadCost: %w", err)
	}

	return variant.OverheadCost, nil
}
