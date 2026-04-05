package main

import (
	"log"

	"github.com/bagusyanuar/genpos-backend/internal/auth/delivery/http"
	"github.com/bagusyanuar/genpos-backend/internal/auth/repository"
	"github.com/bagusyanuar/genpos-backend/internal/auth/usecase"
	"github.com/bagusyanuar/genpos-backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Load Configuration
	conf := config.LoadConfig()

	// 2. Initialize Database
	db := config.InitDB(conf)

	// 3. Initialize App
	app := fiber.New(fiber.Config{
		AppName: "GenPOS Backend",
	})
	app.Use(logger.New())
	app.Use(recover.New())

	// 4. Module Auth
	authRepo := repository.NewAuthRepository(db)
	authUC := usecase.NewAuthUsecase(authRepo)
	authHandler := delivery.NewAuthHandler(authUC)

	// 5. Register Routes
	api := app.Group("/api/v1")
	authHandler.Register(api)

	// 6. Listen
	log.Printf("Server starting on port %s", conf.AppPort)
	if err := app.Listen(":" + conf.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
