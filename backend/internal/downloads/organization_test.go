package downloads

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/qbittorrent"
)

// mockConfigService is a simple mock for config.Service
type mockConfigService struct {
	configs map[string]string
}

func (m *mockConfigService) Get(ctx context.Context, key string) (string, error) {
	if val, ok := m.configs[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("config key not found: %s", key)
}

func (m *mockConfigService) GetAll(ctx context.Context) (map[string]string, error) {
	return m.configs, nil
}

func (m *mockConfigService) Set(ctx context.Context, key, value string) error {
	m.configs[key] = value
	return nil
}

func newMockConfigService(configs map[string]string) *mockConfigService {
	return &mockConfigService{configs: configs}
}

// mockQBClient is a simple mock for qbittorrent.Client
type mockQBClient struct {
	files []*qbittorrent.TorrentFile
	err   error
}

func (m *mockQBClient) GetTorrentFiles(ctx context.Context, hash string) ([]*qbittorrent.TorrentFile, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.files, nil
}

func (m *mockQBClient) Login(ctx context.Context) error {
	return nil
}

func (m *mockQBClient) AddTorrent(ctx context.Context, magnetLink, torrentURL, category string) (string, error) {
	return "", nil
}

func (m *mockQBClient) AddTorrentFromFile(ctx context.Context, torrentData []byte, category string) (string, error) {
	return "", nil
}

func (m *mockQBClient) GetTorrentStatus(ctx context.Context, hash string) (string, float64, error) {
	return "", 0, nil
}

func (m *mockQBClient) DeleteTorrent(ctx context.Context, hash string, deleteFiles bool) error {
	return nil
}

// newTestOrganizationService creates an OrganizationService for testing with mock dependencies
func newTestOrganizationService(qbClient *mockQBClient, configService *mockConfigService) *OrganizationService {
	return &OrganizationService{
		qbClient:      qbClient,
		configService: configService,
	}
}

func TestOrganize(t *testing.T) {
	tests := []struct {
		name           string
		download       *models.Download
		configs        map[string]string
		sourceFiles    map[string]string // filename -> content
		qbFiles        []*qbittorrent.TorrentFile
		qbError        error
		wantErr        bool
		wantErrContain string
		verifyFn       func(t *testing.T, destBase string, dl *models.Download)
	}{
		{
			name: "successful copy operation with series",
			download: &models.Download{
				ID:       "test-1",
				Title:    "Book One",
				Author:   "Jane Doe",
				Series:   "Epic Series",
				QBitHash: "abc123",
			},
			configs: map[string]string{
				"paths.destination":        "", // will be set to tempdir
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "copy",
			},
			sourceFiles: map[string]string{
				"Book One.m4b": "audio content here",
			},
			qbFiles: []*qbittorrent.TorrentFile{
				{
					Name: "Book One.m4b",
					Path: "", // will be set to actual file path
					Size: 1024,
				},
			},
			wantErr: false,
			verifyFn: func(t *testing.T, destBase string, dl *models.Download) {
				// Check organized path is set correctly
				expectedPath := filepath.Join(destBase, "Jane Doe", "Epic Series", "Book One")
				if dl.OrganizedPath != expectedPath {
					t.Errorf("OrganizedPath = %v, want %v", dl.OrganizedPath, expectedPath)
				}

				// Check file exists at destination
				destFile := filepath.Join(expectedPath, "Book One.m4b")
				if _, err := os.Stat(destFile); os.IsNotExist(err) {
					t.Errorf("destination file does not exist: %s", destFile)
				}

				// Check file content
				content, err := os.ReadFile(destFile)
				if err != nil {
					t.Fatalf("failed to read destination file: %v", err)
				}
				if string(content) != "audio content here" {
					t.Errorf("file content = %v, want 'audio content here'", string(content))
				}
			},
		},
		{
			name: "successful copy operation without series",
			download: &models.Download{
				ID:       "test-2",
				Title:    "Standalone Book",
				Author:   "John Smith",
				Series:   "", // No series
				QBitHash: "def456",
			},
			configs: map[string]string{
				"paths.destination":        "", // will be set to tempdir
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "copy",
			},
			sourceFiles: map[string]string{
				"Standalone Book.m4b": "standalone content",
			},
			qbFiles: []*qbittorrent.TorrentFile{
				{
					Name: "Standalone Book.m4b",
					Path: "",
					Size: 2048,
				},
			},
			wantErr: false,
			verifyFn: func(t *testing.T, destBase string, dl *models.Download) {
				// Should use no_series_template
				expectedPath := filepath.Join(destBase, "John Smith", "Standalone Book")
				if dl.OrganizedPath != expectedPath {
					t.Errorf("OrganizedPath = %v, want %v", dl.OrganizedPath, expectedPath)
				}

				// Check file exists
				destFile := filepath.Join(expectedPath, "Standalone Book.m4b")
				if _, err := os.Stat(destFile); os.IsNotExist(err) {
					t.Errorf("destination file does not exist: %s", destFile)
				}
			},
		},
		{
			name: "successful move operation",
			download: &models.Download{
				ID:       "test-3",
				Title:    "Move Test",
				Author:   "Author Name",
				Series:   "Series Name",
				QBitHash: "ghi789",
			},
			configs: map[string]string{
				"paths.destination":        "",
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "move",
			},
			sourceFiles: map[string]string{
				"Move Test.m4b": "move content",
			},
			qbFiles: []*qbittorrent.TorrentFile{
				{
					Name: "Move Test.m4b",
					Path: "",
					Size: 512,
				},
			},
			wantErr: false,
			verifyFn: func(t *testing.T, destBase string, dl *models.Download) {
				expectedPath := filepath.Join(destBase, "Author Name", "Series Name", "Move Test")
				destFile := filepath.Join(expectedPath, "Move Test.m4b")

				// Check file exists at destination
				if _, err := os.Stat(destFile); os.IsNotExist(err) {
					t.Errorf("destination file does not exist: %s", destFile)
				}
			},
		},
		{
			name: "path sanitization with special characters",
			download: &models.Download{
				ID:       "test-4",
				Title:    "Book: With Special? Characters*",
				Author:   "Author/With\\Slashes",
				Series:   `Series"With"Quotes`,
				QBitHash: "jkl012",
			},
			configs: map[string]string{
				"paths.destination":        "",
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "copy",
			},
			sourceFiles: map[string]string{
				"test.m4b": "content",
			},
			qbFiles: []*qbittorrent.TorrentFile{
				{
					Name: "test.m4b",
					Path: "",
					Size: 256,
				},
			},
			wantErr: false,
			verifyFn: func(t *testing.T, destBase string, dl *models.Download) {
				// Special characters should be sanitized (replaced with hyphens)
				expectedPath := filepath.Join(destBase, "Author-With-Slashes", "Series-With-Quotes", "Book- With Special- Characters-")
				if dl.OrganizedPath != expectedPath {
					t.Errorf("OrganizedPath = %v, want %v", dl.OrganizedPath, expectedPath)
				}

				// File should still exist
				destFile := filepath.Join(expectedPath, "test.m4b")
				if _, err := os.Stat(destFile); os.IsNotExist(err) {
					t.Errorf("destination file does not exist: %s", destFile)
				}
			},
		},
		{
			name: "remote path mapping with mount point",
			download: &models.Download{
				ID:       "test-5",
				Title:    "Remote Book",
				Author:   "Remote Author",
				Series:   "Remote Series",
				QBitHash: "mno345",
			},
			configs: map[string]string{
				"paths.destination":        "",
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "copy",
				"paths.local_mount":        "MOUNT_POINT_PLACEHOLDER", // will be set to source dir during test
			},
			sourceFiles: map[string]string{
				"Remote Book.m4b": "remote content",
			},
			qbFiles: []*qbittorrent.TorrentFile{
				{
					Name: "Remote Book.m4b",
					Path: "qbittorrent/downloads/Remote Book.m4b", // Remote path WITHOUT leading slash
					Size: 4096,
				},
			},
			wantErr: false,
			verifyFn: func(t *testing.T, destBase string, dl *models.Download) {
				expectedPath := filepath.Join(destBase, "Remote Author", "Remote Series", "Remote Book")
				destFile := filepath.Join(expectedPath, "Remote Book.m4b")

				if _, err := os.Stat(destFile); os.IsNotExist(err) {
					t.Errorf("destination file does not exist: %s", destFile)
				}
			},
		},
		{
			name: "multiple files in torrent",
			download: &models.Download{
				ID:       "test-6",
				Title:    "Multi File Book",
				Author:   "Multi Author",
				Series:   "Multi Series",
				QBitHash: "pqr678",
			},
			configs: map[string]string{
				"paths.destination":        "",
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "copy",
			},
			sourceFiles: map[string]string{
				"part1.m4b": "part one",
				"part2.m4b": "part two",
				"cover.jpg": "image data",
			},
			qbFiles: []*qbittorrent.TorrentFile{
				{Name: "part1.m4b", Path: "", Size: 1000},
				{Name: "part2.m4b", Path: "", Size: 1000},
				{Name: "cover.jpg", Path: "", Size: 100},
			},
			wantErr: false,
			verifyFn: func(t *testing.T, destBase string, dl *models.Download) {
				expectedPath := filepath.Join(destBase, "Multi Author", "Multi Series", "Multi File Book")

				// All files should exist
				for _, filename := range []string{"part1.m4b", "part2.m4b", "cover.jpg"} {
					destFile := filepath.Join(expectedPath, filename)
					if _, err := os.Stat(destFile); os.IsNotExist(err) {
						t.Errorf("destination file does not exist: %s", destFile)
					}
				}
			},
		},
		{
			name: "error when qBittorrent client fails",
			download: &models.Download{
				ID:       "test-7",
				Title:    "Error Book",
				Author:   "Error Author",
				Series:   "Error Series",
				QBitHash: "error123",
			},
			configs: map[string]string{
				"paths.destination":        "",
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "copy",
			},
			qbError:        fmt.Errorf("qBittorrent connection failed"),
			wantErr:        true,
			wantErrContain: "failed to get torrent files",
		},
		{
			name: "error when source file missing",
			download: &models.Download{
				ID:       "test-8",
				Title:    "Missing File",
				Author:   "Missing Author",
				Series:   "Missing Series",
				QBitHash: "missing123",
			},
			configs: map[string]string{
				"paths.destination":        "",
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "copy",
			},
			sourceFiles: map[string]string{}, // No files created
			qbFiles: []*qbittorrent.TorrentFile{
				{
					Name: "nonexistent.m4b",
					Path: "/tmp/nonexistent.m4b",
					Size: 1024,
				},
			},
			wantErr:        true,
			wantErrContain: "source file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directories
			srcDir := t.TempDir()
			destDir := t.TempDir()

			// Set destination in configs
			tt.configs["paths.destination"] = destDir

			// Create source files
			for filename, content := range tt.sourceFiles {
				srcPath := filepath.Join(srcDir, filename)
				if err := os.WriteFile(srcPath, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}

				// Update qbFiles paths to point to actual source files
				for _, qbFile := range tt.qbFiles {
					if qbFile.Name == filename {
						// For remote path mapping test, keep relative path
						if tt.configs["paths.local_mount"] == "MOUNT_POINT_PLACEHOLDER" {
							// Create subdirectory structure for remote path test
							remoteDir := filepath.Join(srcDir, "qbittorrent", "downloads")
							if err := os.MkdirAll(remoteDir, 0755); err != nil {
								t.Fatalf("failed to create remote directory: %v", err)
							}
							remotePath := filepath.Join(remoteDir, filename)
							if err := os.WriteFile(remotePath, []byte(content), 0644); err != nil {
								t.Fatalf("failed to create remote file: %v", err)
							}
							// Keep the relative path in qbFile
						} else {
							qbFile.Path = srcPath
						}
					}
				}
			}

			// Handle remote path mapping test case
			if tt.configs["paths.local_mount"] == "MOUNT_POINT_PLACEHOLDER" {
				tt.configs["paths.local_mount"] = srcDir
			}

			// Create mocks
			mockConfig := newMockConfigService(tt.configs)
			mockQB := &mockQBClient{
				files: tt.qbFiles,
				err:   tt.qbError,
			}

			// Create service using the mocks
			// We need to work around Go's strict typing by using a test-specific constructor
			svc := newTestOrganizationService(mockQB, mockConfig)

			// Execute
			ctx := context.Background()
			err := svc.Organize(ctx, tt.download)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Organize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrContain != "" {
				if err == nil || !contains(err.Error(), tt.wantErrContain) {
					t.Errorf("Organize() error = %v, want error containing %v", err, tt.wantErrContain)
				}
				return
			}

			// Run verification function
			if !tt.wantErr && tt.verifyFn != nil {
				tt.verifyFn(t, destDir, tt.download)
			}
		})
	}
}

func TestOrganize_NestedDirectoryCreation(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	// Create source file
	srcFile := filepath.Join(srcDir, "test.m4b")
	if err := os.WriteFile(srcFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	download := &models.Download{
		ID:       "nested-test",
		Title:    "Test Book",
		Author:   "Test Author",
		Series:   "Test Series",
		QBitHash: "nested123",
	}

	configs := map[string]string{
		"paths.destination":        destDir,
		"paths.template":           "{author}/{series}/{title}",
		"paths.no_series_template": "{author}/{title}",
		"paths.operation":          "copy",
	}

	mockConfig := newMockConfigService(configs)
	mockQB := &mockQBClient{
		files: []*qbittorrent.TorrentFile{
			{Name: "test.m4b", Path: srcFile, Size: 1024},
		},
	}

	svc := newTestOrganizationService(mockQB, mockConfig)

	ctx := context.Background()
	err := svc.Organize(ctx, download)

	if err != nil {
		t.Fatalf("Organize() failed: %v", err)
	}

	// Verify nested directory structure was created
	expectedPath := filepath.Join(destDir, "Test Author", "Test Series", "Test Book")
	if download.OrganizedPath != expectedPath {
		t.Errorf("OrganizedPath = %v, want %v", download.OrganizedPath, expectedPath)
	}

	// Check all parent directories exist
	if _, err := os.Stat(filepath.Join(destDir, "Test Author")); os.IsNotExist(err) {
		t.Error("author directory not created")
	}
	if _, err := os.Stat(filepath.Join(destDir, "Test Author", "Test Series")); os.IsNotExist(err) {
		t.Error("series directory not created")
	}
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("title directory not created")
	}
}

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name           string
		setupSrc       func(string) error
		wantErr        bool
		wantErrContain string
	}{
		{
			name: "successful copy",
			setupSrc: func(srcPath string) error {
				return os.WriteFile(srcPath, []byte("test content"), 0644)
			},
			wantErr: false,
		},
		{
			name: "source file does not exist",
			setupSrc: func(srcPath string) error {
				// Don't create file
				return nil
			},
			wantErr:        true,
			wantErrContain: "failed to open source file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcDir := t.TempDir()
			destDir := t.TempDir()

			srcPath := filepath.Join(srcDir, "source.txt")
			destPath := filepath.Join(destDir, "dest.txt")

			if err := tt.setupSrc(srcPath); err != nil {
				t.Fatalf("setupSrc failed: %v", err)
			}

			err := copyFile(srcPath, destPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrContain != "" {
				if err == nil || !contains(err.Error(), tt.wantErrContain) {
					t.Errorf("copyFile() error = %v, want error containing %v", err, tt.wantErrContain)
				}
				return
			}

			if !tt.wantErr {
				// Verify content
				content, err := os.ReadFile(destPath)
				if err != nil {
					t.Fatalf("failed to read destination file: %v", err)
				}
				if string(content) != "test content" {
					t.Errorf("content = %v, want 'test content'", string(content))
				}
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
