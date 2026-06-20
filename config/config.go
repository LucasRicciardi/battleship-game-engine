package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string
	Expiration    int
	RefreshExp    int
	Realm         string
	TokenLookup   string
	AuthHeaderPrefix string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 15),
			IdleTimeout:  getEnvInt("SERVER_IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			Host:     getEnvString("DATABASE_HOST", "localhost"),
			Port:     getEnvInt("DATABASE_PORT", 5432),
			User:     getEnvString("DATABASE_USER", "battleship"),
			Password: getEnvString("DATABASE_PASSWORD", "battleship"),
			Name:     getEnvString("DATABASE_NAME", "battleship"),
			SSLMode:  getEnvString("DATABASE_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:             getEnvString("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration:         getEnvInt("JWT_EXPIRATION", 86400), // 24 hours
			RefreshExp:         getEnvInt("JWT_REFRESH_EXPIRATION", 86400),
			Realm:              getEnvString("JWT_REALM", "battleship-game-engine"),
			TokenLookup:        getEnvString("JWT_TOKEN_LOOKUP", "header: Authorization: Bearer"),
			AuthHeaderPrefix:   getEnvString("JWT_AUTH_HEADER_PREFIX", "Bearer"),
		},
		Redis: RedisConfig{
			Host:     getEnvString("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnvString("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
	}
}

// getEnvString returns environment variable value or default
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns environment variable value as int or default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defaultValue
}
