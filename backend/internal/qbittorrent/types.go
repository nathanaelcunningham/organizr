package qbittorrent

type TorrentInfo struct {
	Hash       string  `json:"hash"`
	Name       string  `json:"name"`
	State      string  `json:"state"`
	Progress   float64 `json:"progress"`
	SavePath   string  `json:"save_path"`
	Downloaded int64   `json:"downloaded"`
	Size       int64   `json:"size"`
	AddedOn    int64   `json:"added_on"`
}

type TorrentFile struct {
	Name string
	Path string
	Size int64
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
