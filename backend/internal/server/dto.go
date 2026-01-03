package server

import (
	"strings"
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

type searchResultDTO struct {
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

// Provider DTO converters

func providerConfigToDTO(config *models.ProviderConfig) ProviderConfigDTO {
	return ProviderConfigDTO{
		ProviderType: config.ProviderType,
		DisplayName:  config.DisplayName,
		Enabled:      config.Enabled,
		Config:       sanitizeConfig(config.ConfigJSON),
		CreatedAt:    config.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    config.UpdatedAt.Format(time.RFC3339),
	}
}

func providerConfigListToDTO(configs []*models.ProviderConfig) []ProviderConfigDTO {
	dtos := make([]ProviderConfigDTO, len(configs))
	for i, config := range configs {
		dtos[i] = providerConfigToDTO(config)
	}
	return dtos
}

func providerTypeToDTO(t *models.ProviderType) ProviderTypeDTO {
	schema := make([]ProviderConfigFieldDTO, len(t.ConfigSchema))
	for i, field := range t.ConfigSchema {
		schema[i] = ProviderConfigFieldDTO{
			Name:        field.Name,
			DisplayName: field.DisplayName,
			Type:        field.Type,
			Required:    field.Required,
			Default:     field.Default,
			Description: field.Description,
		}
	}

	return ProviderTypeDTO{
		Type:         t.Type,
		DisplayName:  t.DisplayName,
		Description:  t.Description,
		RequiresAuth: t.RequiresAuth,
		ConfigSchema: schema,
	}
}

func providerTypeListToDTO(types []*models.ProviderType) []ProviderTypeDTO {
	dtos := make([]ProviderTypeDTO, len(types))
	for i, t := range types {
		dtos[i] = providerTypeToDTO(t)
	}
	return dtos
}

// sanitizeConfig removes sensitive fields from config for API responses
func sanitizeConfig(config map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	sensitiveKeys := []string{"secret", "apiKey", "password", "token"}

	for key, value := range config {
		isSensitive := false
		for _, sensitiveKey := range sensitiveKeys {
			if strings.Contains(strings.ToLower(key), strings.ToLower(sensitiveKey)) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			sanitized[key] = "***REDACTED***"
		} else {
			sanitized[key] = value
		}
	}

	return sanitized
}
