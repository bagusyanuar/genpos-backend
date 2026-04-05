package domain

import (
	"context"
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
}
