package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nathanael/organizr/internal/models"
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

	download := &models.Download{
		Title:      req.Title,
		Author:     req.Author,
		Series:     req.Series,
		TorrentURL: req.TorrentURL,
		MagnetLink: req.MagnetLink,
		CreatedAt:  time.Now(),
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

	if req.Value == "" {
		respondWithError(w, http.StatusBadRequest, "Config value cannot be empty", nil)
		return
	}

	if err := s.configService.Set(r.Context(), key, req.Value); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update config", err)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	provider := r.URL.Query().Get("provider")

	results, err := s.searchService.Search(r.Context(), query, provider)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Search failed", err)
		return
	}

	respondWithJSON(w, http.StatusOK, SearchResponse{
		Results: searchResultsToDTOList(results),
		Count:   len(results),
	})
}

func (s *Server) handleListProviders(w http.ResponseWriter, r *http.Request) {
	providers := s.searchService.ListActiveProviders()

	respondWithJSON(w, http.StatusOK, ListProvidersResponse{
		Providers: providers,
	})
}

// Provider configuration handlers

func (s *Server) handleListProviderConfigs(w http.ResponseWriter, r *http.Request) {
	configs, err := s.searchService.ListProviders(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list providers", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ListProvidersConfigResponse{
		Providers: providerConfigListToDTO(configs),
	})
}

func (s *Server) handleListProviderTypes(w http.ResponseWriter, r *http.Request) {
	types, err := s.searchService.GetAvailableTypes(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list provider types", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ListProviderTypesResponse{
		Types: providerTypeListToDTO(types),
	})
}

func (s *Server) handleGetProviderConfig(w http.ResponseWriter, r *http.Request) {
	providerType := chi.URLParam(r, "type")
	if err := validateProviderType(providerType); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider type", err)
		return
	}

	config, err := s.searchService.GetProvider(r.Context(), providerType)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Provider not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetProviderConfigResponse{
		Provider: providerConfigToDTO(config),
	})
}

func (s *Server) handleCreateProvider(w http.ResponseWriter, r *http.Request) {
	var req CreateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validateCreateProviderRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	config := &models.ProviderConfig{
		ProviderType: req.ProviderType,
		DisplayName:  req.DisplayName,
		Enabled:      req.Enabled,
		ConfigJSON:   req.Config,
	}

	if err := s.searchService.CreateProvider(r.Context(), config); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create provider", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, GetProviderConfigResponse{
		Provider: providerConfigToDTO(config),
	})
}

func (s *Server) handleUpdateProvider(w http.ResponseWriter, r *http.Request) {
	providerType := chi.URLParam(r, "type")
	if err := validateProviderType(providerType); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider type", err)
		return
	}

	var req UpdateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validateUpdateProviderRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	config := &models.ProviderConfig{
		ProviderType: providerType,
		DisplayName:  req.DisplayName,
		Enabled:      req.Enabled,
		ConfigJSON:   req.Config,
	}

	if err := s.searchService.UpdateProvider(r.Context(), config); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to update provider", err)
		return
	}

	respondWithJSON(w, http.StatusOK, GetProviderConfigResponse{
		Provider: providerConfigToDTO(config),
	})
}

func (s *Server) handleDeleteProvider(w http.ResponseWriter, r *http.Request) {
	providerType := chi.URLParam(r, "type")
	if err := validateProviderType(providerType); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider type", err)
		return
	}

	if err := s.searchService.DeleteProvider(r.Context(), providerType); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete provider", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleToggleProvider(w http.ResponseWriter, r *http.Request) {
	providerType := chi.URLParam(r, "type")
	if err := validateProviderType(providerType); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider type", err)
		return
	}

	var req ToggleProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := s.searchService.ToggleProvider(r.Context(), providerType, req.Enabled); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to toggle provider", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleTestProviderConnection(w http.ResponseWriter, r *http.Request) {
	providerType := chi.URLParam(r, "type")
	if err := validateProviderType(providerType); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider type", err)
		return
	}

	err := s.searchService.TestProvider(r.Context(), providerType)
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
