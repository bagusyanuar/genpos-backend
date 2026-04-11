package http

import (
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/fileupload"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MediaHandler struct {
	uploader fileupload.FileUploader
	conf     *config.Config
}

func NewMediaHandler(uploader fileupload.FileUploader, conf *config.Config) *MediaHandler {
	return &MediaHandler{
		uploader: uploader,
		conf:     conf,
	}
}

func (h *MediaHandler) Register(router fiber.Router, authMiddleware fiber.Handler) {
	group := router.Group("/uploads")
	group.Use(authMiddleware)

	group.Post("/", h.Upload)
}

func (h *MediaHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("file field is required"))
	}

	// Determine subdir from query or use "general"
	subDir := c.Query("folder", "general")

	// Allowed extensions (you can make this configurable)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".webp", ".pdf"}

	url, err := h.uploader.Upload(file, subDir, allowedExts)
	if err != nil {
		config.Log.Error("failed to upload file", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	res := ToUploadResponse(url)
	return c.Status(fiber.StatusCreated).JSON(response.Success(res, "file uploaded successfully"))
}
