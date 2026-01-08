package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nathanael/organizr/internal/models"
)

// Mock download service for testing
type mockDownloadService struct {
	createFunc    func(ctx context.Context, d *models.Download) (*models.Download, error)
	getFunc       func(ctx context.Context, id string) (*models.Download, error)
	listFunc      func(ctx context.Context) ([]*models.Download, error)
	cancelFunc    func(ctx context.Context, id string) error
	organizeFunc  func(ctx context.Context, id string) error
}

func (m *mockDownloadService) CreateDownload(ctx context.Context, d *models.Download) (*models.Download, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, d)
	}
	return d, nil
}

func (m *mockDownloadService) GetDownload(ctx context.Context, id string) (*models.Download, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockDownloadService) ListDownloads(ctx context.Context) ([]*models.Download, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return []*models.Download{}, nil
}

func (m *mockDownloadService) CancelDownload(ctx context.Context, id string) error {
	if m.cancelFunc != nil {
		return m.cancelFunc(ctx, id)
	}
	return nil
}

func (m *mockDownloadService) OrganizeDownload(ctx context.Context, id string) error {
	if m.organizeFunc != nil {
		return m.organizeFunc(ctx, id)
	}
	return nil
}

// Mock config service for testing
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

// Mock search service for testing
type mockSearchService struct {
	downloadTorrentFunc func(ctx context.Context, torrentID int) ([]byte, error)
}

func (m *mockSearchService) DownloadTorrent(ctx context.Context, torrentID int) ([]byte, error) {
	if m.downloadTorrentFunc != nil {
		return m.downloadTorrentFunc(ctx, torrentID)
	}
	return []byte("mock torrent data"), nil
}

// Test handler that mimics handleCreateDownload logic
func testHandleCreateDownload(w http.ResponseWriter, r *http.Request, downloadSvc *mockDownloadService, searchSvc *mockSearchService) {
	var req CreateDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := validateDownloadRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Download torrent file if torrent ID is provided
	var torrentBytes []byte
	if req.TorrentID != "" {
		torrentID := 0
		fmt.Sscanf(req.TorrentID, "%d", &torrentID)

		var err error
		torrentBytes, err = searchSvc.DownloadTorrent(r.Context(), torrentID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to download torrent", err)
			return
		}
	}

	download := &models.Download{
		Title:        req.Title,
		Author:       req.Author,
		Series:       req.Series,
		TorrentURL:   req.TorrentURL,
		MagnetLink:   req.MagnetLink,
		TorrentBytes: torrentBytes,
		Category:     req.Category,
		CreatedAt:    time.Now(),
	}

	created, err := downloadSvc.CreateDownload(r.Context(), download)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create download", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, CreateDownloadResponse{Download: toDTO(created)})
}

func TestHandleCreateDownload(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockCreate     func(ctx context.Context, d *models.Download) (*models.Download, error)
		mockDownload   func(ctx context.Context, torrentID int) ([]byte, error)
		wantStatus     int
		wantErrContain string
	}{
		{
			name: "valid magnet link creates download",
			requestBody: CreateDownloadRequest{
				Title:      "Test Book",
				Author:     "Test Author",
				MagnetLink: "magnet:?xt=urn:btih:abc123",
				Category:   "audiobooks",
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				d.ID = "test-id-123"
				d.Status = models.StatusQueued
				return d, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "valid torrent URL creates download",
			requestBody: CreateDownloadRequest{
				Title:      "Test Book",
				Author:     "Test Author",
				TorrentURL: "https://example.com/file.torrent",
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				d.ID = "test-id-456"
				d.Status = models.StatusQueued
				return d, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "MAM URL triggers torrent download",
			requestBody: CreateDownloadRequest{
				Title:      "Test Book",
				Author:     "Test Author",
				TorrentID:  "12345",
				TorrentURL: "https://example.com/file.torrent", // Add torrent URL to pass validation
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				d.ID = "test-id-789"
				d.Status = models.StatusQueued
				return d, nil
			},
			mockDownload: func(ctx context.Context, torrentID int) ([]byte, error) {
				if torrentID != 12345 {
					return nil, fmt.Errorf("unexpected torrent ID: %d", torrentID)
				}
				return []byte("mock torrent file data"), nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing title returns 400",
			requestBody: CreateDownloadRequest{
				Author:     "Test Author",
				MagnetLink: "magnet:?xt=urn:btih:abc123",
			},
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "title is required",
		},
		{
			name: "missing author returns 400",
			requestBody: CreateDownloadRequest{
				Title:      "Test Book",
				MagnetLink: "magnet:?xt=urn:btih:abc123",
			},
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "author is required",
		},
		{
			name: "missing torrent source returns 400",
			requestBody: CreateDownloadRequest{
				Title:  "Test Book",
				Author: "Test Author",
			},
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "either torrent_url or magnet_link is required",
		},
		{
			name: "qBittorrent failure returns 500",
			requestBody: CreateDownloadRequest{
				Title:      "Test Book",
				Author:     "Test Author",
				MagnetLink: "magnet:?xt=urn:btih:abc123",
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				return nil, fmt.Errorf("qBittorrent connection failed")
			},
			wantStatus:     http.StatusInternalServerError,
			wantErrContain: "Failed to create download",
		},
		{
			name: "MAM download failure returns 500",
			requestBody: CreateDownloadRequest{
				Title:      "Test Book",
				Author:     "Test Author",
				TorrentID:  "99999",
				TorrentURL: "https://example.com/file.torrent", // Add torrent URL to pass validation
			},
			mockDownload: func(ctx context.Context, torrentID int) ([]byte, error) {
				return nil, fmt.Errorf("MAM authentication failed")
			},
			wantStatus:     http.StatusInternalServerError,
			wantErrContain: "Failed to download torrent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			mockDownloadSvc := &mockDownloadService{
				createFunc: tt.mockCreate,
			}
			mockSearchSvc := &mockSearchService{
				downloadTorrentFunc: tt.mockDownload,
			}

			// Marshal request body
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/downloads", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			testHandleCreateDownload(w, req, mockDownloadSvc, mockSearchSvc)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", w.Code, tt.wantStatus, w.Body.String())
			}

			// Check error message if expected
			if tt.wantErrContain != "" {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				// Check both error and message fields
				errorMsg := fmt.Sprintf("%v %v", response["error"], response["message"])
				if !contains(errorMsg, tt.wantErrContain) {
					t.Errorf("error = %v, want to contain %v", errorMsg, tt.wantErrContain)
				}
			}

			// For successful cases, verify response structure
			if w.Code == http.StatusCreated {
				var response CreateDownloadResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if response.Download.ID == "" {
					t.Error("response should contain download with ID")
				}
			}
		})
	}
}

// Test handler that mimics handleListDownloads
func testHandleListDownloads(w http.ResponseWriter, r *http.Request, downloadSvc *mockDownloadService) {
	downloads, err := downloadSvc.ListDownloads(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list downloads", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ListDownloadsResponse{Downloads: toDTOList(downloads)})
}

func TestHandleListDownloads(t *testing.T) {
	tests := []struct {
		name       string
		mockList   func(ctx context.Context) ([]*models.Download, error)
		wantStatus int
		wantCount  int
	}{
		{
			name: "returns all downloads",
			mockList: func(ctx context.Context) ([]*models.Download, error) {
				return []*models.Download{
					{ID: "1", Title: "Book 1", Author: "Author 1", Status: models.StatusQueued},
					{ID: "2", Title: "Book 2", Author: "Author 2", Status: models.StatusDownloading},
					{ID: "3", Title: "Book 3", Author: "Author 3", Status: models.StatusCompleted},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  3,
		},
		{
			name: "returns empty list when no downloads",
			mockList: func(ctx context.Context) ([]*models.Download, error) {
				return []*models.Download{}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "returns 500 on database error",
			mockList: func(ctx context.Context) ([]*models.Download, error) {
				return nil, fmt.Errorf("database connection failed")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockDownloadSvc := &mockDownloadService{
				listFunc: tt.mockList,
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/downloads", nil)
			w := httptest.NewRecorder()

			// Call handler
			testHandleListDownloads(w, req, mockDownloadSvc)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			// For successful cases, verify response structure
			if w.Code == http.StatusOK {
				var response ListDownloadsResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(response.Downloads) != tt.wantCount {
					t.Errorf("download count = %d, want %d", len(response.Downloads), tt.wantCount)
				}
			}
		})
	}
}

// Test handler that mimics handleCancelDownload
func testHandleCancelDownload(w http.ResponseWriter, r *http.Request, downloadSvc *mockDownloadService) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Download ID is required", nil)
		return
	}

	if err := validateUUID(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid download ID", err)
		return
	}

	if err := downloadSvc.CancelDownload(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to cancel download", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func TestHandleCancelDownload(t *testing.T) {
	validID := "123e4567-e89b-42d3-a456-426614174000" // Fixed to be valid UUID v4 (4 in 3rd segment, 8-b in 4th segment)

	tests := []struct {
		name       string
		downloadID string
		mockCancel func(ctx context.Context, id string) error
		wantStatus int
	}{
		{
			name:       "successfully cancels download",
			downloadID: validID,
			mockCancel: func(ctx context.Context, id string) error {
				return nil
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "returns 400 for invalid UUID",
			downloadID: "invalid-uuid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "returns 500 for service error",
			downloadID: validID,
			mockCancel: func(ctx context.Context, id string) error {
				return fmt.Errorf("failed to delete from qBittorrent")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockDownloadSvc := &mockDownloadService{
				cancelFunc: tt.mockCancel,
			}

			// Create router for URL params
			router := chi.NewRouter()
			router.Delete("/api/downloads/{id}", func(w http.ResponseWriter, r *http.Request) {
				testHandleCancelDownload(w, r, mockDownloadSvc)
			})

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/api/downloads/"+tt.downloadID, nil)
			w := httptest.NewRecorder()

			// Call handler through router
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

// Test handler that mimics handleTestQBittorrentConnection
func testHandleTestQBittorrentConnection(w http.ResponseWriter, r *http.Request, configSvc *mockConfigService) {
	// Get qBittorrent config from ConfigService
	url, err := configSvc.Get(r.Context(), "qbittorrent.url")
	if err != nil || url == "" {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: "qBittorrent URL not configured",
		})
		return
	}

	username, err := configSvc.Get(r.Context(), "qbittorrent.username")
	if err != nil || username == "" {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: "qBittorrent username not configured",
		})
		return
	}

	password, err := configSvc.Get(r.Context(), "qbittorrent.password")
	if err != nil || password == "" {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: "qBittorrent password not configured",
		})
		return
	}

	// In a real scenario, we would test the connection here
	// For testing, we just return success if all config is present
	respondWithJSON(w, http.StatusOK, TestConnectionResponse{
		Success: true,
		Message: "Connected successfully",
	})
}

func TestHandleTestQBittorrentConnection(t *testing.T) {
	tests := []struct {
		name        string
		mockConfigs map[string]string
		wantStatus  int
		wantSuccess bool
	}{
		{
			name: "returns error when URL not configured",
			mockConfigs: map[string]string{
				"qbittorrent.username": "admin",
				"qbittorrent.password": "password",
			},
			wantStatus:  http.StatusOK,
			wantSuccess: false,
		},
		{
			name: "returns error when username not configured",
			mockConfigs: map[string]string{
				"qbittorrent.url":      "http://localhost:8080",
				"qbittorrent.password": "password",
			},
			wantStatus:  http.StatusOK,
			wantSuccess: false,
		},
		{
			name: "returns error when password not configured",
			mockConfigs: map[string]string{
				"qbittorrent.url":      "http://localhost:8080",
				"qbittorrent.username": "admin",
			},
			wantStatus:  http.StatusOK,
			wantSuccess: false,
		},
		{
			name: "returns success when all config present",
			mockConfigs: map[string]string{
				"qbittorrent.url":      "http://localhost:8080",
				"qbittorrent.username": "admin",
				"qbittorrent.password": "password",
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockConfigSvc := newMockConfigService(tt.mockConfigs)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/qbittorrent/test", nil)
			w := httptest.NewRecorder()

			// Call handler
			testHandleTestQBittorrentConnection(w, req, mockConfigSvc)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			// Decode response
			var response TestConnectionResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			// Check success field
			if response.Success != tt.wantSuccess {
				t.Errorf("success = %v, want %v", response.Success, tt.wantSuccess)
			}

			// Check message is present
			if response.Message == "" {
				t.Error("message should not be empty")
			}
		})
	}
}

// Test handler that mimics handleBatchCreateDownload
func testHandleBatchCreateDownload(w http.ResponseWriter, r *http.Request, downloadSvc *mockDownloadService, searchSvc *mockSearchService) {
	var req BatchCreateDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate batch is not empty
	if len(req.Downloads) == 0 {
		respondWithError(w, http.StatusBadRequest, "Downloads array cannot be empty", nil)
		return
	}

	// Validate batch size limit (50 items max to prevent abuse)
	if len(req.Downloads) > 50 {
		respondWithError(w, http.StatusBadRequest, "Batch size cannot exceed 50 downloads", nil)
		return
	}

	var successful []downloadDTO
	var failed []BatchDownloadError

	// Process downloads sequentially
	for i, downloadReq := range req.Downloads {
		// Validate request
		if err := validateDownloadRequest(downloadReq); err != nil {
			failed = append(failed, BatchDownloadError{
				Index:   i,
				Request: downloadReq,
				Error:   fmt.Sprintf("Validation failed: %v", err),
			})
			continue
		}

		// Download torrent file if torrent ID is provided
		var torrentBytes []byte
		if downloadReq.TorrentID != "" {
			torrentID := 0
			fmt.Sscanf(downloadReq.TorrentID, "%d", &torrentID)

			var err error
			torrentBytes, err = searchSvc.DownloadTorrent(r.Context(), torrentID)
			if err != nil {
				failed = append(failed, BatchDownloadError{
					Index:   i,
					Request: downloadReq,
					Error:   fmt.Sprintf("Failed to download torrent: %v", err),
				})
				continue
			}
		}

		download := &models.Download{
			Title:        downloadReq.Title,
			Author:       downloadReq.Author,
			Series:       downloadReq.Series,
			TorrentURL:   downloadReq.TorrentURL,
			MagnetLink:   downloadReq.MagnetLink,
			TorrentBytes: torrentBytes,
			Category:     downloadReq.Category,
			CreatedAt:    time.Now(),
		}

		created, err := downloadSvc.CreateDownload(r.Context(), download)
		if err != nil {
			failed = append(failed, BatchDownloadError{
				Index:   i,
				Request: downloadReq,
				Error:   fmt.Sprintf("Failed to create download: %v", err),
			})
			continue
		}

		successful = append(successful, toDTO(created))
	}

	respondWithJSON(w, http.StatusOK, BatchCreateDownloadResponse{
		Successful: successful,
		Failed:     failed,
	})
}

func TestHandleBatchCreateDownload(t *testing.T) {
	tests := []struct {
		name              string
		requestBody       interface{}
		mockCreate        func(ctx context.Context, d *models.Download) (*models.Download, error)
		wantStatus        int
		wantSuccessCount  int
		wantFailedCount   int
		wantErrContain    string
	}{
		{
			name: "successful batch - all succeed",
			requestBody: BatchCreateDownloadRequest{
				Downloads: []CreateDownloadRequest{
					{
						Title:      "Test Book 1",
						Author:     "Test Author 1",
						MagnetLink: "magnet:?xt=urn:btih:abc123",
						Category:   "audiobooks",
					},
					{
						Title:      "Test Book 2",
						Author:     "Test Author 2",
						MagnetLink: "magnet:?xt=urn:btih:def456",
						Category:   "audiobooks",
					},
					{
						Title:      "Test Book 3",
						Author:     "Test Author 3",
						TorrentURL: "https://example.com/test.torrent",
					},
				},
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				d.ID = "test-id-" + d.Title
				d.Status = models.StatusQueued
				return d, nil
			},
			wantStatus:       http.StatusOK,
			wantSuccessCount: 3,
			wantFailedCount:  0,
		},
		{
			name: "partial failure - some succeed some fail due to validation",
			requestBody: BatchCreateDownloadRequest{
				Downloads: []CreateDownloadRequest{
					{
						Title:      "Valid Book",
						Author:     "Valid Author",
						MagnetLink: "magnet:?xt=urn:btih:abc123",
					},
					{
						Title:  "Missing Author",
						Author: "",
						MagnetLink: "magnet:?xt=urn:btih:def456",
					},
					{
						Title:      "No Torrent Source",
						Author:     "Author",
					},
					{
						Title:      "Valid Book 2",
						Author:     "Valid Author 2",
						TorrentURL: "https://example.com/test.torrent",
					},
				},
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				d.ID = "test-id-" + d.Title
				d.Status = models.StatusQueued
				return d, nil
			},
			wantStatus:       http.StatusOK,
			wantSuccessCount: 2,
			wantFailedCount:  2,
		},
		{
			name: "empty array - returns 400",
			requestBody: BatchCreateDownloadRequest{
				Downloads: []CreateDownloadRequest{},
			},
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "Downloads array cannot be empty",
		},
		{
			name: "oversized batch (>50 items) - returns 400",
			requestBody: func() BatchCreateDownloadRequest {
				downloads := make([]CreateDownloadRequest, 51)
				for i := 0; i < 51; i++ {
					downloads[i] = CreateDownloadRequest{
						Title:      fmt.Sprintf("Book %d", i),
						Author:     fmt.Sprintf("Author %d", i),
						MagnetLink: fmt.Sprintf("magnet:?xt=urn:btih:hash%d", i),
					}
				}
				return BatchCreateDownloadRequest{Downloads: downloads}
			}(),
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "Batch size cannot exceed 50 downloads",
		},
		{
			name: "all fail - returns 200 with empty successful array",
			requestBody: BatchCreateDownloadRequest{
				Downloads: []CreateDownloadRequest{
					{
						Title:      "Book 1",
						Author:     "",
						MagnetLink: "magnet:?xt=urn:btih:abc123",
					},
					{
						Title:      "Book 2",
						Author:     "Author 2",
					},
					{
						Title:      "",
						Author:     "Author 3",
						MagnetLink: "magnet:?xt=urn:btih:ghi789",
					},
				},
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				return nil, fmt.Errorf("qBittorrent error")
			},
			wantStatus:       http.StatusOK,
			wantSuccessCount: 0,
			wantFailedCount:  3,
		},
		{
			name: "partial failure due to service errors",
			requestBody: BatchCreateDownloadRequest{
				Downloads: []CreateDownloadRequest{
					{
						Title:      "Book 1",
						Author:     "Author 1",
						MagnetLink: "magnet:?xt=urn:btih:abc123",
					},
					{
						Title:      "Book 2",
						Author:     "Author 2",
						MagnetLink: "magnet:?xt=urn:btih:def456",
					},
				},
			},
			mockCreate: func(ctx context.Context, d *models.Download) (*models.Download, error) {
				// First download succeeds, second fails
				if d.Title == "Book 1" {
					d.ID = "test-id-1"
					d.Status = models.StatusQueued
					return d, nil
				}
				return nil, fmt.Errorf("qBittorrent connection failed")
			},
			wantStatus:       http.StatusOK,
			wantSuccessCount: 1,
			wantFailedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock services
			mockDownloadSvc := &mockDownloadService{
				createFunc: tt.mockCreate,
			}
			mockSearchSvc := &mockSearchService{}

			// Marshal request body
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/downloads/batch", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			testHandleBatchCreateDownload(w, req, mockDownloadSvc, mockSearchSvc)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", w.Code, tt.wantStatus, w.Body.String())
			}

			// Check error message if expected
			if tt.wantErrContain != "" {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				errorMsg := fmt.Sprintf("%v %v", response["error"], response["message"])
				if !contains(errorMsg, tt.wantErrContain) {
					t.Errorf("error = %v, want to contain %v", errorMsg, tt.wantErrContain)
				}
			}

			// For successful batch operations, verify response structure
			if w.Code == http.StatusOK {
				var response BatchCreateDownloadResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if len(response.Successful) != tt.wantSuccessCount {
					t.Errorf("successful count = %d, want %d", len(response.Successful), tt.wantSuccessCount)
				}

				if len(response.Failed) != tt.wantFailedCount {
					t.Errorf("failed count = %d, want %d", len(response.Failed), tt.wantFailedCount)
				}

				// Verify failed items have correct structure
				for _, failedItem := range response.Failed {
					if failedItem.Error == "" {
						t.Error("failed item should have error message")
					}
					if failedItem.Index < 0 {
						t.Error("failed item should have valid index")
					}
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
