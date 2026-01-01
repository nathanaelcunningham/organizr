package models

type SearchResult struct {
	Title      string
	Author     string
	TorrentURL string
	MagnetLink string
	Size       string
	Seeders    int
	Provider   string
}
