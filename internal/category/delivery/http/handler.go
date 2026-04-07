package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/category/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/bagusyanuar/genpos-backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	uc   domain.CategoryUsecase
	conf *config.Config
}

func NewCategoryHandler(uc domain.CategoryUsecase, conf *config.Config) *CategoryHandler {
	return &CategoryHandler{
		uc:   uc,
		conf: conf,
	}
}

func (h *CategoryHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/categories")
	group.Use(authMiddleware)

	group.Get("/", h.Find)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

func (h *CategoryHandler) Find(c *fiber.Ctx) error {
	var query FindCategoryQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid query parameters"))
	}

	filter := domain.CategoryFilter{
		Search:          query.Search,
		Type:            domain.CategoryType(query.Type),
		IsActive:        query.IsActive,
		PaginationParam: query.PaginationParam,
	}

	if query.ParentID != "" {
		parentID, err := uuid.Parse(query.ParentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid parent_id format"))
		}
		filter.ParentID = &parentID
	}

	categories, total, err := h.uc.Find(c.Context(), filter)
	if err != nil {
		config.Log.Error("handler.Find error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		var parentID *string
		if cat.ParentID != nil {
			pID := cat.ParentID.String()
			parentID = &pID
		}

		res[i] = CategoryResponse{
			ID:          cat.ID.String(),
			ParentID:    parentID,
			Level:       cat.Level,
			Name:        cat.Name,
			Slug:        cat.Slug,
			Description: cat.Description,
			Type:        string(cat.Type),
			ImageURL:    cat.ImageURL,
			SortOrder:   cat.SortOrder,
			IsActive:    cat.IsActive,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		}
	}

	pagination := response.Pagination{
		CurrentPage: query.GetPage(),
		Limit:       query.GetLimit(),
		TotalData:   total,
		TotalPage:   int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(res, pagination, "categories found successfully"))
}

func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	categoryID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid category id format"))
	}

	category, err := h.uc.FindByID(c.Context(), categoryID)
	if err != nil {
		config.Log.Error("handler.GetByID error", zap.Error(err), zap.String("id", idStr))
		return c.Status(fiber.StatusNotFound).JSON(response.Error("category not found"))
	}

	var parentID *string
	if category.ParentID != nil {
		pID := category.ParentID.String()
		parentID = &pID
	}

	res := CategoryResponse{
		ID:          category.ID.String(),
		ParentID:    parentID,
		Level:       category.Level,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Type:        string(category.Type),
		ImageURL:    category.ImageURL,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, "category fetched successfully"))
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails("validation error", errs))
	}

	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Type:        domain.CategoryType(req.Type),
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		IsActive:    *req.IsActive,
	}

	if req.ParentID != nil {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid parent_id format"))
		}
		category.ParentID = &parentID
	}

	if err := h.uc.Create(c.Context(), category); err != nil {
		config.Log.Error("handler.Create error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	var parentID *string
	if category.ParentID != nil {
		pID := category.ParentID.String()
		parentID = &pID
	}

	res := CategoryResponse{
		ID:          category.ID.String(),
		ParentID:    parentID,
		Level:       category.Level,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Type:        string(category.Type),
		ImageURL:    category.ImageURL,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, "category created successfully"))
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	categoryID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid category id format"))
	}

	var req UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails("validation error", errs))
	}

	category := &domain.Category{
		ID:          categoryID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Type:        domain.CategoryType(req.Type),
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		IsActive:    *req.IsActive,
	}

	if req.ParentID != nil {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid parent_id format"))
		}
		category.ParentID = &parentID
	}

	if err := h.uc.Update(c.Context(), category); err != nil {
		config.Log.Error("handler.Update error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil, "category updated successfully"))
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	categoryID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid category id format"))
	}

	if err := h.uc.Delete(c.Context(), categoryID); err != nil {
		config.Log.Error("handler.Delete error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil, "category deleted successfully"))
}
