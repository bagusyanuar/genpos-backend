package delivery

import (
	"github.com/bagusyanuar/genpos-backend/internal/auth/domain"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	token, err := h.uc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

