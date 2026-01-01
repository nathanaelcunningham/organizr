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
