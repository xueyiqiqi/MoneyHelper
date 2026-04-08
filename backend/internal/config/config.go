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
	Path string
}

type AIConfig struct {
	APIKey string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./database.db"),
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