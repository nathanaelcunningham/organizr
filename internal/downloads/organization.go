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

type OrganizationService struct {
	qbClient      *qbittorrent.Client
	configService *config.Service
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

	// Parse template
	path := fileutil.ParseTemplate(pathTemplate, map[string]string{
		"author": dl.Author,
		"series": dl.Series,
		"title":  dl.Title,
	})

	// Sanitize path components
	path = fileutil.SanitizePath(path)

	// Build full destination
	fullPath := filepath.Join(destBase, path)

	// Get torrent files from qBittorrent
	files, err := o.qbClient.GetTorrentFiles(ctx, dl.QBitHash)
	if err != nil {
		return fmt.Errorf("failed to get torrent files: %w", err)
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
