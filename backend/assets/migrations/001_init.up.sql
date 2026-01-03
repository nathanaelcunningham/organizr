-- Downloads table (includes metadata directly, no cart)
CREATE TABLE IF NOT EXISTS downloads (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    series TEXT,
    torrent_url TEXT,
    magnet_link TEXT,
    qbit_hash TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'queued',
    progress REAL DEFAULT 0.0,
    download_path TEXT,
    organized_path TEXT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    organized_at TIMESTAMP
);

-- Configuration table
CREATE TABLE IF NOT EXISTS configs (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_downloads_status ON downloads(status);
CREATE INDEX IF NOT EXISTS idx_downloads_qbit_hash ON downloads(qbit_hash);
CREATE INDEX IF NOT EXISTS idx_downloads_created_at ON downloads(created_at DESC);

-- Default configuration
INSERT OR IGNORE INTO configs (key, value, description) VALUES
    ('qbittorrent.url', 'http://localhost:8080', 'qBittorrent Web UI URL'),
    ('qbittorrent.username', 'admin', 'qBittorrent username'),
    ('qbittorrent.password', 'adminpass', 'qBittorrent password'),
    ('paths.destination', '/audiobooks', 'Base directory for organized audiobooks'),
    ('paths.template', '{author}/{series}/{title}', 'Path template with series'),
    ('paths.no_series_template', '{author}/{title}', 'Path template without series'),
    ('paths.operation', 'copy', 'File operation: copy or move'),
    ('monitor.interval_seconds', '30', 'Monitor polling interval'),
    ('monitor.auto_organize', 'true', 'Auto-organize on download completion');
