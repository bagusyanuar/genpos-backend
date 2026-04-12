package container

import (
	productHttp "github.com/bagusyanuar/genpos-backend/internal/product/delivery/http"
	"github.com/bagusyanuar/genpos-backend/internal/product/repository"
	"github.com/bagusyanuar/genpos-backend/internal/product/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireProductModule(db *gorm.DB, conf *config.Config) {
	productRepo := repository.NewProductRepository(db)
	variantRepo := repository.NewProductVariantRepository(db)
	branchRepo := repository.NewBranchProductRepository(db)

	productUC := usecase.NewProductUsecase(productRepo, variantRepo, branchRepo, c.Uploader)
	productHandler := productHttp.NewProductHandler(productUC, c.Uploader, conf)

	c.ProductUC = productUC
	c.ProductHandler = productHandler
}
