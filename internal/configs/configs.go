package configs

import (
	"os"
	"strconv"
)

const ConfigFilePath = "configs/config.json"

// DBConfig holds the database configuration
type DBConfig struct {
	Provider string `json:"provider"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// Config holds the entire config structure
type Config struct {
	DB DBConfig `json:"db"`
}

// LoadConfig loads the configuration from config.json file
func LoadConfig() *Config {
	// // Open config.json
	// file, err := os.Open(ConfigFilePath)
	// if err != nil {
	// 	log.Fatalf("Failed to open config file: %v", err)
	// }
	// defer file.Close()

	// // Read all bytes
	// bytes, err := io.ReadAll(file)
	// if err != nil {
	// 	log.Fatalf("Failed to read config file: %v", err)
	// }

	// // Unmarshal JSON into Config struct
	// var config Config
	// if err := json.Unmarshal(bytes, &config); err != nil {
	// 	log.Fatalf("Failed to parse config file: %v", err)
	// }
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	return &Config{
		DB: DBConfig{
			Provider: os.Getenv("DB_PROVIDER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			DBName:   os.Getenv("DB_PROVIDER"),
		},
	}
}
