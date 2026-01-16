# System Architecture

## Overview

Organizr is a client-server application that automates audiobook downloading and organization through qBittorrent integration. The system follows a two-tier architecture with clear separation between backend API services and frontend UI.

**Core Purpose:** Automatically download audiobooks from MyAnonamouse (MAM), monitor download progress, and organize completed files using customizable path templates.

## Architectural Pattern

**Layered Architecture** with dependency injection:

1. **API Layer** (backend/internal/server/)
   - HTTP handlers expose REST endpoints
   - Chi router with middleware (CORS, logging)
   - Request/response DTOs and error formatting
   - Swagger documentation generation

2. **Service Layer** (backend/internal/*)
   - `DownloadService` - Orchestrates torrent submission and download creation
   - `Monitor` - Background goroutine polling qBittorrent (30s interval)
   - `OrganizationService` - Post-download file organization
   - `ConfigService` - Environment override + database fallback
   - `MAMService` - MyAnonamouse torrent search
   - `QBittorrentClient` - qBittorrent Web API integration

3. **Persistence Layer** (backend/internal/persistence/)
   - Repository interfaces define data access contracts
   - SQLite implementation with WAL mode for concurrency
   - Connection pooling (25 max connections)
   - SQL migrations in assets/migrations/

4. **Frontend Layer**
   - React 19 SPA with React Router v7
   - Zustand state management (separate stores per domain)
   - API client layer wraps HTTP requests
   - TailwindCSS for styling

## System Boundaries

### Backend Service (Go)
- **Responsibility:** API endpoints, download orchestration, qBittorrent integration, file organization
- **Port:** 8080 (configurable via BACKEND_PORT)
- **Storage:** SQLite database with WAL mode
- **External Dependencies:** qBittorrent Web API, MyAnonamouse API

### Frontend Service (React/Nginx)
- **Responsibility:** User interface, search, download monitoring, configuration
- **Port:** 8081 (configurable via FRONTEND_PORT)
- **API Communication:** REST API via /api proxy
- **State Management:** Zustand stores with periodic polling

### External Integrations
- **qBittorrent:** Torrent client for downloads (Web API on port 8080 default)
- **MyAnonamouse:** Torrent search and metadata (HTTPS API)
- **Filesystem:** Download monitoring and file organization

## Design Decisions

### 1. SQLite with WAL Mode
**Rationale:** Embedded database eliminates external dependencies while supporting concurrent reads during background monitoring. WAL mode prevents write locks from blocking monitor goroutine.

**Trade-off:** Limited to single-host deployment, but acceptable for target use case (personal/home server).

### 2. Background Monitor Goroutine
**Rationale:** Separate goroutine polls qBittorrent status every 30s, updates download progress, and triggers organization when complete.

**Trade-off:** Polling creates slight latency (max 30s) vs push-based notifications, but simpler than webhook management and works with any qBittorrent setup.

### 3. Repository Pattern for Persistence
**Rationale:** Interfaces allow testability, potential database switching, and clear separation between data access and business logic.

**Trade-off:** Adds abstraction layer, but improves testability and maintainability.

### 4. Environment Variables Override Config Database
**Rationale:** Docker/container-friendly configuration with database fallback for runtime changes via UI.

**Trade-off:** Two sources of truth for configuration, but ConfigService provides unified access.

### 5. Template-Based Path Organization
**Rationale:** Flexible file organization supporting various folder structures via `{author}/{series}/{title}` templates.

**Trade-off:** Template validation complexity, but provides user control over organization.

### 6. Zustand Over Redux (Frontend)
**Rationale:** Lightweight state management with less boilerplate than Redux, sufficient for application size.

**Trade-off:** Less tooling/ecosystem than Redux, but better DX for this scale.

## Constraints

### Technical Constraints
- **SQLite Concurrency:** Maximum 25 concurrent connections, write serialization even with WAL mode
- **qBittorrent API:** Cookie-based authentication requires session management, potential rate limiting
- **Filesystem Access:** Backend must have read/write access to qBittorrent download paths
- **Docker Volume Mapping:** Path prefixes must be configured when qBittorrent runs in different container

### Operational Constraints
- **Single Host:** SQLite limits to single-host deployment
- **Network Access:** Requires connectivity to qBittorrent and MyAnonamouse
- **MAM API Key:** Requires valid MyAnonamouse account and API secret

### Security Constraints
- **No Authentication:** Current version has no auth layer (trusted network assumption)
- **API Secrets:** qBittorrent credentials and MAM secret stored in environment/config
- **CORS:** Configured for development, must be restricted in production

## Data Flow

### Download Creation Flow
1. User searches MAM via frontend → `GET /api/search?q=<query>`
2. Backend calls MAM API, returns results with series detection
3. User selects download(s) → `POST /api/downloads` or `POST /api/downloads/batch`
4. DownloadService:
   - Validates request (URL or magnet link)
   - Submits torrent to qBittorrent (file upload or magnet link)
   - Creates Download entity with status=queued
   - Saves to SQLite with qBittorrent hash
5. Returns created download(s) to frontend

### Background Monitor Flow (30s interval)
1. Monitor goroutine queries all downloads with status=queued|downloading
2. Fetches qBittorrent status for each hash
3. Updates progress and status in database
4. When status=completed:
   - Triggers OrganizationService
   - Applies path template (`{author}/{series}/{title}`)
   - Copies/moves files to destination
   - Updates status=organized, sets organized_at
5. On error: Sets status=failed, stores error_message

### Manual Organization Flow
1. User clicks "Organize" → `POST /api/downloads/{id}/organize`
2. OrganizationService:
   - Validates download is complete
   - Reads path template from config
   - Sanitizes and expands template variables
   - Executes file operation (copy or move)
   - Updates organized_path and status
3. Returns updated download to frontend

### Configuration Flow
1. Frontend loads config → `GET /api/config`
2. ConfigService merges environment variables + database values
3. User updates setting → `PUT /api/config/{key}`
4. Validation and save to database
5. Environment variables remain override source
