package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	configs "github.com/haakaashs/todos-backend/internal/configs"
	_ "github.com/lib/pq"
)

const (
	maxRetries    = 10
	retryInterval = 3 * time.Second
)

var config *configs.Config

// ensureDatabase checks if the target database exists, and creates it if not
func ensureDatabase(db *sql.DB) error {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", config.DB.DBName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if !exists {
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DB.DBName))
		if err != nil {
			return fmt.Errorf("error creating database %s: %w", config.DB.DBName, err)
		}
		log.Printf("Database %s created", config.DB.DBName)
	} else {
		log.Printf("Database %s already exists", config.DB.DBName)
	}

	return nil
}

// ensureTables creates tables if they do not exist
func ensureTables(db *sql.DB) error {
	tableSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := db.Exec(tableSQL)
	if err != nil {
		return fmt.Errorf("error creating table todos: %w", err)
	}

	log.Println("Tables ensured")
	return nil
}

// Prerequisite connects to Postgres, ensures database and tables exist, with retries
func prerequisite(db *sql.DB) {
	for i := 1; i <= maxRetries; i++ {
		err := ensureDatabase(db)
		if err != nil {
			log.Printf("Attempt %d/%d: Waiting for Postgres to be ready... (%v)", i, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		err = ensureTables(db)
		if err != nil {
			log.Printf("Attempt %d/%d: Waiting for tables to be creatable... (%v)", i, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		log.Println("Database and tables ready")
		return
	}

	log.Fatal("Failed to ensure database and tables after maximum retries")
}

// getDSN constructs the Data Source Name for connecting to Database
func getDSN(config *configs.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.SSLMode,
	)
}

// InitializeDB initializes the database connection
func InitializeDB() *sql.DB {

	config = configs.LoadConfig()
	db, err := sql.Open(config.DB.Provider, getDSN(config))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Database successfully")
	prerequisite(db)
	return db
}
