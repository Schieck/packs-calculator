package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN string
}

type AuthConfig struct {
	JWTSecret   string
	AuthSecret  string
	TokenExpiry time.Duration
	Issuer      string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DB_DSN", "postgres://packer:secret@localhost:5432/packs?sslmode=disable"),
		},
		Auth: AuthConfig{
			JWTSecret:   getEnv("JWT_SECRET", "change-me"),
			AuthSecret:  getEnv("AUTH_SECRET", "default-auth-secret"),
			TokenExpiry: getEnvDuration("TOKEN_EXPIRY", "24h"),
			Issuer:      getEnv("ISSUER", "packs-calculator"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 24 * time.Hour // Fallback to 24 hours
}
