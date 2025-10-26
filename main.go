package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	DB "goExpenseTracker/config/DB"
	swagger "goExpenseTracker/config/swagger"
	Logger "goExpenseTracker/internal/middlewear"
)

func main() {
	setupDB()

	router := gin.New()
	router.Use(Logger.Logger())
	router.Use(gin.Recovery())
	swagger.SetupSwagger(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)

}

func setupDB() {
	db, err := DB.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
}
