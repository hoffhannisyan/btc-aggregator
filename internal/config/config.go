package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	PollInterval   time.Duration
	RequestTimeout time.Duration
	MaxRetries     int
	HTTPPort       int
}

func Load() (*Config, error) {
	cfg := &Config{
		PollInterval:   10 * time.Second,
		RequestTimeout: 5 * time.Second,
		MaxRetries:     3,
		HTTPPort:       8080,
	}

	if v := os.Getenv("POLL_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid POLL_INTERVAL: %w", err)
		}
		cfg.PollInterval = d
	}

	if v := os.Getenv("REQUEST_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid REQUEST_TIMEOUT: %w", err)
		}
		cfg.RequestTimeout = d
	}

	if v := os.Getenv("MAX_RETRIES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_RETRIES: %w", err)
		}
		cfg.MaxRetries = n
	}

	if v := os.Getenv("HTTP_PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
		}
		cfg.HTTPPort = n
	}

	return cfg, nil
}