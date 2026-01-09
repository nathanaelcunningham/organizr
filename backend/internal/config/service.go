package config

import (
	"context"
	"fmt"
	"os"

	"github.com/nathanael/organizr/internal/persistence"
)

type Service struct {
	repo persistence.ConfigRepository
}

func NewService(repo persistence.ConfigRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, key string) (string, error) {
	// Check environment variable first
	if envKey := getEnvKey(key); envKey != "" {
		if envVal := os.Getenv(envKey); envVal != "" {
			return envVal, nil
		}
	}

	// Fall back to database
	value, err := s.repo.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get config %s: %w", key, err)
	}
	return value, nil
}

func (s *Service) GetAll(ctx context.Context) (map[string]string, error) {
	// Get database configs
	configs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all configs: %w", err)
	}

	// Override with environment variables
	for dbKey, envKey := range envKeyMap {
		if envVal := os.Getenv(envKey); envVal != "" {
			configs[dbKey] = envVal
		}
	}

	return configs, nil
}

func (s *Service) Set(ctx context.Context, key, value string) error {
	if err := s.repo.Set(ctx, key, value); err != nil {
		return fmt.Errorf("failed to set config %s: %w", key, err)
	}
	return nil
}
