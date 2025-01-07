package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
	"log"
	"os"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	// Подключение к базе данных
	connStr := os.Getenv("POSTGRES")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	fmt.Println("Database connection:", db)
	return db, nil
}
