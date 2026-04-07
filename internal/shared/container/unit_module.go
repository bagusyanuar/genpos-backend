package container

import (
	unitHttp "github.com/bagusyanuar/genpos-backend/internal/unit/delivery/http"
	unitRepository "github.com/bagusyanuar/genpos-backend/internal/unit/repository"
	unitUsecase "github.com/bagusyanuar/genpos-backend/internal/unit/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"gorm.io/gorm"
)

func (c *Container) wireUnitModule(db *gorm.DB, conf *config.Config) {
	unitRepo := unitRepository.NewUnitRepository(db)
	c.UnitUC = unitUsecase.NewUnitUsecase(unitRepo)
	c.UnitHandler = unitHttp.NewUnitHandler(c.UnitUC, conf)
}
