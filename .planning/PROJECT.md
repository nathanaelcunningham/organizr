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

### Active

(None - v1.0 MVP shipped)

### Out of Scope

- Multiple torrent sites - MyAnonamouse only for v1
- Audiobook playback - Audiobookshelf handles that
- Metadata editing - no tagging, cover art, or file modification
- Multi-user support - single-user tool

## Context

**Current State (v1.0 shipped 2026-01-07):**
- **Codebase**: 8,968 LOC (5,060 Go backend + 3,908 TypeScript frontend)
- **Tech Stack**: Go (Chi router, SQLite with WAL), React (TypeScript, Vite, Zustand), qBittorrent Web API
- **Features**: Full qBittorrent integration with MAM authenticated downloads, background monitoring with 3-second polling, configurable folder templates with real-time preview, auto-organization on completion, comprehensive test coverage
- **Testing**: 22 frontend tests (Vitest), comprehensive backend tests (handlers, monitor, organization), zero race conditions, manual E2E documentation with 7 scenarios
- **Status**: Production-ready, ready for deployment and real-world usage

**Established Patterns:**
- Repository pattern for data access
- Service layer with interface-based dependency injection
- Background monitoring via goroutines with resilience
- Frontend polling for real-time updates (3-second intervals, auto-stops when idle)
- Template validation with real-time preview
- All-or-nothing copy operations with automatic cleanup
- Comprehensive error handling with user-friendly messages

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

---
*Last updated: 2026-01-07 after v1.0 milestone*
