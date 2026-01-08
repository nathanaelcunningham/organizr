# Project Milestones: Organizr

## v1.1 Enhancements (Shipped: 2026-01-08)

**Delivered:** Enhanced search and download experience with MAM series detection, batch operations for multiple downloads, and series number organization templates.

**Phases completed:** 7-9 (plus 7.1 inserted) — 7 plans total

**Key accomplishments:**

- MAM Series Detection - Structured series data parsing (ID, Name, Number), search results grouped by series with numerical ordering, support for books in multiple series
- Series Download Fix - Corrected series field to send name only (not "name #number"), preventing incorrect folder organization
- Batch Operations - Backend batch endpoint with 50-item limit and partial success handling, frontend multi-select with floating action bar, series group batch selection with indeterminate checkbox state
- Series Number Organization - Added {series_number} template variable for folder paths (e.g., "Author/Series/1 - Title"), complete integration from MAM → download → organization

**Stats:**

- 46 files created/modified
- +3,636 lines added, -264 lines removed
- 4 phases (3 planned + 7.1 inserted), 7 plans, ~15 tasks
- 2 days from first commit to ship (2026-01-07 → 2026-01-08)

**Git range:** `f1306b8` → `718be00`

**What's next:** User testing with real MAM series searches and batch downloads. Consider v2.0 with advanced features (download scheduling, bandwidth limits, notifications).

---

## v1.0 MVP (Shipped: 2026-01-07)

**Delivered:** Complete audiobook automation from torrent submission through organized files with qBittorrent integration, monitoring, and Audiobookshelf-compatible folder structures.

**Phases completed:** 1-6 (plus 1.1 inserted) — 8 plans total

**Key accomplishments:**

- qBittorrent Integration - Full Web API integration with cookie-based authentication, torrent file upload, MAM authenticated downloads, and hash retrieval
- Download Monitoring - Background monitor with 3-second polling, resilience to qBittorrent unavailability, auto-organization on completion
- Configuration System - Template validation, real-time path preview with 500ms debounce, user-friendly error messages
- File Organization Engine - Automated folder creation with Author/Series/Book structure, disk space validation, partial failure recovery
- Frontend Integration - Real-time status updates, auto-organization toggle, comprehensive test coverage (22 tests), enhanced UX
- End-to-End Testing - Comprehensive backend tests (handlers + monitor), zero race conditions, manual E2E documentation with 7 scenarios

**Stats:**

- 50 files created/modified
- 8,968 lines total (5,060 Go + 3,908 TypeScript)
- 7 phases (6 planned + 1.1 inserted), 8 plans, ~35 tasks
- 6 days from first commit to ship (2026-01-01 → 2026-01-07)

**Git range:** `9dad827` → `70c0475`

**What's next:** Production deployment and real-world usage validation. Monitor for user feedback and potential bug fixes.

---
