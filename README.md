# Organizr

Audiobook download and organization service that integrates with qBittorrent.

## Features

- REST API for managing audiobook downloads
- Automatic organization of completed downloads into configurable directory structures
- Background monitoring of qBittorrent downloads
- Metadata-based file organization with customizable path templates
- SQLite-based configuration and download tracking

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

### Configuration

- `GET /api/config` - Get all configuration
- `GET /api/config/{key}` - Get specific config value
- `PUT /api/config/{key}` - Update config value

### Health

- `GET /api/health` - Service health check

## Example Usage

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

### List downloads:

```bash
curl http://localhost:8080/api/downloads
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

## Development

```bash
# Build
make build

# Run
make run

# Clean
make clean

# Test
make test

# Tidy dependencies
make tidy
```

## License

MIT
