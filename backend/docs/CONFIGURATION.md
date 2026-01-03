# Configuration Guide

Organizr stores all configuration in the SQLite database. This document explains how to configure the service.

## Initial Setup

On first run, Organizr creates a database with default configuration values. You can update these via the REST API.

## Configuration Options

### qBittorrent Connection

Configure connection to your qBittorrent instance:

```bash
# Set qBittorrent URL
curl -X PUT http://localhost:8080/api/config/qbittorrent.url \
  -H "Content-Type: application/json" \
  -d '{"value": "http://192.168.1.100:8080"}'

# Set username
curl -X PUT http://localhost:8080/api/config/qbittorrent.username \
  -H "Content-Type: application/json" \
  -d '{"value": "your-username"}'

# Set password
curl -X PUT http://localhost:8080/api/config/qbittorrent.password \
  -H "Content-Type: application/json" \
  -d '{"value": "your-password"}'
```

### File Organization Paths

Configure where and how files are organized:

```bash
# Set destination directory
curl -X PUT http://localhost:8080/api/config/paths.destination \
  -H "Content-Type: application/json" \
  -d '{"value": "/mnt/nas/audiobooks"}'

# Set path template for books with series
curl -X PUT http://localhost:8080/api/config/paths.template \
  -H "Content-Type: application/json" \
  -d '{"value": "{author}/{series}/{title}"}'

# Set path template for books without series
curl -X PUT http://localhost:8080/api/config/paths.no_series_template \
  -H "Content-Type: application/json" \
  -d '{"value": "{author}/{title}"}'

# Set file operation (copy or move)
curl -X PUT http://localhost:8080/api/config/paths.operation \
  -H "Content-Type: application/json" \
  -d '{"value": "copy"}'
```

**Path Template Variables:**
- `{author}` - Book author name
- `{series}` - Series name (if provided)
- `{title}` - Book title

**Example Results:**

With series (Template: `{author}/{series}/{title}`):
```
/audiobooks/Stephen King/The Dark Tower/The Gunslinger/
```

Without series (Template: `{author}/{title}`):
```
/audiobooks/Stephen King/The Stand/
```

### File Operations

Choose between copying or moving files:

- **`copy`** (default): Leaves original files in qBittorrent download directory
  - Pros: Keeps seeding, safer
  - Cons: Uses more disk space

- **`move`**: Moves files from qBittorrent to organized location
  - Pros: Saves disk space
  - Cons: Stops seeding, can't recover if organization fails

### Monitor Settings

Configure the background monitor that watches for completed downloads:

```bash
# Set poll interval (seconds)
curl -X PUT http://localhost:8080/api/config/monitor.interval_seconds \
  -H "Content-Type: application/json" \
  -d '{"value": "60"}'

# Enable/disable auto-organization
curl -X PUT http://localhost:8080/api/config/monitor.auto_organize \
  -H "Content-Type: application/json" \
  -d '{"value": "true"}'
```

**Recommendations:**
- `interval_seconds`: 30-60 seconds for most use cases
- `auto_organize`: Keep as `true` unless you want manual control

## Viewing Current Configuration

Get all configuration:

```bash
curl http://localhost:8080/api/config
```

Get specific value:

```bash
curl http://localhost:8080/api/config/paths.destination
```

## Environment-Specific Configurations

### Home Server

```bash
# Local qBittorrent instance
curl -X PUT http://localhost:8080/api/config/qbittorrent.url \
  -d '{"value": "http://localhost:8080"}'

# Local storage
curl -X PUT http://localhost:8080/api/config/paths.destination \
  -d '{"value": "/home/user/audiobooks"}'

# Copy to preserve seeding
curl -X PUT http://localhost:8080/api/config/paths.operation \
  -d '{"value": "copy"}'
```

### NAS Setup

```bash
# Remote qBittorrent
curl -X PUT http://localhost:8080/api/config/qbittorrent.url \
  -d '{"value": "http://nas.local:8080"}'

# NAS mount point
curl -X PUT http://localhost:8080/api/config/paths.destination \
  -d '{"value": "/mnt/nas/media/audiobooks"}'

# Move to save space
curl -X PUT http://localhost:8080/api/config/paths.operation \
  -d '{"value": "move"}'
```

### Docker Setup

When running in Docker, ensure:
1. qBittorrent URL is accessible from container
2. Destination path is mounted as volume
3. qBittorrent download path is accessible (for file operations)

```bash
# Example Docker run
docker run -d \
  -p 8080:8080 \
  -v /path/to/audiobooks:/audiobooks \
  -v /path/to/downloads:/downloads \
  organizr:latest
```

## Security Considerations

1. **qBittorrent Password**: Store securely, consider using environment variables in production
2. **Network**: Restrict qBittorrent Web UI access to trusted networks
3. **File Permissions**: Ensure Organizr has appropriate read/write permissions
4. **API Access**: Add authentication middleware for production use

## Backup

The entire configuration is stored in `organizr.db`. Back up this file to preserve:
- All configuration settings
- Download history and status
- Organized file paths

```bash
# Backup database
cp organizr.db organizr.db.backup

# Restore from backup
cp organizr.db.backup organizr.db
```

## Troubleshooting

### Configuration not taking effect

1. Verify configuration was saved:
   ```bash
   curl http://localhost:8080/api/config/your-key
   ```

2. Check logs for errors

3. Restart the service (some changes may require restart)

### Invalid path templates

Path templates must:
- Use valid placeholders: `{author}`, `{series}`, `{title}`
- Not include invalid filesystem characters (these are auto-sanitized)
- Not start or end with `/`

Valid:
- `{author}/{series}/{title}`
- `{author} - {title}`
- `Audiobooks/{author}/{title}`

Invalid:
- `/{author}/{title}` (starts with /)
- `{author}/{series}` (missing title)
- `{invalid}` (unknown placeholder)

### Permission errors

Ensure the user running Organizr has:
- Read access to qBittorrent download directory
- Write access to destination directory

```bash
# Check permissions
ls -la /path/to/audiobooks

# Fix if needed
sudo chown -R your-user:your-group /path/to/audiobooks
chmod -R 755 /path/to/audiobooks
```
