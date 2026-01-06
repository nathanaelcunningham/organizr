package server

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// UUID v4 pattern
	uuidPattern = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)
)

// validateDownloadRequest validates a download creation request
func validateDownloadRequest(req CreateDownloadRequest) error {
	// Validate required fields
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required and cannot be empty")
	}

	if strings.TrimSpace(req.Author) == "" {
		return fmt.Errorf("author is required and cannot be empty")
	}

	// Validate torrent source
	if req.TorrentURL == "" && req.MagnetLink == "" {
		return fmt.Errorf("either torrent_url or magnet_link is required")
	}

	// Validate lengths
	if len(req.Title) > 500 {
		return fmt.Errorf("title must be 500 characters or less")
	}

	if len(req.Author) > 200 {
		return fmt.Errorf("author must be 200 characters or less")
	}

	if len(req.Series) > 200 {
		return fmt.Errorf("series must be 200 characters or less")
	}

	return nil
}

// validateUUID validates a UUID string
func validateUUID(id string) error {
	if !uuidPattern.MatchString(id) {
		return fmt.Errorf("invalid UUID format")
	}
	return nil
}

// validateConfigKey validates a configuration key
func validateConfigKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("config key cannot be empty")
	}

	// Only allow alphanumeric, dots, underscores, and hyphens
	validKey := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validKey.MatchString(key) {
		return fmt.Errorf("config key can only contain letters, numbers, dots, underscores, and hyphens")
	}

	return nil
}
