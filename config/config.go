package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	portInt, _ := strconv.Atoi(port)

	databaseURL := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/peeingdog?sslmode=disable")

	return &Config{
		Port:        portInt,
		DatabaseURL: databaseURL,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
