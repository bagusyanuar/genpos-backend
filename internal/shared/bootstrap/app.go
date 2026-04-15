package bootstrap

import (
	"errors"

	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/internal/shared/container"
	"github.com/bagusyanuar/genpos-backend/internal/shared/middleware"
	"github.com/bagusyanuar/genpos-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
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
	app.Use(requestid.New())   // 1. Generate Request ID
	app.Use(middleware.Logger(conf)) // 2. Log using custom Zap middleware
	app.Use(recover.New())     // 3. Panic recovery

	// Static files
	app.Static("/public", "./public")

	// Register Routes
	api := app.Group("/api/v1")

	// Middlewares
	jwtMiddleware := middleware.JWTProtected(conf)

	deps.AuthHandler.Register(api, jwtMiddleware)
	deps.UserHandler.Register(api, jwtMiddleware)
	deps.BranchHandler.Register(api, jwtMiddleware)
	deps.UnitHandler.Register(api, jwtMiddleware)
	deps.CategoryHandler.Register(api, jwtMiddleware)
	deps.MaterialHandler.Register(api, jwtMiddleware)
	deps.InventoryHandler.Register(api, jwtMiddleware)
	deps.ProductHandler.Register(api, jwtMiddleware)
	deps.RecipeHandler.Register(api, jwtMiddleware)
	deps.MediaHandler.Register(api, jwtMiddleware)

	// Start Server
	config.Log.Info("Server is starting...", zap.String("port", conf.AppPort))
	if err := app.Listen(":" + conf.AppPort); err != nil {
		config.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
