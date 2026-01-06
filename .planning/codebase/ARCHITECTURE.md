# Architecture

**Analysis Date:** 2026-01-06

## Pattern Overview

**Overall:** Full-Stack Monolith with Layered Architecture (Go Backend + React Frontend)

**Key Characteristics:**
- Clear separation between backend (Go) and frontend (React)
- REST API communication over HTTP/JSON
- Traditional layered architecture: handler → service → persistence
- Component-based frontend with centralized state management
- Background monitoring via goroutines

## Layers

**Backend Layers:**

**HTTP Layer:**
- Purpose: Request routing, middleware, CORS, validation
- Contains: `backend/internal/server/routes.go`, `backend/internal/server/handlers.go`, `backend/internal/server/validation.go`
- Depends on: Service layer for business logic
- Used by: Frontend API clients

**Service Layer:**
- Purpose: Business logic orchestration
- Contains: `backend/internal/config/service.go`, `backend/internal/downloads/service.go`, `backend/internal/search/search_service.go`
- Depends on: Data access layer (repository interfaces)
- Used by: HTTP handlers

**Data Access Layer:**
- Purpose: Database operations and persistence
- Contains: `backend/internal/persistence/interfaces.go`, `backend/internal/persistence/sqlite/`
- Depends on: SQLite database
- Used by: Service layer

**Domain Models:**
- Purpose: Core entities and types
- Contains: `backend/internal/models/download.go`, `backend/internal/models/search.go`
- Depends on: Nothing (pure data structures)
- Used by: All layers

**External Integration Layer:**
- Purpose: Third-party service communication
- Contains: `backend/internal/qbittorrent/client.go`, `backend/internal/search/providers/mam.go`
- Depends on: Go net/http
- Used by: Service layer

**Frontend Layers:**

**Page Layer:**
- Purpose: Route handlers and page-level state coordination
- Contains: `frontend/src/pages/SearchPage.tsx`, `frontend/src/pages/DownloadsPage.tsx`, `frontend/src/pages/ConfigPage.tsx`
- Depends on: Component layer, stores, API clients
- Used by: React Router

**Component Layer:**
- Purpose: Reusable UI components
- Contains:
  - Layout: `frontend/src/components/layout/` (Layout, Header, Sidebar)
  - Common: `frontend/src/components/common/` (Button, Input, Card, Modal, etc.)
  - Feature: `frontend/src/components/search/`, `frontend/src/components/downloads/`, `frontend/src/components/config/`
- Depends on: Type definitions
- Used by: Page layer

**State Management Layer:**
- Purpose: Global application state
- Contains: `frontend/src/stores/useDownloadStore.ts`, `frontend/src/stores/useSearchStore.ts`, `frontend/src/stores/useConfigStore.ts`, `frontend/src/stores/useNotificationStore.ts`
- Depends on: API client layer
- Used by: Pages and components

**API Client Layer:**
- Purpose: HTTP communication with backend
- Contains: `frontend/src/api/client.ts` (base client), `frontend/src/api/downloads.ts`, `frontend/src/api/search.ts`, `frontend/src/api/config.ts`
- Depends on: Type definitions, environment config
- Used by: Stores

**Type System:**
- Purpose: TypeScript interfaces and types
- Contains: `frontend/src/types/download.ts`, `frontend/src/types/search.ts`, `frontend/src/types/config.ts`, `frontend/src/types/api.ts`
- Depends on: Nothing
- Used by: All layers

## Data Flow

**Search and Download Flow:**

1. User enters search query in `SearchPage`
2. `useSearchStore.search()` called
3. `searchApi.search()` makes HTTP request to backend
4. Backend: `GET /api/search` → `Server.handleSearch()`
5. `MAMService.Search()` calls `MyAnonamouseProvider.Search()`
6. External API returns search results
7. Results transformed to domain models
8. Response sent to frontend
9. Store updates state, SearchResults component re-renders
10. User clicks download button
11. `useDownloadStore.createDownload()` called
12. `POST /api/downloads` → `Server.handleCreateDownload()`
13. `DownloadService.CreateDownload()` validates and stores in database
14. `qBittorrentClient.AddTorrent()` adds to torrent client
15. Download status returned to frontend

**Background Monitoring Flow:**

1. `main.go` initializes `Monitor.Run()` in background goroutine
2. Monitor periodically polls qBittorrent for download status
3. Database updated with current progress
4. On completion, auto-organization triggered if enabled
5. Files moved to organized paths based on templates
6. Frontend polls `GET /api/downloads` every 3 seconds
7. Store updates download status
8. UI reflects current state

**State Management:**
- Backend: Stateless request handling (no in-memory session)
- Frontend: Zustand stores with polling for real-time updates
- Database: All persistent state in SQLite

## Key Abstractions

**Repository Pattern:**
- Purpose: Interface-based data access for dependency inversion
- Examples: `DownloadRepository`, `ConfigRepository` in `backend/internal/persistence/interfaces.go`
- Implementations: `backend/internal/persistence/sqlite/downloads.go`, `backend/internal/persistence/sqlite/config.go`
- Pattern: Dependency injection via constructor

**Service Pattern:**
- Purpose: Business logic separation from HTTP layer
- Examples: `ConfigService` (`backend/internal/config/service.go`), `DownloadService` (`backend/internal/downloads/service.go`), `MAMService` (`backend/internal/search/search_service.go`)
- Pattern: Struct with dependencies injected via `New*Service()` constructors

**Provider Pattern:**
- Purpose: External integration abstraction
- Examples: `MyAnonamouseProvider` (`backend/internal/search/providers/mam.go`)
- Pattern: Pluggable search providers (future extensibility)

**Monitor Pattern:**
- Purpose: Background task orchestration
- Examples: `Monitor` (`backend/internal/downloads/monitor.go`)
- Pattern: Goroutine with context cancellation

**Client Pattern:**
- Purpose: External service wrappers
- Examples: `qBittorrentClient` (`backend/internal/qbittorrent/client.go`)
- Pattern: Cookie jar for session management, typed request/response

**Zustand Store Pattern:**
- Purpose: Centralized frontend state management
- Examples: All stores in `frontend/src/stores/`
- Pattern: Immutable state updates, computed getters, polling mechanism

**API Client Pattern:**
- Purpose: Centralized HTTP abstraction
- Examples: Base client in `frontend/src/api/client.ts`, domain clients in `frontend/src/api/downloads.ts`, etc.
- Pattern: Unified error handling with `APIClientError`, timeout support

## Entry Points

**Backend Entry:**
- Location: `backend/cmd/api/main.go`
- Triggers: Binary execution
- Responsibilities:
  - Database initialization
  - Service setup and dependency wiring
  - HTTP server start
  - Background monitor initialization
  - Graceful shutdown handling

**Frontend Entry:**
- Location: `frontend/src/main.tsx`
- Triggers: Browser loads index.html
- Responsibilities: Wraps App in StrictMode, mounts React to DOM
- Router: `frontend/src/App.tsx` - BrowserRouter with routes: `/search`, `/downloads`, `/config`, `/404`

## Error Handling

**Strategy:**
- Backend: Throw errors, catch at handler boundaries, return structured JSON errors
- Frontend: Catch in API client, wrap in `APIClientError`, display via notification store

**Patterns:**
- Backend: `fmt.Errorf(..., %w, err)` for error context, deferred cleanup
- Frontend: try/catch in async operations, error state in stores
- HTTP handlers: Centralized error response in `backend/internal/server/errors.go`

**Concerns:**
- Multiple ignored errors in I/O operations (`backend/internal/qbittorrent/client.go:57,144`, `backend/internal/search/providers/mam.go:140`)
- Context misuse in background operations (`backend/internal/downloads/monitor.go:96`)

## Cross-Cutting Concerns

**Logging:**
- Backend: Standard Go log package, console output
- Frontend: Development console.log in `frontend/src/api/client.ts`
- No structured logging framework

**Validation:**
- Backend: Manual validation in `backend/internal/server/validation.go`, regex patterns for UUIDs
- Frontend: React Hook Form for form validation
- API boundary: Input validation at HTTP handlers

**Configuration:**
- Backend: Database-driven config repository (`backend/internal/config/service.go`)
- Frontend: Environment variables via Vite (`frontend/src/utils/env.ts`)
- No centralized config management

**Authentication:**
- None implemented (security gap)

---

*Architecture analysis: 2026-01-06*
*Update when major patterns change*
