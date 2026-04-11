package container

import (
	materialHttp "github.com/bagusyanuar/genpos-backend/internal/material/delivery/http"
	materialRepo "github.com/bagusyanuar/genpos-backend/internal/material/repository"
	materialUC "github.com/bagusyanuar/genpos-backend/internal/material/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireMaterialModule(db *gorm.DB, conf *config.Config) {
	repo := materialRepo.NewMaterialRepository(db)
	uomRepo := materialRepo.NewMaterialUOMRepository(db)
	
	uc := materialUC.NewMaterialUsecase(repo, uomRepo)
	uomUC := materialUC.NewMaterialUOMUsecase(uomRepo)

	handler := materialHttp.NewMaterialHandler(uc, uomUC, conf)

	c.MaterialUC = uc
	c.MaterialHandler = handler
}
