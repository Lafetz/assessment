package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const defaultPort = 8080

var logLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type Config struct {
	Port     int
	LogLevel slog.Level
	Env      string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	portStr := os.Getenv("PORT")
	port := defaultPort

	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		} else {
			fmt.Printf("Invalid PORT value '%s', defaulting to %d\n", portStr, defaultPort)
		}
	} else {
		fmt.Printf("PORT not set, defaulting to %d\n", defaultPort)
	}

	logLevelStr := os.Getenv("LOG_LEVEL")
	level, exists := logLevels[logLevelStr]
	if !exists {
		fmt.Printf("Invalid LOG_LEVEL '%s', defaulting to 'info'\n", logLevelStr)
		level = logLevels["info"]
	}

	env := os.Getenv("ENV")
	if env != "development" && env != "production" {
		fmt.Printf("Invalid ENV '%s', defaulting to 'development'\n", env)
		env = "development"
	} else if env == "" {
		fmt.Printf("ENV not set, defaulting to 'development'\n")
		env = "development"
	}

	return &Config{
		Port:     port,
		LogLevel: level,
		Env:      env,
	}
}
