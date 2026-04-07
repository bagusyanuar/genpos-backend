package container

import (
	categoryHttp "github.com/bagusyanuar/genpos-backend/internal/category/delivery/http"
	categoryRepository "github.com/bagusyanuar/genpos-backend/internal/category/repository"
	categoryUsecase "github.com/bagusyanuar/genpos-backend/internal/category/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireCategoryModule(db *gorm.DB, conf *config.Config) {
	categoryRepo := categoryRepository.NewCategoryRepository(db)
	c.CategoryUC = categoryUsecase.NewCategoryUsecase(categoryRepo)
	c.CategoryHandler = categoryHttp.NewCategoryHandler(c.CategoryUC, conf)
}
