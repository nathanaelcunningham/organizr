package downloads

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanael/organizr/internal/config"
	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/persistence"
	"github.com/nathanael/organizr/internal/qbittorrent"
)

type Service struct {
	db            *sql.DB
	qbClient      *qbittorrent.Client
	downloadRepo  persistence.DownloadRepository
	configService *config.Service
}

func NewService(db *sql.DB, qbClient *qbittorrent.Client, downloadRepo persistence.DownloadRepository, configService *config.Service) *Service {
	return &Service{
		db:            db,
		qbClient:      qbClient,
		downloadRepo:  downloadRepo,
		configService: configService,
	}
}

func (s *Service) CreateDownload(ctx context.Context, d *models.Download) (*models.Download, error) {
	// Validate input
	if d.Title == "" || d.Author == "" {
		return nil, fmt.Errorf("title and author are required")
	}
	if d.TorrentURL == "" && d.MagnetLink == "" {
		return nil, fmt.Errorf("either torrent URL or magnet link is required")
	}

	// Generate ID
	d.ID = uuid.New().String()

	// Add to qBittorrent
	hash, err := s.qbClient.AddTorrent(ctx, d.MagnetLink, d.TorrentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to add torrent to qBittorrent: %w", err)
	}

	d.QBitHash = hash
	d.Status = models.StatusQueued

	// Save to database
	if err := s.downloadRepo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("failed to save download: %w", err)
	}

	return d, nil
}

func (s *Service) GetDownload(ctx context.Context, id string) (*models.Download, error) {
	download, err := s.downloadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get download: %w", err)
	}
	return download, nil
}

func (s *Service) ListDownloads(ctx context.Context) ([]*models.Download, error) {
	downloads, err := s.downloadRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list downloads: %w", err)
	}
	return downloads, nil
}

func (s *Service) CancelDownload(ctx context.Context, id string) error {
	download, err := s.downloadRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get download: %w", err)
	}

	// Delete from qBittorrent
	if err := s.qbClient.DeleteTorrent(ctx, download.QBitHash, false); err != nil {
		return fmt.Errorf("failed to delete torrent from qBittorrent: %w", err)
	}

	// Delete from database
	if err := s.downloadRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete download: %w", err)
	}

	return nil
}

func (s *Service) OrganizeDownload(ctx context.Context, id string) error {
	download, err := s.downloadRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get download: %w", err)
	}

	// Create organization service and organize
	orgService := NewOrganizationService(s.qbClient, s.configService)
	if err := orgService.Organize(ctx, download); err != nil {
		return fmt.Errorf("failed to organize download: %w", err)
	}

	// Update status and path in database
	if err := s.downloadRepo.UpdateStatus(ctx, id, models.StatusOrganized); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	if err := s.downloadRepo.UpdateOrganizedPath(ctx, id, download.OrganizedPath); err != nil {
		return fmt.Errorf("failed to update organized path: %w", err)
	}

	return nil
}
