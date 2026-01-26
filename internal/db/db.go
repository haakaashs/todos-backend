package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

// Prerequisite ensures database and tables exist
func prerequisite(db *sql.DB) {
	err := ensureDatabase(db)
	if err != nil {
		log.Default().Fatalf("Failed to ensure database exists: %v", err)
	}

	err = ensureTables(db)
	if err != nil {
		log.Default().Fatalf("Failed to ensure tables exist: %v", err)
	}

	log.Println("Database and tables ready")
}

// getDSN constructs the DSN for connecting to the default "postgres" database
func getDSN(config *configs.Config) string {
	// Connect to default "postgres" database to create the target DB
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.DBName,
		config.DB.SSLMode,
	)
}

func InitializeDB() *sql.DB {
	config = configs.LoadConfig()

	// Step 1: connect to default DB (postgres)
	adminDB, err := sql.Open(config.DB.Provider, getDSN(config))
	if err != nil {
		log.Fatal(err)
	}

	// wait until the DB is ready
	for range maxRetries {
		if err := adminDB.Ping(); err == nil {
			log.Println("Connected to Postgres for DB creation")
			break
		}
		log.Printf("Waiting for Postgres to be ready... (%v)", err)
		time.Sleep(retryInterval)
	}

	// Step 2: ensure target database exists
	config.DB.DBName = os.Getenv("DB_NAME")
	prerequisite(adminDB)
	defer adminDB.Close()

	// Step 3: now connect to the actual database
	db, err := sql.Open(config.DB.Provider, getDSN(config))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to target database successfully")
	ensureTables(db)
	return db
}
