package usecase

import (
	"context"
	"fmt"

	inventoryDomain "github.com/bagusyanuar/genpos-backend/internal/inventory/domain"
	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/fileupload"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type materialUsecase struct {
	materialRepo  domain.MaterialRepository
	uomRepo       domain.MaterialUOMRepository
	inventoryRepo inventoryDomain.InventoryRepository
	auditRepo     domain.MaterialAuditRepository
	uploader      fileupload.FileUploader
}

func NewMaterialUsecase(
	materialRepo domain.MaterialRepository,
	uomRepo domain.MaterialUOMRepository,
	inventoryRepo inventoryDomain.InventoryRepository,
	auditRepo domain.MaterialAuditRepository,
	uploader fileupload.FileUploader,
) domain.MaterialUsecase {
	return &materialUsecase{
		materialRepo:  materialRepo,
		uomRepo:       uomRepo,
		inventoryRepo: inventoryRepo,
		auditRepo:     auditRepo,
		uploader:      uploader,
	}
}

func (u *materialUsecase) Create(ctx context.Context, material *domain.Material, uoms []domain.MaterialUOM) error {
	// 1. Validation for MaterialType
	if material.MaterialType != "RAW" && material.MaterialType != "SEMI_FINISHED" {
		return fmt.Errorf("invalid material type: %s", material.MaterialType)
	}

	// 2. Validation for UOMs: Must have exactly 1 default
	hasDefault := false
	for _, uom := range uoms {
		if uom.IsDefault {
			if hasDefault {
				return fmt.Errorf("multiple default UOMs provided")
			}
			hasDefault = true
		}
	}
	if !hasDefault {
		return fmt.Errorf("default UOM (Base Unit) must be provided")
	}

	// 3. Start Transaction
	tx := u.materialRepo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. Save Material in transaction
	if err := tx.WithContext(ctx).Create(material).Error; err != nil {
		tx.Rollback()
		config.Log.Error("failed to create material in transaction", zap.Error(err))
		return fmt.Errorf("material_uc.Create.Material: %w", err)
	}

	// 5. Prepare UOMs with MaterialID and Save
	for i := range uoms {
		uoms[i].MaterialID = material.ID
	}

	if err := tx.WithContext(ctx).Create(&uoms).Error; err != nil {
		tx.Rollback()
		config.Log.Error("failed to create material uoms in transaction", zap.Error(err))
		return fmt.Errorf("material_uc.Create.UOMs: %w", err)
	}

	// 6. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("material_uc.Create.Commit: %w", err)
	}

	return nil
}

func (u *materialUsecase) Update(ctx context.Context, material *domain.Material) error {
	// 1. Validation for MaterialType
	if material.MaterialType != "" && material.MaterialType != "RAW" && material.MaterialType != "SEMI_FINISHED" {
		return fmt.Errorf("invalid material type: %s", material.MaterialType)
	}

	// 2. Find existing material to handle image cleanup
	existing, err := u.materialRepo.FindByID(ctx, material.ID)
	if err != nil {
		return fmt.Errorf("material_uc.Update.FindByID: %w", err)
	}

	// 3. Update Material
	if err := u.materialRepo.Update(ctx, material); err != nil {
		return fmt.Errorf("material_uc.Update.Repo: %w", err)
	}

	// 4. Image Cleanup: Delete old image ONLY if DB update succeeded and URL changed
	if material.ImageURL != nil && existing.ImageURL != nil && *material.ImageURL != *existing.ImageURL {
		if err := u.uploader.Delete(*existing.ImageURL); err != nil {
			config.Log.Warn("failed to delete old image after successful update",
				zap.Error(err),
				zap.String("url", *existing.ImageURL),
			)
		}
	}

	return nil
}

func (u *materialUsecase) UpdateImage(ctx context.Context, id uuid.UUID, imageURL string) error {
	// 1. Find existing material
	existing, err := u.materialRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("material_uc.UpdateImage.FindByID: %w", err)
	}

	// 2. Update Image URL in DB first
	oldURL := ""
	if existing.ImageURL != nil {
		oldURL = *existing.ImageURL
	}

	existing.ImageURL = &imageURL
	if err := u.materialRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("material_uc.UpdateImage.Repo: %w", err)
	}

	// 3. Delete old image ONLY if DB update succeeded
	if oldURL != "" && oldURL != imageURL {
		if err := u.uploader.Delete(oldURL); err != nil {
			config.Log.Warn("failed to delete old image after successful patch",
				zap.Error(err),
				zap.String("url", oldURL),
			)
		}
	}

	return nil
}

func (u *materialUsecase) FindByID(ctx context.Context, id uuid.UUID) (*domain.Material, error) {
	material, err := u.materialRepo.FindByID(ctx, id)
	if err != nil {
		config.Log.Error("failed to find material by id",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, fmt.Errorf("material_uc.FindByID: %w", err)
	}
	return material, nil
}

func (u *materialUsecase) Find(ctx context.Context, filter domain.MaterialFilter) ([]domain.Material, int64, error) {
	materials, total, err := u.materialRepo.Find(ctx, filter)
	if err != nil {
		config.Log.Error("failed to find materials",
			zap.Error(err),
			zap.String("search", filter.Search),
		)
		return nil, 0, fmt.Errorf("material_uc.Find: %w", err)
	}

	return materials, total, nil
}

func (u *materialUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.materialRepo.Delete(ctx, id); err != nil {
		config.Log.Error("failed to delete material",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return fmt.Errorf("material_uc.Delete: %w", err)
	}
	return nil
}

func (u *materialUsecase) RecalibrateUOM(ctx context.Context, materialID uuid.UUID, targetUOMID uuid.UUID, userID uuid.UUID) error {
	// 1. Fetch current multipliers to find target CF
	uoms, err := u.uomRepo.Find(ctx, materialID)
	if err != nil {
		return fmt.Errorf("material_uc.RecalibrateUOM.FindUOMs: %w", err)
	}

	var targetUOM *domain.MaterialUOM
	for i := range uoms {
		if uoms[i].ID == targetUOMID {
			targetUOM = &uoms[i]
			break
		}
	}

	if targetUOM == nil {
		return fmt.Errorf("target UOM not found for this material")
	}

	if targetUOM.IsDefault {
		return fmt.Errorf("target UOM is already the default (Base Unit)")
	}

	// 2. Fetch Material to get current base cost
	material, err := u.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		return fmt.Errorf("material_uc.RecalibrateUOM.FindMaterial: %w", err)
	}

	// 3. Mathematical Recalibration
	cf := targetUOM.Multiplier
	if cf <= 0 {
		return fmt.Errorf("invalid conversion factor: %.4f", cf)
	}

	// 4. Execute Transaction
	tx := u.materialRepo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4a. Normalise Multipliers (Optimized: Single query)
	if err := u.uomRepo.RecalibrateUOMs(ctx, tx, materialID, cf, targetUOMID); err != nil {
		tx.Rollback()
		return fmt.Errorf("material_uc.RecalibrateUOM.UpdateUOMs: %w", err)
	}

	// 4b. Convert Inventories Stock & Min Stock
	if err := u.inventoryRepo.RecalibrateStock(ctx, tx, materialID, cf); err != nil {
		tx.Rollback()
		return fmt.Errorf("material_uc.RecalibrateUOM.UpdateStock: %w", err)
	}

	// 4c. Convert Material Base Cost
	newBaseCost := material.BaseCost * cf
	if err := tx.WithContext(ctx).Model(material).Update("base_cost", newBaseCost).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("material_uc.RecalibrateUOM.UpdateCost: %w", err)
	}

	// 4d. Record Audit Trail in MaterialAudit
	audit := &domain.MaterialAudit{
		MaterialID: materialID,
		Action:     "RECALIBRATE",
		Note:       fmt.Sprintf("Base UOM changed. CF: %.4f. Cost: %.2f -> %.2f", cf, material.BaseCost, newBaseCost),
		CreatedBy:  userID,
	}
	if err := tx.WithContext(ctx).Create(audit).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("material_uc.RecalibrateUOM.AuditLog: %w", err)
	}

	// 5. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("material_uc.RecalibrateUOM.Commit: %w", err)
	}

	config.Log.Info("Material UOM Recalibrated successfully",
		zap.String("material_id", materialID.String()),
		zap.String("target_uom_id", targetUOMID.String()),
		zap.Float64("cf", cf),
		zap.String("user_id", userID.String()),
	)

	return nil
}
