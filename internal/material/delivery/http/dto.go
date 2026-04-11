package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/google/uuid"
)

type MaterialResponse struct {
	ID           uuid.UUID `json:"id"`
	CategoryID   *uuid.UUID `json:"category_id"`
	SKU          string    `json:"sku"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	MaterialType string    `json:"material_type"`
	ImageURL     *string   `json:"image_url"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateMaterialUOMRequest struct {
	UnitID     uuid.UUID `json:"unit_id" validate:"required"`
	Multiplier float64   `json:"multiplier" validate:"required,gt=0"`
	IsDefault  bool      `json:"is_default"`
}

type CreateMaterialRequest struct {
	CategoryID   *uuid.UUID                 `json:"category_id"`
	SKU          string                     `json:"sku" validate:"required"`
	Name         string                     `json:"name" validate:"required"`
	Description  *string                    `json:"description"`
	MaterialType string                     `json:"material_type" validate:"required"`
	ImageURL     *string                    `json:"image_url"`
	IsActive     bool                       `json:"is_active" validate:"required"`
	UOMs         []CreateMaterialUOMRequest `json:"uoms" validate:"required,min=1"`
}

func (r *CreateMaterialRequest) ToEntity() *domain.Material {
	return &domain.Material{
		CategoryID:   r.CategoryID,
		SKU:          r.SKU,
		Name:         r.Name,
		Description:  r.Description,
		MaterialType: r.MaterialType,
		ImageURL:     r.ImageURL,
		IsActive:     r.IsActive,
	}
}

func (r *CreateMaterialUOMRequest) ToEntity() domain.MaterialUOM {
	return domain.MaterialUOM{
		UnitID:     r.UnitID,
		Multiplier: r.Multiplier,
		IsDefault:  r.IsDefault,
	}
}

func ToMaterialResponse(m domain.Material) MaterialResponse {
	return MaterialResponse{
		ID:           m.ID,
		CategoryID:   m.CategoryID,
		SKU:          m.SKU,
		Name:         m.Name,
		Description:  m.Description,
		MaterialType: m.MaterialType,
		ImageURL:     m.ImageURL,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
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

type UpdateMaterialUOMRequest struct {
	ID         *uuid.UUID `json:"id"`
	UnitID     uuid.UUID  `json:"unit_id" validate:"required"`
	Multiplier float64    `json:"multiplier" validate:"required,gt=0"`
	IsDefault  bool       `json:"is_default"`
}

func (r *UpdateMaterialUOMRequest) ToEntity() domain.MaterialUOM {
	uom := domain.MaterialUOM{
		UnitID:     r.UnitID,
		Multiplier: r.Multiplier,
		IsDefault:  r.IsDefault,
	}
	if r.ID != nil {
		uom.ID = *r.ID
	}
	return uom
}

type SyncMaterialUOMRequest struct {
	UOMs []UpdateMaterialUOMRequest `json:"uoms" validate:"required,min=1"`
}
