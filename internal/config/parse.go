package config

import (
	"log"
	"os"
	"strconv"

	"github.com/zonder12120/brandscout-quotebook/pkg/env"
)

const (
	envFilePath        = "config/.env"
	defaultQuotesLimit = 1000
)

func MustLoad() *App {
	cfg, err := parseConfig(envFilePath)
	if err != nil {
		log.Printf("WARNING: config error: %v", err)
	}
	return cfg
}

func parseConfig(filePath string) (*App, error) {
	_ = env.LoadEnv(filePath)

	port := os.Getenv("PORT")
	logLevel := os.Getenv("LOG_LEVEL")

	if port == "" {
		port = "8080"
	}
	if logLevel == "" {
		logLevel = "info"
	}

	quotesLimit := defaultQuotesLimit
	if envLimit := os.Getenv("QUOTES_LIMIT"); envLimit != "" {
		if v, err := strconv.Atoi(envLimit); err == nil {
			quotesLimit = v
		}
	}

	return &App{
		Port:        port,
		QuotesLimit: quotesLimit,
		LogLevel:    logLevel,
	}, nil
}
