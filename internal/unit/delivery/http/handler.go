package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/unit/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/bagusyanuar/genpos-backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UnitHandler struct {
	uc   domain.UnitUsecase
	conf *config.Config
}

func NewUnitHandler(uc domain.UnitUsecase, conf *config.Config) *UnitHandler {
	return &UnitHandler{
		uc:   uc,
		conf: conf,
	}
}

func (h *UnitHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/units")
	group.Use(authMiddleware)

	group.Get("/", h.Find)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

func (h *UnitHandler) Find(c *fiber.Ctx) error {
	var query FindUnitQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid query parameters"))
	}

	filter := domain.UnitFilter{
		Search:          query.Search,
		PaginationParam: query.PaginationParam,
	}

	units, total, err := h.uc.Find(c.Context(), filter)
	if err != nil {
		config.Log.Error("handler.Find error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := make([]UnitResponse, len(units))
	for i, u := range units {
		res[i] = UnitResponse{
			ID:        u.ID.String(),
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	pagination := response.Pagination{
		CurrentPage: query.GetPage(),
		Limit:       query.GetLimit(),
		TotalData:   total,
		TotalPage:   int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(res, pagination, "units found successfully"))
}

func (h *UnitHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	unitID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid unit id format"))
	}

	unit, err := h.uc.FindByID(c.Context(), unitID)
	if err != nil {
		config.Log.Error("handler.GetByID error", zap.Error(err), zap.String("id", idStr))
		return c.Status(fiber.StatusNotFound).JSON(response.Error("unit not found"))
	}

	res := UnitResponse{
		ID:        unit.ID.String(),
		Name:      unit.Name,
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, "unit fetched successfully"))
}

func (h *UnitHandler) Create(c *fiber.Ctx) error {
	var req CreateUnitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails("validation error", errs))
	}

	unit := &domain.Unit{
		Name: req.Name,
	}

	if err := h.uc.Create(c.Context(), unit); err != nil {
		config.Log.Error("handler.Create error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := UnitResponse{
		ID:        unit.ID.String(),
		Name:      unit.Name,
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, "unit created successfully"))
}

func (h *UnitHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	unitID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid unit id format"))
	}

	var req UpdateUnitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails("validation error", errs))
	}

	unit, err := h.uc.FindByID(c.Context(), unitID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error("unit not found"))
	}

	unit.Name = req.Name

	if err := h.uc.Update(c.Context(), unit); err != nil {
		config.Log.Error("handler.Update error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := UnitResponse{
		ID:        unit.ID.String(),
		Name:      unit.Name,
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, "unit updated successfully"))
}

func (h *UnitHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	unitID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid unit id format"))
	}

	if err := h.uc.Delete(c.Context(), unitID); err != nil {
		config.Log.Error("handler.Delete error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil, "unit deleted successfully"))
}
