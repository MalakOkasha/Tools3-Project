package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Postgres driver
)

var DB *sql.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	// Get database connection details from environment variables
	DBUser := os.Getenv("DB_USER")         // Get DB user from environment variable
	DBPassword := os.Getenv("DB_PASSWORD") // Get DB password from environment variable
	DBName := os.Getenv("DB_NAME")         // Get DB name from environment variable
	DBHost := os.Getenv("DB_HOST")         // Get DB host from environment variable
	DBPort := os.Getenv("DB_PORT")         // Get DB port from environment variable
	SSLMode := "disable"                   // SSL mode

	if DBUser == "" || DBPassword == "" || DBName == "" || DBHost == "" || DBPort == "" {
		log.Fatal("Missing required environment variables for database connection")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		DBUser, DBPassword, DBName, DBHost, DBPort, SSLMode)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	fmt.Println("Successfully connected to the database")
}
