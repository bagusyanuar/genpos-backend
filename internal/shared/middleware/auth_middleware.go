package middleware

import (
	"strings"

	"github.com/bagusyanuar/genpos-backend/internal/config"
	"github.com/bagusyanuar/genpos-backend/pkg/jwt"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// JWTProtected is a middleware that protects routes with a JWT access token.
func JWTProtected(conf *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get Authorization Header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("missing authorization header"))
		}

		// 2. Exact Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid authorization format"))
		}

		tokenString := parts[1]

		// 3. Parse and Validate Token
		claims, err := jwt.ParseToken(tokenString, conf.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid or expired token"))
		}

		// 4. Set Claims to Context for Handlers
		c.Locals("user_id", claims.Subject)
		c.Locals("email", claims.Email)
		c.Locals("roles", claims.Roles)

		return c.Next()
	}
}
