---
phase: 15-refactoring-opportunities
plan: 01
subsystem: infrastructure
tags: [error-handling, code-quality, linting, errcheck, golang]

# Dependency graph
requires:
  - phase: 14-code-quality-tools
    provides: golangci-lint configuration and CI integration
provides:
  - Comprehensive error handling for all defer Close() calls
  - Proper error checking for all I/O operations
  - Error handling for database operations in error paths
  - errcheck linter enabled in CI pipeline
affects: [all-phases]

# Tech tracking
tech-stack:
  added: []
  patterns: [defer-func-error-checking, best-effort-logging-in-error-paths]

key-files:
  created: []
  modified:
    - backend/internal/downloads/organization.go
    - backend/internal/persistence/sqlite/config.go
    - backend/internal/persistence/sqlite/downloads.go
    - backend/internal/persistence/sqlite/downloads_test.go
    - backend/internal/qbittorrent/client.go
    - backend/internal/search/providers/mam.go
    - backend/internal/search/search_service_integration_test.go
    - backend/cmd/api/main.go
    - backend/internal/downloads/monitor.go
    - backend/internal/downloads/monitor_test.go
    - backend/internal/server/errors.go
    - backend/internal/server/handlers_test.go
    - .golangci.yml

key-decisions:
  - "Use defer func() pattern to check close errors and log them appropriately"
  - "File operations log close errors as they may indicate I/O issues"
  - "Database operations log close errors as they may indicate connection issues"
  - "HTTP response bodies log close errors as they may indicate network issues"
  - "Best-effort error logging in error paths (log update failures, don't fail the failure handler)"
  - "Explicitly ignore parse errors in test code with _ = ... and comments explaining why"

patterns-established:
  - "defer func() { if err := x.Close(); err != nil { log.Printf(...) } }()"
  - "Best-effort database updates in error paths with logging"
  - "Document ignored errors with explicit _ = and comment"

issues-created: []

# Metrics
duration: 9min
completed: 2026-01-09
---

# Phase 15 Plan 1: Error Handling Cleanup Summary

**Eliminated 74 errcheck violations, enabled errcheck linter in CI for zero-tolerance error handling policy**

## Performance

- **Duration:** 9 min
- **Started:** 2026-01-09T20:15:51Z
- **Completed:** 2026-01-09T20:25:43Z
- **Tasks:** 3
- **Files modified:** 13

## Accomplishments

- Fixed all defer Close() error handling (40+ violations across file I/O, database operations, and HTTP responses)
- Fixed ignored I/O errors in HTTP response handling (io.ReadAll calls)
- Fixed update error handling in monitor error paths with best-effort logging
- Fixed fmt.Sscanf and json.Encoder.Encode error handling
- Enabled errcheck linter in .golangci.yml configuration
- Removed all errcheck exclusions - zero tolerance for unchecked errors going forward

## Task Commits

Each task was committed atomically:

1. **Task 1: Fix defer Close() error handling** - `db78b37` (refactor)
2. **Task 2: Fix ignored I/O errors** - (completed within Task 1)
3. **Task 3: Fix remaining violations and enable errcheck** - `aa78fe4` (refactor)

**Plan metadata:** (to be committed)

## Files Created/Modified

- `backend/internal/downloads/organization.go` - Added defer func() error checking for file Close() calls in copyFile
- `backend/internal/persistence/sqlite/config.go` - Added defer func() error checking for rows.Close()
- `backend/internal/persistence/sqlite/downloads.go` - Added defer func() error checking for rows.Close() in GetActive and List
- `backend/internal/persistence/sqlite/downloads_test.go` - Added defer func() error checking for db.Close()
- `backend/internal/qbittorrent/client.go` - Added defer func() error checking for all HTTP response Body.Close() calls (6 locations)
- `backend/internal/search/providers/mam.go` - Added defer func() error checking for HTTP Body.Close() + fixed io.ReadAll error handling (3 locations)
- `backend/internal/search/search_service_integration_test.go` - Added defer func() error checking for db.Close()
- `backend/cmd/api/main.go` - Added defer func() error checking for db.Close()
- `backend/internal/downloads/monitor.go` - Fixed UpdateError/UpdateStatus calls in error paths with proper error logging
- `backend/internal/downloads/monitor_test.go` - Fixed UpdateError/UpdateStatus and fmt.Sscanf with explicit error handling
- `backend/internal/server/errors.go` - Added error checking for json.Encoder.Encode() calls with logging
- `backend/internal/server/handlers_test.go` - Fixed fmt.Sscanf error handling with explicit ignoring and comments
- `.golangci.yml` - Enabled errcheck linter, removed exclusions, cleaned up config

## Decisions Made

1. **defer func() pattern for Close() errors**: Using `defer func() { if err := x.Close(); err != nil { log.Printf(...) } }()` pattern allows proper error checking while maintaining cleanup guarantee. Errors are logged because close failures may indicate underlying I/O, network, or database issues that should be visible in logs.

2. **Best-effort logging in error paths**: In error handlers (e.g., organization failure in monitor), database update failures are logged but don't override the original error. Pattern: `if updateErr := repo.Update(...); updateErr != nil { log.Printf(...) }`. This prevents cascading failures where error handling itself fails.

3. **Explicit error ignoring with comments**: In test code where parse errors are truly unimportant (e.g., fmt.Sscanf in tests), use `_, _ = fmt.Sscanf(...)` with a comment explaining why it's safe to ignore. This documents the decision and satisfies errcheck.

4. **Zero tolerance policy**: Enabled errcheck in CI with no exclusions. All errors must be either handled or explicitly documented as safely ignorable. This prevents error handling debt from accumulating.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - all errcheck violations were straightforward to fix with consistent patterns.

## Next Phase Readiness

- All error handling technical debt cleared
- errcheck linter now running in pre-commit hooks and CI
- Codebase has zero errcheck violations
- Ready for phase 15-02 (code quality improvements)

---
*Phase: 15-refactoring-opportunities*
*Completed: 2026-01-09*
