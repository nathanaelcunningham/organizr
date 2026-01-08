---
phase: 09-series-number-organization
plan: 02
subsystem: file-organization
tags: [template-system, path-preview, series-numbers]

# Dependency graph
requires:
  - phase: 09-01
    provides: SeriesNumber field in Download model and database
  - phase: 03-01
    provides: Template validation and path preview infrastructure
  - phase: 04-01
    provides: Organization service with template parsing

provides:
  - {series_number} template variable support in validation
  - series_number included in organization path generation
  - series_number support in path preview endpoint

affects: [09-03-frontend-integration]

# Tech tracking
tech-stack:
  added: []
  patterns: []

key-files:
  created: []
  modified:
    - backend/internal/fileutil/template_test.go
    - backend/internal/downloads/organization.go
    - backend/internal/downloads/organization_test.go
    - backend/internal/server/handlers.go
    - backend/internal/server/request_types.go

key-decisions:
  - "Empty series_number replaced with empty string in templates (preserves user's template structure)"

patterns-established: []

issues-created: []

# Metrics
duration: 3min
completed: 2026-01-08
---

# Phase 9 Plan 2: Template and Organization Summary

**{series_number} template variable fully integrated - validation, organization, and path preview all support series numbers**

## Performance

- **Duration:** 3 min
- **Started:** 2026-01-08T13:57:37Z
- **Completed:** 2026-01-08T14:00:48Z
- **Tasks:** 3
- **Files modified:** 5

## Accomplishments

- Template validation supports {series_number} variable with test coverage
- Organization service includes series_number in sanitized template vars
- Path preview endpoint validates and uses series_number for real-time preview
- Empty series_number handled gracefully (replaced with empty string)

## Task Commits

Each task was committed atomically:

1. **Task 1: Add series_number to template validation** - `acc4bb3` (test)
2. **Task 2: Update organization logic to use series_number in templates** - `683880a` (feat)
3. **Task 3: Update path preview endpoint to support series_number** - `b521fc0` (feat)

**Plan metadata:** (will be added in docs commit)

## Files Created/Modified

- `backend/internal/fileutil/template_test.go` - Added 3 series_number validation test cases
- `backend/internal/downloads/organization.go` - Added series_number to sanitizedVars map
- `backend/internal/downloads/organization_test.go` - Added 2 test cases for series_number paths
- `backend/internal/server/handlers.go` - Added series_number to allowedVars and sanitized vars
- `backend/internal/server/request_types.go` - Added SeriesNumber field to PreviewPathRequest

## Decisions Made

**Empty series_number handling:** When SeriesNumber is empty string, it's replaced with empty string in templates. This preserves the user's template structure (e.g., "{series_number} - {title}" becomes " - Title") rather than collapsing or removing the placeholder. Users have control over their template format.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## Next Phase Readiness

Backend fully ready for {series_number} template variable:
- Template validation accepts it
- Organization uses it in folder paths
- Preview endpoint displays it in real-time

Ready for 09-03-PLAN.md (Frontend Integration) - add series_number input to settings UI and search result display.

---
*Phase: 09-series-number-organization*
*Completed: 2026-01-08*
