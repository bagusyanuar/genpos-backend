package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/recipe/domain"
	"github.com/google/uuid"
)

type RecipeItemRequest struct {
	MaterialID   uuid.UUID `json:"material_id" validate:"required"`
	UomID        uuid.UUID `json:"uom_id" validate:"required"`
	Quantity     float64   `json:"quantity" validate:"required,gt=0"`
	SubtotalCost float64   `json:"subtotal_cost"`
}

type SyncRecipeRequest struct {
	Recipes []RecipeItemRequest `json:"recipes" validate:"required,dive"`
}

func (r RecipeItemRequest) ToEntity() domain.Recipe {
	return domain.Recipe{
		MaterialID:   r.MaterialID,
		UomID:        r.UomID,
		Quantity:     r.Quantity,
		SubtotalCost: r.SubtotalCost,
	}
}

type RecipeResponse struct {
	ID               uuid.UUID `json:"id"`
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	MaterialID       uuid.UUID `json:"material_id"`
	UomID            uuid.UUID `json:"uom_id"`
	Quantity         float64   `json:"quantity"`
	SubtotalCost     float64   `json:"subtotal_cost"`
}

func ToRecipeResponse(r domain.Recipe) RecipeResponse {
	return RecipeResponse{
		ID:               r.ID,
		ProductVariantID: r.ProductVariantID,
		MaterialID:       r.MaterialID,
		UomID:            r.UomID,
		Quantity:         r.Quantity,
		SubtotalCost:     r.SubtotalCost,
	}
}

func ToRecipeListResponse(recipes []domain.Recipe) []RecipeResponse {
	res := make([]RecipeResponse, 0)
	for _, r := range recipes {
		res = append(res, ToRecipeResponse(r))
	}
	return res
}

type COGSResponse struct {
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	EstimatedCOGS    float64   `json:"estimated_cogs"`
}
