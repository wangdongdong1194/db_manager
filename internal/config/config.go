package config

import (
	"errors"
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
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getenvDefault("DB_PORT", "3306"),
		DBUser:     os.Getenv("DB_USER"),
		DBPass:     os.Getenv("DB_PASS"),
		ServerPort: os.Getenv("SERVER_PORT"),
		GinMode:    os.Getenv("GIN_MODE"),
	}
	if cfg.GinMode == "" {
		cfg.GinMode = "release"
	}

	cfg.TrustedProxies = parseTrustedProxies(os.Getenv("TRUSTED_PROXIES"))
	if len(cfg.TrustedProxies) == 0 {
		cfg.TrustedProxies = []string{"127.0.0.1", "::1"}
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPass == "" || cfg.ServerPort == "" {
		return Config{}, errors.New("missing required env: DB_HOST, DB_USER, DB_PASS or SERVER_PORT")
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
