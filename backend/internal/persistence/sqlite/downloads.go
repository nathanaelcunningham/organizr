package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nathanael/organizr/internal/models"
)

type DownloadRepository struct {
	db *sql.DB
}

func NewDownloadRepository(db *sql.DB) *DownloadRepository {
	return &DownloadRepository{db: db}
}

func (r *DownloadRepository) Create(ctx context.Context, d *models.Download) error {
	query := `
		INSERT INTO downloads (id, title, author, series, torrent_url, magnet_link, category, qbit_hash, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.Title, d.Author, d.Series, d.TorrentURL, d.MagnetLink, d.Category, d.QBitHash, d.Status, d.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert download: %w", err)
	}

	return nil
}

func (r *DownloadRepository) GetByID(ctx context.Context, id string) (*models.Download, error) {
	query := `
		SELECT id, title, author, series, torrent_url, magnet_link, category, qbit_hash, status, progress,
		       download_path, organized_path, error_message, created_at, completed_at, organized_at
		FROM downloads
		WHERE id = ?
	`

	var d models.Download
	var completedAt, organizedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.Title, &d.Author, &d.Series, &d.TorrentURL, &d.MagnetLink, &d.Category, &d.QBitHash,
		&d.Status, &d.Progress, &d.DownloadPath, &d.OrganizedPath, &d.ErrorMessage,
		&d.CreatedAt, &completedAt, &organizedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("download not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query download: %w", err)
	}

	if completedAt.Valid {
		d.CompletedAt = &completedAt.Time
	}
	if organizedAt.Valid {
		d.OrganizedAt = &organizedAt.Time
	}

	return &d, nil
}

func (r *DownloadRepository) GetActive(ctx context.Context) ([]*models.Download, error) {
	query := `
		SELECT id, title, author, series, qbit_hash, status, progress
		FROM downloads
		WHERE status IN ('queued', 'downloading', 'completed')
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active downloads: %w", err)
	}
	defer rows.Close()

	var downloads []*models.Download
	for rows.Next() {
		var d models.Download
		if err := rows.Scan(&d.ID, &d.Title, &d.Author, &d.Series, &d.QBitHash, &d.Status, &d.Progress); err != nil {
			return nil, fmt.Errorf("failed to scan download: %w", err)
		}
		downloads = append(downloads, &d)
	}

	return downloads, nil
}

func (r *DownloadRepository) List(ctx context.Context) ([]*models.Download, error) {
	query := `
		SELECT id, title, author, series, qbit_hash, status, progress, created_at
		FROM downloads
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query downloads: %w", err)
	}
	defer rows.Close()

	var downloads []*models.Download
	for rows.Next() {
		var d models.Download
		if err := rows.Scan(&d.ID, &d.Title, &d.Author, &d.Series, &d.QBitHash, &d.Status, &d.Progress, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan download: %w", err)
		}
		downloads = append(downloads, &d)
	}

	return downloads, nil
}

func (r *DownloadRepository) UpdateStatus(ctx context.Context, id string, status models.DownloadStatus) error {
	query := `UPDATE downloads SET status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

func (r *DownloadRepository) UpdateProgress(ctx context.Context, id string, progress float64) error {
	query := `UPDATE downloads SET progress = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, progress, id)
	if err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}
	return nil
}

func (r *DownloadRepository) UpdateError(ctx context.Context, id string, errorMsg string) error {
	query := `UPDATE downloads SET error_message = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, errorMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update error: %w", err)
	}
	return nil
}

func (r *DownloadRepository) UpdateOrganizedPath(ctx context.Context, id string, path string) error {
	query := `UPDATE downloads SET organized_path = ?, organized_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, path, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update organized path: %w", err)
	}
	return nil
}

func (r *DownloadRepository) UpdateCompleted(ctx context.Context, id string) error {
	query := `UPDATE downloads SET completed_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update completed time: %w", err)
	}
	return nil
}

func (r *DownloadRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM downloads WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete download: %w", err)
	}
	return nil
}
