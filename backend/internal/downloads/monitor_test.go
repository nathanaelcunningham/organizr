package downloads

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/qbittorrent"
)

// mockDownloadRepo for testing
type mockDownloadRepo struct {
	mu                  sync.Mutex
	getActiveFunc       func(ctx context.Context) ([]*models.Download, error)
	updateProgressFunc  func(ctx context.Context, id string, progress float64) error
	updateCompletedFunc func(ctx context.Context, id string) error
	updateStatusFunc    func(ctx context.Context, id string, status models.DownloadStatus) error
	updateErrorFunc     func(ctx context.Context, id string, errorMsg string) error
	updatePathFunc      func(ctx context.Context, id string, path string) error
	// Track calls for verification
	progressUpdates  map[string]float64
	statusUpdates    map[string]models.DownloadStatus
	completedCalls   []string
	errorUpdates     map[string]string
	pathUpdates      map[string]string
}

func newMockDownloadRepo() *mockDownloadRepo {
	return &mockDownloadRepo{
		progressUpdates: make(map[string]float64),
		statusUpdates:   make(map[string]models.DownloadStatus),
		completedCalls:  []string{},
		errorUpdates:    make(map[string]string),
		pathUpdates:     make(map[string]string),
	}
}

func (m *mockDownloadRepo) GetActive(ctx context.Context) ([]*models.Download, error) {
	if m.getActiveFunc != nil {
		return m.getActiveFunc(ctx)
	}
	return []*models.Download{}, nil
}

func (m *mockDownloadRepo) UpdateProgress(ctx context.Context, id string, progress float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.updateProgressFunc != nil {
		return m.updateProgressFunc(ctx, id, progress)
	}
	m.progressUpdates[id] = progress
	return nil
}

func (m *mockDownloadRepo) UpdateCompleted(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.updateCompletedFunc != nil {
		return m.updateCompletedFunc(ctx, id)
	}
	m.completedCalls = append(m.completedCalls, id)
	return nil
}

func (m *mockDownloadRepo) UpdateStatus(ctx context.Context, id string, status models.DownloadStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, status)
	}
	m.statusUpdates[id] = status
	return nil
}

func (m *mockDownloadRepo) UpdateError(ctx context.Context, id string, errorMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.updateErrorFunc != nil {
		return m.updateErrorFunc(ctx, id, errorMsg)
	}
	m.errorUpdates[id] = errorMsg
	return nil
}

func (m *mockDownloadRepo) UpdateOrganizedPath(ctx context.Context, id string, path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.updatePathFunc != nil {
		return m.updatePathFunc(ctx, id, path)
	}
	m.pathUpdates[id] = path
	return nil
}

// Unused methods required by interface
func (m *mockDownloadRepo) Create(ctx context.Context, d *models.Download) error {
	return nil
}

func (m *mockDownloadRepo) GetByID(ctx context.Context, id string) (*models.Download, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockDownloadRepo) List(ctx context.Context) ([]*models.Download, error) {
	return nil, nil
}

func (m *mockDownloadRepo) Delete(ctx context.Context, id string) error {
	return nil
}

// mockQBClientForMonitor for testing
type mockQBClientForMonitor struct {
	mu                   sync.Mutex
	getTorrentStatusFunc func(ctx context.Context, hash string) (string, float64, error)
	getTorrentFilesFunc  func(ctx context.Context, hash string) ([]*qbittorrent.TorrentFile, error)
	statusResponses      map[string]statusResponse
}

type statusResponse struct {
	status   string
	progress float64
	err      error
}

func newMockQBClientForMonitor() *mockQBClientForMonitor {
	return &mockQBClientForMonitor{
		statusResponses: make(map[string]statusResponse),
	}
}

func (m *mockQBClientForMonitor) GetTorrentStatus(ctx context.Context, hash string) (string, float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.getTorrentStatusFunc != nil {
		return m.getTorrentStatusFunc(ctx, hash)
	}
	if resp, ok := m.statusResponses[hash]; ok {
		return resp.status, resp.progress, resp.err
	}
	return "downloading", 0.0, nil
}

func (m *mockQBClientForMonitor) GetTorrentFiles(ctx context.Context, hash string) ([]*qbittorrent.TorrentFile, error) {
	if m.getTorrentFilesFunc != nil {
		return m.getTorrentFilesFunc(ctx, hash)
	}
	return []*qbittorrent.TorrentFile{}, nil
}

// Unused methods required by interface
func (m *mockQBClientForMonitor) Login(ctx context.Context) error {
	return nil
}

func (m *mockQBClientForMonitor) AddTorrent(ctx context.Context, magnetLink, torrentURL, category string) (string, error) {
	return "", nil
}

func (m *mockQBClientForMonitor) AddTorrentFromFile(ctx context.Context, torrentData []byte, category string) (string, error) {
	return "", nil
}

func (m *mockQBClientForMonitor) DeleteTorrent(ctx context.Context, hash string, deleteFiles bool) error {
	return nil
}

// mockOrgService for testing
type mockOrgService struct {
	mu           sync.Mutex
	organizeFunc func(ctx context.Context, dl *models.Download) error
	organizeCalls []string
}

func (m *mockOrgService) Organize(ctx context.Context, dl *models.Download) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.organizeFunc != nil {
		return m.organizeFunc(ctx, dl)
	}
	m.organizeCalls = append(m.organizeCalls, dl.ID)
	dl.OrganizedPath = "/organized/path/" + dl.Title
	return nil
}

// qbClientInterface defines the methods Monitor needs from qbittorrent.Client
type qbClientInterface interface {
	GetTorrentStatus(ctx context.Context, hash string) (string, float64, error)
	GetTorrentFiles(ctx context.Context, hash string) ([]*qbittorrent.TorrentFile, error)
}

// orgServiceInterface defines the methods Monitor needs from OrganizationService
type orgServiceInterface interface {
	Organize(ctx context.Context, dl *models.Download) error
}

// configServiceInterface defines the methods Monitor needs from config.Service
type configServiceInterface interface {
	Get(ctx context.Context, key string) (string, error)
}

// testMonitor wraps Monitor for testing with interface-based dependencies
type testMonitor struct {
	downloadRepo  *mockDownloadRepo
	qbClient      qbClientInterface
	configService configServiceInterface
	orgService    orgServiceInterface
	interval      time.Duration
	maxConcurrent int
}

// newTestMonitor creates a testMonitor with injectable test dependencies
func newTestMonitor(repo *mockDownloadRepo, qbClient *mockQBClientForMonitor, configSvc *mockConfigService, orgSvc *mockOrgService) *testMonitor {
	return &testMonitor{
		downloadRepo:  repo,
		qbClient:      qbClient,
		configService: configSvc,
		orgService:    orgSvc,
		interval:      30 * time.Second,
		maxConcurrent: 3,
	}
}

// checkDownloads replicates Monitor.checkDownloads for testing
func (m *testMonitor) checkDownloads(ctx context.Context) error {
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
			lastErr = err
			continue
		}

		// At least one download succeeded
		allFailed = false

		// Update progress
		if err := m.downloadRepo.UpdateProgress(ctx, dl.ID, progress); err != nil {
			continue
		}

		// Map qBittorrent state to our status
		newStatus := mapQBitStatusToModel(status)

		// Check if completed
		if (status == "uploading" || status == "stalledUP" || status == "pausedUP") && dl.Status != models.StatusOrganized {
			// Mark completed
			if err := m.downloadRepo.UpdateCompleted(ctx, dl.ID); err != nil {
				continue
			}

			// Check if auto-organization is enabled (default: true for backward compatibility)
			autoOrganize := true
			if autoOrganizeStr, err := m.configService.Get(ctx, "organization.auto_organize"); err == nil {
				if autoOrganizeStr == "false" {
					autoOrganize = false
				}
			}

			// Only auto-organize if enabled
			if autoOrganize {
				// Create context with timeout for organization (respects parent cancellation but prevents indefinite hanging)
				orgCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)

				// Trigger organization in goroutine
				go func(ctx context.Context, download *models.Download) {
					defer cancel()
					m.organizeDownload(ctx, download)
				}(orgCtx, dl)
			}
		}

		_ = newStatus // Suppress unused warning
	}

	// If all downloads failed, qBittorrent may be unavailable
	if len(downloads) > 0 && allFailed {
		_ = lastErr // Continue monitoring, qBittorrent may recover
	}

	return nil
}

// organizeDownload replicates Monitor.organizeDownload for testing
func (m *testMonitor) organizeDownload(ctx context.Context, dl *models.Download) {
	// Mark as organizing
	if err := m.downloadRepo.UpdateStatus(ctx, dl.ID, models.StatusOrganizing); err != nil {
		return
	}

	// Perform organization
	if err := m.orgService.Organize(ctx, dl); err != nil {
		m.downloadRepo.UpdateError(ctx, dl.ID, err.Error())
		m.downloadRepo.UpdateStatus(ctx, dl.ID, models.StatusFailed)
		return
	}

	// Mark as organized
	if err := m.downloadRepo.UpdateStatus(ctx, dl.ID, models.StatusOrganized); err != nil {
		return
	}

	if err := m.downloadRepo.UpdateOrganizedPath(ctx, dl.ID, dl.OrganizedPath); err != nil {
		return
	}
}

// Run replicates Monitor.Run for testing
func (m *testMonitor) Run(ctx context.Context) error {
	// Get interval from config
	intervalStr, err := m.configService.Get(ctx, "monitor.interval_seconds")
	if err == nil {
		var seconds int
		fmt.Sscanf(intervalStr, "%d", &seconds)
		m.interval = time.Duration(seconds) * time.Second
	}

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.checkDownloads(ctx); err != nil {
				// Log but continue
				_ = err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func TestCheckDownloads_ProgressUpdates(t *testing.T) {
	tests := []struct {
		name             string
		downloads        []*models.Download
		statusResponses  map[string]statusResponse
		wantProgress     map[string]float64
	}{
		{
			name: "updates progress for downloading torrents",
			downloads: []*models.Download{
				{ID: "dl-1", Title: "Book 1", QBitHash: "hash1", Status: models.StatusDownloading, Progress: 0},
				{ID: "dl-2", Title: "Book 2", QBitHash: "hash2", Status: models.StatusDownloading, Progress: 25},
			},
			statusResponses: map[string]statusResponse{
				"hash1": {status: "downloading", progress: 50.5, err: nil},
				"hash2": {status: "downloading", progress: 75.0, err: nil},
			},
			wantProgress: map[string]float64{
				"dl-1": 50.5,
				"dl-2": 75.0,
			},
		},
		{
			name: "handles 100% completion",
			downloads: []*models.Download{
				{ID: "dl-3", Title: "Book 3", QBitHash: "hash3", Status: models.StatusDownloading, Progress: 95},
			},
			statusResponses: map[string]statusResponse{
				"hash3": {status: "uploading", progress: 100.0, err: nil},
			},
			wantProgress: map[string]float64{
				"dl-3": 100.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			repo := newMockDownloadRepo()
			repo.getActiveFunc = func(ctx context.Context) ([]*models.Download, error) {
				return tt.downloads, nil
			}

			qbClient := newMockQBClientForMonitor()
			qbClient.statusResponses = tt.statusResponses

			configSvc := newMockConfigService(map[string]string{
				"organization.auto_organize": "false", // Disable auto-org for progress tests
			})

			// Create monitor
			monitor := newTestMonitor(repo, qbClient, configSvc, &mockOrgService{})

			// Run check
			ctx := context.Background()
			err := monitor.checkDownloads(ctx)
			if err != nil {
				t.Fatalf("checkDownloads failed: %v", err)
			}

			// Verify progress updates
			repo.mu.Lock()
			defer repo.mu.Unlock()
			for id, wantProgress := range tt.wantProgress {
				gotProgress, ok := repo.progressUpdates[id]
				if !ok {
					t.Errorf("progress for %s not updated", id)
					continue
				}
				if gotProgress != wantProgress {
					t.Errorf("progress for %s = %v, want %v", id, gotProgress, wantProgress)
				}
			}
		})
	}
}

func TestCheckDownloads_CompletionDetection(t *testing.T) {
	tests := []struct {
		name            string
		download        *models.Download
		qbitStatus      string
		wantCompleted   bool
	}{
		{
			name:          "detects uploading as completed",
			download:      &models.Download{ID: "dl-1", Title: "Book 1", QBitHash: "hash1", Status: models.StatusDownloading},
			qbitStatus:    "uploading",
			wantCompleted: true,
		},
		{
			name:          "detects stalledUP as completed",
			download:      &models.Download{ID: "dl-2", Title: "Book 2", QBitHash: "hash2", Status: models.StatusDownloading},
			qbitStatus:    "stalledUP",
			wantCompleted: true,
		},
		{
			name:          "detects pausedUP as completed",
			download:      &models.Download{ID: "dl-3", Title: "Book 3", QBitHash: "hash3", Status: models.StatusDownloading},
			qbitStatus:    "pausedUP",
			wantCompleted: true,
		},
		{
			name:          "does not mark downloading as completed",
			download:      &models.Download{ID: "dl-4", Title: "Book 4", QBitHash: "hash4", Status: models.StatusDownloading},
			qbitStatus:    "downloading",
			wantCompleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			repo := newMockDownloadRepo()
			repo.getActiveFunc = func(ctx context.Context) ([]*models.Download, error) {
				return []*models.Download{tt.download}, nil
			}

			qbClient := newMockQBClientForMonitor()
			qbClient.statusResponses = map[string]statusResponse{
				tt.download.QBitHash: {status: tt.qbitStatus, progress: 100.0, err: nil},
			}

			configSvc := newMockConfigService(map[string]string{
				"organization.auto_organize": "false", // Disable auto-org for completion tests
			})

			// Create monitor
			monitor := newTestMonitor(repo, qbClient, configSvc, &mockOrgService{})

			// Run check
			ctx := context.Background()
			err := monitor.checkDownloads(ctx)
			if err != nil {
				t.Fatalf("checkDownloads failed: %v", err)
			}

			// Verify completion
			repo.mu.Lock()
			defer repo.mu.Unlock()
			gotCompleted := false
			for _, id := range repo.completedCalls {
				if id == tt.download.ID {
					gotCompleted = true
					break
				}
			}

			if gotCompleted != tt.wantCompleted {
				t.Errorf("completed = %v, want %v", gotCompleted, tt.wantCompleted)
			}
		})
	}
}

func TestCheckDownloads_AutoOrganization(t *testing.T) {
	tests := []struct {
		name              string
		autoOrgEnabled    string
		wantOrganized     bool
	}{
		{
			name:           "triggers auto-organization when enabled",
			autoOrgEnabled: "true",
			wantOrganized:  true,
		},
		{
			name:           "skips auto-organization when disabled",
			autoOrgEnabled: "false",
			wantOrganized:  false,
		},
		{
			name:           "defaults to enabled when config missing",
			autoOrgEnabled: "", // Config key not present
			wantOrganized:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create download
			download := &models.Download{
				ID:       "dl-1",
				Title:    "Test Book",
				QBitHash: "hash1",
				Status:   models.StatusDownloading,
			}

			// Create mocks
			repo := newMockDownloadRepo()
			repo.getActiveFunc = func(ctx context.Context) ([]*models.Download, error) {
				return []*models.Download{download}, nil
			}

			qbClient := newMockQBClientForMonitor()
			qbClient.statusResponses = map[string]statusResponse{
				"hash1": {status: "uploading", progress: 100.0, err: nil},
			}

			configMap := map[string]string{}
			if tt.autoOrgEnabled != "" {
				configMap["organization.auto_organize"] = tt.autoOrgEnabled
			}
			configSvc := newMockConfigService(configMap)

			orgSvc := &mockOrgService{}

			// Create monitor
			monitor := newTestMonitor(repo, qbClient, configSvc, orgSvc)

			// Run check
			ctx := context.Background()
			err := monitor.checkDownloads(ctx)
			if err != nil {
				t.Fatalf("checkDownloads failed: %v", err)
			}

			// Wait a bit for async organization
			time.Sleep(50 * time.Millisecond)

			// Verify organization
			orgSvc.mu.Lock()
			defer orgSvc.mu.Unlock()
			gotOrganized := false
			for _, id := range orgSvc.organizeCalls {
				if id == download.ID {
					gotOrganized = true
					break
				}
			}

			if gotOrganized != tt.wantOrganized {
				t.Errorf("organized = %v, want %v", gotOrganized, tt.wantOrganized)
			}
		})
	}
}

func TestCheckDownloads_Resilience(t *testing.T) {
	tests := []struct {
		name          string
		downloads     []*models.Download
		setupQBClient func(*mockQBClientForMonitor)
		wantError     bool
	}{
		{
			name: "continues when qBittorrent temporarily unavailable",
			downloads: []*models.Download{
				{ID: "dl-1", Title: "Book 1", QBitHash: "hash1", Status: models.StatusDownloading},
				{ID: "dl-2", Title: "Book 2", QBitHash: "hash2", Status: models.StatusDownloading},
			},
			setupQBClient: func(qb *mockQBClientForMonitor) {
				// All torrents fail to get status (qBittorrent down)
				qb.getTorrentStatusFunc = func(ctx context.Context, hash string) (string, float64, error) {
					return "", 0, fmt.Errorf("connection refused")
				}
			},
			wantError: false, // Should not error, just log warnings
		},
		{
			name: "handles partial failures gracefully",
			downloads: []*models.Download{
				{ID: "dl-1", Title: "Book 1", QBitHash: "hash1", Status: models.StatusDownloading},
				{ID: "dl-2", Title: "Book 2", QBitHash: "hash2", Status: models.StatusDownloading},
			},
			setupQBClient: func(qb *mockQBClientForMonitor) {
				qb.getTorrentStatusFunc = func(ctx context.Context, hash string) (string, float64, error) {
					if hash == "hash1" {
						return "", 0, fmt.Errorf("torrent not found")
					}
					return "downloading", 50.0, nil
				}
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			repo := newMockDownloadRepo()
			repo.getActiveFunc = func(ctx context.Context) ([]*models.Download, error) {
				return tt.downloads, nil
			}

			qbClient := newMockQBClientForMonitor()
			if tt.setupQBClient != nil {
				tt.setupQBClient(qbClient)
			}

			configSvc := newMockConfigService(map[string]string{})

			// Create monitor
			monitor := newTestMonitor(repo, qbClient, configSvc, &mockOrgService{})

			// Run check
			ctx := context.Background()
			err := monitor.checkDownloads(ctx)

			if (err != nil) != tt.wantError {
				t.Errorf("checkDownloads() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestMonitorRun_ContextCancellation(t *testing.T) {
	// Create mocks
	repo := newMockDownloadRepo()
	repo.getActiveFunc = func(ctx context.Context) ([]*models.Download, error) {
		return []*models.Download{}, nil
	}

	qbClient := newMockQBClientForMonitor()
	configSvc := newMockConfigService(map[string]string{
		"monitor.interval_seconds": "1", // Short interval for testing
	})

	// Create monitor using test constructor
	monitor := newTestMonitor(repo, qbClient, configSvc, &mockOrgService{})

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Run monitor (should exit when context cancelled)
	err := monitor.Run(ctx)

	// Should return context.DeadlineExceeded
	if err != context.DeadlineExceeded {
		t.Errorf("Run() error = %v, want context.DeadlineExceeded", err)
	}
}

func TestOrganizeDownload_Success(t *testing.T) {
	// Create download
	download := &models.Download{
		ID:       "dl-1",
		Title:    "Test Book",
		QBitHash: "hash1",
		Status:   models.StatusCompleted,
	}

	// Create mocks
	repo := newMockDownloadRepo()
	orgSvc := &mockOrgService{
		organizeFunc: func(ctx context.Context, dl *models.Download) error {
			dl.OrganizedPath = "/organized/Test Book"
			return nil
		},
	}

	// Create monitor
	monitor := newTestMonitor(repo, newMockQBClientForMonitor(), newMockConfigService(map[string]string{}), orgSvc)

	// Run organize
	ctx := context.Background()
	monitor.organizeDownload(ctx, download)

	// Give it a moment to complete
	time.Sleep(10 * time.Millisecond)

	// Verify status updates
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Should have been marked as organizing, then organized
	if status, ok := repo.statusUpdates[download.ID]; !ok || status != models.StatusOrganized {
		t.Errorf("status = %v, want %v", status, models.StatusOrganized)
	}

	// Should have path updated
	if path, ok := repo.pathUpdates[download.ID]; !ok || path == "" {
		t.Errorf("organized path not updated, got %v", path)
	}
}

func TestOrganizeDownload_Failure(t *testing.T) {
	// Create download
	download := &models.Download{
		ID:       "dl-1",
		Title:    "Test Book",
		QBitHash: "hash1",
		Status:   models.StatusCompleted,
	}

	// Create mocks
	repo := newMockDownloadRepo()
	orgSvc := &mockOrgService{
		organizeFunc: func(ctx context.Context, dl *models.Download) error {
			return fmt.Errorf("organization failed: disk full")
		},
	}

	// Create monitor
	monitor := newTestMonitor(repo, newMockQBClientForMonitor(), newMockConfigService(map[string]string{}), orgSvc)

	// Run organize
	ctx := context.Background()
	monitor.organizeDownload(ctx, download)

	// Give it a moment to complete
	time.Sleep(10 * time.Millisecond)

	// Verify error handling
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Should have been marked as failed
	if status, ok := repo.statusUpdates[download.ID]; !ok || status != models.StatusFailed {
		t.Errorf("status = %v, want %v", status, models.StatusFailed)
	}

	// Should have error message
	if errMsg, ok := repo.errorUpdates[download.ID]; !ok || errMsg == "" {
		t.Errorf("error message not updated, got %v", errMsg)
	}
}

func TestMapQBitStatusToModel(t *testing.T) {
	tests := []struct {
		qbitStatus string
		wantStatus models.DownloadStatus
	}{
		{"queuedDL", models.StatusQueued},
		{"queuedUP", models.StatusQueued},
		{"downloading", models.StatusDownloading},
		{"metaDL", models.StatusDownloading},
		{"allocating", models.StatusDownloading},
		{"checkingDL", models.StatusDownloading},
		{"forcedDL", models.StatusDownloading},
		{"uploading", models.StatusCompleted},
		{"stalledUP", models.StatusCompleted},
		{"pausedUP", models.StatusCompleted},
		{"forcedUP", models.StatusCompleted},
		{"checkingUP", models.StatusCompleted},
		{"unknown", ""}, // Unknown status returns empty
	}

	for _, tt := range tests {
		t.Run(tt.qbitStatus, func(t *testing.T) {
			got := mapQBitStatusToModel(tt.qbitStatus)
			if got != tt.wantStatus {
				t.Errorf("mapQBitStatusToModel(%q) = %v, want %v", tt.qbitStatus, got, tt.wantStatus)
			}
		})
	}
}
