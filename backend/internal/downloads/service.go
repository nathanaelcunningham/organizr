package downloads

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nathanael/organizr/internal/config"
	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/persistence"
	"github.com/nathanael/organizr/internal/qbittorrent"
	"github.com/nathanael/organizr/internal/search"
)

type Service struct {
	db            *sql.DB
	qbClient      *qbittorrent.Client
	downloadRepo  persistence.DownloadRepository
	configService *config.Service
	mamService    *search.MAMService
}

func NewService(db *sql.DB, qbClient *qbittorrent.Client, downloadRepo persistence.DownloadRepository, configService *config.Service, mamService *search.MAMService) *Service {
	return &Service{
		db:            db,
		qbClient:      qbClient,
		downloadRepo:  downloadRepo,
		configService: configService,
		mamService:    mamService,
	}
}

func (s *Service) CreateDownload(ctx context.Context, d *models.Download) (*models.Download, error) {
	// Validate input
	if d.Title == "" || d.Author == "" {
		return nil, fmt.Errorf("title and author are required")
	}
	if d.TorrentURL == "" && d.MagnetLink == "" && len(d.TorrentBytes) == 0 {
		return nil, fmt.Errorf("either torrent URL, magnet link, or torrent bytes is required")
	}

	// Generate ID
	d.ID = uuid.New().String()

	var hash string
	var err error

	// Determine the download method
	if len(d.TorrentBytes) > 0 {
		// Use torrent bytes (from MAM or direct upload)
		hash, err = s.qbClient.AddTorrentFromFile(ctx, d.TorrentBytes, d.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to add torrent from file to qBittorrent: %w", err)
		}
	} else if d.TorrentURL != "" && strings.Contains(d.TorrentURL, "/tor/download.php") {
		// MAM URL - validate format first
		if !strings.Contains(d.TorrentURL, "myanonamouse.net") {
			return nil, fmt.Errorf("invalid MAM torrent URL: must be from myanonamouse.net")
		}

		// Extract torrent ID from URL
		torrentID, err := extractTorrentIDFromURL(d.TorrentURL)
		if err != nil {
			return nil, fmt.Errorf("invalid MAM torrent URL: %w", err)
		}

		// Download torrent file from MAM
		torrentData, err := s.mamService.DownloadTorrent(ctx, torrentID)
		if err != nil {
			// Categorize MAM errors
			if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
				return nil, fmt.Errorf("MAM authentication failed: please check your credentials in settings")
			} else if strings.Contains(err.Error(), "404") {
				return nil, fmt.Errorf("torrent not found on MAM (ID: %d)", torrentID)
			}
			return nil, fmt.Errorf("failed to download torrent from MAM: %w", err)
		}

		// Add torrent from file data
		hash, err = s.qbClient.AddTorrentFromFile(ctx, torrentData, d.Category)
		if err != nil {
			// Categorize qBittorrent errors
			if strings.Contains(err.Error(), "authentication failed") || strings.Contains(err.Error(), "invalid username or password") {
				return nil, fmt.Errorf("qBittorrent authentication failed: check username and password in settings")
			} else if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
				// Get base URL from qbClient if possible, otherwise use generic message
				return nil, fmt.Errorf("qBittorrent not reachable: ensure qBittorrent is running and Web UI is enabled")
			}
			return nil, fmt.Errorf("failed to add torrent to qBittorrent: %w", err)
		}
	} else {
		// Use magnet link or direct URL
		hash, err = s.qbClient.AddTorrent(ctx, d.MagnetLink, d.TorrentURL, d.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to add torrent to qBittorrent: %w", err)
		}
	}

	d.QBitHash = hash
	d.Status = models.StatusQueued

	// Save to database
	if err := s.downloadRepo.Create(ctx, d); err != nil {
		// Check for unique constraint violations
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate") {
			return nil, fmt.Errorf("download with this ID already exists: %w", err)
		}
		return nil, fmt.Errorf("failed to save download to database: %w", err)
	}

	return d, nil
}

// extractTorrentIDFromURL extracts the torrent ID from a MAM download URL
// Example: https://www.myanonamouse.net/tor/download.php?tid=12345
func extractTorrentIDFromURL(url string) (int, error) {
	// Validate URL format
	if !strings.Contains(url, "tid=") {
		return 0, fmt.Errorf("missing required 'tid' parameter")
	}

	// Find the tid parameter
	parts := strings.Split(url, "tid=")
	if len(parts) < 2 {
		return 0, fmt.Errorf("torrent ID not found in URL")
	}

	// Get the ID part (might have other params after it)
	idPart := parts[1]

	// Remove any additional query parameters
	if idx := strings.Index(idPart, "&"); idx != -1 {
		idPart = idPart[:idx]
	}

	// Validate ID is not empty
	if idPart == "" {
		return 0, fmt.Errorf("torrent ID is empty")
	}

	// Convert to int
	id, err := strconv.Atoi(idPart)
	if err != nil {
		return 0, fmt.Errorf("invalid torrent ID format (expected number): %w", err)
	}

	return id, nil
}

func (s *Service) GetDownload(ctx context.Context, id string) (*models.Download, error) {
	download, err := s.downloadRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("download not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get download from database: %w", err)
	}
	return download, nil
}

func (s *Service) ListDownloads(ctx context.Context) ([]*models.Download, error) {
	downloads, err := s.downloadRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list downloads from database: %w", err)
	}
	return downloads, nil
}

func (s *Service) CancelDownload(ctx context.Context, id string) error {
	download, err := s.downloadRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("download not found: %s", id)
		}
		return fmt.Errorf("failed to get download from database: %w", err)
	}

	// Delete from qBittorrent
	if download.QBitHash != "" {
		if err := s.qbClient.DeleteTorrent(ctx, download.QBitHash, false); err != nil {
			// Log but don't fail if torrent already removed from qBittorrent
			if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "404") {
				return fmt.Errorf("failed to delete torrent from qBittorrent: %w", err)
			}
		}
	}

	// Delete from database
	if err := s.downloadRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete download from database: %w", err)
	}

	return nil
}

func (s *Service) OrganizeDownload(ctx context.Context, id string) error {
	download, err := s.downloadRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("download not found: %s", id)
		}
		return fmt.Errorf("failed to get download from database: %w", err)
	}

	// Create organization service and organize
	orgService := NewOrganizationService(s.qbClient, s.configService)
	if err := orgService.Organize(ctx, download); err != nil {
		return fmt.Errorf("failed to organize download files: %w", err)
	}

	// Update status and path in database
	if err := s.downloadRepo.UpdateStatus(ctx, id, models.StatusOrganized); err != nil {
		return fmt.Errorf("failed to update download status in database: %w", err)
	}

	if err := s.downloadRepo.UpdateOrganizedPath(ctx, id, download.OrganizedPath); err != nil {
		return fmt.Errorf("failed to update organized path in database: %w", err)
	}

	return nil
}
