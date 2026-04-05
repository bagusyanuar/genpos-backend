package repository

import (
	"log"

	"github.com/bagusyanuar/genpos-backend/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userSeeder struct{}

// NewUserSeeder provides a new instance of the UserSeeder.
func NewUserSeeder() *userSeeder {
	return &userSeeder{}
}

// Run executes the database seeding for the users table.
func (s *userSeeder) Run(db *gorm.DB) error {
	log.Println("Seeding users...")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminUser := &domain.User{
		Email:    "admin@genpos.com",
		Username: "admin",
		Password: string(hashedPassword),
	}

	// Use GORM's FirstOrCreate for idempotency (upsert style)
	var existingUser domain.User
	result := db.Where("email = ?", adminUser.Email).FirstOrCreate(&existingUser, adminUser)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Default admin user created: %s\n", adminUser.Email)
	} else {
		log.Printf("User %s already exists, skipping.\n", adminUser.Email)
	}

	return nil
}
