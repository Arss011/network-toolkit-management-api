package database

import (
	"log"

	"toolkit-management/config"

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

	err = db.AutoMigrate()
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("✅ Database connected, migrated successfully")
	return db
}
