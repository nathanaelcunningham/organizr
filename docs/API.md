# Organizr API Documentation

Base URL: `http://localhost:8080/api`

All requests and responses use `application/json` content type unless otherwise specified.

## Error Responses

All API errors follow a standardized format:

```json
{
  "error": "Error type",
  "message": "Detailed error message",
  "code": 400
}
```

Common HTTP status codes:
- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Request successful with no response body
- `400 Bad Request` - Invalid input or validation error
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Downloads

### Create Download

Create a new download and add it to qBittorrent.

**Endpoint:** `POST /api/downloads`

**Request Body:**
```json
{
  "title": "The Gunslinger",
  "author": "Stephen King",
  "series": "The Dark Tower",
  "torrent_url": "https://example.com/torrent.torrent",
  "magnet_link": "magnet:?xt=urn:btih:..."
}
```

**Fields:**
- `title` (string, required): Book title (max 500 characters)
- `author` (string, required): Author name (max 200 characters)
- `series` (string, optional): Series name (max 200 characters)
- `torrent_url` (string, optional): Direct torrent file URL
- `magnet_link` (string, optional): Magnet link

**Note:** Either `torrent_url` or `magnet_link` must be provided.

**Response:** `201 Created`
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

---

### List Downloads

Get all downloads with their current status.

**Endpoint:** `GET /api/downloads`

**Response:** `200 OK`
```json
{
  "downloads": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "The Gunslinger",
      "author": "Stephen King",
      "series": "The Dark Tower",
      "status": "downloading",
      "progress": 45.5,
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

**Download Statuses:**
- `queued` - Added to qBittorrent, waiting to start
- `downloading` - Currently downloading
- `completed` - Download finished, pending organization
- `organizing` - Files being organized
- `organized` - Fully complete and organized
- `failed` - Error occurred

---

### Get Download

Get details for a specific download.

**Endpoint:** `GET /api/downloads/{id}`

**Parameters:**
- `id` (UUID): Download ID

**Response:** `200 OK`
```json
{
  "download": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "The Gunslinger",
    "author": "Stephen King",
    "series": "The Dark Tower",
    "status": "organized",
    "progress": 100,
    "organized_path": "/audiobooks/Stephen King/The Dark Tower/The Gunslinger",
    "created_at": "2026-01-01T00:00:00Z",
    "completed_at": "2026-01-01T01:30:00Z",
    "organized_at": "2026-01-01T01:31:00Z"
  }
}
```

---

### Cancel Download

Cancel an active download and remove it from qBittorrent.

**Endpoint:** `DELETE /api/downloads/{id}`

**Parameters:**
- `id` (UUID): Download ID

**Response:** `204 No Content`

---

### Manually Organize Download

Manually trigger file organization for a completed download.

**Endpoint:** `POST /api/downloads/{id}/organize`

**Parameters:**
- `id` (UUID): Download ID

**Response:** `200 OK`

---

## Search

### Search Torrents

Search for torrents across configured providers.

**Endpoint:** `GET /api/search`

**Query Parameters:**
- `q` (string, required): Search query (min 2 characters)
- `provider` (string, optional): Specific provider name

**Examples:**
```bash
# Search all providers
GET /api/search?q=dark+tower+king

# Search specific provider
GET /api/search?q=dark+tower&provider=AudiobookBay
```

**Response:** `200 OK`
```json
{
  "results": [
    {
      "title": "The Dark Tower: The Gunslinger",
      "author": "Stephen King",
      "torrent_url": "https://example.com/torrent.torrent",
      "magnet_link": "magnet:?xt=urn:btih:...",
      "size": "450 MB",
      "seeders": 42,
      "provider": "AudiobookBay"
    }
  ],
  "count": 1
}
```

---

### List Providers

Get list of available search providers.

**Endpoint:** `GET /api/search/providers`

**Response:** `200 OK`
```json
{
  "providers": [
    "AudiobookBay",
    "MyCustomProvider"
  ]
}
```

---

## Configuration

### Get All Configuration

Retrieve all configuration key-value pairs.

**Endpoint:** `GET /api/config`

**Response:** `200 OK`
```json
{
  "configs": {
    "qbittorrent.url": "http://localhost:8080",
    "qbittorrent.username": "admin",
    "qbittorrent.password": "adminpass",
    "paths.destination": "/audiobooks",
    "paths.template": "{author}/{series}/{title}",
    "paths.operation": "copy",
    "monitor.interval_seconds": "30",
    "monitor.auto_organize": "true"
  }
}
```

---

### Get Configuration Value

Get a specific configuration value.

**Endpoint:** `GET /api/config/{key}`

**Parameters:**
- `key` (string): Configuration key

**Response:** `200 OK`
```json
{
  "key": "paths.destination",
  "value": "/audiobooks"
}
```

---

### Update Configuration

Update a configuration value.

**Endpoint:** `PUT /api/config/{key}`

**Parameters:**
- `key` (string): Configuration key

**Request Body:**
```json
{
  "value": "/mnt/audiobooks"
}
```

**Response:** `200 OK`

**Configuration Keys:**

| Key | Description | Default | Type |
|-----|-------------|---------|------|
| `qbittorrent.url` | qBittorrent Web UI URL | `http://localhost:8080` | URL |
| `qbittorrent.username` | qBittorrent username | `admin` | string |
| `qbittorrent.password` | qBittorrent password | `adminpass` | string |
| `paths.destination` | Base directory for organized files | `/audiobooks` | path |
| `paths.template` | Path template with series | `{author}/{series}/{title}` | template |
| `paths.no_series_template` | Path template without series | `{author}/{title}` | template |
| `paths.operation` | File operation type | `copy` | `copy` or `move` |
| `monitor.interval_seconds` | Monitor poll interval | `30` | integer |
| `monitor.auto_organize` | Auto-organize on completion | `true` | `true` or `false` |

**Path Template Variables:**
- `{author}` - Book author
- `{series}` - Series name (if provided)
- `{title}` - Book title

---

## Health

### Health Check

Check service health status.

**Endpoint:** `GET /api/health`

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "database": "ok",
  "qbittorrent": "unknown",
  "monitor": "running"
}
```

---

## Rate Limiting

Currently, there is no rate limiting implemented. For production use, consider adding rate limiting middleware.

## Authentication

Currently, there is no authentication required. For production use, consider adding authentication and authorization.

## CORS

CORS is enabled for all origins by default. Configure `AllowedOrigins` in `internal/server/server.go` for production use.
