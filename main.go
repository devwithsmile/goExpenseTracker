package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	DB "goExpenseTracker/config/DB"
	"goExpenseTracker/internal/handlers"
	Logger "goExpenseTracker/internal/middlewears"
	"goExpenseTracker/internal/models"
	"goExpenseTracker/internal/repositories"
	"goExpenseTracker/internal/routes"
	"goExpenseTracker/internal/services"

	_ "goExpenseTracker/docs" // Swagger docs

	"gorm.io/gorm"
)

// @title Go Expense Tracker API
// @version 1.0
// @description RESTful API for managing expenses and categories
// @termsOfService http://example.com/terms/

// @contact.name Developer Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	app := bootstrapApp()
	port := getPort()
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// bootstrapApp sets up the full app (DB + Dependencies + Routes)
func bootstrapApp() *gin.Engine {
	db, err := DB.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database tables
	if err := db.AutoMigrate(&models.Category{}, &models.Expense{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database tables migrated successfully!")

	// Initialize repositories, services, handlers
	categoryHandler, expenseHandler := initializeDependencies(db)

	// Create Gin router and attach middleware
	router := gin.New()
	router.Use(Logger.Logger(), gin.Recovery())

	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Mount routes
	setupRoutes(router, categoryHandler, expenseHandler)

	return router
}

// initializeDependencies wires repositories → services → handlers
func initializeDependencies(db *gorm.DB) (*handlers.CategoryHandler, *handlers.ExpenseHandler) {
	// Category dependencies
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Expense dependencies (with category repo for relationship mapping)
	expenseRepo := repositories.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo, categoryRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	return categoryHandler, expenseHandler
}

// setupRoutes registers all route groups
func setupRoutes(router *gin.Engine, categoryHandler *handlers.CategoryHandler, expenseHandler *handlers.ExpenseHandler) {
	api := router.Group("/api")
	{
		routes.SetupCategoryRoutes(api, categoryHandler)
		routes.SetupExpenseRoutes(api, expenseHandler)
	}
}

// getPort retrieves the server port or defaults to 8080
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
