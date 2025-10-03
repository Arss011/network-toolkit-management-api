package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	Environment string
	DB          *gorm.DB
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env." + getEnv("GO_ENV", "development"))
	_ = godotenv.Load()

	env := getEnv("GO_ENV", "development")

	config := &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "root"),
		DBName:      getEnv("DB_NAME", "toolkit_db"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: env,
	}

	if err := config.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	return config
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) InitDB() error {
	dsn := c.GetDBConnectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	c.DB = db
	return nil
}

func (c *Config) CloseDB() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
