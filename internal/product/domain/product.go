package domain

import (
	"context"
	"time"

	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description *string        `gorm:"type:text" json:"description"`
	ImageURL    *string        `gorm:"type:text" json:"image_url"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type ProductVariant struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	SKU          string         `gorm:"type:varchar(100);not null;index" json:"sku"`
	Price        float64        `gorm:"type:decimal(15,2);not null;default:0" json:"price"`
	OverheadCost float64        `gorm:"type:decimal(15,2);not null;default:0" json:"overhead_cost"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (v *ProductVariant) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return
}

type BranchProduct struct {
	BranchID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"branch_id"`
	ProductID uuid.UUID `gorm:"type:uuid;primaryKey" json:"product_id"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductFilter struct {
	Search     string     `json:"search"`
	CategoryID *uuid.UUID `json:"category_id"`
	BranchID   *uuid.UUID `json:"branch_id"`
	IsActive   *bool      `json:"is_active"`
	request.PaginationParam
}

type ProductRepository interface {
	Find(ctx context.Context, filter ProductFilter) ([]Product, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDB() *gorm.DB
}

type ProductVariantRepository interface {
	CreateBatch(ctx context.Context, variants []ProductVariant) error
	UpdateBatch(ctx context.Context, variants []ProductVariant) error
	DeleteByProductID(ctx context.Context, productID uuid.UUID) error
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]ProductVariant, error)
	FindByID(ctx context.Context, id uuid.UUID) (*ProductVariant, error)
}

type BranchProductRepository interface {
	Assign(ctx context.Context, branchID uuid.UUID, productIDs []uuid.UUID) error
	Unassign(ctx context.Context, branchID uuid.UUID, productIDs []uuid.UUID) error
	FindByBranch(ctx context.Context, branchID uuid.UUID, filter ProductFilter) ([]Product, int64, error)
}

type ProductUsecase interface {
	Find(ctx context.Context, filter ProductFilter) ([]Product, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Create(ctx context.Context, product *Product, variants []ProductVariant, branchIDs []uuid.UUID) error
	Update(ctx context.Context, product *Product, variants []ProductVariant, branchIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateImage(ctx context.Context, id uuid.UUID, imageURL string) error
	AssignToBranch(ctx context.Context, branchID uuid.UUID, productIDs []uuid.UUID) error
}
