package configs

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	t.Setenv("DB_PROVIDER", "postgres")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpassword")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_SSLMODE", "disable")

	config := LoadConfig()

	if config.DB.Provider != "postgres" {
		t.Errorf("Expected provider to be 'postgres', got '%s'", config.DB.Provider)
	}
	if config.DB.Host != "localhost" {
		t.Errorf("Expected host to be 'localhost', got '%s'", config.DB.Host)
	}
	if config.DB.Port != 5432 {
		t.Errorf("Expected port to be 5432, got %d", config.DB.Port)
	}
	if config.DB.User != "testuser" {
		t.Errorf("Expected user to be 'testuser', got '%s'", config.DB.User)
	}
	if config.DB.Password != "testpassword" {
		t.Errorf("Expected password to be 'testpassword', got '%s'", config.DB.Password)
	}
	if config.DB.DBName != "postgres" {
		t.Errorf("Expected dbname to be 'postgres', got '%s'", config.DB.DBName)
	}
	if config.DB.SSLMode != "disable" {
		t.Errorf("Expected sslmode to be 'disable', got '%s'", config.DB.SSLMode)
	}
}
