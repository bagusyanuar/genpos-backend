package container

import (
	authHttp "github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	branchHttp "github.com/bagusyanuar/genpos-backend/internal/branch/delivery/http"
	branchDomain "github.com/bagusyanuar/genpos-backend/internal/branch/domain"
	categoryHttp "github.com/bagusyanuar/genpos-backend/internal/category/delivery/http"
	categoryDomain "github.com/bagusyanuar/genpos-backend/internal/category/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	unitHttp "github.com/bagusyanuar/genpos-backend/internal/unit/delivery/http"
	unitDomain "github.com/bagusyanuar/genpos-backend/internal/unit/domain"
	userHttp "github.com/bagusyanuar/genpos-backend/internal/user/delivery/http"
	userDomain "github.com/bagusyanuar/genpos-backend/internal/user/domain"
	materialHttp "github.com/bagusyanuar/genpos-backend/internal/material/delivery/http"
	materialDomain "github.com/bagusyanuar/genpos-backend/internal/material/domain"
	inventoryHttp "github.com/bagusyanuar/genpos-backend/internal/inventory/delivery/http"
	inventoryDomain "github.com/bagusyanuar/genpos-backend/internal/inventory/domain"
	mediaHttp "github.com/bagusyanuar/genpos-backend/internal/media/delivery/http"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler   *authHttp.AuthHandler
	UserUC        userDomain.UserUsecase
	UserHandler   *userHttp.UserHandler
	BranchUC      branchDomain.BranchUsecase
	BranchHandler *branchHttp.BranchHandler
	UnitUC        unitDomain.UnitUsecase
	UnitHandler   *unitHttp.UnitHandler
	CategoryUC      categoryDomain.CategoryUsecase
	CategoryHandler *categoryHttp.CategoryHandler
	MaterialUC      materialDomain.MaterialUsecase
	MaterialHandler *materialHttp.MaterialHandler
	InventoryUC      inventoryDomain.InventoryUsecase
	InventoryHandler *inventoryHttp.InventoryHandler
	MediaHandler     *mediaHttp.MediaHandler
}

func NewContainer(db *gorm.DB, conf *config.Config) *Container {
	c := &Container{}

	// Wiring modules (delegated to modular files)
	userRepo := c.wireUserModule(db, conf)
	c.wireAuthModule(userRepo, conf)
	c.wireBranchModule(db, conf)
	c.wireUnitModule(db, conf)
	c.wireCategoryModule(db, conf)
	c.wireMaterialModule(db, conf)
	c.wireInventoryModule(db, conf)
	c.wireMediaModule(db, conf)

	return c
}
