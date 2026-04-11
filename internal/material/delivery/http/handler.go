package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/material/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/request"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MaterialHandler struct {
	uc     domain.MaterialUsecase
	uomUC  domain.MaterialUOMUsecase
	conf   *config.Config
}

func NewMaterialHandler(uc domain.MaterialUsecase, uomUC domain.MaterialUOMUsecase, conf *config.Config) *MaterialHandler {
	return &MaterialHandler{
		uc:     uc,
		uomUC:  uomUC,
		conf:   conf,
	}
}

func (h *MaterialHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/materials")
	group.Use(authMiddleware)

	group.Get("/", h.Find)
	group.Get("/:id", h.FindByID)
	group.Get("/:id/uoms", h.FindUOMs)
}

func (h *MaterialHandler) FindByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid material id format"))
	}

	material, err := h.uc.FindByID(c.Context(), id)
	if err != nil {
		config.Log.Error("handler.FindByID material error", zap.Error(err), zap.String("id", idStr))
		return c.Status(fiber.StatusNotFound).JSON(response.Error(err.Error()))
	}

	res := ToMaterialResponse(*material)
	return c.Status(fiber.StatusOK).JSON(response.Success(res, "material found successfully"))
}

type FindMaterialQuery struct {
	Search   string `query:"search"`
	request.PaginationParam
}

func (h *MaterialHandler) Find(c *fiber.Ctx) error {
	var query FindMaterialQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid query parameters"))
	}

	filter := domain.MaterialFilter{
		Search:          query.Search,
		PaginationParam: query.PaginationParam,
	}

	materials, total, err := h.uc.Find(c.Context(), filter)
	if err != nil {
		config.Log.Error("handler.Find material error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := ToMaterialListResponse(materials)

	pagination := response.Pagination{
		CurrentPage: query.GetPage(),
		Limit:       query.GetLimit(),
		TotalData:   total,
		TotalPage:   int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(res, pagination, "materials found successfully"))
}

func (h *MaterialHandler) FindUOMs(c *fiber.Ctx) error {
	idStr := c.Params("id")
	materialID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid material id format"))
	}

	uoms, err := h.uomUC.Find(c.Context(), materialID)
	if err != nil {
		config.Log.Error("handler.FindUOMs error", zap.Error(err), zap.String("material_id", idStr))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := ToMaterialUOMListResponse(uoms)
	return c.Status(fiber.StatusOK).JSON(response.Success(res, "material UOMs found successfully"))
}
