package models

type SeriesInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

type SearchResult struct {
	// Basic info
	ID         string
	Title      string
	Author     string
	TorrentURL string
	MagnetLink string
	Provider   string

	// Metadata
	Series      []SeriesInfo
	Narrator    string
	Category    string
	FileType    string
	Language    string
	Tags        []string
	Description string
	Added       string

	// Torrent stats
	Size           string
	Seeders        int
	Leechers       int
	NumFiles       int
	TimesCompleted int

	// Special flags
	Freeleech    bool
	FreeleechVIP bool
	VIP          bool
}
