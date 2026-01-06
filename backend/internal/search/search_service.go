package search

import (
	"context"
	"fmt"

	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/persistence"
	"github.com/nathanael/organizr/internal/search/providers"
)

// MAMService handles torrent search and download operations for MyAnonamouse
type MAMService struct {
	configRepo persistence.ConfigRepository
	provider   *providers.MyAnonamouseProvider
}

// NewMAMService creates a new MAM service
func NewMAMService(configRepo persistence.ConfigRepository) *MAMService {
	return &MAMService{
		configRepo: configRepo,
	}
}

// initProvider lazy-loads the MAM configuration
func (s *MAMService) initProvider(ctx context.Context) error {
	if s.provider != nil {
		return nil
	}

	baseURL, err := s.configRepo.Get(ctx, "mam.baseurl")
	if err != nil {
		return fmt.Errorf("MAM base URL not configured: %w", err)
	}

	secret, err := s.configRepo.Get(ctx, "mam.secret")
	if err != nil {
		return fmt.Errorf("MAM secret not configured: %w", err)
	}

	s.provider = providers.NewMyAnonamouseProvider(baseURL, secret)
	return nil
}

// Search performs a torrent search on MyAnonamouse
func (s *MAMService) Search(ctx context.Context, query string) ([]*models.SearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if err := s.initProvider(ctx); err != nil {
		return nil, err
	}

	return s.provider.Search(ctx, query)
}

// TestConnection validates MyAnonamouse credentials
func (s *MAMService) TestConnection(ctx context.Context) error {
	if err := s.initProvider(ctx); err != nil {
		return err
	}

	return s.provider.TestConnection(ctx)
}

// DownloadTorrent fetches the torrent file bytes for a given torrent ID
func (s *MAMService) DownloadTorrent(ctx context.Context, torrentID int) ([]byte, error) {
	if err := s.initProvider(ctx); err != nil {
		return nil, err
	}

	return s.provider.DownloadTorrent(ctx, torrentID)
}
