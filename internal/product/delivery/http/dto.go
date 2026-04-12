package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/product/domain"
	"github.com/google/uuid"
)

type CreateProductVariantRequest struct {
	ID       *uuid.UUID `json:"id"`
	Name     string     `json:"name" validate:"required"`
	SKU      string     `json:"sku" validate:"required"`
	Price    float64    `json:"price" validate:"required,gte=0"`
	IsActive bool       `json:"is_active" validate:"required"`
}

type CreateProductRequest struct {
	CategoryID  uuid.UUID                     `json:"category_id" validate:"required"`
	Name        string                        `json:"name" validate:"required"`
	Description *string                       `json:"description"`
	IsActive    bool                          `json:"is_active" validate:"required"`
	Variants    []CreateProductVariantRequest `json:"variants" validate:"required,min=1"`
	BranchIDs   []uuid.UUID                   `json:"branch_ids"`
}

func (r *CreateProductRequest) ToEntity() *domain.Product {
	return &domain.Product{
		CategoryID:  r.CategoryID,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

func (r *CreateProductVariantRequest) ToEntity() domain.ProductVariant {
	variant := domain.ProductVariant{
		Name:     r.Name,
		SKU:      r.SKU,
		Price:    r.Price,
		IsActive: r.IsActive,
	}

	if r.ID != nil {
		variant.ID = *r.ID
	}

	return variant
}

type ProductVariantResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	Price     float64   `json:"price"`
	IsActive  bool      `json:"is_active"`
}

type ProductResponse struct {
	ID          uuid.UUID                `json:"id"`
	CategoryID  uuid.UUID                `json:"category_id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"description"`
	ImageURL    *string                  `json:"image_url"`
	IsActive    bool                     `json:"is_active"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

func ToProductResponse(p domain.Product) ProductResponse {
	variants := make([]ProductVariantResponse, 0)
	for _, v := range p.Variants {
		variants = append(variants, ProductVariantResponse{
			ID:       v.ID,
			Name:     v.Name,
			SKU:      v.SKU,
			Price:    v.Price,
			IsActive: v.IsActive,
		})
	}

	return ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
		Variants:    variants,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToProductListResponse(products []domain.Product) []ProductResponse {
	res := make([]ProductResponse, 0)
	for _, p := range products {
		res = append(res, ToProductResponse(p))
	}
	return res
}
