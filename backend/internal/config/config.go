package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AI       AIConfig
	Redis    RedisConfig
	JWT      JWTConfig
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
}

type AIConfig struct {
	BaseURL        string
	APIKey         string
	Model          string
	TimeoutSeconds int
}

type RedisConfig struct {
	Addr string
}

type JWTConfig struct {
	Secret string
}

func LoadEnvFileIfPresent(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return godotenv.Load(path)
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "abc123"),
			Database: getEnv("DB_NAME", "financial_assistant"),
		},
		AI: AIConfig{
			BaseURL:        getEnv("AI_BASE_URL", "https://tianjiajie.xin"),
			APIKey:         getEnv("AI_API_KEY", ""),
			Model:          getEnv("AI_MODEL", "claude-3-5-sonnet-20241022"),
			TimeoutSeconds: getEnvAsInt("AI_TIMEOUT_SECONDS", 30),
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return defaultValue
}
