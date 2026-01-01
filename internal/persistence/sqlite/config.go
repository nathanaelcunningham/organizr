package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) Get(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM configs WHERE key = ?`

	var value string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&value)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("config key not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to query config: %w", err)
	}

	return value, nil
}

func (r *ConfigRepository) Set(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO configs (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = ?
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, key, value, now, value, now)

	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	return nil
}

func (r *ConfigRepository) GetAll(ctx context.Context) (map[string]string, error) {
	query := `SELECT key, value FROM configs`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs: %w", err)
	}
	defer rows.Close()

	configs := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs[key] = value
	}

	return configs, nil
}
