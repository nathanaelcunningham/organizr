# Organizr

A Go-based REST API service for managing audiobook downloads through qBittorrent with automatic file organization.

## Features

- 📥 REST API for managing audiobook downloads
- 🔍 Pluggable torrent search provider system
- 📁 Automatic organization of completed downloads into configurable directory structures
- 🔄 Background monitoring of qBittorrent downloads
- 🏷️ Metadata-based file organization with customizable path templates
- 💾 SQLite-based configuration and download tracking
- ✅ Comprehensive input validation and error handling

## Architecture

- **Go backend** with Chi router
- **Direct SQL** with SQLite (no ORMs)
- **qBittorrent integration** via Web API
- **Background monitor** for automatic file organization
- **Configurable path templates** for organizing audiobooks

## Project Structure

```
organizr/
├── cmd/api/              # Application entry point
├── internal/
│   ├── models/           # Domain models
│   ├── server/           # HTTP server, routes, handlers
│   ├── persistence/      # Database access layer
│   ├── downloads/        # Download service logic
│   ├── config/           # Configuration service
│   ├── qbittorrent/      # qBittorrent API client
│   └── fileutil/         # File utilities
└── assets/migrations/    # Database migrations

## Setup

1. **Build the project**:
   ```bash
   make build
   ```

2. **Run the service**:
   ```bash
   make run
   ```

   The service will start on port 8080 and create an `organizr.db` SQLite database.

3. **Configure qBittorrent connection**:
   ```bash
   curl -X PUT http://localhost:8080/api/config/qbittorrent.url \
     -H "Content-Type: application/json" \
     -d '{"value": "http://your-qbittorrent:8080"}'
   ```

## API Endpoints

### Downloads

- `POST /api/downloads` - Create a new download
- `GET /api/downloads` - List all downloads
- `GET /api/downloads/{id}` - Get download details
- `DELETE /api/downloads/{id}` - Cancel download
- `POST /api/downloads/{id}/organize` - Manually trigger organization

### Search

- `GET /api/search?q=<query>` - Search all providers
- `GET /api/search?q=<query>&provider=<name>` - Search specific provider
- `GET /api/search/providers` - List available providers

### Configuration

- `GET /api/config` - Get all configuration
- `GET /api/config/{key}` - Get specific config value
- `PUT /api/config/{key}` - Update config value

### Health

- `GET /api/health` - Service health check

## Example Usage

### Search for torrents:

```bash
curl "http://localhost:8080/api/search?q=dark+tower+king"
```

### Create a download:

```bash
curl -X POST http://localhost:8080/api/downloads \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Gunslinger",
    "author": "Stephen King",
    "series": "The Dark Tower",
    "magnet_link": "magnet:?xt=urn:btih:..."
  }'
```

**Response:**
```json
{
  "download": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "The Gunslinger",
    "author": "Stephen King",
    "series": "The Dark Tower",
    "status": "queued",
    "progress": 0,
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

### List downloads:

```bash
curl http://localhost:8080/api/downloads
```

### Get download status:

```bash
curl http://localhost:8080/api/downloads/550e8400-e29b-41d4-a716-446655440000
```

### Update configuration:

```bash
curl -X PUT http://localhost:8080/api/config/paths.destination \
  -H "Content-Type: application/json" \
  -d '{"value": "/audiobooks"}'
```

## Configuration

All configuration is stored in the SQLite database. Default values:

- `qbittorrent.url`: `http://localhost:8080`
- `qbittorrent.username`: `admin`
- `qbittorrent.password`: `adminpass`
- `paths.destination`: `/audiobooks`
- `paths.template`: `{author}/{series}/{title}`
- `paths.no_series_template`: `{author}/{title}`
- `paths.operation`: `copy` (or `move`)
- `monitor.interval_seconds`: `30`
- `monitor.auto_organize`: `true`

## Search Providers

Organizr supports pluggable search providers. To add your own provider:

1. **Implement the Provider interface** in `internal/search/providers/`:

```go
package providers

import (
    "context"
    "github.com/nathanael/organizr/internal/models"
)

type MyCustomProvider struct {
    apiKey string
}

func NewMyCustomProvider(apiKey string) *MyCustomProvider {
    return &MyCustomProvider{apiKey: apiKey}
}

func (p *MyCustomProvider) Name() string {
    return "MyCustomProvider"
}

func (p *MyCustomProvider) Search(ctx context.Context, query string) ([]*models.SearchResult, error) {
    // Your search implementation here
    return results, nil
}
```

2. **Register your provider** in `cmd/api/main.go`:

```go
providers := []search.Provider{
    providers.NewMyCustomProvider("your-api-key"),
}
searchService := search.NewService(providers)
```

See `internal/search/providers/example.go` for a complete implementation example.

## File Organization

Organizr automatically organizes completed downloads using path templates:

**With series:**
- Template: `{author}/{series}/{title}`
- Result: `/audiobooks/Stephen King/The Dark Tower/The Gunslinger/`

**Without series:**
- Template: `{author}/{title}`
- Result: `/audiobooks/Stephen King/The Gunslinger/`

Path components are automatically sanitized to remove invalid filesystem characters.

## Development

```bash
# Build
make build

# Run (with hot reload)
make run

# Clean build artifacts
make clean

# Tidy dependencies
make tidy
```

## Troubleshooting

### Connection to qBittorrent fails

Ensure qBittorrent Web UI is enabled and accessible:
1. Open qBittorrent → Tools → Options → Web UI
2. Enable "Web User Interface"
3. Note the port (default 8080)
4. Update Organizr configuration with correct URL and credentials

### Downloads not organizing automatically

Check monitor configuration:
```bash
curl http://localhost:8080/api/config/monitor.auto_organize
curl http://localhost:8080/api/config/monitor.interval_seconds
```

Ensure `monitor.auto_organize` is set to `true` and interval is reasonable (e.g., `30` seconds).

### Permission denied when organizing files

Ensure the Organizr process has write permissions to the destination directory:
```bash
chmod 755 /audiobooks
```

## Error Responses

All API errors return a standardized JSON format:

```json
{
  "error": "Validation failed",
  "message": "title is required and cannot be empty",
  "code": 400
}
```

## License

MIT
