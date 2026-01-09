# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-01-08 after v1.1 milestone)

**Core value:** Perfect folder structures and file placement every time - audiobooks land exactly where Audiobookshelf expects them, automatically.
**Current focus:** Developer experience improvements — code quality, documentation, and maintainability (v1.2)

## Current Position

Milestone: v1.2 Developer Experience
Phase: 15 of 15 (Refactoring Opportunities)
Plan: 2 of 2 in current phase
Status: Phase complete
Last activity: 2026-01-09 - Completed 15-02-PLAN.md

Progress: ██████████ 100% (9/9 plans complete in v1.2)

## Performance Metrics

**Velocity (v1.1 milestone):**
- Plans completed: 7
- Average duration: ~6 minutes
- Total execution time: ~0.7 hours
- Milestone duration: 2 days

**Overall Velocity (all milestones):**
- Total plans completed: 15 (8 in v1.0 + 7 in v1.1)
- Cumulative execution time: ~12.5 hours
- Two milestones shipped: v1.0 MVP, v1.1 Enhancements

**Recent Trend:**
- v1.1 plans very efficient: 3-7 min each for focused changes
- Small incremental plans on established codebase = high velocity

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

### Roadmap Evolution

- v1.0 MVP shipped (2026-01-07): 8 plans across 7 phases
- v1.1 Enhancements shipped (2026-01-08): 7 plans across 4 phases
- v1.2 Developer Experience created (2026-01-08): Code quality, documentation, and maintainability, 6 phases (10-15)

### Deferred Issues

None

### Blockers/Concerns Carried Forward

None

### Technical Debt

- Phase 07-03 plan created but not executed (pre-download confirmation modal) - could be addressed in future milestone if needed

## Session Continuity

Last session: 2026-01-09
Stopped at: Completed 15-02-PLAN.md (final plan in Phase 15)
Resume file: None
Next action: v1.2 milestone complete - all 9 plans finished. Ready for milestone completion (/gsd:complete-milestone)
