package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/recipe/domain"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/bagusyanuar/genpos-backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RecipeHandler struct {
	uc domain.RecipeUsecase
}

func NewRecipeHandler(uc domain.RecipeUsecase) *RecipeHandler {
	return &RecipeHandler{uc: uc}
}

func (h *RecipeHandler) Sync(c *fiber.Ctx) error {
	variantIDStr := c.Params("variant_id")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid variant id format"))
	}

	var req SyncRecipeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationError(err))
	}

	recipes := make([]domain.Recipe, 0)
	for _, r := range req.Recipes {
		recipes = append(recipes, r.ToEntity())
	}

	if err := h.uc.SyncRecipe(c.Context(), variantID, recipes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil, "recipe synced successfully"))
}

func (h *RecipeHandler) GetByVariantID(c *fiber.Ctx) error {
	variantIDStr := c.Params("variant_id")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid variant id format"))
	}

	recipes, err := h.uc.GetByVariantID(c.Context(), variantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(ToRecipeListResponse(recipes), "recipes found successfully"))
}

func (h *RecipeHandler) CalculateCOGS(c *fiber.Ctx) error {
	variantIDStr := c.Params("variant_id")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid variant id format"))
	}

	cogs, err := h.uc.CalculateEstimatedCOGS(c.Context(), variantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := COGSResponse{
		ProductVariantID: variantID,
		EstimatedCOGS:    cogs,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, "cogs calculated successfully"))
}

func (h *RecipeHandler) Register(router fiber.Router, auth fiber.Handler) {
	recipe := router.Group("/recipes", auth)
	recipe.Post("/variant/:variant_id", h.Sync)
	recipe.Get("/variant/:variant_id", h.GetByVariantID)
	recipe.Get("/variant/:variant_id/cogs", h.CalculateCOGS)
}
