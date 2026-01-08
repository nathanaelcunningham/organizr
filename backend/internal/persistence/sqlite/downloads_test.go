package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/nathanael/organizr/internal/models"
)

func TestDownloadRepository_NullHandling(t *testing.T) {
	// Create in-memory database
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create downloads table
	schema := `
		CREATE TABLE downloads (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			series TEXT,
			series_number TEXT,
			torrent_url TEXT,
			magnet_link TEXT,
			category TEXT,
			qbit_hash TEXT UNIQUE NOT NULL,
			status TEXT NOT NULL DEFAULT 'queued',
			progress REAL DEFAULT 0.0,
			download_path TEXT,
			organized_path TEXT,
			error_message TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP,
			organized_at TIMESTAMP
		);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	repo := NewDownloadRepository(db)
	ctx := context.Background()

	// Test 1: Create download with minimal fields (nullable fields are NULL)
	download := &models.Download{
		ID:        "test-id-1",
		Title:     "Test Audiobook",
		Author:    "Test Author",
		QBitHash:  "testhash123",
		Status:    models.StatusQueued,
		CreatedAt: time.Now(),
	}

	if err := repo.Create(ctx, download); err != nil {
		t.Fatalf("Failed to create download: %v", err)
	}

	// Test 2: Retrieve download - should not fail on NULL fields
	retrieved, err := repo.GetByID(ctx, "test-id-1")
	if err != nil {
		t.Fatalf("Failed to get download by ID: %v", err)
	}

	// Verify basic fields
	if retrieved.Title != download.Title {
		t.Errorf("Expected title %q, got %q", download.Title, retrieved.Title)
	}
	if retrieved.Author != download.Author {
		t.Errorf("Expected author %q, got %q", download.Author, retrieved.Author)
	}

	// Verify nullable fields are empty strings (not errors)
	if retrieved.Series != "" {
		t.Errorf("Expected empty series, got %q", retrieved.Series)
	}
	if retrieved.DownloadPath != "" {
		t.Errorf("Expected empty download_path, got %q", retrieved.DownloadPath)
	}
	if retrieved.OrganizedPath != "" {
		t.Errorf("Expected empty organized_path, got %q", retrieved.OrganizedPath)
	}
	if retrieved.ErrorMessage != "" {
		t.Errorf("Expected empty error_message, got %q", retrieved.ErrorMessage)
	}

	// Test 3: Create download with series and series number
	download2 := &models.Download{
		ID:           "test-id-2",
		Title:        "Book with Series",
		Author:       "Another Author",
		Series:       "My Series",
		SeriesNumber: "1",
		QBitHash:     "testhash456",
		Status:       models.StatusQueued,
		CreatedAt:    time.Now(),
	}

	if err := repo.Create(ctx, download2); err != nil {
		t.Fatalf("Failed to create download 2: %v", err)
	}

	retrieved2, err := repo.GetByID(ctx, "test-id-2")
	if err != nil {
		t.Fatalf("Failed to get download 2 by ID: %v", err)
	}

	if retrieved2.Series != download2.Series {
		t.Errorf("Expected series %q, got %q", download2.Series, retrieved2.Series)
	}

	if retrieved2.SeriesNumber != download2.SeriesNumber {
		t.Errorf("Expected series_number %q, got %q", download2.SeriesNumber, retrieved2.SeriesNumber)
	}

	// Test 4: List downloads - should handle mixed NULL/non-NULL values
	allDownloads, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list downloads: %v", err)
	}

	if len(allDownloads) != 2 {
		t.Errorf("Expected 2 downloads, got %d", len(allDownloads))
	}

	// Test 5: GetActive - should handle mixed NULL/non-NULL values
	activeDownloads, err := repo.GetActive(ctx)
	if err != nil {
		t.Fatalf("Failed to get active downloads: %v", err)
	}

	if len(activeDownloads) != 2 {
		t.Errorf("Expected 2 active downloads, got %d", len(activeDownloads))
	}

	t.Log("✓ All NULL handling tests passed")
}
