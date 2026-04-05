package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/auth/domain"
	"github.com/bagusyanuar/genpos-backend/internal/config"
	"github.com/bagusyanuar/genpos-backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo domain.AuthRepository
	conf *config.Config
}

func NewAuthUsecase(repo domain.AuthRepository, conf *config.Config) domain.AuthUsecase {
	return &authUsecase{
		repo: repo,
		conf: conf,
	}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (domain.TokenPair, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.TokenPair{}, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domain.TokenPair{}, errors.New("invalid email or password")
	}

	// TODO: roles should be fetched from DB
	roles := []string{"admin"}

	// 1. Generate JWT Access Token
	accessToken, err := jwt.GenerateToken(
		user.ID.String(),
		user.Email,
		roles,
		u.conf.JWTSecret,
		u.conf.JWTIssuer,
		time.Duration(u.conf.JWTExpiration)*time.Minute,
	)
	if err != nil {
		return domain.TokenPair{}, errors.New("failed to generate access token")
	}

	// 2. Generate JWT Refresh Token (longer lived, potentially different roles/scope)
	refreshToken, err := jwt.GenerateToken(
		user.ID.String(),
		user.Email,
		roles,
		u.conf.JWTRefreshSecret,
		u.conf.JWTIssuer,
		time.Duration(u.conf.JWTRefreshExpiration)*time.Hour*24,
	)
	if err != nil {
		return domain.TokenPair{}, errors.New("failed to generate refresh token")
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	// 1. Parse and Validate Refresh Token
	claims, err := jwt.ParseToken(refreshToken, u.conf.JWTRefreshSecret)
	if err != nil {
		return domain.TokenPair{}, errors.New("invalid refresh token")
	}

	// 2. Extract User ID and Fetch User
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return domain.TokenPair{}, errors.New("invalid token payload")
	}

	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return domain.TokenPair{}, errors.New("user not found")
	}

	// 3. Generate New Token Pair (Rotation)
	// TODO: roles from DB
	roles := []string{"admin"}

	accessToken, err := jwt.GenerateToken(
		user.ID.String(),
		user.Email,
		roles,
		u.conf.JWTSecret,
		u.conf.JWTIssuer,
		time.Duration(u.conf.JWTExpiration)*time.Minute,
	)
	if err != nil {
		return domain.TokenPair{}, errors.New("failed to generate access token")
	}

	newRefreshToken, err := jwt.GenerateToken(
		user.ID.String(),
		user.Email,
		roles,
		u.conf.JWTRefreshSecret,
		u.conf.JWTIssuer,
		time.Duration(u.conf.JWTRefreshExpiration)*time.Hour*24,
	)
	if err != nil {
		return domain.TokenPair{}, errors.New("failed to generate refresh token")
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
