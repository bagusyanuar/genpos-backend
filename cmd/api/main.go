package main

import (
	"github.com/bagusyanuar/genpos-backend/internal/config"
	"github.com/bagusyanuar/genpos-backend/internal/shared/bootstrap"
	"github.com/bagusyanuar/genpos-backend/internal/shared/container"
)

func main() {
	// 1. Load Configuration
	conf := config.LoadConfig()

	// 2. Initialize Database
	db := config.InitDB(conf)

	// 3. Initialize Dependency Container
	deps := container.NewContainer(db)

	// 4. Start Application
	bootstrap.Start(conf, deps)
}
