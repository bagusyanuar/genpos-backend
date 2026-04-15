package container

import (
	materialRepo "github.com/bagusyanuar/genpos-backend/internal/material/repository"
	productRepo "github.com/bagusyanuar/genpos-backend/internal/product/repository"
	recipeHttp "github.com/bagusyanuar/genpos-backend/internal/recipe/delivery/http"
	recipeRepo "github.com/bagusyanuar/genpos-backend/internal/recipe/repository"
	recipeUC "github.com/bagusyanuar/genpos-backend/internal/recipe/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireRecipeModule(db *gorm.DB, _ *config.Config) {
	repo := recipeRepo.NewRecipeRepository(db)
	matRepo := materialRepo.NewMaterialRepository(db)
	variantRepo := productRepo.NewProductVariantRepository(db)

	uc := recipeUC.NewRecipeUsecase(repo, matRepo, variantRepo)
	handler := recipeHttp.NewRecipeHandler(uc)

	c.RecipeUC = uc
	c.RecipeHandler = handler
}
