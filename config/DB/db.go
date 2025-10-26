package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, relying on environment variables")
	}
	fmt.Println(".env found")
}

func ConnectPostgres() (*gorm.DB, error) {
	// Support both naming conventions
	user := os.Getenv("DB_USER")
	if user == "" {
		user = os.Getenv("POSTGRES_USER")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = os.Getenv("POSTGRES_PASSWORD")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = os.Getenv("POSTGRES_DB")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kolkata",
		host, user, password, dbName, port, sslmode,
	)

	// Open a GORM DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to Postgres: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully!")
	return db, nil
}
