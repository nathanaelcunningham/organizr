-- Add mount point configuration for remote qBittorrent setups
INSERT OR IGNORE INTO configs (key, value, description) VALUES
    ('paths.local_mount', '', 'Local mount point for qBittorrent downloads (e.g., /Volumes/data for network share or Docker volume)');
