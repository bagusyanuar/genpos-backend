package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/branch/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/bagusyanuar/genpos-backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BranchHandler struct {
	uc   domain.BranchUsecase
	conf *config.Config
}

func NewBranchHandler(uc domain.BranchUsecase, conf *config.Config) *BranchHandler {
	return &BranchHandler{
		uc:   uc,
		conf: conf,
	}
}

func (h *BranchHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/branches")
	group.Use(authMiddleware)

	group.Get("/", h.Find)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

func (h *BranchHandler) Find(c *fiber.Ctx) error {
	var query FindBranchQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid query parameters"))
	}

	filter := domain.BranchFilter{
		Search:          query.Search,
		PaginationParam: query.PaginationParam,
	}

	branches, total, err := h.uc.Find(c.Context(), filter)
	if err != nil {
		config.Log.Error("handler.Find error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := make([]BranchResponse, len(branches))
	for i, b := range branches {
		res[i] = BranchResponse{
			ID:        b.ID.String(),
			Name:      b.Name,
			IsDefault: b.IsDefault,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
	}

	pagination := response.Pagination{
		CurrentPage: query.GetPage(),
		Limit:       query.GetLimit(),
		TotalData:   total,
		TotalPage:   int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(res, pagination, "branches found successfully"))
}

func (h *BranchHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	branchID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid branch id format"))
	}

	branch, err := h.uc.FindByID(c.Context(), branchID)
	if err != nil {
		config.Log.Error("handler.GetByID error", zap.Error(err), zap.String("id", idStr))
		return c.Status(fiber.StatusNotFound).JSON(response.Error("branch not found"))
	}

	res := BranchResponse{
		ID:        branch.ID.String(),
		Name:      branch.Name,
		IsDefault: branch.IsDefault,
		CreatedAt: branch.CreatedAt,
		UpdatedAt: branch.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, "branch fetched successfully"))
}

func (h *BranchHandler) Create(c *fiber.Ctx) error {
	var req CreateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails("validation error", errs))
	}

	branch := &domain.Branch{
		Name:      req.Name,
		IsDefault: req.IsDefault,
	}

	if err := h.uc.Create(c.Context(), branch); err != nil {
		config.Log.Error("handler.Create error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := BranchResponse{
		ID:        branch.ID.String(),
		Name:      branch.Name,
		IsDefault: branch.IsDefault,
		CreatedAt: branch.CreatedAt,
		UpdatedAt: branch.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, "branch created successfully"))
}

func (h *BranchHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	branchID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid branch id format"))
	}

	var req UpdateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails("validation error", errs))
	}

	branch, err := h.uc.FindByID(c.Context(), branchID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error("branch not found"))
	}

	branch.Name = req.Name
	branch.IsDefault = req.IsDefault

	if err := h.uc.Update(c.Context(), branch); err != nil {
		config.Log.Error("handler.Update error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := BranchResponse{
		ID:        branch.ID.String(),
		Name:      branch.Name,
		IsDefault: branch.IsDefault,
		CreatedAt: branch.CreatedAt,
		UpdatedAt: branch.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, "branch updated successfully"))
}

func (h *BranchHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	branchID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid branch id format"))
	}

	if err := h.uc.Delete(c.Context(), branchID); err != nil {
		config.Log.Error("handler.Delete error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil, "branch deleted successfully"))
}
