package testutil

import (
	"testing"

	"github.com/nathanael/organizr/internal/models"
)

func TestNewTestDownload(t *testing.T) {
	t.Run("creates download with defaults", func(t *testing.T) {
		dl := NewTestDownload()

		if dl.ID == "" {
			t.Error("expected ID to be set")
		}

		if dl.Title != "Test Download" {
			t.Errorf("expected title 'Test Download', got %s", dl.Title)
		}

		if dl.Author != "Test Author" {
			t.Errorf("expected author 'Test Author', got %s", dl.Author)
		}

		if dl.Status != models.StatusQueued {
			t.Errorf("expected status queued, got %s", dl.Status)
		}

		if dl.Progress != 0 {
			t.Errorf("expected progress 0, got %f", dl.Progress)
		}

		if dl.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}
	})

	t.Run("applies functional options", func(t *testing.T) {
		dl := NewTestDownload(
			WithTitle("Custom Title"),
			WithAuthor("Custom Author"),
			WithProgress(50),
			WithStatus(models.StatusDownloading),
		)

		if dl.Title != "Custom Title" {
			t.Errorf("expected title 'Custom Title', got %s", dl.Title)
		}

		if dl.Author != "Custom Author" {
			t.Errorf("expected author 'Custom Author', got %s", dl.Author)
		}

		if dl.Progress != 50 {
			t.Errorf("expected progress 50, got %f", dl.Progress)
		}

		if dl.Status != models.StatusDownloading {
			t.Errorf("expected status downloading, got %s", dl.Status)
		}
	})

	t.Run("applies all available options", func(t *testing.T) {
		dl := NewTestDownload(
			WithID("test-id-123"),
			WithTitle("Custom Book"),
			WithAuthor("John Doe"),
			WithSeries("Test Series"),
			WithSeriesNumber("1"),
			WithStatus(models.StatusCompleted),
			WithProgress(100),
			WithMagnetLink("magnet:?xt=urn:btih:test"),
			WithTorrentURL("https://example.com/test.torrent"),
			WithCategory("Audiobooks"),
			WithQBitHash("abc123"),
			WithOrganizedPath("/path/to/organized"),
		)

		if dl.ID != "test-id-123" {
			t.Errorf("expected ID 'test-id-123', got %s", dl.ID)
		}

		if dl.Title != "Custom Book" {
			t.Errorf("expected title 'Custom Book', got %s", dl.Title)
		}

		if dl.Series != "Test Series" {
			t.Errorf("expected series 'Test Series', got %s", dl.Series)
		}

		if dl.SeriesNumber != "1" {
			t.Errorf("expected series number '1', got %s", dl.SeriesNumber)
		}

		if dl.QBitHash != "abc123" {
			t.Errorf("expected hash 'abc123', got %s", dl.QBitHash)
		}

		if dl.OrganizedPath != "/path/to/organized" {
			t.Errorf("expected organized path '/path/to/organized', got %s", dl.OrganizedPath)
		}
	})
}

func TestNewTestSearchResult(t *testing.T) {
	t.Run("creates search result with defaults", func(t *testing.T) {
		sr := NewTestSearchResult()

		if sr.ID == "" {
			t.Error("expected ID to be set")
		}

		if sr.Title != "Test Book" {
			t.Errorf("expected title 'Test Book', got %s", sr.Title)
		}

		if sr.Author != "Test Author" {
			t.Errorf("expected author 'Test Author', got %s", sr.Author)
		}

		if sr.Provider != "MyAnonamouse" {
			t.Errorf("expected provider 'MyAnonamouse', got %s", sr.Provider)
		}

		if sr.Category != "Audiobooks" {
			t.Errorf("expected category 'Audiobooks', got %s", sr.Category)
		}

		if sr.Seeders != 10 {
			t.Errorf("expected 10 seeders, got %d", sr.Seeders)
		}
	})

	t.Run("applies functional options", func(t *testing.T) {
		series := []models.SeriesInfo{
			{ID: "1", Name: "Test Series", Number: "1"},
		}

		sr := NewTestSearchResult(
			WithResultID("custom-id"),
			WithResultTitle("Custom Book"),
			WithResultAuthor("Jane Doe"),
			WithResultSeries(series),
			WithResultProvider("TestProvider"),
			WithSeeders(100),
			WithFreeleech(true),
		)

		if sr.ID != "custom-id" {
			t.Errorf("expected ID 'custom-id', got %s", sr.ID)
		}

		if sr.Title != "Custom Book" {
			t.Errorf("expected title 'Custom Book', got %s", sr.Title)
		}

		if sr.Author != "Jane Doe" {
			t.Errorf("expected author 'Jane Doe', got %s", sr.Author)
		}

		if len(sr.Series) != 1 || sr.Series[0].Name != "Test Series" {
			t.Errorf("expected series 'Test Series', got %+v", sr.Series)
		}

		if sr.Provider != "TestProvider" {
			t.Errorf("expected provider 'TestProvider', got %s", sr.Provider)
		}

		if sr.Seeders != 100 {
			t.Errorf("expected 100 seeders, got %d", sr.Seeders)
		}

		if !sr.Freeleech {
			t.Error("expected freeleech to be true")
		}
	})
}

func TestNewTestConfig(t *testing.T) {
	t.Run("creates config with defaults", func(t *testing.T) {
		config := NewTestConfig(nil)

		expectedKeys := []string{
			"qbittorrent.url",
			"qbittorrent.username",
			"qbittorrent.password",
			"mam.baseurl",
			"mam.secret",
			"output.path",
			"output.template",
		}

		for _, key := range expectedKeys {
			if _, ok := config[key]; !ok {
				t.Errorf("expected config to have key %s", key)
			}
		}

		if config["qbittorrent.url"] != "http://localhost:8080" {
			t.Errorf("expected qbittorrent.url to be 'http://localhost:8080', got %s", config["qbittorrent.url"])
		}
	})

	t.Run("applies overrides", func(t *testing.T) {
		overrides := map[string]string{
			"qbittorrent.url":      "http://custom:9090",
			"custom.key":           "custom value",
		}

		config := NewTestConfig(overrides)

		if config["qbittorrent.url"] != "http://custom:9090" {
			t.Errorf("expected custom qbittorrent.url, got %s", config["qbittorrent.url"])
		}

		if config["custom.key"] != "custom value" {
			t.Errorf("expected custom.key to be 'custom value', got %s", config["custom.key"])
		}

		// Check defaults still present
		if config["qbittorrent.username"] != "admin" {
			t.Error("expected default username to still be present")
		}
	})

	t.Run("handles empty overrides", func(t *testing.T) {
		config := NewTestConfig(map[string]string{})

		if len(config) != 7 {
			t.Errorf("expected 7 config keys, got %d", len(config))
		}
	})
}
