# Bug Fixes Summary

## Issue 1: Series Information Not Appearing (FIXED ✅)

### Problem
The `formatSeriesInfo` function in `internal/search/providers/mam.go` was failing to parse series information from the MAM API because it expected an array of strings, but MAM actually returns mixed types.

**Example MAM Response:**
```json
{"30281": ["Awaken Online", "10", 10.000000]}
```

The array contains:
- `[0]` - Series name (string)
- `[1]` - Book number (string)
- `[2]` - Numeric value (float)

### Solution
Changed the parsing to use `[]interface{}` instead of `[]string` and added proper type assertions.

**Fixed Code:**
```go
seriesMap := make(map[string][]interface{})  // Was: []string
if name, ok := s[0].(string); ok {
    seriesStr = name
}
```

**Result:** Series information now displays correctly as "Awaken Online (10)"

---

## Issue 2: GetDownloadByID Failing on NULL Fields (FIXED ✅)

### Problem
The `GetByID` method in `internal/persistence/sqlite/downloads.go` was failing when trying to read downloads because nullable database columns (`download_path`, `organized_path`, `series`, etc.) contained NULL values.

**Error:** SQL scan fails when trying to read NULL into a Go `string` type.

### Root Cause
Several database columns are nullable:
- `download_path TEXT` - not set until download starts
- `organized_path TEXT` - not set until organization completes
- `series TEXT` - optional field
- `error_message TEXT` - only set on errors
- `category TEXT` - optional field
- `torrent_url TEXT` - may be NULL if using magnet link
- `magnet_link TEXT` - may be NULL if using torrent URL

But the repository was scanning directly into string fields without NULL handling.

### Solution
Updated all scan operations to use `sql.NullString` for nullable fields:

**Before:**
```go
err := rows.Scan(&d.ID, &d.Title, &d.Author, &d.Series, &d.DownloadPath, ...)
```

**After:**
```go
var series, downloadPath sql.NullString
err := rows.Scan(&d.ID, &d.Title, &d.Author, &series, &downloadPath, ...)
if series.Valid {
    d.Series = series.String
}
if downloadPath.Valid {
    d.DownloadPath = downloadPath.String
}
```

### Files Modified
- `internal/persistence/sqlite/downloads.go`:
  - `GetByID()` - handles 7 nullable fields
  - `GetActive()` - handles nullable series
  - `List()` - handles nullable series

### Testing
Created comprehensive test in `internal/persistence/sqlite/downloads_test.go`:
- Tests downloads with NULL fields
- Tests downloads with populated fields
- Tests mixed scenarios in List and GetActive
- All tests pass ✅

---

## Additional Fixes

### Removed Debug Code
- Removed debug `fmt.Printf` statements from `mam.go` (lines 80-82)

### Fixed Unused Import
- Removed unused `strconv` import from `cmd/test-mam-raw/main.go`

---

## Verification

### Build Status
```bash
✅ go build ./...        # All packages compile
✅ go test ./...         # All tests pass
```

### Test Results
```bash
✅ TestFormatSeriesInfo            # Series parsing works
✅ TestDownloadRepository_NullHandling  # NULL handling works
```

---

## Impact

### Before
- Series information was never displayed (appeared empty)
- GetDownloadByID would fail with SQL scan errors
- Frontend couldn't display download details

### After
- Series information displays correctly: "Awaken Online (10)"
- GetDownloadByID works with both new and existing downloads
- Frontend can properly display all download information
- No breaking changes to API or database schema

---

## Next Steps

No additional changes needed. Both issues are fully resolved and tested.
