# Organizr

## What This Is

A tool for automating audiobook torrent downloads and organization. Searches MyAnonamouse by author, title, or series; sends torrents to qBittorrent; monitors download progress; and automatically organizes completed files into Audiobookshelf-compatible folder structures (Author/Series/Book). Eliminates manual file management for audiobook collectors.

## Core Value

Perfect folder structures and file placement every time - audiobooks land exactly where Audiobookshelf expects them, automatically.

## Requirements

### Validated

- ✓ Search interface with MyAnonamouse integration - existing
- ✓ Frontend application (React + TypeScript + Vite) - existing
- ✓ Backend API (Go + Chi router + SQLite) - existing
- ✓ Download tracking data model and persistence - existing
- ✓ qBittorrent client wrapper code - existing
- ✓ qBittorrent integration - send torrents to qBittorrent Web API — v1.0
- ✓ Download monitoring - poll qBittorrent for progress and detect completion — v1.0
- ✓ Configurable folder structure - support Author/Series/Book organization (and other patterns) — v1.0
- ✓ Auto-organization on completion - create folders and copy files to final destination — v1.0
- ✓ Real-time UI updates - frontend reflects current download status — v1.0
- ✓ MAM series detection - parse and display series information from search results — v1.1
- ✓ Batch download operations - add multiple torrents simultaneously — v1.1
- ✓ Series number organization - support {series_number} in folder templates — v1.1
- ✓ Architecture documentation - ADR, contribution guidelines, codebase reference — v1.2
- ✓ API standardization - typed error helpers, OpenAPI/Swagger documentation — v1.2
- ✓ Testing infrastructure - test utilities, fixtures, coverage reporting — v1.2
- ✓ Developer documentation - README, troubleshooting, architecture diagrams, deployment guide — v1.2
- ✓ Code quality automation - linting, formatting, pre-commit hooks, CI workflows — v1.2
- ✓ Error handling cleanup - zero errcheck violations, proper Close() handling — v1.2

### Active

- [ ] Improved user experience - better visual feedback and status display
- [ ] Performance optimizations - reduce API calls and improve responsiveness

### Out of Scope

- Multiple torrent sites - MyAnonamouse only for v1
- Audiobook playback - Audiobookshelf handles that
- Metadata editing - no tagging, cover art, or file modification
- Multi-user support - single-user tool

## Context

**Current State (v1.2 shipped 2026-01-09):**
- **Codebase**: ~12,291 LOC (Go backend + TypeScript frontend)
- **Tech Stack**: Go (Chi router, SQLite with WAL), React (TypeScript, Vite, Zustand), qBittorrent Web API
- **Features**: Full qBittorrent integration, background monitoring, configurable folder templates with {series_number} support, auto-organization, MAM series detection with grouped display, batch downloads with multi-select UI
- **Testing**: Comprehensive backend tests (handlers, monitor, organization, batch), frontend tests (Vitest), test utilities and fixtures, 60% coverage baseline, zero race conditions
- **Documentation**: Complete README, ADR, CONTRIBUTING guide, API docs (Swagger), architecture diagrams, deployment guide, troubleshooting guide
- **Quality**: golangci-lint + Prettier automation, pre-commit hooks, CI quality gates, zero errcheck violations
- **Status**: Production-ready with professional developer experience and maintainable codebase

**Established Patterns:**
- Repository pattern for data access
- Service layer with interface-based dependency injection
- Background monitoring via goroutines with resilience
- Frontend polling for real-time updates (3-second intervals, auto-stops when idle)
- Template validation with real-time preview
- All-or-nothing copy operations with automatic cleanup
- Comprehensive error handling with user-friendly messages
- Structured data models over string concatenation (SeriesInfo with ID, Name, Number)
- Partial success patterns for batch operations
- Client-side grouping and sorting for search results
- Typed error helpers for consistent API responses
- Test utilities and fixtures for reduced boilerplate
- Architecture Decision Record for documenting technical choices
- Pre-commit hooks and CI workflows for automated quality checks
- Zero-tolerance error handling policy (errcheck enabled)

## Constraints

- **Torrent Site**: MyAnonamouse (MAM) - private tracker with specific API requirements
- **Download Client**: qBittorrent Web API - must integrate via HTTP API (cookie-based authentication)
- **File Organization**: Configurable folder structure templates - must support Author/Series/Book pattern as default
- **Tech Stack**: Go backend + React frontend (established, no changes)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| MyAnonamouse as single torrent source | Focused scope for v1, MAM is primary use case | ✓ Good - Clean integration, can add more sources in future |
| qBittorrent Web API integration | Standard API, cookie-based auth, widely used | ✓ Good - Robust integration with retry logic and resilience |
| Configurable folder templates | Author/Series/Book preferred but needs flexibility | ✓ Good - Template validation prevents errors, real-time preview excellent UX |
| Background monitoring with polling | Matches existing architecture patterns | ✓ Good - 3-second polling, auto-stops when idle, resilient to qBittorrent unavailability |
| MAM authenticated downloads | MAM URLs trigger torrent file download before qBittorrent upload | ✓ Good - Required for private tracker authentication |
| Hash retrieval via timestamp | Query qBittorrent API sorted by added_on timestamp | ✓ Good - Reliable for all torrent sources |
| Context timeout for organization | 5-minute timeout for large file operations | ✓ Good - Prevents indefinite hanging during shutdown |
| Interface-based testing | Refactored OrganizationService for dependency injection | ✓ Good - Enables comprehensive mocking and test coverage |
| Race detection requirement | All concurrency code must pass go test -race | ✓ Good - Catches data races early, zero races in v1.0 |
| Structured series data | SeriesInfo[] with ID, Name, Number instead of concatenated string | ✓ Good - Enables sorting, grouping, and better data fidelity |
| Sequential batch processing | Process batch downloads sequentially rather than concurrently | ✓ Good - Prevents overwhelming qBittorrent, more predictable behavior |
| Partial success pattern | Batch operations return 200 OK with separate success/failed arrays | ✓ Good - Graceful partial failure handling, clear user feedback |
| First series as primary | Books in multiple series use first for organization | ✓ Good - Simple rule, typically the main series is listed first |
| Empty series_number handling | Replace with empty string in templates | ✓ Good - Preserves user's template structure, simple and predictable |
| Typed error helpers | Function-based error helpers over map/enum approach | ✓ Good - Type safety and clear call site documentation |
| 60% coverage baseline | Established 60% thresholds for both backend and frontend | ✓ Good - Realistic quality baseline, can increase over time |
| Functional options for Go fixtures | Use functional options pattern for test fixtures | ✓ Good - More idiomatic than struct overrides |
| Partial&lt;T&gt; for TypeScript fixtures | Use Partial&lt;T&gt; pattern for type-safe test data overrides | ✓ Good - Familiar to TS developers |
| V8 coverage provider | Vitest uses v8 provider (faster than c8) | ✓ Good - Built-in support, better performance |
| Mermaid for architecture diagrams | GitHub-native rendering, version control friendly | ✓ Good - Easy to maintain, automatically rendered |
| golangci-lint pragmatic ruleset | govet, staticcheck, unused, misspell, goimports | ✓ Good - Focus on correctness, errcheck added in v1.2 |
| Prettier for frontend formatting | Single quotes, no semicolons, 100 char width | ✓ Good - Separate from ESLint for performance |
| Husky for pre-commit hooks | Standard in Node.js ecosystem | ✓ Good - Simple setup, reliable execution |
| Parallel CI quality jobs | Quality checks run in parallel with tests | ✓ Good - Fast feedback and fail-fast behavior |
| Pre-commit scope | Format + lint + type-check (skip tests) | ✓ Good - Tests too slow for pre-commit, remain in CI |
| defer func() pattern for Close() errors | Use defer func() with proper error logging | ✓ Good - Proper error checking with cleanup guarantee |
| Zero-tolerance error handling policy | errcheck enabled in CI with no exclusions | ✓ Good - Prevents error handling debt from accumulating |

---
*Last updated: 2026-01-09 after v1.2 milestone*
