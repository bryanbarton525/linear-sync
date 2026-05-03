package config

import (
	"errors"
	"os"
)

// Config holds the application configuration
type Config struct {
	APIKey      string
	TeamID      string
	DatabaseURL string
}

// Load reads configuration from environment variables and validates them
func Load() (*Config, error) {
	apiKey := os.Getenv("LINEAR_API_KEY")
	if apiKey == "" {
		return nil, errors.New("LINEAR_API_KEY is required")
	}

	teamID := os.Getenv("LINEAR_TEAM_ID")
	if teamID == "" {
		return nil, errors.New("LINEAR_TEAM_ID is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return &Config{
		APIKey:      apiKey,
		TeamID:      teamID,
		DatabaseURL: databaseURL,
	}, nil
}
