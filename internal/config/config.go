package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SqlitePath string
	ServerPort string
	JWTSecret  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		SqlitePath: os.Getenv("SQLITE_PATH"),
		ServerPort: os.Getenv("SERVER_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	if cfg.SqlitePath == "" {
		return nil, fmt.Errorf("SQLITE_PATH is required")
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
