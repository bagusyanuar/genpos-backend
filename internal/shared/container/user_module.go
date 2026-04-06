package container

import (
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	userHttp "github.com/bagusyanuar/genpos-backend/internal/user/delivery/http"
	userDomain "github.com/bagusyanuar/genpos-backend/internal/user/domain"
	userRepository "github.com/bagusyanuar/genpos-backend/internal/user/repository"
	userUsecase "github.com/bagusyanuar/genpos-backend/internal/user/usecase"
	"gorm.io/gorm"
)

func (c *Container) wireUserModule(db *gorm.DB, conf *config.Config) userDomain.UserRepository {
	userRepo := userRepository.NewUserRepository(db)
	c.UserUC = userUsecase.NewUserUsecase(userRepo, conf)
	c.UserHandler = userHttp.NewUserHandler(c.UserUC, conf)
	return userRepo
}
