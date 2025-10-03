package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"toolkit-management/config"
	"toolkit-management/internal/handlers"
	"toolkit-management/internal/repositories"
	"toolkit-management/internal/services"
	"toolkit-management/pkg/auth"
	"toolkit-management/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	// Set Mode Gin
	if cfg.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init Database
	db := database.InitDB(cfg)

	//seed user admin
	database.SeedAdminUser(db)

	if cfg.Environment == "development" {
		database.SeedTestData(db)
	}

	// Init Auth Service
	authConfig := auth.AuthConfig{
		SecretKey:     "secrect-key-rahasia",
		TokenDuration: time.Hour * time.Duration(24),
	}
	authService := auth.NewAuthService(authConfig)

	// init repo & service
	userRepo := repositories.NewUserRepository(db)
	toolkitRepo := repositories.NewToolkitRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	loanRepo := repositories.NewLoanRepository(db)

	userService := services.NewUserService(userRepo)
	toolkitService := services.NewToolkitService(toolkitRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	loanService := services.NewLoanService(loanRepo, toolkitRepo)

	// init handler
	userHandler := handlers.NewUserHandler(userService)
	toolkitHandler := handlers.NewToolkitHandler(toolkitService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	loanHandler := handlers.NewLoanHandler(loanService)

	// Setup Router
	router := gin.Default()

	// Config CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// Define route
	api := router.Group("/api")
	{
		// Public routes
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "Service is healthy",
			})
		})
		api.POST("/auth/login", userHandler.Login)

		// Protected routes
		protected := api.Group("")
		protected.Use(authService.RequireAuth())
		{
			// Current user
			protected.GET("/auth/me", userHandler.GetCurrentUser)

			// User routes - Admin only
			users := protected.Group("/users")
			users.Use(authService.RequireAdmin())
			{
				users.POST("", userHandler.Create)
				users.GET("", userHandler.GetAll)
				users.POST("/search", userHandler.GetAll)
				users.GET("/:id", userHandler.GetByID)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
			}

			// Toolkit routes - Admin & user
			toolkits := protected.Group("/toolkits")
			{
				// Admin only
				toolkitsAdmin := toolkits.Group("")
				toolkitsAdmin.Use(authService.RequireAdmin())
				{
					toolkitsAdmin.POST("", toolkitHandler.Create)
					toolkitsAdmin.PUT("/:id", toolkitHandler.Update)
					toolkitsAdmin.DELETE("/:id", toolkitHandler.Delete)
					toolkitsAdmin.PATCH("/:id/stock", toolkitHandler.UpdateStock)
				}

				// All authenticated users
				toolkits.GET("", toolkitHandler.GetAll)
				toolkits.POST("/search", toolkitHandler.GetAll)
				toolkits.GET("/:id", toolkitHandler.GetByID)
			}

			// Category routes - Admin only
			categories := protected.Group("/categories")
			categories.Use(authService.RequireAdmin())
			{
				categories.POST("", categoryHandler.Create)
				categories.GET("", categoryHandler.GetAll)
				categories.GET("/:id", categoryHandler.GetByID)
				categories.PUT("/:id", categoryHandler.Update)
				categories.DELETE("/:id", categoryHandler.Delete)
				categories.GET("/tree", categoryHandler.GetTree)
			}

			// Loan routes
			loans := protected.Group("/loans")
			{
				loans.POST("", loanHandler.Create)
				loans.GET("", loanHandler.GetAll)
				loans.GET("/:id", loanHandler.GetByID)
				loans.PUT("/:id", loanHandler.Update)
				loans.DELETE("/:id", loanHandler.Delete)
			}
		}
	}

	log.Printf("Starting server on port %s...", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
