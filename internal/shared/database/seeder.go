package database

import (
	"log"

	"gorm.io/gorm"
)

// Seeder defines the common interface for all database seeders.
type Seeder interface {
	Run(db *gorm.DB) error
}

// RunSeeders executes all registered seeders in the system.
func RunSeeders(db *gorm.DB, registeredSeeders []Seeder) error {
	log.Println("Starting database seeding...")

	for _, seeder := range registeredSeeders {
		if err := seeder.Run(db); err != nil {
			log.Printf("Seeder failed: %v\n", err)
			return err
		}
	}

	log.Println("Database seeding completed successfully.")
	return nil
}
