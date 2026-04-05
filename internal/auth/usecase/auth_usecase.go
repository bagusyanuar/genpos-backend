package usecase

import (
	"context"
	"errors"

	"github.com/bagusyanuar/genpos-backend/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo domain.AuthRepository
}

func NewAuthUsecase(repo domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{repo: repo}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Placeholder for JWT generation
	token := "placeholder-jwt-token"
	return token, nil
}
