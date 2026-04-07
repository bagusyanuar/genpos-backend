package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
)

type FindUnitQuery struct {
	Search string `query:"search"`
	request.PaginationParam
}

type UnitResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUnitRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateUnitRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}
