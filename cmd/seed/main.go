package main

import (
	"log"

	"github.com/bagusyanuar/genpos-backend/internal/config"
	"github.com/bagusyanuar/genpos-backend/internal/shared/database"
	"github.com/bagusyanuar/genpos-backend/internal/user/repository"
)

func main() {
	// 1. Load configuration
	conf := config.LoadConfig()

	// 2. Initialize Database
	db := config.InitDB(conf)

	// 3. Register Seeders
	seeders := []database.Seeder{
		repository.NewUserSeeder(),
		// register more seeders here
	}

	// 4. Run Seeders
	if err := database.RunSeeders(db, seeders); err != nil {
		log.Fatalf("Failed to seed database: %v\n", err)
	}

	log.Println("Seeding process finished.")
}
