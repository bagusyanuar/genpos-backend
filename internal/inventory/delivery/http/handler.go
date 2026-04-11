package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/inventory/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InventoryHandler struct {
	uc   domain.InventoryUsecase
	conf *config.Config
}

func NewInventoryHandler(uc domain.InventoryUsecase, conf *config.Config) *InventoryHandler {
	return &InventoryHandler{
		uc:   uc,
		conf: conf,
	}
}

func (h *InventoryHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/inventories")
	group.Use(authMiddleware)

	group.Get("/", h.Find)
	group.Get("/summary", h.GetSummary)
}

type FindInventoryQuery struct {
	BranchID   string `query:"branch_id"`
	MaterialID string `query:"material_id"`
	Search     string `query:"search"`
	request.PaginationParam
}

func (h *InventoryHandler) Find(c *fiber.Ctx) error {
	var query FindInventoryQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid query parameters"))
	}

	filter := domain.InventoryFilter{
		Search:          query.Search,
		PaginationParam: query.PaginationParam,
	}

	// Multi-tenancy check
	if query.BranchID != "" {
		branchID, err := uuid.Parse(query.BranchID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid branch_id format"))
		}
		filter.BranchID = branchID
	} else {
		ctxBranchID := c.Locals("branch_id")
		if ctxBranchID != nil {
			if id, ok := ctxBranchID.(uuid.UUID); ok {
				filter.BranchID = id
			} else if idStr, ok := ctxBranchID.(string); ok {
				id, _ := uuid.Parse(idStr)
				filter.BranchID = id
			}
		}
	}

	if filter.BranchID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("branch_id is required"))
	}

	// Material Filter
	if query.MaterialID != "" {
		materialID, err := uuid.Parse(query.MaterialID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid material_id format"))
		}
		filter.MaterialID = materialID
	}

	inventories, total, err := h.uc.Find(c.Context(), filter)
	if err != nil {
		config.Log.Error("handler.Find inventories error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := ToInventoryListResponse(inventories)

	pagination := response.Pagination{
		CurrentPage: query.GetPage(),
		Limit:       query.GetLimit(),
		TotalData:   total,
		TotalPage:   int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(res, pagination, "inventories found successfully"))
}

func (h *InventoryHandler) GetSummary(c *fiber.Ctx) error {
	var query FindInventoryQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid query parameters"))
	}

	filter := domain.InventoryFilter{
		Search:          query.Search,
		PaginationParam: query.PaginationParam,
	}

	var branchID uuid.UUID

	// Multi-tenancy check (HQ view allowed if no branch_id)
	if query.BranchID != "" {
		id, err := uuid.Parse(query.BranchID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid branch_id format"))
		}
		branchID = id
	} else {
		ctxBranchID := c.Locals("branch_id")
		if ctxBranchID != nil {
			if id, ok := ctxBranchID.(uuid.UUID); ok {
				branchID = id
			} else if idStr, ok := ctxBranchID.(string); ok {
				id, _ := uuid.Parse(idStr)
				branchID = id
			}
		}
	}

	// Material Filter
	if query.MaterialID != "" {
		materialID, err := uuid.Parse(query.MaterialID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid material_id format"))
		}
		filter.MaterialID = materialID
	}

	views, total, err := h.uc.GetSummary(c.Context(), branchID, filter)
	if err != nil {
		config.Log.Error("handler.GetSummary error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	pagination := response.Pagination{
		CurrentPage: query.GetPage(),
		Limit:       query.GetLimit(),
		TotalData:   total,
		TotalPage:   int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(views, pagination, "inventory summary found successfully"))
}
