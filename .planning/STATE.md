# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-01-09 after v1.2 milestone)

**Core value:** Perfect folder structures and file placement every time - audiobooks land exactly where Audiobookshelf expects them, automatically.
**Current focus:** v1.3 Production Deployment - Docker containerization for Unraid

## Current Position

Milestone: v1.3 Production Deployment
Phase: 21 of 21 (Configurable Port Variables)
Plan: 1 of 1 in current phase
Status: Phase complete
Last activity: 2026-01-09 - Completed 21-01-PLAN.md

Progress: ██████████ 100% (21 of 21 phases complete)

## Performance Metrics

**Velocity (v1.2 milestone):**
- Plans completed: 11
- Average duration: ~4 minutes
- Total execution time: ~0.7 hours
- Milestone duration: 2 days

**Overall Velocity (all milestones):**
- Total plans completed: 26 (8 in v1.0 + 7 in v1.1 + 11 in v1.2)
- Cumulative execution time: ~13 hours
- Three milestones shipped: v1.0 MVP, v1.1 Enhancements, v1.2 Developer Experience

**Recent Trend:**
- v1.2 plans highly efficient: 3-9 min each for documentation and tooling
- Documentation/infrastructure work = consistent velocity

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- **Documentation hierarchy**: README serves as navigation hub with quick starts and links to comprehensive guides in docs/ (prevents README bloat while ensuring discoverability)
- **Deployment section placement**: Positioned between Quick Start and Screenshots for production user visibility
- **MAM authenticated downloads**: MAM URLs trigger torrent file download before qBittorrent upload (required for private tracker authentication)
- **Category support**: Categories are optional parameters passed through to qBittorrent for torrent organization
- **Hash retrieval**: Query qBittorrent API after upload sorted by added_on timestamp (reliable for all torrent sources)
- **Retry strategy**: Max 3 attempts with 500ms delay for torrent info query to handle qBittorrent processing delay
- **Error handling philosophy**: User-friendly messages without exposing internals
- **Context timeout for organization**: 5-minute timeout allows large file operations to complete while preventing indefinite hanging during shutdown
- **Monitor resilience**: Continue monitoring when all downloads fail (qBittorrent may be unavailable), log warnings but don't stop
- **Path mapping for remote deployments**: Single mount point config prepended to qBittorrent paths (supports Docker and network shares)
- **Template validation approach**: Sanitize individual variables before template parsing to preserve directory structure while cleaning filenames
- **Preview debounce timing**: 500ms debounce for preview API calls balances responsiveness with server load
- **Copy operation strategy**: All-or-nothing pattern with automatic cleanup on failure (delete partial copies)
- **Move operation strategy**: Partial success acceptable (atomic per file, already moved files stay moved)
- **Disk space checking**: 10% buffer added to required space for filesystem overhead
- **Interface-based testing**: Refactored OrganizationService to use interfaces for dependency injection and mockability
- **Handler test strategy**: HTTP handler tests focus on request/response layer only, mock all service dependencies
- **Concurrency test approach**: Use channels and synchronization primitives instead of time.Sleep for deterministic tests
- **Race detection requirement**: All concurrency-related code must pass go test -race
- **Structured series data**: Series field changed from concatenated string to []SeriesInfo array with ID, Name, Number fields
- **Docker base image choice**: Alpine over scratch for ca-certificates (HTTPS) and sqlite-libs (database runtime)
- **nginx for frontend serving**: nginx over Node.js for static file serving (more efficient, production-standard)
- **Non-root container security**: Both containers run as non-root user (uid/gid 1001) for security
- **nginx non-privileged port**: nginx listens on port 8080 instead of 80 for non-root compatibility
- **nginx cache directory**: Use /tmp instead of /var/cache/nginx for non-root user write access
- **Series Number as string**: Keep Number field as string to accommodate various formats ("1", "Book 1", "1.5") - frontend handles parsing
- **Empty series array**: Return empty array instead of null for books without series (consistent API responses)
- **Sequential batch processing**: Process batch downloads sequentially rather than concurrently to avoid overwhelming qBittorrent
- **Batch size limit**: 50-item maximum for batch operations to prevent abuse and maintain system stability
- **Partial success pattern**: Batch operations return 200 OK with separate successful/failed arrays for graceful partial failure handling
- **Selection ID strategy**: Use `result.id || result.title` as unique identifier for multi-select to handle results without ID fields
- **Notification granularity**: Three notification states (all success, partial, all failed) provide clear user feedback for batch operations
- **Indeterminate checkbox state**: Series groups show indeterminate checkbox when some (but not all) books are selected
- **Empty series_number in templates**: Replace with empty string rather than removing placeholder (preserves user's template structure)
- **First series as primary**: Books can have multiple series - use first series for folder organization (most relevant)
- **Typed error helpers**: Function-based error helpers over map/enum approach for type safety and clear call site documentation
- **60% coverage baseline**: Established 60% thresholds for both backend and frontend as realistic quality baseline (can increase over time)
- **Functional options for Go fixtures**: Use functional options pattern for test fixtures (more idiomatic than struct overrides)
- **Partial&lt;T&gt; for TypeScript fixtures**: Use Partial&lt;T&gt; pattern for type-safe test data overrides (familiar to TS developers)
- **V8 coverage provider**: Vitest uses v8 provider (faster than c8, built-in support)
- **Mermaid for architecture diagrams**: GitHub-native rendering, version control friendly (system overview, workflow, backend, frontend diagrams)
- **systemd for backend service**: Standard on modern Linux, reliable process supervision with automatic restarts
- **nginx reverse proxy pattern**: Serve frontend static files and proxy API requests to backend
- **golangci-lint pragmatic ruleset**: govet, staticcheck, unused, misspell, goimports (errcheck deferred to phase 15)
- **Prettier for frontend formatting**: Single quotes, no semicolons, 100 char width - separate from ESLint for performance
- **errcheck deferral**: 31 defer Close() violations require refactoring - addressed in phase 15 (refactoring opportunities)
- **Husky for pre-commit hooks**: Standard in Node.js ecosystem, simple setup for running quality checks before commit
- **Parallel CI quality jobs**: Quality checks (lint, format, type-check) run in parallel with tests for fast feedback and fail-fast behavior
- **Pre-commit scope**: Format + lint + type-check (skip tests - too slow), tests remain in CI only
- **defer func() pattern for Close() errors**: Use `defer func() { if err := x.Close(); err != nil { log.Printf(...) } }()` for proper error checking
- **Best-effort logging in error paths**: Database update failures in error handlers are logged but don't override the original error
- **Zero-tolerance error handling policy**: errcheck linter enabled in CI with no exclusions - all errors must be handled or explicitly documented
- **Environment variable for DB path**: ORGANIZR_DB_PATH env var allows configurable database location (defaults to ./organizr.db for backward compatibility)
- **Volume mount at /data**: Mount database volume at /data instead of /app to avoid permission conflicts with non-root user ownership
- **Frontend port mapping**: Map frontend to host port 8081 (container port 8080) to avoid conflict with backend on host port 8080
- **Health check dependencies**: Frontend depends on backend health before starting to ensure database and API readiness
- **godotenv for .env loading**: Simple, focused library for .env file support without config framework overhead
- **Environment variable precedence**: ENV > Database > Defaults pattern enables deployment-time config without database access
- **env_mapping centralization**: Single map of database keys to environment variable names for maintainability
- **Set() database-only behavior**: Config.Set() writes to database only, environment variables are read-only at runtime
- **Volume mount strategy**: Named volumes for development, host path mounts for production (Unraid/NAS deployments)
- **Container path defaults**: PATHS_LOCAL_MOUNT=/downloads, PATHS_DESTINATION=/audiobooks match volume mount paths for zero-config development
- **Port variables in Docker Compose only**: BACKEND_PORT and FRONTEND_PORT control host-to-container mappings but aren't passed to containers (internal ports fixed at 8080)
- **Docker Compose default value syntax**: ${VAR:-default} parameter expansion provides inline defaults without requiring .env file

### Roadmap Evolution

- v1.0 MVP shipped (2026-01-07): 8 plans across 7 phases
- v1.1 Enhancements shipped (2026-01-08): 7 plans across 4 phases
- v1.2 Developer Experience shipped (2026-01-09): 11 plans across 6 phases (10-15)
- v1.3 Production Deployment shipped (2026-01-09): 6 phases (16-21), focus on Docker containerization for Unraid deployment with configurable ports

### Deferred Issues

None

### Blockers/Concerns Carried Forward

None

### Technical Debt

- Phase 07-03 plan created but not executed (pre-download confirmation modal) - could be addressed in future milestone if needed

## Session Continuity

Last session: 2026-01-09T22:32:46Z
Stopped at: Completed 21-01-PLAN.md (v1.3 milestone complete)
Resume file: None
Next action: /gsd:complete-milestone
