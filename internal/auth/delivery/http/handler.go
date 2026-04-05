package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/auth/domain"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	uc domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(router fiber.Router) {
	group := router.Group("/auth")
	group.Post("/login", h.Login)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	token, err := h.uc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(fiber.Map{"token": token}, "login success"))
}

