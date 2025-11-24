package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {   
		fmt.Println("Warning : .env file not found, using system environment variables")
	}

	// Parse JWT expiry
	accessExpiry, _ := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	refreshExpiry, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "3000"),
			Env: getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName: getEnv("DB_NAME", "patungan_db"),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-this-secret"),
			AccessExpiry: accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}