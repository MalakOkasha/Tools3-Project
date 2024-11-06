package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Postgres driver
)

const (
	DBUser     = "postgres"  // The actual PostgreSQL username
	DBPassword = "hijazy_1"  // The actual PostgreSQL password
	DBName     = "PTS"       // The actual PostgreSQL database name
	DBHost     = "localhost" // The database host (e.g., localhost)
	DBPort     = "5432"      // The PostgreSQL port (default is 5432)
	SSLMode    = "disable"   // Change to "require" if SSL is enabled
)

var DB *sql.DB

// ConnectDB initializes the database connection
func ConnectDB() {
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
