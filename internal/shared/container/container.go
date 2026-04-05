package container

import (
	"github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	"github.com/bagusyanuar/genpos-backend/internal/auth/repository"
	"github.com/bagusyanuar/genpos-backend/internal/auth/usecase"
	"gorm.io/gorm"
)

type Container struct {
	AuthHandler *http.AuthHandler
}

func NewContainer(db *gorm.DB) *Container {
	// Auth Module Wiring
	authRepo := repository.NewAuthRepository(db)
	authUC := usecase.NewAuthUsecase(authRepo)
	authHandler := http.NewAuthHandler(authUC)

	return &Container{
		AuthHandler: authHandler,
	}
}
