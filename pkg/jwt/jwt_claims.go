package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}
