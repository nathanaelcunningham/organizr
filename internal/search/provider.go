package search

import (
	"context"

	"github.com/nathanael/organizr/internal/models"
)

// Provider defines the interface for torrent search providers.
// Each provider implements its own API-based search logic.
type Provider interface {
	// Name returns the provider's unique name
	Name() string

	// Search queries the provider's API for audiobook torrents
	Search(ctx context.Context, query string) ([]*models.SearchResult, error)
}
