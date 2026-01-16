# Key Files & Directories

## Entry Points

### Backend Entry Point
**File:** `backend/cmd/api/main.go` (156 lines)

**Purpose:** Application initialization and startup

**Key Responsibilities:**
- Load environment variables from .env file
- Initialize SQLite database with WAL mode (25 connection pool)
- Run database migrations (001-004)
- Create repository instances (DownloadRepository, ConfigRepository)
- Initialize services (ConfigService, MAMService, QBittorrentClient)
- Create DownloadService and start Monitor goroutine (background polling)
- Setup HTTP server with Chi router on port 8080
- Handle graceful shutdown (SIGINT/SIGTERM)

**Critical Initialization Order:**
1. Database connection
2. Migrations
3. Repositories
4. Services
5. HTTP server
6. Monitor goroutine (after server ready)

### Frontend Entry Point
**File:** `frontend/src/main.tsx` (10 lines)

**Purpose:** React application bootstrap

**Key Responsibilities:**
- Create React root with StrictMode
- Render App component with BrowserRouter
- Mount to DOM element #root

**File:** `frontend/src/App.tsx` (30 lines)

**Purpose:** Route configuration

**Key Responsibilities:**
- Define application routes (/, /downloads, /config)
- Layout wrapper with navigation
- 404 handling

## Core Business Logic

### Download Orchestration
**Directory:** `backend/internal/downloads/`

**Files:**
- `service.go` (250 lines) - Download creation, torrent submission to qBittorrent
- `monitor.go` (180 lines) - Background polling (30s interval), status updates
- `organization.go` (200 lines) - File organization with path templates
- `*_test.go` - Test coverage for download workflows

**Key Functions:**
- `CreateDownload()` - Submit torrent and create download record
- `Monitor.Start()` - Background goroutine for polling qBittorrent
- `OrganizeDownload()` - Apply path template and move/copy files

### qBittorrent Integration
**Directory:** `backend/internal/qbittorrent/`

**Files:**
- `client.go` (300 lines) - qBittorrent Web API client
- `types.go` (100 lines) - API request/response types

**Key Operations:**
- Login with cookie authentication
- Upload torrent file (multipart/form-data)
- Add magnet link
- Get torrent status by hash
- Map qBittorrent status to Organizr status

### MyAnonamouse Search
**Directory:** `backend/internal/search/`

**Files:**
- `search_service.go` (80 lines) - Search service wrapper
- `providers/mam.go` (250 lines) - MAM API implementation
- `providers/mam_test.go` - MAM integration tests

**Key Operations:**
- Search torrents by query
- Parse series information from results
- Handle MAM authentication (API secret)
- Map MAM response to SearchResult model

### Configuration Management
**Directory:** `backend/internal/config/`

**Files:**
- `service.go` (150 lines) - Config service with env override + DB fallback
- `env_mapping.go` (50 lines) - Environment variable to config key mappings
- `service_test.go` - Config service tests

**Key Behavior:**
- Environment variables override database values
- Database provides runtime-configurable defaults
- Hierarchical keys (e.g., "qbittorrent.url", "paths.template")

### Path Templates & Sanitization
**Directory:** `backend/internal/fileutil/`

**Files:**
- `template.go` (120 lines) - Template parsing and expansion
- `sanitizer.go` (80 lines) - Path sanitization for filesystems
- `*_test.go` - Template and sanitization tests

**Supported Variables:**
- `{author}` - Author name
- `{series}` - Series name
- `{title}` - Book title
- `{series_number}` - Series position

**Critical Logic:**
- Validate template syntax
- Sanitize paths (remove invalid chars, handle spaces)
- Support separate templates for series vs non-series books
- Handle path prefix for Docker volume mapping

## API Layer

### HTTP Handlers
**Directory:** `backend/internal/server/`

**Files:**
- `server.go` (100 lines) - HTTP server setup with Chi router
- `routes.go` (80 lines) - Route registration and grouping
- `handlers.go` (500 lines) - API endpoint implementations
- `errors.go` (60 lines) - Error response helpers
- `dto.go` (150 lines) - Data transfer objects (request/response types)
- `handlers_test.go` - Handler integration tests

**API Routes:**
```
POST   /api/downloads              # Create download
POST   /api/downloads/batch        # Batch create
GET    /api/downloads              # List downloads
GET    /api/downloads/{id}         # Get download
DELETE /api/downloads/{id}         # Cancel download
POST   /api/downloads/{id}/organize # Manual organize

GET    /api/config                 # Get all config
GET    /api/config/{key}           # Get single config
PUT    /api/config/{key}           # Update config
POST   /api/config/preview-path    # Preview template

GET    /api/search?q=<query>       # Search MAM
POST   /api/search/test            # Test MAM

GET    /api/qbittorrent/test       # Test qBittorrent
GET    /api/health                 # Health check
```

## Data Models & Persistence

### Domain Models
**Directory:** `backend/internal/models/`

**Files:**
- `download.go` (100 lines) - Download entity with statuses
- `config.go` (30 lines) - Config key-value model
- `search.go` (80 lines) - SearchResult, SeriesInfo models

**Key Model: Download**
- 6 statuses: queued, downloading, completed, organizing, organized, failed
- Tracks qBittorrent hash, progress, paths, timestamps
- Unique index on qbit_hash

### Repository Layer
**Directory:** `backend/internal/persistence/`

**Files:**
- `interfaces.go` (80 lines) - Repository interface definitions
- `sqlite/db.go` (100 lines) - SQLite connection with WAL mode
- `sqlite/downloads.go` (350 lines) - Download CRUD operations
- `sqlite/config.go` (150 lines) - Config key-value operations
- `sqlite/downloads_test.go` - Repository tests

**Critical Configuration:**
- WAL mode for concurrent reads
- Connection pool: 25 max open, 25 max idle
- Prepared statements for common queries

### Database Migrations
**Directory:** `backend/assets/migrations/`

**Files:**
- `001_init.up.sql` - Initial schema (downloads, configs tables)
- `002_add_category.up.sql` - Add category field
- `003_add_path_prefix.up.sql` - Add path prefix for Docker
- `004_add_series_number.up.sql` - Add series_number field

**Migration Strategy:**
- Sequential numbering (001, 002, 003, ...)
- Only .up.sql files (no down migrations)
- Run on startup via main.go

## Frontend State Management

### Zustand Stores
**Directory:** `frontend/src/stores/`

**Files:**
- `useDownloadStore.ts` (200 lines) - Downloads, polling logic
- `useSearchStore.ts` (120 lines) - Search results, query
- `useConfigStore.ts` (100 lines) - Settings state
- `useNotificationStore.ts` (60 lines) - Toast notifications
- `useDownloadStore.test.ts` (12,761 lines) - Store tests

**Key Store: useDownloadStore**
- Manages downloads list, loading, error state
- Periodic polling (configurable interval)
- CRUD operations (create, fetch, delete)
- Computed getters (activeDownloads, completedDownloads, failedDownloads)

### API Client Layer
**Directory:** `frontend/src/api/`

**Files:**
- `client.ts` (120 lines) - HTTP request wrapper with error handling
- `downloads.ts` (80 lines) - Download API endpoints
- `search.ts` (40 lines) - Search API endpoint
- `config.ts` (60 lines) - Config API endpoints
- `qbittorrent.ts` (30 lines) - Test connection

**Key Pattern:**
- Centralized `apiRequest()` function
- `APIClientError` with statusCode and apiError fields
- Query parameter building
- Timeout support

## Frontend Components

### Page Components
**Directory:** `frontend/src/pages/`

**Files:**
- `SearchPage.tsx` (150 lines) - MAM search UI
- `DownloadsPage.tsx` (120 lines) - Download list with polling
- `ConfigPage.tsx` (100 lines) - Settings form
- `NotFoundPage.tsx` (40 lines) - 404 page

### Feature Components
**Directories:**
- `frontend/src/components/search/` - Search form, series grouping
- `frontend/src/components/downloads/` - Download table, progress bars
- `frontend/src/components/config/` - Settings forms
- `frontend/src/components/layout/` - Layout wrapper with nav
- `frontend/src/components/common/` - Notifications, Modals, Loading

## Configuration Files

### Backend Configuration
**Files:**
- `backend/go.mod` - Go module dependencies
- `backend/Makefile` - Build, test, lint targets
- `backend/.air.toml` - Air hot reload configuration (dev)
- `backend/Dockerfile` - Multi-stage Docker build
- `.golangci.yml` - Go linter configuration

### Frontend Configuration
**Files:**
- `frontend/package.json` - npm dependencies and scripts
- `frontend/vite.config.ts` - Vite build configuration
- `frontend/vitest.config.ts` - Test configuration (60% coverage)
- `frontend/tsconfig.json` - TypeScript configuration
- `frontend/Dockerfile` - Multi-stage Docker build
- `frontend/nginx.conf` - Nginx server configuration

### Infrastructure
**Files:**
- `docker-compose.yml` - Service orchestration (backend, frontend)
- `.env.example` - Environment variable template
- `.dockerignore` - Docker build exclusions
- `.gitignore` - Version control exclusions

## Documentation

**Files:**
- `README.md` - Project overview and quick start
- `CONTRIBUTING.md` - Contribution guidelines (deleted in git status)
- `backend/README.md` - Backend-specific documentation
- `backend/docs/API.md` - API documentation
- `backend/docs/CONFIGURATION.md` - Configuration reference
- `frontend/README.md` - Frontend-specific documentation
- `docs/DEPLOYMENT.md` - Deployment guide (Docker, Unraid, bare metal)
- `docs/TROUBLESHOOTING.md` - Common issues and solutions
- `docs/architecture/ADR.md` - Architecture decision records
- `docs/architecture/DIAGRAMS.md` - System diagrams

## Common Pitfalls

### Template Validation Issues
**Problem:** Invalid template syntax causes organization failures

**Files to Check:**
- `backend/internal/fileutil/template.go` - Template parsing logic
- `backend/internal/fileutil/sanitizer.go` - Path sanitization

**Common Causes:**
- Missing closing brace in template (e.g., `{author/`)
- Unsupported variables (only author, series, title, series_number allowed)
- Path sanitization edge cases (unicode, special chars)

**How to Debug:**
- Use `/api/config/preview-path` endpoint to test templates
- Check error_message field on failed downloads
- Review sanitization logic for specific character handling

### Path Mapping Confusion
**Problem:** Organized files end up in wrong location or paths don't match

**Files to Check:**
- `backend/internal/downloads/organization.go` - File organization logic
- `backend/internal/config/env_mapping.go` - PATHS_LOCAL_MOUNT mapping

**Common Causes:**
- Docker volume mounts differ from qBittorrent download paths
- Missing or incorrect PATHS_LOCAL_MOUNT prefix
- Destination path not writable by backend user

**How to Debug:**
- Compare qBittorrent download path to backend container paths
- Verify PATHS_LOCAL_MOUNT is set correctly
- Check Docker Compose volume mappings
- Test with manual organize endpoint first

## Performance Considerations

### Monitor Polling Interval
**File:** `backend/internal/downloads/monitor.go`

**Consideration:** Lower intervals increase qBittorrent API load

**Recommendation:** 30s default is reasonable, avoid <10s

### Concurrent Organization
**File:** `backend/internal/downloads/monitor.go`

**Current Limit:** 3 concurrent organization operations

**Consideration:** File I/O can be slow, avoid saturating disk

### SQLite Connection Pool
**File:** `backend/internal/persistence/sqlite/db.go`

**Current Setting:** 25 max connections

**Consideration:** SQLite has write serialization even with WAL, don't over-provision

### Frontend Polling
**File:** `frontend/src/stores/useDownloadStore.ts`

**Consideration:** Unnecessary polling when no active downloads

**Recommendation:** Only poll when activeDownloads.length > 0

## Security Requirements

### API Authentication
**Current State:** No authentication implemented

**Assumption:** Trusted network environment (home server, behind VPN)

**Future Consideration:** Add JWT or session-based auth for internet exposure

### qBittorrent Credentials
**Files:** Environment variables (QBITTORRENT_USERNAME, QBITTORRENT_PASSWORD)

**Storage:** Should be in .env file (not committed) or Docker secrets

### MAM API Secret
**Files:** Environment variable (MAM_SECRET)

**Importance:** Protects MAM account, treat as sensitive credential

### CORS Configuration
**File:** `backend/internal/server/server.go`

**Current:** Configured for development (localhost:5173)

**Production:** Restrict to specific frontend origin only
