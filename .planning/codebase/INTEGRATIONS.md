# External Integrations

**Analysis Date:** 2026-01-06

## APIs & External Services

**qBittorrent Torrent Client:**
- Purpose: Download management, torrent manipulation, progress monitoring
- Integration method: REST API via custom HTTP client
- Implementation: `backend/internal/qbittorrent/client.go`, `backend/internal/qbittorrent/types.go`
- Service integration: `backend/internal/downloads/service.go`, `backend/internal/downloads/monitor.go`
- API endpoints used:
  - `/api/v2/auth/login` - Authentication
  - `/api/v2/torrents/add` - Add torrents via magnet or URL
  - `/api/v2/torrents/info` - Get torrent status and info
  - `/api/v2/torrents/files` - Retrieve torrent file list
  - `/api/v2/torrents/delete` - Remove torrents
- Configuration keys (stored in database):
  - `qbittorrent.url` - Base URL (default: `http://localhost:8080`)
  - `qbittorrent.username` - Web UI username (default: `admin`)
  - `qbittorrent.password` - Web UI password (default: `adminpass`)
- Auth: Cookie-based session management with cookie jar

**MyAnonamouse (MAM) Private Torrent Tracker:**
- Purpose: Search for audiobooks and books with metadata
- Integration method: REST API via custom HTTP client
- Implementation: `backend/internal/search/providers/mam.go`, `backend/internal/search/search_service.go`
- Frontend integration: `frontend/src/api/search.ts`
- API endpoints used:
  - `/tor/js/loadSearchJSONbasic.php` - Search torrents
  - `/tor/download.php?tid={id}` - Download torrent file
  - `/jsonLoad.php` - Test connection/authentication
- Auth: Cookie-based (`mam_id` header)
- Search parameters: Main categories (Audiobooks: 13, Books: 14), search fields (Title, Author, Series, Narrator)
- Results: Up to 100 results per page
- Configuration keys:
  - `mam.baseurl` - Base URL (default: `https://www.myanonamouse.net`)
  - `mam.secret` - API secret/authentication key

## Data Storage

**Databases:**
- SQLite3 - Embedded database
  - Connection: Direct file access at `backend/organizr.db`
  - Driver: mattn/go-sqlite3 v1.14.32 (`backend/go.mod`)
  - Implementation: `backend/internal/persistence/sqlite/db.go`
  - Migrations: SQL-based at `backend/assets/migrations/001_init.up.sql`
  - Schema: Downloads and configuration tables
  - WAL mode: Enabled for concurrency

**File Storage:**
- Local filesystem only
- No cloud storage integrations
- Download paths tracked in database
- Organized paths: Template-based file organization (`backend/internal/downloads/organization.go`, `backend/internal/fileutil/`)

**Caching:**
- None detected

## Authentication & Identity

**Auth Provider:**
- None implemented
- No authentication layer present
- API endpoints are publicly accessible

**OAuth Integrations:**
- None detected

## Monitoring & Observability

**Error Tracking:**
- None detected

**Analytics:**
- None detected

**Logs:**
- Standard output/stderr only (Go log package)
- No structured logging framework

## CI/CD & Deployment

**Hosting:**
- Not detected

**CI Pipeline:**
- Not detected

## Environment Configuration

**Development:**
- Required env vars: `VITE_API_URL` for frontend (`frontend/.env.development`)
- Backend: Database-driven configuration, no environment variables required
- Secrets location: Database configs table
- Missing: No `.env.example` files

**Staging:**
- Not configured

**Production:**
- Not configured

## Webhooks & Callbacks

**Incoming:**
- None detected

**Outgoing:**
- None detected

## Frontend API Communication

**Base Configuration:**
- `frontend/src/utils/env.ts` - Environment-driven API URL
- `frontend/src/api/client.ts` - Generic HTTP request handler with:
  - Request timeout handling (30 second default)
  - Error handling and APIClientError wrapper
  - Development logging
  - Query parameter serialization

**Endpoint Interfaces:**
- `frontend/src/api/search.ts` - Search endpoints
- `frontend/src/api/downloads.ts` - Download management
- `frontend/src/api/config.ts` - Configuration management

## Backend HTTP Server

**CORS Policy:**
- Location: `backend/internal/server/server.go`
- Allowed Origins: `*` (all origins - security concern)
- Allowed Methods: GET, POST, PUT, DELETE, OPTIONS
- Allowed Headers: Accept, Authorization, Content-Type, X-CSRF-Token

**Middleware Stack:**
- Request ID tracking
- Real IP detection
- Request logging
- Panic recovery
- CORS handling

**No Rate Limiting:**
- Not detected

---

*Integration audit: 2026-01-06*
*Update when adding/removing external services*
