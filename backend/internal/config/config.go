package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AI       AIConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	// SQLite 兼容模式（开发用）
	UseSQLite bool
	Path      string
}

type AIConfig struct {
	APIKey string
}

func Load() *Config {
	useSQLite := getEnv("USE_SQLITE", "false") == "true"
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:      getEnv("DB_HOST", "localhost"),
			Port:      getEnv("DB_PORT", "3306"),
			User:      getEnv("DB_USER", "root"),
			Password:  getEnv("DB_PASSWORD", ""),
			Database:  getEnv("DB_NAME", "financial_assistant"),
			UseSQLite: useSQLite,
			Path:      getEnv("DATABASE_PATH", "./database.db"),
		},
		AI: AIConfig{
			APIKey: getEnv("AI_API_KEY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
