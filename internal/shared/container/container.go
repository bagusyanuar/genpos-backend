package container

import (
	authHttp "github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	branchHttp "github.com/bagusyanuar/genpos-backend/internal/branch/delivery/http"
	branchDomain "github.com/bagusyanuar/genpos-backend/internal/branch/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	userHttp "github.com/bagusyanuar/genpos-backend/internal/user/delivery/http"
	userDomain "github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler   *authHttp.AuthHandler
	UserUC        userDomain.UserUsecase
	UserHandler   *userHttp.UserHandler
	BranchUC      branchDomain.BranchUsecase
	BranchHandler *branchHttp.BranchHandler
}

func NewContainer(db *gorm.DB, conf *config.Config) *Container {
	c := &Container{}

	// Wiring modules (delegated to modular files)
	userRepo := c.wireUserModule(db, conf)
	c.wireAuthModule(userRepo, conf)
	c.wireBranchModule(db, conf)

	return c
}
