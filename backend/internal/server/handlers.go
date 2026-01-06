package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/qbittorrent"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:      "healthy",
		Database:    "ok",
		QBittorrent: "unknown",
		Monitor:     "running",
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (s *Server) handleCreateDownload(w http.ResponseWriter, r *http.Request) {
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
		torrentID, err := strconv.Atoi(req.TorrentID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid torrent ID", err)
			return
		}

		torrentBytes, err = s.searchService.DownloadTorrent(r.Context(), torrentID)
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

	created, err := s.downloadService.CreateDownload(r.Context(), download)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create download", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, CreateDownloadResponse{Download: toDTO(created)})
}

func (s *Server) handleListDownloads(w http.ResponseWriter, r *http.Request) {
	downloads, err := s.downloadService.ListDownloads(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list downloads", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ListDownloadsResponse{Downloads: toDTOList(downloads)})
}

func (s *Server) handleGetDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Download ID is required", nil)
		return
	}

	if err := validateUUID(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid download ID", err)
		return
	}

	download, err := s.downloadService.GetDownload(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Download not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetDownloadResponse{Download: toDTO(download)})
}

func (s *Server) handleCancelDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Download ID is required", nil)
		return
	}

	if err := validateUUID(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid download ID", err)
		return
	}

	if err := s.downloadService.CancelDownload(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to cancel download", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleOrganize(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Download ID is required", nil)
		return
	}

	if err := validateUUID(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid download ID", err)
		return
	}

	if err := s.downloadService.OrganizeDownload(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to organize download", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		respondWithError(w, http.StatusBadRequest, "Config key is required", nil)
		return
	}

	if err := validateConfigKey(key); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid config key", err)
		return
	}

	value, err := s.configService.Get(r.Context(), key)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Config not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetConfigResponse{Key: key, Value: value})
}

func (s *Server) handleGetAllConfig(w http.ResponseWriter, r *http.Request) {
	configs, err := s.configService.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get configuration", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetAllConfigResponse{Configs: configs})
}

func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		respondWithError(w, http.StatusBadRequest, "Config key is required", nil)
		return
	}

	if err := validateConfigKey(key); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid config key", err)
		return
	}

	var req UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Config key '%s' cannot be empty", key), nil)
		return
	}

	if err := s.configService.Set(r.Context(), key, req.Value); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update config", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Query parameter 'q' is required", nil)
		return
	}

	if len(query) < 2 {
		respondWithError(w, http.StatusBadRequest, "Query must be at least 2 characters", nil)
		return
	}

	results, err := s.searchService.Search(r.Context(), query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Search failed", err)
		return
	}

	respondWithJSON(w, http.StatusOK, SearchResponse{
		Results: searchResultsToDTOList(results),
		Count:   len(results),
	})
}

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
