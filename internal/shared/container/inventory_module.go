package container

import (
	inventoryHttp "github.com/bagusyanuar/genpos-backend/internal/inventory/delivery/http"
	inventoryRepo "github.com/bagusyanuar/genpos-backend/internal/inventory/repository"
	inventoryUC "github.com/bagusyanuar/genpos-backend/internal/inventory/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireInventoryModule(db *gorm.DB, conf *config.Config) {
	repo := inventoryRepo.NewInventoryRepository(db)
	uc := inventoryUC.NewInventoryUsecase(repo)
	handler := inventoryHttp.NewInventoryHandler(uc, conf)

	c.InventoryUC = uc
	c.InventoryHandler = handler
}
