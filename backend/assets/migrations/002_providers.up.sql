-- Create providers configuration table
CREATE TABLE IF NOT EXISTS providers (
    provider_type TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT 1,
    config_json TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for quick enabled providers lookup
CREATE INDEX IF NOT EXISTS idx_providers_enabled ON providers(enabled);
