package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryType represents the discriminator for product or ingredient categories.
type CategoryType string

const (
	CategoryTypeProduct    CategoryType = "PRODUCT"
	CategoryTypeIngredient CategoryType = "INGREDIENT"
	CategoryTypeAll        CategoryType = "ALL"
)

// Category represents the hierarchical category entity.
type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ParentID    *uuid.UUID     `gorm:"type:uuid" json:"parent_id"`
	Level       int            `gorm:"default:0;not null" json:"level"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Slug        *string        `gorm:"type:varchar(100);unique" json:"slug"`
	Description *string        `gorm:"type:text" json:"description"`
	Type        CategoryType   `gorm:"type:varchar(20);default:'PRODUCT';not null" json:"type"`
	ImageURL    *string        `gorm:"type:text" json:"image_url"`
	SortOrder   int            `gorm:"default:0;not null" json:"sort_order"`
	IsActive    bool           `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate hooks into GORM to auto-generate UUID IDs.
func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

// CategoryFilter represents criteria for category data operations.
type CategoryFilter struct {
	Search   string
	ParentID *uuid.UUID
	Type     CategoryType
	IsActive *bool
	request.PaginationParam
}

// CategoryRepository defines the interface for category data operations.
type CategoryRepository interface {
	Find(ctx context.Context, filter CategoryFilter) ([]*Category, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDB() *gorm.DB
}

// CategoryUsecase defines the interface for category business logic.
type CategoryUsecase interface {
	Find(ctx context.Context, filter CategoryFilter) ([]*Category, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}
