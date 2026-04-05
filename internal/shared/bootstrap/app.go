package bootstrap

import (
	"errors"
	"log"

	"github.com/bagusyanuar/genpos-backend/internal/config"
	"github.com/bagusyanuar/genpos-backend/internal/shared/container"
	"github.com/bagusyanuar/genpos-backend/internal/shared/middleware"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Start(conf *config.Config, deps *container.Container) {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: conf.AppName,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(response.Error(err.Error()))
		},
	})

	// Global Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Register Routes
	api := app.Group("/api/v1")

	// Middlewares
	jwtMiddleware := middleware.JWTProtected(conf)

	deps.AuthHandler.Register(api, jwtMiddleware)

	// Start Server
	log.Printf("Server starting on port %s", conf.AppPort)
	if err := app.Listen(":" + conf.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
