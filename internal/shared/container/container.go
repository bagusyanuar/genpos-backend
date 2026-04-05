package container

import (
	"github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	"github.com/bagusyanuar/genpos-backend/internal/auth/repository"
	"github.com/bagusyanuar/genpos-backend/internal/auth/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/config"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *http.AuthHandler
}

func NewContainer(db *gorm.DB, conf *config.Config) *Container {
	// Auth Module Wiring
	authRepo := repository.NewAuthRepository(db)
	authUC := usecase.NewAuthUsecase(authRepo, conf)
	authHandler := http.NewAuthHandler(authUC, conf)

	return &Container{
		AuthHandler: authHandler,
	}
}
