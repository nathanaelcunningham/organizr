-- Add path prefix configuration for remote qBittorrent setups
INSERT OR IGNORE INTO configs (key, value, description) VALUES
    ('paths.qbittorrent_prefix', '', 'qBittorrent save path prefix (e.g., /mnt/user/downloads) - stripped from file paths'),
    ('paths.local_mount', '', 'Local mount point where qBittorrent downloads are accessible (e.g., /Volumes/downloads)');
