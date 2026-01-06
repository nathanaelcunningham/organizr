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

### Active

- [ ] qBittorrent integration - send torrents to qBittorrent Web API
- [ ] Download monitoring - poll qBittorrent for progress and detect completion
- [ ] Configurable folder structure - support Author/Series/Book organization (and other patterns)
- [ ] Auto-organization on completion - create folders and copy files to final destination
- [ ] Real-time UI updates - frontend reflects current download status

### Out of Scope

- Multiple torrent sites - MyAnonamouse only for v1
- Audiobook playback - Audiobookshelf handles that
- Metadata editing - no tagging, cover art, or file modification
- Multi-user support - single-user tool

## Context

**What already exists:**
- Full-stack application with Go backend and React frontend
- MyAnonamouse search provider with API integration
- qBittorrent client code structure in place
- SQLite database with WAL mode for persistence
- Background monitor pattern via goroutines
- Frontend state management via Zustand stores

**What needs to be built:**
- Complete qBittorrent integration (authentication, torrent submission, status polling)
- Download lifecycle management (monitoring, completion detection, status updates)
- File organization engine (path templates, folder creation, file copying)
- Configuration system for folder patterns and destination paths

**Known patterns to follow:**
- Repository pattern for data access
- Service layer for business logic
- Background monitoring via goroutines
- Frontend polling for real-time updates (3-second intervals)

## Constraints

- **Torrent Site**: MyAnonamouse (MAM) - private tracker with specific API requirements
- **Download Client**: qBittorrent Web API - must integrate via HTTP API (cookie-based authentication)
- **File Organization**: Configurable folder structure templates - must support Author/Series/Book pattern as default
- **Tech Stack**: Go backend + React frontend (established, no changes)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| MyAnonamouse as single torrent source | Focused scope for v1, MAM is primary use case | — Pending |
| qBittorrent Web API integration | Standard API, cookie-based auth, widely used | — Pending |
| Configurable folder templates | Author/Series/Book preferred but needs flexibility | — Pending |
| Background monitoring with polling | Matches existing architecture patterns | — Pending |

---
*Last updated: 2026-01-06 after initialization*
