# MAM Series Debugging Guide

This document explains how to use the debugging tools to investigate why series information isn't appearing in search results.

## Available Tools

### 1. Simple CLI Test Tool (Recommended)
**File:** `cmd/test-mam/main.go`

The easiest way to test MAM search with real API calls.

**Usage:**
```bash
cd backend
go run ./cmd/test-mam -query "wheel of time"
```

**Options:**
- `-query`: Search query (default: "wheel of time")
- `-db`: Path to database (default: "./organizr.db")
- `-debug`: Show additional debug info

**Example output:**
```
🔍 Testing MAM Search API
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📁 Connecting to database: ./organizr.db
✓ Database connected
✓ MAM configuration loaded

📊 SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total results:          50
With series info:       35 (70.0%)
Without series info:    15 (30.0%)
```

### 2. Raw API Response Inspector
**File:** `cmd/test-mam-raw/main.go`

Shows the exact JSON response from MAM API, including the raw `series_info` field.

**Usage:**
```bash
cd backend
go run ./cmd/test-mam-raw -query "mistborn" -limit 5
```

**Options:**
- `-query`: Search query
- `-db`: Path to database
- `-limit`: Number of results to show in detail (default: 3)

**What it shows:**
- Raw API response
- Exact `series_info` JSON from MAM
- How the `formatSeriesInfo` function parses it
- Comparison of what MAM sends vs. what the app displays

**Example output:**
```json
Raw series_info:
  {
    "123": ["Mistborn", "Book 1"],
    "456": ["The Final Empire"]
  }

Parsed series (what the app shows):
  "Mistborn (Book 1), The Final Empire"
```

### 3. Integration Tests
**File:** `internal/search/search_service_integration_test.go`

Go test suite that calls the real API.

**Usage:**
```bash
cd backend
go test -tags=integration -v ./internal/search -run TestMAMSearchIntegration
```

## How to Debug

### Step 1: Run the Simple Test
```bash
cd backend
go run ./cmd/test-mam -query "wheel of time"
```

This will show you:
- If series info is coming back at all
- Percentage of results with series data
- First 5 results with their series field

### Step 2: Inspect Raw API Response
```bash
go run ./cmd/test-mam-raw -query "wheel of time" -limit 3
```

This will show you:
- The exact JSON that MAM returns
- Whether `series_info` field is empty or populated
- How the parsing function transforms it

### Step 3: Analyze the Results

**If series_info is EMPTY in the raw response:**
- MAM doesn't have series metadata for these items
- Try different search queries (popular book series)
- The issue is with MAM's data, not our code

**If series_info HAS DATA but appears empty in the app:**
- There's a bug in the `formatSeriesInfo` function (line 278 in `internal/search/providers/mam.go`)
- Check the parsing logic
- The raw inspector will show exactly what's failing

**If series_info has data AND appears in the CLI tool but NOT in the frontend:**
- The backend is working correctly
- Check the frontend code that displays search results
- Verify the API response is being parsed correctly in the React app

## Common Issues

### 1. "Config key not found: mam.baseurl"
You need to configure MAM credentials first. Run the backend and set them via the settings page.

### 2. "Authentication failed"
Your MAM secret (cookie) may have expired. Log into MAM in your browser, get a new `mam_id` cookie, and update the config.

### 3. No results with series
Some audiobooks on MAM don't have series metadata. Try searching for:
- "wheel of time"
- "stormlight archive"
- "mistborn"
- "harry potter"

These are well-known series that should have metadata.

## Understanding series_info Format

MAM returns series information as a JSON object:

```json
{
  "series_id": ["Series Name", "Book Number"]
}
```

Examples:
```json
// Single series with book number
{"123": ["The Wheel of Time", "Book 1"]}

// Single series without book number
{"456": ["Harry Potter"]}

// Multiple series
{
  "123": ["Series A", "1"],
  "456": ["Series B", "Book 2"]
}
```

Our `formatSeriesInfo` function converts this to:
- `"The Wheel of Time (Book 1)"`
- `"Harry Potter"`
- `"Series A (1), Series B (Book 2)"`

## Saving Results

To save the full JSON response to a file:
```bash
SAVE_RESULTS=1 go run ./cmd/test-mam -query "wheel of time"
```

This creates `mam_search_results.json` with all results.
