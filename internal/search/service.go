package search

import (
	"context"
	"fmt"
	"log"

	"github.com/nathanael/organizr/internal/models"
)

// Service coordinates search operations across multiple providers
type Service struct {
	providers []Provider
}

// NewService creates a new search service with the given providers
func NewService(providers []Provider) *Service {
	return &Service{
		providers: providers,
	}
}

// Search performs a search across providers.
// If providerName is specified, only that provider is queried.
// Otherwise, all providers are queried and results are aggregated.
func (s *Service) Search(ctx context.Context, query string, providerName string) ([]*models.SearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Search specific provider
	if providerName != "" {
		for _, p := range s.providers {
			if p.Name() == providerName {
				return p.Search(ctx, query)
			}
		}
		return nil, fmt.Errorf("provider not found: %s", providerName)
	}

	// Search all providers and aggregate results
	var allResults []*models.SearchResult
	for _, p := range s.providers {
		results, err := p.Search(ctx, query)
		if err != nil {
			log.Printf("Provider %s error: %v", p.Name(), err)
			continue
		}
		allResults = append(allResults, results...)
	}

	if len(allResults) == 0 && len(s.providers) > 0 {
		return nil, fmt.Errorf("all providers failed to return results")
	}

	return allResults, nil
}

// ListProviders returns the names of all registered providers
func (s *Service) ListProviders() []string {
	names := make([]string, len(s.providers))
	for i, p := range s.providers {
		names[i] = p.Name()
	}
	return names
}
