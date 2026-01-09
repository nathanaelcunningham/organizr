package downloads

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nathanael/organizr/internal/config"
	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/persistence"
	"github.com/nathanael/organizr/internal/qbittorrent"
)

type Monitor struct {
	db            *sql.DB
	qbClient      *qbittorrent.Client
	downloadRepo  persistence.DownloadRepository
	orgService    *OrganizationService
	configService *config.Service
	interval      time.Duration
	maxConcurrent int
}

func NewMonitor(db *sql.DB, qbClient *qbittorrent.Client, downloadRepo persistence.DownloadRepository, configService *config.Service) *Monitor {
	return &Monitor{
		db:            db,
		qbClient:      qbClient,
		downloadRepo:  downloadRepo,
		orgService:    NewOrganizationService(qbClient, configService),
		configService: configService,
		interval:      30 * time.Second,
		maxConcurrent: 3,
	}
}

func (m *Monitor) Run(ctx context.Context) error {
	// Get interval from config
	intervalStr, err := m.configService.Get(ctx, "monitor.interval_seconds")
	if err == nil {
		if seconds, err := strconv.Atoi(intervalStr); err == nil {
			m.interval = time.Duration(seconds) * time.Second
		}
	}

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	log.Println("Download monitor started")

	for {
		select {
		case <-ticker.C:
			if err := m.checkDownloads(ctx); err != nil {
				log.Printf("Monitor error: %v", err)
			}
		case <-ctx.Done():
			log.Println("Monitor stopped")
			return ctx.Err()
		}
	}
}

func (m *Monitor) checkDownloads(ctx context.Context) error {
	// Get active downloads
	downloads, err := m.downloadRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active downloads: %w", err)
	}

	// Track if all downloads failed (suggests qBittorrent is down)
	allFailed := true
	var lastErr error

	for _, dl := range downloads {
		// Check status in qBittorrent
		status, progress, err := m.qbClient.GetTorrentStatus(ctx, dl.QBitHash)
		if err != nil {
			log.Printf("Warning: Failed to get status for download %s (%s): %v", dl.ID, dl.Title, err)
			lastErr = err
			continue
		}

		// At least one download succeeded
		allFailed = false

		// Update progress
		if err := m.downloadRepo.UpdateProgress(ctx, dl.ID, progress); err != nil {
			log.Printf("Failed to update progress for download %s: %v", dl.ID, err)
		}

		// Map qBittorrent state to our status
		newStatus := mapQBitStatusToModel(status)

		// Log state transitions (only when state actually changes)
		if newStatus != dl.Status && newStatus != "" {
			log.Printf("Download %s (%s) state changed: %s → %s", dl.ID, dl.Title, dl.Status, newStatus)
		}

		// Check if completed
		if (status == "uploading" || status == "stalledUP" || status == "pausedUP") && dl.Status != models.StatusOrganized {
			log.Printf("Download %s (%s) completed, marking as complete", dl.ID, dl.Title)

			// Mark completed
			if err := m.downloadRepo.UpdateCompleted(ctx, dl.ID); err != nil {
				log.Printf("Failed to mark download as completed: %v", err)
				continue
			}

			// Check if auto-organization is enabled (default: true for backward compatibility)
			autoOrganize := true
			if autoOrganizeStr, err := m.configService.Get(ctx, "organization.auto_organize"); err == nil {
				if autoOrganizeStr == "false" {
					autoOrganize = false
					log.Printf("Auto-organization disabled for download %s (%s), skipping organization", dl.ID, dl.Title)
				}
			}

			// Only auto-organize if enabled
			if autoOrganize {
				log.Printf("Auto-organizing download %s (%s)", dl.ID, dl.Title)

				// Create context with timeout for organization (respects parent cancellation but prevents indefinite hanging)
				orgCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)

				// Trigger organization in goroutine
				go func(ctx context.Context, download *models.Download) {
					defer cancel()
					m.organizeDownload(ctx, download)
				}(orgCtx, dl)
			}
		}
	}

	// If all downloads failed, qBittorrent may be unavailable
	if len(downloads) > 0 && allFailed {
		log.Printf("Warning: qBittorrent may be unavailable - all %d download status checks failed (last error: %v)", len(downloads), lastErr)
		// Don't return error - continue monitoring, qBittorrent may recover
	}

	return nil
}

// mapQBitStatusToModel maps qBittorrent state to our download status
func mapQBitStatusToModel(qbitStatus string) models.DownloadStatus {
	switch qbitStatus {
	case "queuedDL", "queuedUP":
		return models.StatusQueued
	case "downloading", "metaDL", "allocating", "checkingDL", "forcedDL":
		return models.StatusDownloading
	case "uploading", "stalledUP", "pausedUP", "forcedUP", "checkingUP":
		return models.StatusCompleted
	default:
		return "" // Unknown state, don't update
	}
}

func (m *Monitor) organizeDownload(ctx context.Context, dl *models.Download) {
	// Mark as organizing
	if err := m.downloadRepo.UpdateStatus(ctx, dl.ID, models.StatusOrganizing); err != nil {
		log.Printf("Failed to update status for download %s: %v", dl.ID, err)
		return
	}

	// Perform organization
	if err := m.orgService.Organize(ctx, dl); err != nil {
		log.Printf("Failed to organize download %s: %v", dl.ID, err)
		// Update error status in database - log if this also fails
		if updateErr := m.downloadRepo.UpdateError(ctx, dl.ID, err.Error()); updateErr != nil {
			log.Printf("Failed to update download error for %s: %v", dl.ID, updateErr)
		}
		if updateErr := m.downloadRepo.UpdateStatus(ctx, dl.ID, models.StatusFailed); updateErr != nil {
			log.Printf("Failed to update download status for %s: %v", dl.ID, updateErr)
		}
		return
	}

	// Mark as organized
	if err := m.downloadRepo.UpdateStatus(ctx, dl.ID, models.StatusOrganized); err != nil {
		log.Printf("Failed to update status for download %s: %v", dl.ID, err)
		return
	}

	if err := m.downloadRepo.UpdateOrganizedPath(ctx, dl.ID, dl.OrganizedPath); err != nil {
		log.Printf("Failed to update organized path for download %s: %v", dl.ID, err)
	}

	log.Printf("Download %s organized successfully to %s", dl.ID, dl.OrganizedPath)
}
