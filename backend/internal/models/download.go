package models

import "time"

type Download struct {
	ID            string
	Title         string
	Author        string
	Series        string
	SeriesNumber  string
	TorrentURL    string
	MagnetLink    string
	TorrentBytes  []byte
	Category      string
	QBitHash      string
	Status        DownloadStatus
	Progress      float64
	DownloadPath  string
	OrganizedPath string
	ErrorMessage  string
	CreatedAt     time.Time
	CompletedAt   *time.Time
	OrganizedAt   *time.Time
}

type DownloadStatus string

const (
	StatusQueued      DownloadStatus = "queued"
	StatusDownloading DownloadStatus = "downloading"
	StatusCompleted   DownloadStatus = "completed"
	StatusOrganizing  DownloadStatus = "organizing"
	StatusOrganized   DownloadStatus = "organized"
	StatusFailed      DownloadStatus = "failed"
)
