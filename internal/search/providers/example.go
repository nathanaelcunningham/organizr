package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/nathanael/organizr/internal/models"
)

// AudiobookBayProvider is an example API-based provider implementation.
// This demonstrates how to implement the Provider interface for a torrent site
// that provides an API (not HTML scraping).
type AudiobookBayProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewAudiobookBayProvider creates a new AudiobookBay provider instance
func NewAudiobookBayProvider(baseURL, apiKey string) *AudiobookBayProvider {
	return &AudiobookBayProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// Name returns the provider's name
func (p *AudiobookBayProvider) Name() string {
	return "AudiobookBay"
}

// Search queries the AudiobookBay API for audiobook torrents
func (p *AudiobookBayProvider) Search(ctx context.Context, query string) ([]*models.SearchResult, error) {
	// Build API request
	searchURL := fmt.Sprintf("%s/api/search?q=%s&key=%s", p.baseURL, url.QueryEscape(query), p.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse API response
	var apiResp struct {
		Results []struct {
			Title      string `json:"title"`
			Author     string `json:"author"`
			TorrentURL string `json:"torrent_url"`
			MagnetLink string `json:"magnet"`
			Size       string `json:"size"`
			Seeders    int    `json:"seeders"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert API results to domain models
	results := make([]*models.SearchResult, len(apiResp.Results))
	for i, r := range apiResp.Results {
		results[i] = &models.SearchResult{
			Title:      r.Title,
			Author:     r.Author,
			TorrentURL: r.TorrentURL,
			MagnetLink: r.MagnetLink,
			Size:       r.Size,
			Seeders:    r.Seeders,
			Provider:   p.Name(),
		}
	}

	return results, nil
}
