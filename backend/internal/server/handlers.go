package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nathanael/organizr/internal/fileutil"
	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/qbittorrent"
)

// handleHealth godoc
// @Summary Health check
// @Description Check API and dependency health status
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:      "healthy",
		Database:    "ok",
		QBittorrent: "unknown",
		Monitor:     "running",
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// handleCreateDownload godoc
// @Summary Create a new download
// @Description Create a new audiobook download from torrent URL, magnet link, or torrent ID
// @Tags downloads
// @Accept json
// @Produce json
// @Param request body CreateDownloadRequest true "Download request"
// @Success 201 {object} CreateDownloadResponse
// @Failure 400 {object} ErrorResponse "Invalid request body or validation failed"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /downloads [post]
func (s *Server) handleCreateDownload(w http.ResponseWriter, r *http.Request) {
	var req CreateDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithBadRequest(w, "invalid request body", err)
		return
	}

	// Validate request
	if err := validateDownloadRequest(req); err != nil {
		respondWithBadRequest(w, "validation failed", err)
		return
	}

	// Download torrent file if torrent ID is provided
	var torrentBytes []byte
	if req.TorrentID != "" {
		torrentID, err := strconv.Atoi(req.TorrentID)
		if err != nil {
			respondWithValidationError(w, "torrent ID", err)
			return
		}

		torrentBytes, err = s.searchService.DownloadTorrent(r.Context(), torrentID)
		if err != nil {
			respondWithInternalError(w, "download torrent", err)
			return
		}
	}

	download := &models.Download{
		Title:        req.Title,
		Author:       req.Author,
		Series:       req.Series,
		SeriesNumber: req.SeriesNumber,
		TorrentURL:   req.TorrentURL,
		MagnetLink:   req.MagnetLink,
		TorrentBytes: torrentBytes,
		Category:     req.Category,
		CreatedAt:    time.Now(),
	}

	created, err := s.downloadService.CreateDownload(r.Context(), download)
	if err != nil {
		respondWithInternalError(w, "create download", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, CreateDownloadResponse{Download: toDTO(created)})
}

// handleListDownloads godoc
// @Summary List all downloads
// @Description Get a list of all downloads with their status and progress
// @Tags downloads
// @Produce json
// @Success 200 {object} ListDownloadsResponse
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /downloads [get]
func (s *Server) handleListDownloads(w http.ResponseWriter, r *http.Request) {
	downloads, err := s.downloadService.ListDownloads(r.Context())
	if err != nil {
		respondWithInternalError(w, "list downloads", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ListDownloadsResponse{Downloads: toDTOList(downloads)})
}

// handleGetDownload godoc
// @Summary Get a specific download
// @Description Get detailed information about a specific download by ID
// @Tags downloads
// @Produce json
// @Param id path string true "Download ID (UUID)"
// @Success 200 {object} GetDownloadResponse
// @Failure 400 {object} ErrorResponse "Invalid download ID"
// @Failure 404 {object} ErrorResponse "Download not found"
// @Router /downloads/{id} [get]
func (s *Server) handleGetDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithValidationError(w, "download ID", nil)
		return
	}

	if err := validateUUID(id); err != nil {
		respondWithValidationError(w, "download ID", err)
		return
	}

	download, err := s.downloadService.GetDownload(r.Context(), id)
	if err != nil {
		respondWithNotFound(w, "download", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetDownloadResponse{Download: toDTO(download)})
}

// handleCancelDownload godoc
// @Summary Cancel a download
// @Description Cancel an active download and remove it from qBittorrent
// @Tags downloads
// @Param id path string true "Download ID (UUID)"
// @Success 204 "Download cancelled successfully"
// @Failure 400 {object} ErrorResponse "Invalid download ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /downloads/{id} [delete]
func (s *Server) handleCancelDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithValidationError(w, "download ID", nil)
		return
	}

	if err := validateUUID(id); err != nil {
		respondWithValidationError(w, "download ID", err)
		return
	}

	if err := s.downloadService.CancelDownload(r.Context(), id); err != nil {
		respondWithInternalError(w, "cancel download", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleOrganize godoc
// @Summary Organize a completed download
// @Description Move and organize a completed download to the configured destination path
// @Tags downloads
// @Param id path string true "Download ID (UUID)"
// @Success 200 "Download organized successfully"
// @Failure 400 {object} ErrorResponse "Invalid download ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /downloads/{id}/organize [post]
func (s *Server) handleOrganize(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithValidationError(w, "download ID", nil)
		return
	}

	if err := validateUUID(id); err != nil {
		respondWithValidationError(w, "download ID", err)
		return
	}

	if err := s.downloadService.OrganizeDownload(r.Context(), id); err != nil {
		respondWithInternalError(w, "organize download", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleGetConfig godoc
// @Summary Get a configuration value
// @Description Get the value of a specific configuration key
// @Tags config
// @Produce json
// @Param key path string true "Configuration key"
// @Success 200 {object} GetConfigResponse
// @Failure 400 {object} ErrorResponse "Invalid config key"
// @Failure 404 {object} ErrorResponse "Config key not found"
// @Router /config/{key} [get]
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		respondWithValidationError(w, "config key", nil)
		return
	}

	if err := validateConfigKey(key); err != nil {
		respondWithValidationError(w, "config key", err)
		return
	}

	value, err := s.configService.Get(r.Context(), key)
	if err != nil {
		respondWithNotFound(w, "config", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetConfigResponse{Key: key, Value: value})
}

// handleGetAllConfig godoc
// @Summary Get all configuration values
// @Description Get all configuration key-value pairs
// @Tags config
// @Produce json
// @Success 200 {object} GetAllConfigResponse
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /config [get]
func (s *Server) handleGetAllConfig(w http.ResponseWriter, r *http.Request) {
	configs, err := s.configService.GetAll(r.Context())
	if err != nil {
		respondWithInternalError(w, "get configuration", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetAllConfigResponse{Configs: configs})
}

// handleUpdateConfig godoc
// @Summary Update a configuration value
// @Description Update the value of a specific configuration key
// @Tags config
// @Accept json
// @Param key path string true "Configuration key"
// @Param request body UpdateConfigRequest true "New configuration value"
// @Success 204 "Configuration updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or validation failed"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /config/{key} [put]
func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		respondWithValidationError(w, "config key", nil)
		return
	}

	if err := validateConfigKey(key); err != nil {
		respondWithValidationError(w, "config key", err)
		return
	}

	var req UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithBadRequest(w, "invalid request body", err)
		return
	}

	// Required config keys that cannot be empty
	requiredKeys := map[string]bool{
		"qbittorrent.url":          true,
		"qbittorrent.username":     true,
		"paths.destination":        true,
		"paths.template":           true,
		"paths.no_series_template": true,
		"paths.operation":          true,
		"monitor.interval_seconds": true,
		"monitor.auto_organize":    true,
		"mam.baseurl":              true,
	}

	if req.Value == "" && requiredKeys[key] {
		respondWithBadRequest(w, fmt.Sprintf("config key '%s' cannot be empty", key), nil)
		return
	}

	if err := s.configService.Set(r.Context(), key, req.Value); err != nil {
		respondWithInternalError(w, "update config", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleSearch godoc
// @Summary Search for audiobooks
// @Description Search for audiobooks on configured providers (e.g., MyAnonamouse)
// @Tags search
// @Produce json
// @Param q query string true "Search query" minlength(2)
// @Success 200 {object} SearchResponse
// @Failure 400 {object} ErrorResponse "Invalid query parameter"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /search [get]
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithValidationError(w, "query parameter 'q'", nil)
		return
	}

	if len(query) < 2 {
		respondWithValidationError(w, "query length", nil)
		return
	}

	results, err := s.searchService.Search(r.Context(), query)
	if err != nil {
		respondWithInternalError(w, "search", err)
		return
	}

	respondWithJSON(w, http.StatusOK, SearchResponse{
		Results: searchResultsToDTOList(results),
		Count:   len(results),
	})
}

// handleTestConnection godoc
// @Summary Test search provider connection
// @Description Test connectivity to the configured search provider (e.g., MyAnonamouse)
// @Tags search
// @Produce json
// @Success 200 {object} TestConnectionResponse
// @Router /search/test [post]
func (s *Server) handleTestConnection(w http.ResponseWriter, r *http.Request) {
	err := s.searchService.TestConnection(r.Context())
	if err != nil {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, TestConnectionResponse{
		Success: true,
		Message: "Connection successful",
	})
}

// handleTestQBittorrentConnection godoc
// @Summary Test qBittorrent connection
// @Description Test connectivity to the configured qBittorrent instance
// @Tags qbittorrent
// @Produce json
// @Success 200 {object} TestConnectionResponse
// @Router /qbittorrent/test [get]
func (s *Server) handleTestQBittorrentConnection(w http.ResponseWriter, r *http.Request) {
	// Get qBittorrent config from ConfigService
	url, err := s.configService.Get(r.Context(), "qbittorrent.url")
	if err != nil || url == "" {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: "qBittorrent URL not configured",
		})
		return
	}

	username, err := s.configService.Get(r.Context(), "qbittorrent.username")
	if err != nil || username == "" {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: "qBittorrent username not configured",
		})
		return
	}

	password, err := s.configService.Get(r.Context(), "qbittorrent.password")
	if err != nil || password == "" {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: "qBittorrent password not configured",
		})
		return
	}

	// Create qBittorrent client and test connection
	client, err := qbittorrent.NewClient(url, username, password)
	if err != nil {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create client: %v", err),
		})
		return
	}
	if err := client.Login(r.Context()); err != nil {
		respondWithJSON(w, http.StatusOK, TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Connection failed: %v", err),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, TestConnectionResponse{
		Success: true,
		Message: "Connected successfully",
	})
}

// handlePreviewPath godoc
// @Summary Preview path template
// @Description Preview the result of applying a path template with given audiobook metadata
// @Tags config
// @Accept json
// @Produce json
// @Param request body PreviewPathRequest true "Path preview request"
// @Success 200 {object} PreviewPathResponse
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Router /config/preview-path [post]
func (s *Server) handlePreviewPath(w http.ResponseWriter, r *http.Request) {
	var req PreviewPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithBadRequest(w, "invalid request body", err)
		return
	}

	allowedVars := []string{"author", "series", "series_number", "title"}

	// Validate template
	if err := fileutil.ValidateTemplate(req.Template, allowedVars); err != nil {
		respondWithJSON(w, http.StatusOK, PreviewPathResponse{
			Valid: false,
			Error: err.Error(),
		})
		return
	}

	// Sanitize individual variables before parsing template (preserves directory structure)
	vars := map[string]string{
		"author":        fileutil.SanitizePath(req.Author),
		"series":        fileutil.SanitizePath(req.Series),
		"series_number": fileutil.SanitizePath(req.SeriesNumber),
		"title":         fileutil.SanitizePath(req.Title),
	}

	// Parse template with sanitized values (directory separators preserved)
	path := fileutil.ParseTemplate(req.Template, vars)

	respondWithJSON(w, http.StatusOK, PreviewPathResponse{
		Valid: true,
		Path:  path,
	})
}

// handleBatchCreateDownload godoc
// @Summary Create multiple downloads in batch
// @Description Create multiple audiobook downloads in a single request (max 50 items)
// @Tags downloads
// @Accept json
// @Produce json
// @Param request body BatchCreateDownloadRequest true "Batch download request"
// @Success 200 {object} BatchCreateDownloadResponse "Returns successful and failed downloads"
// @Failure 400 {object} ErrorResponse "Invalid request body or batch size exceeded"
// @Router /downloads/batch [post]
func (s *Server) handleBatchCreateDownload(w http.ResponseWriter, r *http.Request) {
	var req BatchCreateDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithBadRequest(w, "invalid request body", err)
		return
	}

	// Validate batch is not empty
	if len(req.Downloads) == 0 {
		respondWithValidationError(w, "downloads array", nil)
		return
	}

	// Validate batch size limit (50 items max to prevent abuse)
	if len(req.Downloads) > 50 {
		respondWithBadRequest(w, "batch size exceeds 50 item limit", nil)
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
			torrentID, err := strconv.Atoi(downloadReq.TorrentID)
			if err != nil {
				failed = append(failed, BatchDownloadError{
					Index:   i,
					Request: downloadReq,
					Error:   fmt.Sprintf("Invalid torrent ID: %v", err),
				})
				continue
			}

			torrentBytes, err = s.searchService.DownloadTorrent(r.Context(), torrentID)
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
			SeriesNumber: downloadReq.SeriesNumber,
			TorrentURL:   downloadReq.TorrentURL,
			MagnetLink:   downloadReq.MagnetLink,
			TorrentBytes: torrentBytes,
			Category:     downloadReq.Category,
			CreatedAt:    time.Now(),
		}

		created, err := s.downloadService.CreateDownload(r.Context(), download)
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

	// Log batch results
	fmt.Printf("Batch download processed: %d successful, %d failed (total: %d)\n",
		len(successful), len(failed), len(req.Downloads))

	respondWithJSON(w, http.StatusOK, BatchCreateDownloadResponse{
		Successful: successful,
		Failed:     failed,
	})
}
