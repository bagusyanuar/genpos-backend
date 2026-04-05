package main

import (
	"github.com/bagusyanuar/genpos-backend/internal/shared/bootstrap"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/internal/shared/container"
	"go.uber.org/zap"
)

func main() {
	// 1. Load Configuration
	conf := config.LoadConfig()

	// 2. Initialize Logger
	config.InitLogger(conf)
	config.Log.Info("Starting Application",
		zap.String("app_name", conf.AppName),
		zap.String("app_version", conf.AppVersion),
		zap.String("app_env", conf.AppEnv),
	)

	// 3. Initialize Database
	db := config.InitDB(conf)

	// 3. Initialize Dependency Container
	deps := container.NewContainer(db, conf)

	// 4. Start Application
	bootstrap.Start(conf, deps)
}
