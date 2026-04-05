package container

import (
	authHttp "github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	authUsecase "github.com/bagusyanuar/genpos-backend/internal/auth/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	userHttp "github.com/bagusyanuar/genpos-backend/internal/user/delivery/http"
	userRepository "github.com/bagusyanuar/genpos-backend/internal/user/repository"
	userUsecase "github.com/bagusyanuar/genpos-backend/internal/user/usecase"
	userDomain "github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *authHttp.AuthHandler
	UserUC      userDomain.UserUsecase
	UserHandler *userHttp.UserHandler
}

func NewContainer(db *gorm.DB, conf *config.Config) *Container {
	// User Module Wiring
	userRepo := userRepository.NewUserRepository(db)

	// Auth Module Wiring
	authUC := authUsecase.NewAuthUsecase(userRepo, conf)
	authHandler := authHttp.NewAuthHandler(authUC, conf)

	// User Module Usecase Wiring
	userUC := userUsecase.NewUserUsecase(userRepo, conf)
	userHandler := userHttp.NewUserHandler(userUC, conf)

	return &Container{
		AuthHandler: authHandler,
		UserUC:      userUC,
		UserHandler: userHandler,
	}
}
