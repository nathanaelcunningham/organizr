package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nathanael/organizr/internal/models"
)

type ProviderRepository struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) Create(ctx context.Context, config *models.ProviderConfig) error {
	configJSON, err := json.Marshal(config.ConfigJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `
		INSERT INTO providers (provider_type, display_name, enabled, config_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err = r.db.ExecContext(ctx, query,
		config.ProviderType,
		config.DisplayName,
		config.Enabled,
		string(configJSON),
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to insert provider: %w", err)
	}

	config.CreatedAt = now
	config.UpdatedAt = now

	return nil
}

func (r *ProviderRepository) GetByType(ctx context.Context, providerType string) (*models.ProviderConfig, error) {
	query := `
		SELECT provider_type, display_name, enabled, config_json, created_at, updated_at
		FROM providers
		WHERE provider_type = ?
	`

	var config models.ProviderConfig
	var configJSON string

	err := r.db.QueryRowContext(ctx, query, providerType).Scan(
		&config.ProviderType,
		&config.DisplayName,
		&config.Enabled,
		&configJSON,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("provider not found: %s", providerType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query provider: %w", err)
	}

	if err := json.Unmarshal([]byte(configJSON), &config.ConfigJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (r *ProviderRepository) List(ctx context.Context) ([]*models.ProviderConfig, error) {
	query := `
		SELECT provider_type, display_name, enabled, config_json, created_at, updated_at
		FROM providers
		ORDER BY display_name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query providers: %w", err)
	}
	defer rows.Close()

	var configs []*models.ProviderConfig
	for rows.Next() {
		var config models.ProviderConfig
		var configJSON string

		if err := rows.Scan(
			&config.ProviderType,
			&config.DisplayName,
			&config.Enabled,
			&configJSON,
			&config.CreatedAt,
			&config.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan provider: %w", err)
		}

		if err := json.Unmarshal([]byte(configJSON), &config.ConfigJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

func (r *ProviderRepository) ListEnabled(ctx context.Context) ([]*models.ProviderConfig, error) {
	query := `
		SELECT provider_type, display_name, enabled, config_json, created_at, updated_at
		FROM providers
		WHERE enabled = 1
		ORDER BY display_name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled providers: %w", err)
	}
	defer rows.Close()

	var configs []*models.ProviderConfig
	for rows.Next() {
		var config models.ProviderConfig
		var configJSON string

		if err := rows.Scan(
			&config.ProviderType,
			&config.DisplayName,
			&config.Enabled,
			&configJSON,
			&config.CreatedAt,
			&config.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan provider: %w", err)
		}

		if err := json.Unmarshal([]byte(configJSON), &config.ConfigJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

func (r *ProviderRepository) Update(ctx context.Context, config *models.ProviderConfig) error {
	configJSON, err := json.Marshal(config.ConfigJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `
		UPDATE providers
		SET display_name = ?, enabled = ?, config_json = ?, updated_at = ?
		WHERE provider_type = ?
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		config.DisplayName,
		config.Enabled,
		string(configJSON),
		now,
		config.ProviderType,
	)

	if err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("provider not found: %s", config.ProviderType)
	}

	config.UpdatedAt = now

	return nil
}

func (r *ProviderRepository) UpdateEnabled(ctx context.Context, providerType string, enabled bool) error {
	query := `UPDATE providers SET enabled = ?, updated_at = ? WHERE provider_type = ?`

	result, err := r.db.ExecContext(ctx, query, enabled, time.Now(), providerType)
	if err != nil {
		return fmt.Errorf("failed to update enabled status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("provider not found: %s", providerType)
	}

	return nil
}

func (r *ProviderRepository) Delete(ctx context.Context, providerType string) error {
	query := `DELETE FROM providers WHERE provider_type = ?`

	result, err := r.db.ExecContext(ctx, query, providerType)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("provider not found: %s", providerType)
	}

	return nil
}

func (r *ProviderRepository) Exists(ctx context.Context, providerType string) (bool, error) {
	query := `SELECT COUNT(*) FROM providers WHERE provider_type = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, providerType).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return count > 0, nil
}
