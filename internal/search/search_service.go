package search

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/persistence"
)

// SearchService handles both search operations and provider management
type SearchService struct {
	repo      persistence.ProviderRepository
	registry  *Registry
	providers []Provider
	mu        sync.RWMutex
}

// NewSearchService creates a new unified search service
func NewSearchService(repo persistence.ProviderRepository, registry *Registry) *SearchService {
	s := &SearchService{
		repo:     repo,
		registry: registry,
	}

	// Load providers at initialization
	if err := s.reloadProviders(context.Background()); err != nil {
		log.Printf("Warning: Failed to load providers during initialization: %v", err)
	}

	return s
}

// reloadProviders loads enabled providers from database (internal method)
func (s *SearchService) reloadProviders(ctx context.Context) error {
	configs, err := s.repo.ListEnabled(ctx)
	if err != nil {
		return fmt.Errorf("failed to list enabled providers: %w", err)
	}

	providers := make([]Provider, 0, len(configs))
	for _, config := range configs {
		provider, err := s.registry.Create(config.ProviderType, config.ConfigJSON)
		if err != nil {
			log.Printf("Warning: Failed to create provider %s: %v", config.ProviderType, err)
			continue
		}
		providers = append(providers, provider)
	}

	s.mu.Lock()
	s.providers = providers
	s.mu.Unlock()

	log.Printf("Loaded %d search providers", len(providers))
	return nil
}

// ============================================================================
// Search Operations
// ============================================================================

// Search performs a search across providers.
// If providerName is specified, only that provider is queried.
// Otherwise, all enabled providers are queried and results are aggregated.
func (s *SearchService) Search(ctx context.Context, query string, providerName string) ([]*models.SearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	s.mu.RLock()
	providers := s.providers
	s.mu.RUnlock()

	// Search specific provider
	if providerName != "" {
		for _, p := range providers {
			if p.Name() == providerName {
				return p.Search(ctx, query)
			}
		}
		return nil, fmt.Errorf("provider not found: %s", providerName)
	}

	// Search all providers and aggregate results
	var allResults []*models.SearchResult
	for _, p := range providers {
		results, err := p.Search(ctx, query)
		if err != nil {
			log.Printf("Provider %s error: %v", p.Name(), err)
			continue
		}
		allResults = append(allResults, results...)
	}

	if len(allResults) == 0 && len(providers) > 0 {
		return nil, fmt.Errorf("all providers failed to return results")
	}

	return allResults, nil
}

// ListActiveProviders returns the names of all currently loaded providers
func (s *SearchService) ListActiveProviders() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, len(s.providers))
	for i, p := range s.providers {
		names[i] = p.Name()
	}
	return names
}

// ============================================================================
// Provider Management (CRUD)
// ============================================================================

// GetAvailableTypes returns all provider types that can be configured
func (s *SearchService) GetAvailableTypes(ctx context.Context) ([]*models.ProviderType, error) {
	return s.registry.GetTypes(), nil
}

// GetType returns details about a specific provider type
func (s *SearchService) GetType(ctx context.Context, providerType string) (*models.ProviderType, error) {
	return s.registry.GetType(providerType)
}

// CreateProvider adds a new provider configuration and reloads providers
func (s *SearchService) CreateProvider(ctx context.Context, config *models.ProviderConfig) error {
	// Validate provider type exists in registry
	if _, err := s.registry.GetType(config.ProviderType); err != nil {
		return fmt.Errorf("invalid provider type: %w", err)
	}

	// Validate configuration
	if err := s.registry.ValidateConfig(config.ProviderType, config.ConfigJSON); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Check if already exists (one instance per type)
	exists, err := s.repo.Exists(ctx, config.ProviderType)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return fmt.Errorf("provider already configured: %s", config.ProviderType)
	}

	// Create in database
	if err := s.repo.Create(ctx, config); err != nil {
		return err
	}

	// Reload providers to pick up the new one
	return s.reloadProviders(ctx)
}

// GetProvider retrieves a provider configuration
func (s *SearchService) GetProvider(ctx context.Context, providerType string) (*models.ProviderConfig, error) {
	return s.repo.GetByType(ctx, providerType)
}

// ListProviders returns all configured providers (both enabled and disabled)
func (s *SearchService) ListProviders(ctx context.Context) ([]*models.ProviderConfig, error) {
	return s.repo.List(ctx)
}

// UpdateProvider modifies a provider configuration and reloads providers
func (s *SearchService) UpdateProvider(ctx context.Context, config *models.ProviderConfig) error {
	// Validate provider type exists in registry
	if _, err := s.registry.GetType(config.ProviderType); err != nil {
		return fmt.Errorf("invalid provider type: %w", err)
	}

	// Validate configuration
	if err := s.registry.ValidateConfig(config.ProviderType, config.ConfigJSON); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Update in database
	if err := s.repo.Update(ctx, config); err != nil {
		return err
	}

	// Reload providers to pick up changes
	return s.reloadProviders(ctx)
}

// DeleteProvider removes a provider configuration and reloads providers
func (s *SearchService) DeleteProvider(ctx context.Context, providerType string) error {
	if err := s.repo.Delete(ctx, providerType); err != nil {
		return err
	}

	// Reload providers to remove the deleted one
	return s.reloadProviders(ctx)
}

// ToggleProvider enables or disables a provider and reloads providers
func (s *SearchService) ToggleProvider(ctx context.Context, providerType string, enabled bool) error {
	if err := s.repo.UpdateEnabled(ctx, providerType, enabled); err != nil {
		return err
	}

	// Reload providers to reflect the change
	return s.reloadProviders(ctx)
}

// TestProvider validates provider credentials
func (s *SearchService) TestProvider(ctx context.Context, providerType string) error {
	// Get config from database
	config, err := s.repo.GetByType(ctx, providerType)
	if err != nil {
		return fmt.Errorf("failed to get provider config: %w", err)
	}

	// Create provider instance
	provider, err := s.registry.Create(config.ProviderType, config.ConfigJSON)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Test connection
	if err := provider.TestConnection(ctx); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}
