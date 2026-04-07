package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
)

type FindCategoryQuery struct {
	Search   string `query:"search"`
	ParentID string `query:"parent_id"`
	Type     string `query:"type"`
	IsActive *bool  `query:"is_active"`
	request.PaginationParam
}

type CategoryResponse struct {
	ID          string    `json:"id"`
	ParentID    *string   `json:"parent_id,omitempty"`
	Level       int       `json:"level"`
	Name        string    `json:"name"`
	Slug        *string   `json:"slug,omitempty"`
	Description *string   `json:"description,omitempty"`
	Type        string    `json:"type"`
	ImageURL    *string   `json:"image_url,omitempty"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	ParentID    *string `json:"parent_id" validate:"omitempty,uuid"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Slug        *string `json:"slug" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Type        string  `json:"type" validate:"required,oneof=PRODUCT INGREDIENT ALL"`
	ImageURL    *string `json:"image_url" validate:"omitempty,url"`
	SortOrder   int     `json:"sort_order" validate:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active" validate:"required"`
}

type UpdateCategoryRequest struct {
	ParentID    *string `json:"parent_id" validate:"omitempty,uuid"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Slug        *string `json:"slug" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Type        string  `json:"type" validate:"required,oneof=PRODUCT INGREDIENT ALL"`
	ImageURL    *string `json:"image_url" validate:"omitempty,url"`
	SortOrder   int     `json:"sort_order" validate:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active" validate:"required"`
}
