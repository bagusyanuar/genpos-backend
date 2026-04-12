package usecase

import (
	"context"
	"fmt"

	"github.com/bagusyanuar/genpos-backend/internal/product/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/fileupload"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	variantRepo domain.ProductVariantRepository
	branchRepo  domain.BranchProductRepository
	uploader    fileupload.FileUploader
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	variantRepo domain.ProductVariantRepository,
	branchRepo domain.BranchProductRepository,
	uploader fileupload.FileUploader,
) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		variantRepo: variantRepo,
		branchRepo:  branchRepo,
		uploader:    uploader,
	}
}

func (u *productUsecase) Create(ctx context.Context, product *domain.Product, variants []domain.ProductVariant, branchIDs []uuid.UUID) error {
	// 1. Validation: Product must have at least 1 variant
	if len(variants) == 0 {
		return fmt.Errorf("product must have at least one variant")
	}

	// 2. Start Transaction
	tx := u.productRepo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. Save Product (Header)
	if err := tx.WithContext(ctx).Create(product).Error; err != nil {
		tx.Rollback()
		config.Log.Error("failed to create product header", zap.Error(err))
		return fmt.Errorf("product_uc.Create.Header: %w", err)
	}

	// 4. Prepare & Save Variants
	for i := range variants {
		variants[i].ProductID = product.ID
	}

	if err := tx.WithContext(ctx).Create(&variants).Error; err != nil {
		tx.Rollback()
		config.Log.Error("failed to create product variants", zap.Error(err))
		return fmt.Errorf("product_uc.Create.Variants: %w", err)
	}

	// 5. Assign to Branches if provided
	if len(branchIDs) > 0 {
		branchProducts := make([]domain.BranchProduct, 0, len(branchIDs))
		for _, bID := range branchIDs {
			branchProducts = append(branchProducts, domain.BranchProduct{
				BranchID:  bID,
				ProductID: product.ID,
				IsActive:  true,
			})
		}
		if err := tx.WithContext(ctx).Create(&branchProducts).Error; err != nil {
			tx.Rollback()
			config.Log.Error("failed to assign product to branches", zap.Error(err))
			return fmt.Errorf("product_uc.Create.BranchAssignment: %w", err)
		}
	}

	// 6. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("product_uc.Create.Commit: %w", err)
	}

	return nil
}

func (u *productUsecase) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product_uc.FindByID: %w", err)
	}
	return product, nil
}

func (u *productUsecase) Find(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	return u.productRepo.Find(ctx, filter)
}

func (u *productUsecase) Update(ctx context.Context, product *domain.Product, variants []domain.ProductVariant, branchIDs []uuid.UUID) error {
	// 1. Validation: Product must have at least 1 variant
	if len(variants) == 0 {
		return fmt.Errorf("product must have at least one variant")
	}

	// 2. Find existing for image cleanup and variant diffing
	existing, err := u.productRepo.FindByID(ctx, product.ID)
	if err != nil {
		return fmt.Errorf("product_uc.Update.FindByID: %w", err)
	}

	// 3. Start Transaction
	tx := u.productRepo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. Update Product (Header)
	if err := tx.WithContext(ctx).Model(product).
		Select("category_id", "name", "description", "is_active", "updated_at").
		Updates(product).Error; err != nil {
		tx.Rollback()
		config.Log.Error("failed to update product header", zap.Error(err))
		return fmt.Errorf("product_uc.Update.Header: %w", err)
	}

	// 5. Optimized Variant Sync (Upsert/Diff)
	incomingVariantIDs := make([]uuid.UUID, 0, len(variants))
	for i := range variants {
		variants[i].ProductID = product.ID
		if variants[i].ID != uuid.Nil {
			if err := tx.WithContext(ctx).Model(&variants[i]).
				Select("name", "sku", "price", "is_active", "updated_at").
				Updates(&variants[i]).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("product_uc.Update.UpdateVariant(%s): %w", variants[i].ID, err)
			}
			incomingVariantIDs = append(incomingVariantIDs, variants[i].ID)
		} else {
			if err := tx.WithContext(ctx).Create(&variants[i]).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("product_uc.Update.CreateVariant: %w", err)
			}
			incomingVariantIDs = append(incomingVariantIDs, variants[i].ID)
		}
	}

	if err := tx.WithContext(ctx).
		Where("product_id = ? AND id NOT IN ?", product.ID, incomingVariantIDs).
		Delete(&domain.ProductVariant{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("product_uc.Update.CleanupVariants: %w", err)
	}

	// 6. Branch Assignment Sync
	if err := tx.WithContext(ctx).
		Where("product_id = ?", product.ID).
		Delete(&domain.BranchProduct{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("product_uc.Update.ClearBranchAssignments: %w", err)
	}

	if len(branchIDs) > 0 {
		branchProducts := make([]domain.BranchProduct, 0, len(branchIDs))
		for _, bID := range branchIDs {
			branchProducts = append(branchProducts, domain.BranchProduct{
				BranchID:  bID,
				ProductID: product.ID,
				IsActive:  true,
			})
		}
		if err := tx.WithContext(ctx).Create(&branchProducts).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("product_uc.Update.RecreateBranchAssignments: %w", err)
		}
	}

	// 7. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("product_uc.Update.Commit: %w", err)
	}

	// 8. Image Cleanup
	if product.ImageURL != nil && existing.ImageURL != nil && *product.ImageURL != *existing.ImageURL {
		if err := u.uploader.Delete(*existing.ImageURL); err != nil {
			config.Log.Warn("failed to delete old image after product update", zap.Error(err))
		}
	}

	return nil
}

func (u *productUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.productRepo.Delete(ctx, id)
}

func (u *productUsecase) UpdateImage(ctx context.Context, id uuid.UUID, imageURL string) error {
	// 1. Find existing product
	existing, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product_uc.UpdateImage.FindByID: %w", err)
	}

	// 2. Prepare update
	oldURL := ""
	if existing.ImageURL != nil {
		oldURL = *existing.ImageURL
	}

	existing.ImageURL = &imageURL
	if err := u.productRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("product_uc.UpdateImage.Repo: %w", err)
	}

	// 3. Delete old image if exists and changed
	if oldURL != "" && oldURL != imageURL {
		if err := u.uploader.Delete(oldURL); err != nil {
			config.Log.Warn("failed to delete old product image", zap.Error(err), zap.String("url", oldURL))
		}
	}

	return nil
}

func (u *productUsecase) AssignToBranch(ctx context.Context, branchID uuid.UUID, productIDs []uuid.UUID) error {
	return u.branchRepo.Assign(ctx, branchID, productIDs)
}
