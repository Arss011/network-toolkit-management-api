package database

import (
	"log"

	"toolkit-management/config"
	"toolkit-management/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnectionString()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.Toolkit{},
		&models.Loan{},
		&models.Category{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed default admin user
	SeedAdminUser(db)

	if cfg.IsDevelopment() {
		SeedTestData(db)
	}

	log.Println("Database connected, migrated, and seeded successfully")
	return db
}
