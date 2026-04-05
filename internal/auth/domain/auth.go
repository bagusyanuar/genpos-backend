package domain

import (
	"context"

	userDomain "github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"github.com/google/uuid"
)

// TokenPair holds the access and refresh token pair.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// AuthUsecase defines the interface for authentication business logic.
type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (TokenPair, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*userDomain.User, error)
}
