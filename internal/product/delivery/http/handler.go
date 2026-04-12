package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/product/domain"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/fileupload"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/bagusyanuar/genpos-backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	uc       domain.ProductUsecase
	uploader fileupload.FileUploader
	conf     *config.Config
}

func NewProductHandler(uc domain.ProductUsecase, uploader fileupload.FileUploader, conf *config.Config) *ProductHandler {
	return &ProductHandler{
		uc:       uc,
		uploader: uploader,
		conf:     conf,
	}
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationError(err))
	}

	product := req.ToEntity()
	variants := make([]domain.ProductVariant, 0)
	for _, v := range req.Variants {
		variants = append(variants, v.ToEntity())
	}

	if err := h.uc.Create(c.Context(), product, variants, req.BranchIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(ToProductResponse(*product), "product created successfully"))
}

func (h *ProductHandler) PatchImage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid product id format"))
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("image file is required"))
	}

	// Upload to products folder
	url, err := h.uploader.Upload(file, "products", []string{".jpg", ".jpeg", ".png", ".webp"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	if err := h.uc.UpdateImage(c.Context(), id, url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(fiber.Map{"image_url": url}, "image updated successfully"))
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid product id format"))
	}

	product, err := h.uc.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(ToProductResponse(*product), "product found successfully"))
}

func (h *ProductHandler) Find(c *fiber.Ctx) error {
	var filter domain.ProductFilter
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	products, total, err := h.uc.Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	pagination := response.Pagination{
		CurrentPage: filter.Page,
		Limit:       filter.Limit,
		TotalData:   total,
		TotalPage:   int((total + int64(filter.Limit) - 1) / int64(filter.Limit)),
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPagination(ToProductListResponse(products), pagination, "products found successfully"))
}

func (h *ProductHandler) Register(router fiber.Router, auth fiber.Handler) {
	product := router.Group("/products", auth)
	product.Post("/", h.Create)
	product.Get("/", h.Find)
	product.Get("/:id", h.GetByID)
	product.Put("/:id", h.Update)
	product.Delete("/:id", h.Delete)
	product.Patch("/:id/image", h.PatchImage)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid product id format"))
	}

	var req CreateProductRequest // Reuse CreateProductRequest for PUT
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ValidationError(err))
	}

	product := req.ToEntity()
	product.ID = id
	variants := make([]domain.ProductVariant, 0)
	for _, v := range req.Variants {
		variants = append(variants, v.ToEntity())
	}

	if err := h.uc.Update(c.Context(), product, variants, req.BranchIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(ToProductResponse(*product), "product updated successfully"))
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid product id format"))
	}

	if err := h.uc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil, "product deleted successfully"))
}
