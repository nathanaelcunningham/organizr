package persistence

import (
	"context"

	"github.com/nathanael/organizr/internal/models"
)

type DownloadRepository interface {
	Create(ctx context.Context, d *models.Download) error
	GetByID(ctx context.Context, id string) (*models.Download, error)
	GetActive(ctx context.Context) ([]*models.Download, error)
	List(ctx context.Context) ([]*models.Download, error)
	UpdateStatus(ctx context.Context, id string, status models.DownloadStatus) error
	UpdateProgress(ctx context.Context, id string, progress float64) error
	UpdateError(ctx context.Context, id string, errorMsg string) error
	UpdateOrganizedPath(ctx context.Context, id string, path string) error
	UpdateCompleted(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type ConfigRepository interface {
	Get(ctx context.Context, key string) (string, error)
	GetAll(ctx context.Context) (map[string]string, error)
	Set(ctx context.Context, key, value string) error
}

type ProviderRepository interface {
	Create(ctx context.Context, config *models.ProviderConfig) error
	GetByType(ctx context.Context, providerType string) (*models.ProviderConfig, error)
	List(ctx context.Context) ([]*models.ProviderConfig, error)
	ListEnabled(ctx context.Context) ([]*models.ProviderConfig, error)
	Update(ctx context.Context, config *models.ProviderConfig) error
	UpdateEnabled(ctx context.Context, providerType string, enabled bool) error
	Delete(ctx context.Context, providerType string) error
	Exists(ctx context.Context, providerType string) (bool, error)
}
