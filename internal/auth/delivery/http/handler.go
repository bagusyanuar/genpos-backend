package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/auth/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	uc   domain.AuthUsecase
	conf *config.Config
}

func NewAuthHandler(uc domain.AuthUsecase, conf *config.Config) *AuthHandler {
	return &AuthHandler{
		uc:   uc,
		conf: conf,
	}
}

func (h *AuthHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/auth")
	group.Post("/login", h.Login)
	group.Post("/refresh", h.RefreshToken)

	// Protected
	group.Get("/me", authMiddleware, h.Me)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	idStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid user id in token"))
	}

	user, err := h.uc.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(err.Error()))
	}

	res := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	return c.Status(fiber.StatusOK).JSON(response.Success(res, "user profile retrieved"))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	tokenPair, err := h.uc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(err.Error()))
	}

	// Set Refresh Token in Cookie (HttpOnly, Secure, SameSite)
	h.setRefreshTokenCookie(c, tokenPair.RefreshToken)

	res := LoginResponse{AccessToken: tokenPair.AccessToken}
	return c.Status(fiber.StatusOK).JSON(response.Success(res, "login success"))
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error("refresh token missing"))
	}

	tokenPair, err := h.uc.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		// Clear invalid cookie
		c.ClearCookie("refresh_token")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(err.Error()))
	}

	// Set New Refresh Token in Cookie (Rotation)
	h.setRefreshTokenCookie(c, tokenPair.RefreshToken)

	res := LoginResponse{AccessToken: tokenPair.AccessToken}
	return c.Status(fiber.StatusOK).JSON(response.Success(res, "refresh success"))
}

func (h *AuthHandler) setRefreshTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  time.Now().Add(time.Duration(h.conf.JWTRefreshExpiration) * time.Hour * 24),
		HTTPOnly: true,
		Secure:   h.conf.AppEnv == "production",
		SameSite: "Lax",
		Path:     "/",
	})
}
