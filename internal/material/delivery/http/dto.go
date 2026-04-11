package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/google/uuid"
)

type MaterialResponse struct {
	ID        uuid.UUID `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToMaterialResponse(m domain.Material) MaterialResponse {
	return MaterialResponse{
		ID:        m.ID,
		SKU:       m.SKU,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToMaterialListResponse(materials []domain.Material) []MaterialResponse {
	res := make([]MaterialResponse, 0)
	for _, m := range materials {
		res = append(res, ToMaterialResponse(m))
	}
	return res
}

type MaterialUOMResponse struct {
	ID         uuid.UUID `json:"id"`
	MaterialID uuid.UUID `json:"material_id"`
	UnitID     uuid.UUID `json:"unit_id"`
	Multiplier float64   `json:"multiplier"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToMaterialUOMResponse(m domain.MaterialUOM) MaterialUOMResponse {
	return MaterialUOMResponse{
		ID:         m.ID,
		MaterialID: m.MaterialID,
		UnitID:     m.UnitID,
		Multiplier: m.Multiplier,
		IsDefault:  m.IsDefault,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func ToMaterialUOMListResponse(uoms []domain.MaterialUOM) []MaterialUOMResponse {
	res := make([]MaterialUOMResponse, 0)
	for _, uom := range uoms {
		res = append(res, ToMaterialUOMResponse(uom))
	}
	return res
}
