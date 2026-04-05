package container

import (
	authHttp "github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	authUsecase "github.com/bagusyanuar/genpos-backend/internal/auth/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/config"
	userRepository "github.com/bagusyanuar/genpos-backend/internal/user/repository"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *authHttp.AuthHandler
}

func NewContainer(db *gorm.DB, conf *config.Config) *Container {
	// User Module Wiring
	userRepo := userRepository.NewUserRepository(db)

	// Auth Module Wiring
	authUC := authUsecase.NewAuthUsecase(userRepo, conf)
	authHandler := authHttp.NewAuthHandler(authUC, conf)

	return &Container{
		AuthHandler: authHandler,
	}
}
