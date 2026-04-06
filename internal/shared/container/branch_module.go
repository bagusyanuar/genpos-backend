package container

import (
	branchHttp "github.com/bagusyanuar/genpos-backend/internal/branch/delivery/http"
	branchRepository "github.com/bagusyanuar/genpos-backend/internal/branch/repository"
	branchUsecase "github.com/bagusyanuar/genpos-backend/internal/branch/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireBranchModule(db *gorm.DB, conf *config.Config) {
	branchRepo := branchRepository.NewBranchRepository(db)
	c.BranchUC = branchUsecase.NewBranchUsecase(branchRepo)
	c.BranchHandler = branchHttp.NewBranchHandler(c.BranchUC, conf)
}
