package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"toolkit-management/internal/models"
)

// SeedAdminUser trigger if emty
func SeedAdminUser(db *gorm.DB) {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Error during initial user count check: %v", err)
		return
	}

	if userCount == 0 {
		log.Println("Admin seed: Creating default admin user.")

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Admin seed: Failed to hash password: %v", err)
			return
		}

		adminUser := &models.User{
			Username:    "admin",
			Email:       "admin@gmail.com",
			FullName:    "Administrator",
			Password:    string(hashedPassword),
			Role:        "admin",
			Department:  "IT",
			PhoneNumber: "+6281234567890",
			IsActive:    true,
		}

		if err := db.Create(adminUser).Error; err != nil {
			log.Printf("Admin seed: Failed to create admin user: %v", err)
			return
		}

		log.Println("Admin seed: Default admin user created successfully.")

	} else {
		log.Println("Admin seed: Skipping, users found in database.")
	}
}

func SeedTestData(db *gorm.DB) {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)

	if userCount > 0 {
		log.Println("Test data seed: Creating sample demo user.")

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

		sampleUser := &models.User{
			Username:    "demo",
			Email:       "demo@gmail.com",
			FullName:    "Demo User",
			Password:    string(hashedPassword),
			Role:        "user",
			Department:  "Operations",
			PhoneNumber: "+6280987654321",
			IsActive:    true,
		}

		var existingUser models.User
		result := db.Where("username = ?", sampleUser.Username).First(&existingUser)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(sampleUser).Error; err != nil {
				log.Printf("Test data seed: Failed to create demo user: %v", err)
				return
			}
			log.Println("Test data seed: Sample demo user created successfully.")
		} else if result.Error != nil {
			log.Printf("Test data seed: Error checking for existing demo user: %v", result.Error)
		} else {
			log.Println("Test data seed: Sample demo user already exists, skipping.")
		}
	}
}
