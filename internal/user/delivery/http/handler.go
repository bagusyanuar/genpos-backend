package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/internal/user/domain"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/bagusyanuar/genpos-backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	uc   domain.UserUsecase
	conf *config.Config
}

func NewUserHandler(uc domain.UserUsecase, conf *config.Config) *UserHandler {
	return &UserHandler{
		uc:   uc,
		conf: conf,
	}
}

func (h *UserHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/users")
	group.Use(authMiddleware)

	group.Get("/", h.Find)
	group.Get("/:id", h.GetByID)
	group.Post("/", h.Create)
}

func (h *UserHandler) Find(c *fiber.Ctx) error {
	var query FindUserQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid query parameters"))
	}

	filter := domain.UserFilter{
		Search:          query.Search,
		PaginationParam: query.PaginationParam,
	}

	users, total, err := h.uc.Find(c.Context(), filter)
	if err != nil {
		config.Log.Error("handler.Find error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = UserResponse{
			ID:        u.ID.String(),
			Email:     u.Email,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
		}
	}

	pagination := response.Pagination{
		CurrentPage: query.GetPage(),
		Limit:       query.GetLimit(),
		TotalData:   total,
		TotalPage:   int((total + int64(query.GetLimit()) - 1) / int64(query.GetLimit())),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(res, pagination, "users found successfully"))
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid user id format"))
	}

	user, err := h.uc.FindByID(c.Context(), userID)
	if err != nil {
		config.Log.Error("handler.GetByID error", zap.Error(err), zap.String("id", idStr))
		return c.Status(fiber.StatusNotFound).JSON(response.Error("user not found"))
	}

	res := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, "user fetched successfully"))
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails("validation error", errs))
	}

	user := &domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	if err := h.uc.Create(c.Context(), user); err != nil {
		config.Log.Error("handler.Create error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, "user created successfully"))
}
