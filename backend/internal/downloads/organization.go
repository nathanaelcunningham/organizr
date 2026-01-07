package downloads

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"

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

	// Pre-organization validation: check source files exist and are readable
	var totalSize int64
	for _, file := range files {
		info, err := os.Stat(file.Path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("source file does not exist: %s", file.Path)
			}
			return fmt.Errorf("source file is not accessible: %s: %w", file.Path, err)
		}
		totalSize += info.Size()
	}

	// Check available disk space at destination
	var stat syscall.Statfs_t
	if err := syscall.Statfs(destBase, &stat); err != nil {
		return fmt.Errorf("failed to check disk space at destination: %w", err)
	}

	// Available space = block size * available blocks
	availableSpace := int64(stat.Bavail) * int64(stat.Bsize)

	// Return error if insufficient space (with buffer of 10% for filesystem overhead)
	requiredSpace := int64(float64(totalSize) * 1.1)
	if availableSpace < requiredSpace {
		return fmt.Errorf("insufficient disk space: need %s, only %s available",
			formatBytes(requiredSpace), formatBytes(availableSpace))
	}

	// Log organization start
	log.Printf("Organizing %d files (%.2f MB) from torrent %s to %s",
		len(files), float64(totalSize)/(1024*1024), dl.QBitHash, fullPath)

	// Create destination directory
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Track successfully copied files for cleanup on partial failure
	var copiedFiles []string

	// Defer cleanup function to handle panic during copy operations
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic during organization: %v, cleaning up %d files", r, len(copiedFiles))
			for _, path := range copiedFiles {
				if err := os.Remove(path); err != nil {
					log.Printf("Failed to clean up file %s: %v", path, err)
				}
			}
			panic(r) // Re-panic after cleanup
		}
	}()

	// Copy or move files
	for i, file := range files {
		srcPath := file.Path
		destPath := filepath.Join(fullPath, filepath.Base(file.Name))

		if operation == "move" {
			// Move operation: atomic per file, partial success is acceptable
			if err := os.Rename(srcPath, destPath); err != nil {
				log.Printf("Failed to move file %s (%s -> %s): %v", file.Name, srcPath, destPath, err)
				return fmt.Errorf("failed to move file %s (%s -> %s): %w", file.Name, srcPath, destPath, err)
			}
		} else {
			// Copy operation: all-or-nothing, clean up on failure
			if err := copyFile(srcPath, destPath); err != nil {
				// Cleanup: delete all previously copied files
				log.Printf("Failed to copy file %s (%s -> %s): %v, cleaning up %d previously copied files",
					file.Name, srcPath, destPath, err, len(copiedFiles))
				for _, path := range copiedFiles {
					if removeErr := os.Remove(path); removeErr != nil {
						log.Printf("Failed to clean up file %s: %v", path, removeErr)
					}
				}
				return fmt.Errorf("failed to copy file %s (%s -> %s): %w", file.Name, srcPath, destPath, err)
			}
			// Track successfully copied file
			copiedFiles = append(copiedFiles, destPath)
			log.Printf("Successfully copied file %d/%d: %s", i+1, len(files), file.Name)
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

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
