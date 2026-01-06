package server

import (
	"time"

	"github.com/nathanael/organizr/internal/models"
)

type downloadDTO struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	Series        string     `json:"series,omitempty"`
	Category      string     `json:"category,omitempty"`
	Status        string     `json:"status"`
	Progress      float64    `json:"progress"`
	OrganizedPath string     `json:"organized_path,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	OrganizedAt   *time.Time `json:"organized_at,omitempty"`
}

func toDTO(d *models.Download) downloadDTO {
	return downloadDTO{
		ID:            d.ID,
		Title:         d.Title,
		Author:        d.Author,
		Series:        d.Series,
		Category:      d.Category,
		Status:        string(d.Status),
		Progress:      d.Progress,
		OrganizedPath: d.OrganizedPath,
		ErrorMessage:  d.ErrorMessage,
		CreatedAt:     d.CreatedAt,
		CompletedAt:   d.CompletedAt,
		OrganizedAt:   d.OrganizedAt,
	}
}

func toDTOList(downloads []*models.Download) []downloadDTO {
	dtos := make([]downloadDTO, len(downloads))
	for i, d := range downloads {
		dtos[i] = toDTO(d)
	}
	return dtos
}

type searchResultDTO struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	TorrentURL string `json:"torrent_url,omitempty"`
	MagnetLink string `json:"magnet_link,omitempty"`
	Size       string `json:"size"`
	Seeders    int    `json:"seeders"`
	Provider   string `json:"provider"`
}

func searchResultToDTO(s *models.SearchResult) searchResultDTO {
	return searchResultDTO{
		ID:         s.ID,
		Title:      s.Title,
		Author:     s.Author,
		TorrentURL: s.TorrentURL,
		MagnetLink: s.MagnetLink,
		Size:       s.Size,
		Seeders:    s.Seeders,
		Provider:   s.Provider,
	}
}

func searchResultsToDTOList(results []*models.SearchResult) []searchResultDTO {
	dtos := make([]searchResultDTO, len(results))
	for i, r := range results {
		dtos[i] = searchResultToDTO(r)
	}
	return dtos
}
