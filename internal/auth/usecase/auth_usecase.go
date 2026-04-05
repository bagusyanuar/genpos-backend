package usecase

import (
	"context"
	"errors"
	"time"

	authDomain "github.com/bagusyanuar/genpos-backend/internal/auth/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	userDomain "github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"github.com/bagusyanuar/genpos-backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo userDomain.UserRepository
	conf     *config.Config
}

func NewAuthUsecase(userRepo userDomain.UserRepository, conf *config.Config) authDomain.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		conf:     conf,
	}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (authDomain.TokenPair, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return authDomain.TokenPair{}, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return authDomain.TokenPair{}, errors.New("invalid email or password")
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
		return authDomain.TokenPair{}, errors.New("failed to generate access token")
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
		return authDomain.TokenPair{}, errors.New("failed to generate refresh token")
	}

	return authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (authDomain.TokenPair, error) {
	// 1. Parse and Validate Refresh Token
	claims, err := jwt.ParseToken(refreshToken, u.conf.JWTRefreshSecret)
	if err != nil {
		return authDomain.TokenPair{}, errors.New("invalid refresh token")
	}

	// 2. Extract User ID and Fetch User
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return authDomain.TokenPair{}, errors.New("invalid token payload")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return authDomain.TokenPair{}, errors.New("user not found")
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
		return authDomain.TokenPair{}, errors.New("failed to generate access token")
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
		return authDomain.TokenPair{}, errors.New("failed to generate refresh token")
	}

	return authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
