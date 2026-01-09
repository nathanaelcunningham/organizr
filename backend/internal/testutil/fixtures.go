package testutil

import (
	"maps"
	"time"

	"github.com/google/uuid"
	"github.com/nathanael/organizr/internal/models"
)

// DownloadOption is a functional option for configuring test Download fixtures.
type DownloadOption func(*models.Download)

// NewTestDownload creates a Download with sensible defaults and optional overrides.
func NewTestDownload(opts ...DownloadOption) *models.Download {
	dl := &models.Download{
		ID:        uuid.New().String(),
		Title:     "Test Download",
		Author:    "Test Author",
		Status:    models.StatusQueued,
		Progress:  0,
		CreatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(dl)
	}

	return dl
}

// WithTitle sets the download title.
func WithTitle(title string) DownloadOption {
	return func(d *models.Download) {
		d.Title = title
	}
}

// WithAuthor sets the download author.
func WithAuthor(author string) DownloadOption {
	return func(d *models.Download) {
		d.Author = author
	}
}

// WithSeries sets the download series.
func WithSeries(series string) DownloadOption {
	return func(d *models.Download) {
		d.Series = series
	}
}

// WithSeriesNumber sets the download series number.
func WithSeriesNumber(seriesNumber string) DownloadOption {
	return func(d *models.Download) {
		d.SeriesNumber = seriesNumber
	}
}

// WithStatus sets the download status.
func WithStatus(status models.DownloadStatus) DownloadOption {
	return func(d *models.Download) {
		d.Status = status
	}
}

// WithProgress sets the download progress.
func WithProgress(progress float64) DownloadOption {
	return func(d *models.Download) {
		d.Progress = progress
	}
}

// WithID sets the download ID.
func WithID(id string) DownloadOption {
	return func(d *models.Download) {
		d.ID = id
	}
}

// WithMagnetLink sets the download magnet link.
func WithMagnetLink(magnetLink string) DownloadOption {
	return func(d *models.Download) {
		d.MagnetLink = magnetLink
	}
}

// WithTorrentURL sets the download torrent URL.
func WithTorrentURL(torrentURL string) DownloadOption {
	return func(d *models.Download) {
		d.TorrentURL = torrentURL
	}
}

// WithCategory sets the download category.
func WithCategory(category string) DownloadOption {
	return func(d *models.Download) {
		d.Category = category
	}
}

// WithQBitHash sets the qBittorrent hash.
func WithQBitHash(hash string) DownloadOption {
	return func(d *models.Download) {
		d.QBitHash = hash
	}
}

// WithOrganizedPath sets the organized path.
func WithOrganizedPath(path string) DownloadOption {
	return func(d *models.Download) {
		d.OrganizedPath = path
	}
}

// SearchResultOption is a functional option for configuring test SearchResult fixtures.
type SearchResultOption func(*models.SearchResult)

// NewTestSearchResult creates a SearchResult with sensible defaults and optional overrides.
func NewTestSearchResult(opts ...SearchResultOption) *models.SearchResult {
	sr := &models.SearchResult{
		ID:         uuid.New().String(),
		Title:      "Test Book",
		Author:     "Test Author",
		TorrentURL: "https://example.com/torrent.torrent",
		MagnetLink: "magnet:?xt=urn:btih:test",
		Provider:   "MyAnonamouse",
		Category:   "Audiobooks",
		FileType:   "M4B",
		Language:   "English",
		Size:       "100 MB",
		Seeders:    10,
		Leechers:   2,
	}

	for _, opt := range opts {
		opt(sr)
	}

	return sr
}

// WithResultID sets the search result ID.
func WithResultID(id string) SearchResultOption {
	return func(sr *models.SearchResult) {
		sr.ID = id
	}
}

// WithResultTitle sets the search result title.
func WithResultTitle(title string) SearchResultOption {
	return func(sr *models.SearchResult) {
		sr.Title = title
	}
}

// WithResultAuthor sets the search result author.
func WithResultAuthor(author string) SearchResultOption {
	return func(sr *models.SearchResult) {
		sr.Author = author
	}
}

// WithResultSeries sets the search result series information.
func WithResultSeries(series []models.SeriesInfo) SearchResultOption {
	return func(sr *models.SearchResult) {
		sr.Series = series
	}
}

// WithResultProvider sets the search result provider.
func WithResultProvider(provider string) SearchResultOption {
	return func(sr *models.SearchResult) {
		sr.Provider = provider
	}
}

// WithSeeders sets the number of seeders.
func WithSeeders(seeders int) SearchResultOption {
	return func(sr *models.SearchResult) {
		sr.Seeders = seeders
	}
}

// WithFreeleech sets the freeleech flag.
func WithFreeleech(freeleech bool) SearchResultOption {
	return func(sr *models.SearchResult) {
		sr.Freeleech = freeleech
	}
}

// NewTestConfig creates a config map with sensible defaults and optional overrides.
func NewTestConfig(overrides map[string]string) map[string]string {
	config := map[string]string{
		"qbittorrent.url":      "http://localhost:8080",
		"qbittorrent.username": "admin",
		"qbittorrent.password": "adminpass",
		"mam.baseurl":          "https://www.myanonamouse.net",
		"mam.secret":           "test-secret-key",
		"output.path":          "/tmp/audiobooks",
		"output.template":      "{author}/{series}/{title}",
	}

	maps.Copy(config, overrides)

	return config
}
