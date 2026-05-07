package config

import (
	"errors"
	"os"
)

type Config struct {
	APIKey      string
	TeamID      string
	DatabaseURL string
}

func Load() (*Config, error) {
	apiKey := os.Getenv("LINEAR_API_KEY")
	if apiKey == "" {
		return nil, errors.New("LINEAR_API_KEY is not set")
	}
	teamID := os.Getenv("LINEAR_TEAM_ID")
	if teamID == "" {
		return nil, errors.New("LINEAR_TEAM_ID is not set")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	return &Config{
		APIKey:      apiKey,
		TeamID:      teamID,
		DatabaseURL: databaseURL,
	}, nil
}
