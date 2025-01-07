package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
	"log"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	connStr := "postgresql://test_db_mvaz_user:NgcEd82NG6iHSgqfwhkhPukcnsBHC0c4@dpg-cttsnalumphs73ei09c0-a.oregon-postgres.render.com/test_db_mvaz"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//DB = db

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	return db, nil
}
