package server

type CreateDownloadRequest struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Series     string `json:"series"`
	TorrentURL string `json:"torrent_url,omitempty"`
	MagnetLink string `json:"magnet_link,omitempty"`
}

type UpdateConfigRequest struct {
	Value string `json:"value"`
}

type CreateDownloadResponse struct {
	Download downloadDTO `json:"download"`
}

type ListDownloadsResponse struct {
	Downloads []downloadDTO `json:"downloads"`
}

type GetDownloadResponse struct {
	Download downloadDTO `json:"download"`
}

type GetConfigResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetAllConfigResponse struct {
	Configs map[string]string `json:"configs"`
}

type HealthResponse struct {
	Status      string `json:"status"`
	Database    string `json:"database"`
	QBittorrent string `json:"qbittorrent"`
	Monitor     string `json:"monitor"`
}

type SearchResponse struct {
	Results []searchResultDTO `json:"results"`
	Count   int               `json:"count"`
}

type ListProvidersResponse struct {
	Providers []string `json:"providers"`
}

// Provider configuration requests/responses
type CreateProviderRequest struct {
	ProviderType string                 `json:"provider_type"`
	DisplayName  string                 `json:"display_name"`
	Enabled      bool                   `json:"enabled"`
	Config       map[string]interface{} `json:"config"`
}

type UpdateProviderRequest struct {
	DisplayName string                 `json:"display_name"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
}

type ToggleProviderRequest struct {
	Enabled bool `json:"enabled"`
}

type ProviderConfigDTO struct {
	ProviderType string                 `json:"provider_type"`
	DisplayName  string                 `json:"display_name"`
	Enabled      bool                   `json:"enabled"`
	Config       map[string]interface{} `json:"config"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
}

type ProviderTypeDTO struct {
	Type         string                   `json:"type"`
	DisplayName  string                   `json:"display_name"`
	Description  string                   `json:"description"`
	RequiresAuth bool                     `json:"requires_auth"`
	ConfigSchema []ProviderConfigFieldDTO `json:"config_schema"`
}

type ProviderConfigFieldDTO struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

type ListProvidersConfigResponse struct {
	Providers []ProviderConfigDTO `json:"providers"`
}

type GetProviderConfigResponse struct {
	Provider ProviderConfigDTO `json:"provider"`
}

type ListProviderTypesResponse struct {
	Types []ProviderTypeDTO `json:"types"`
}

type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
