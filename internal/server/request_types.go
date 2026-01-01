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
