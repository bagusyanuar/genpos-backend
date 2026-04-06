package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
)

type BranchResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FindBranchQuery struct {
	Search string `json:"search" query:"search"`
	request.PaginationParam
}

type CreateBranchRequest struct {
	Name      string `json:"name" validate:"required,min=3"`
	IsDefault bool   `json:"is_default"`
}

type UpdateBranchRequest struct {
	Name      string `json:"name" validate:"required,min=3"`
	IsDefault bool   `json:"is_default"`
}
