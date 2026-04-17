package container

import (
	inventoryRepo "github.com/bagusyanuar/genpos-backend/internal/inventory/repository"
	materialHttp "github.com/bagusyanuar/genpos-backend/internal/material/delivery/http"
	materialRepo "github.com/bagusyanuar/genpos-backend/internal/material/repository"
	materialUC "github.com/bagusyanuar/genpos-backend/internal/material/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireMaterialModule(db *gorm.DB, conf *config.Config) {
	repo := materialRepo.NewMaterialRepository(db)
	uomRepo := materialRepo.NewMaterialUOMRepository(db)
	invRepo := inventoryRepo.NewInventoryRepository(db)
	
	uc := materialUC.NewMaterialUsecase(repo, uomRepo, invRepo, c.Uploader)
	uomUC := materialUC.NewMaterialUOMUsecase(uomRepo)

	handler := materialHttp.NewMaterialHandler(uc, uomUC, conf)

	c.MaterialUC = uc
	c.MaterialHandler = handler
}
