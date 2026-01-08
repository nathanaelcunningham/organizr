# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-01-07 after v1.0 milestone)

**Core value:** Perfect folder structures and file placement every time - audiobooks land exactly where Audiobookshelf expects them, automatically.
**Current focus:** v1.1 Enhancements — MAM series support, batch operations, series number templates

## Current Position

Milestone: v1.1 Enhancements
Phase: 8 of 9 (Batch Operations)
Plan: 2 of 2 in current phase
Status: Phase complete
Last activity: 2026-01-08 - Completed 08-02-PLAN.md

Progress: ████░░░░░░ 28%

## Performance Metrics

**Velocity:**
- Total plans completed: 13
- Average duration: ~53 minutes
- Total execution time: 11.6 hours

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 1 | 2 | 1.0h | 30m |
| 1.1 | 1 | 0.1h | 6m |
| 2 | 1 | 6.0h | 357m |
| 3 | 1 | 0.4h | 22m |
| 4 | 1 | 0.8h | 45m |
| 5 | 1 | 1.5h | 90m |
| 6 | 1 | 0.3h | 16m |
| 7 | 2 | 0.4h | 11m |
| 7.1 | 1 | 0.05h | 3m |
| 8 | 2 | 0.17h | 5m |

**Recent Trend:**
- Last 5 plans: 07-02 (17m), 07.1-01 (3m), 08-01 (3m), 08-02 (7m)
- Trend: Small focused plans very efficient (3-7 min each)

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

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
- **Series Number as string**: Keep Number field as string to accommodate various formats ("1", "Book 1", "1.5") - frontend handles parsing
- **Empty series array**: Return empty array instead of null for books without series (consistent API responses)
- **Sequential batch processing**: Process batch downloads sequentially rather than concurrently to avoid overwhelming qBittorrent
- **Batch size limit**: 50-item maximum for batch operations to prevent abuse and maintain system stability
- **Partial success pattern**: Batch operations return 200 OK with separate successful/failed arrays for graceful partial failure handling
- **Selection ID strategy**: Use `result.id || result.title` as unique identifier for multi-select to handle results without ID fields
- **Notification granularity**: Three notification states (all success, partial, all failed) provide clear user feedback for batch operations
- **Indeterminate checkbox state**: Series groups show indeterminate checkbox when some (but not all) books are selected

### Roadmap Evolution

- **Phase 1.1 inserted after Phase 1** (2026-01-06): qBittorrent Connection Testing (URGENT)
  - Reason: ISS-001 discovered during Phase 1 testing - users need diagnostic tool for qBittorrent connectivity
  - Impact: Addresses testing gaps before building Phase 2 monitoring (which depends on reliable qBittorrent connection)
- **Milestone v1.1 created**: Enhancements focus, 3 phases (Phase 7-9)
  - Theme: MAM series detection, batch operations, series number organization
- **Phase 7.1 inserted after Phase 7** (2026-01-07): Fix Series Download Field (URGENT)
  - Reason: Download requests send series with number ("Discworld #1") instead of name only ("Discworld"), breaking folder organization
  - Impact: Critical fix for file organization - must be resolved before Phase 8 batch operations to avoid compounding the issue

### Deferred Issues

None - ISS-001 resolved in Phase 1.1

### Blockers/Concerns Carried Forward

None.

## Session Continuity

Last session: 2026-01-08 04:09
Stopped at: Completed 08-02-PLAN.md (Phase 8 complete)
Resume file: None
Next action: Plan Phase 9 (Series Number Organization) — /gsd:plan-phase 9
