package container

import (
	authHttp "github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	authUsecase "github.com/bagusyanuar/genpos-backend/internal/auth/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	userDomain "github.com/bagusyanuar/genpos-backend/internal/user/domain"
)

func (c *Container) wireAuthModule(userRepo userDomain.UserRepository, conf *config.Config) {
	authUC := authUsecase.NewAuthUsecase(userRepo, conf)
	c.AuthHandler = authHttp.NewAuthHandler(authUC, conf)
}
