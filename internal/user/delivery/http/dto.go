package http

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type FindUserQuery struct {
	Search string `json:"search" query:"search"`
	request.PaginationParam
}
