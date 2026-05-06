package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	ServerPort     string
	GinMode        string
	TrustedProxies []string
}

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		GinMode:    os.Getenv("GIN_MODE"),
	}
	if cfg.GinMode == "" {
		cfg.GinMode = "release"
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	return cfg, nil
}

func loadDotEnv() {
	if exePath, err := os.Executable(); err == nil {
		envPath := filepath.Join(filepath.Dir(exePath), ".env")
		if err := godotenv.Load(envPath); err == nil {
			return
		}
	}
	_ = godotenv.Load(".env")
}

func parseTrustedProxies(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	proxies := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			proxies = append(proxies, p)
		}
	}
	return proxies
}

func getenvDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
