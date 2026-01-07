package downloads

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nathanael/organizr/internal/config"
	"github.com/nathanael/organizr/internal/fileutil"
	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/qbittorrent"
)

// qbittorrentClient interface defines the methods we need from qbittorrent.Client
type qbittorrentClient interface {
	GetTorrentFiles(ctx context.Context, hash string) ([]*qbittorrent.TorrentFile, error)
}

// configService interface defines the methods we need from config.Service
type configService interface {
	Get(ctx context.Context, key string) (string, error)
}

type OrganizationService struct {
	qbClient      qbittorrentClient
	configService configService
}

func NewOrganizationService(qbClient *qbittorrent.Client, configService *config.Service) *OrganizationService {
	return &OrganizationService{
		qbClient:      qbClient,
		configService: configService,
	}
}

func (o *OrganizationService) Organize(ctx context.Context, dl *models.Download) error {
	// Get config
	destBase, err := o.configService.Get(ctx, "paths.destination")
	if err != nil {
		return fmt.Errorf("failed to get destination path: %w", err)
	}

	// Ensure base destination directory exists
	if err := os.MkdirAll(destBase, 0755); err != nil {
		return fmt.Errorf("failed to create base destination directory %s: %w", destBase, err)
	}

	template, err := o.configService.Get(ctx, "paths.template")
	if err != nil {
		template = "{author}/{series}/{title}"
	}

	noSeriesTemplate, err := o.configService.Get(ctx, "paths.no_series_template")
	if err != nil {
		noSeriesTemplate = "{author}/{title}"
	}

	operation, err := o.configService.Get(ctx, "paths.operation")
	if err != nil {
		operation = "copy"
	}

	// Choose template based on whether series exists
	pathTemplate := template
	if dl.Series == "" {
		pathTemplate = noSeriesTemplate
	}

	// Sanitize variables BEFORE template parsing (preserves directory structure)
	sanitizedVars := map[string]string{
		"author": fileutil.SanitizePath(dl.Author),
		"series": fileutil.SanitizePath(dl.Series),
		"title":  fileutil.SanitizePath(dl.Title),
	}

	// Parse template with sanitized variables
	path := fileutil.ParseTemplate(pathTemplate, sanitizedVars)

	// Build full destination
	fullPath := filepath.Join(destBase, path)

	// Get torrent files from qBittorrent
	files, err := o.qbClient.GetTorrentFiles(ctx, dl.QBitHash)
	if err != nil {
		return fmt.Errorf("failed to get torrent files: %w", err)
	}

	// Get mount point configuration for remote qBittorrent setups
	mountPoint, _ := o.configService.Get(ctx, "paths.local_mount")

	// Prepend mount point if configured (for network shares or Docker volumes)
	if mountPoint != "" {
		for _, file := range files {
			// qBittorrent reports paths relative to its filesystem
			// Prepend the local mount point to access them
			file.Path = filepath.Join(mountPoint, file.Path)
		}
	}

	// Create destination directory
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy or move files
	for _, file := range files {
		srcPath := file.Path
		destPath := filepath.Join(fullPath, filepath.Base(file.Name))

		if operation == "move" {
			if err := os.Rename(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to move file %s: %w", file.Name, err)
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", file.Name, err)
			}
		}
	}

	dl.OrganizedPath = fullPath
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Sync to ensure data is written to disk
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}
