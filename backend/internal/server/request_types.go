package server

type CreateDownloadRequest struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Series       string `json:"series"`
	SeriesNumber string `json:"series_number,omitempty"`
	TorrentID    string `json:"torrent_id,omitempty"`
	TorrentURL   string `json:"torrent_url,omitempty"`
	MagnetLink   string `json:"magnet_link,omitempty"`
	Category     string `json:"category,omitempty"`
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

type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type PreviewPathRequest struct {
	Template     string `json:"template"`
	Author       string `json:"author"`
	Series       string `json:"series,omitempty"`
	SeriesNumber string `json:"series_number,omitempty"`
	Title        string `json:"title"`
}

type PreviewPathResponse struct {
	Valid bool   `json:"valid"`
	Path  string `json:"path,omitempty"`
	Error string `json:"error,omitempty"`
}

type BatchCreateDownloadRequest struct {
	Downloads []CreateDownloadRequest `json:"downloads"`
}

type BatchDownloadError struct {
	Index   int                   `json:"index"`
	Request CreateDownloadRequest `json:"request"`
	Error   string                `json:"error"`
}

type BatchCreateDownloadResponse struct {
	Successful []downloadDTO        `json:"successful"`
	Failed     []BatchDownloadError `json:"failed"`
}
