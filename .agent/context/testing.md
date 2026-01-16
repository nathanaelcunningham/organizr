# Testing Approach

## Testing Philosophy

**Pragmatic testing with focus on critical paths**

- Prioritize testing integration points (API handlers, repositories, external services)
- Test complex business logic (path templates, series detection, organization)
- Test error handling and edge cases
- Current coverage: ~20% backend (pragmatic approach), 60% frontend threshold

**What to test:**
- API endpoint behavior (handlers_test.go)
- Repository operations (downloads_test.go, config_test.go)
- Service orchestration (downloads service, monitor)
- External integrations (MAM API, qBittorrent client)
- Complex utilities (template parsing, path sanitization)
- State management (Zustand stores)

**What not to test:**
- Simple getters/setters
- Trivial DTOs or models
- Third-party library behavior
- Obvious one-liners

## Test Frameworks

### Backend Testing
**Framework:** Go standard library `testing` package

**Key Packages:**
- `testing` - Core test framework
- `net/http/httptest` - HTTP handler testing (mock requests/responses)
- `context` - Context-aware testing
- `database/sql` - In-memory SQLite for repository tests

**Test File Pattern:** `*_test.go` (same directory as code under test)

**Example Test Structure:**
```go
func TestCreateDownload(t *testing.T) {
    // Setup
    repo := setupTestRepo(t)
    service := downloads.NewDownloadService(repo, mockQBitClient, mockConfig)

    // Execute
    download, err := service.CreateDownload(ctx, &models.CreateDownloadRequest{
        Title: "Test Book",
        // ...
    })

    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if download.Status != models.StatusQueued {
        t.Errorf("expected status queued, got %v", download.Status)
    }
}
```

### Frontend Testing
**Framework:** Vitest v4.0.0

**Key Packages:**
- `vitest` - Test runner (Jest-compatible API)
- `@testing-library/react` - React component testing utilities
- `@testing-library/jest-dom` - DOM matchers
- `happy-dom` - Lightweight DOM implementation (faster than jsdom)

**Test File Pattern:** `*.test.ts` or `*.test.tsx`

**Coverage Thresholds:** 60% for lines, functions, branches, statements
- Configured in vitest.config.ts
- Run `npm run test:coverage` to check

**Example Test Structure:**
```typescript
import { renderHook, waitFor } from '@testing-library/react'
import { useDownloadStore } from './useDownloadStore'

describe('useDownloadStore', () => {
  it('fetches downloads on mount', async () => {
    const { result } = renderHook(() => useDownloadStore())

    await waitFor(() => {
      expect(result.current.downloads).toHaveLength(2)
    })
  })
})
```

## Test Organization

### Backend Test Files
```
backend/internal/
├── server/
│   ├── handlers.go
│   └── handlers_test.go          # API endpoint tests
├── persistence/sqlite/
│   ├── downloads.go
│   └── downloads_test.go         # Repository CRUD tests
├── downloads/
│   ├── service.go
│   ├── service_test.go           # Service orchestration tests
│   ├── monitor.go
│   ├── monitor_test.go           # Monitor polling tests
│   ├── organization.go
│   └── organization_test.go      # File organization tests
├── search/providers/
│   ├── mam.go
│   └── mam_test.go               # MAM API integration tests
├── config/
│   ├── service.go
│   └── service_test.go           # Config service tests
└── fileutil/
    ├── template.go
    ├── template_test.go          # Template parsing tests
    ├── sanitizer.go
    └── sanitizer_test.go         # Path sanitization tests
```

**Naming Convention:**
- Test file: `<name>_test.go` (e.g., `downloads_test.go`)
- Test function: `Test<FunctionName>` (e.g., `TestCreateDownload`)
- Subtests: `t.Run("subtest name", func(t *testing.T) { ... })`

### Frontend Test Files
```
frontend/src/
├── stores/
│   ├── useDownloadStore.ts
│   └── useDownloadStore.test.ts  # Zustand store tests
├── utils/
│   ├── groupSeries.ts
│   └── groupSeries.test.ts       # Utility function tests
└── components/
    ├── SearchForm.tsx
    └── SearchForm.test.tsx       # Component tests (future)
```

**Naming Convention:**
- Test file: `<name>.test.ts` or `<name>.test.tsx`
- Test suite: `describe('<ComponentName>', () => { ... })`
- Test case: `it('should do something', () => { ... })`

## Test Patterns

### Backend Unit Tests

**Repository Tests (with in-memory SQLite):**
```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to create test db: %v", err)
    }

    // Run migrations
    // ...

    return db
}

func TestGetDownload(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := sqlite.NewDownloadRepository(db)

    // Test implementation
}
```

**Handler Tests (with httptest):**
```go
func TestHandleCreateDownload(t *testing.T) {
    // Create test server with mocked dependencies
    server := &Server{
        downloadService: mockDownloadService,
    }

    // Create test request
    body := `{"title":"Test","torrent_url":"http://example.com/test.torrent"}`
    req := httptest.NewRequest(http.MethodPost, "/api/downloads", strings.NewReader(body))
    w := httptest.NewRecorder()

    // Execute handler
    server.CreateDownload(w, req)

    // Assert response
    if w.Code != http.StatusCreated {
        t.Errorf("expected status 201, got %d", w.Code)
    }
}
```

**Service Tests (with mocked dependencies):**
```go
type mockDownloadRepo struct {
    downloads []*models.Download
}

func (m *mockDownloadRepo) CreateDownload(ctx context.Context, d *models.Download) error {
    m.downloads = append(m.downloads, d)
    return nil
}

func TestDownloadService_CreateDownload(t *testing.T) {
    mockRepo := &mockDownloadRepo{}
    mockQBit := &mockQBitClient{}
    service := downloads.NewDownloadService(mockRepo, mockQBit, nil)

    // Test implementation
}
```

### Frontend Unit Tests

**Zustand Store Tests:**
```typescript
describe('useDownloadStore', () => {
  beforeEach(() => {
    // Reset store state
    useDownloadStore.setState({
      downloads: [],
      loading: false,
      error: null,
    })
  })

  it('sets loading state when fetching', async () => {
    const { result } = renderHook(() => useDownloadStore())

    act(() => {
      result.current.fetchDownloads()
    })

    expect(result.current.loading).toBe(true)

    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })
  })
})
```

**Utility Function Tests:**
```typescript
describe('groupSeries', () => {
  it('groups results by series', () => {
    const results: SearchResult[] = [
      { id: '1', title: 'Book 1', series: [{ id: 's1', name: 'Series A', number: '1' }] },
      { id: '2', title: 'Book 2', series: [{ id: 's1', name: 'Series A', number: '2' }] },
    ]

    const groups = groupSeries(results)

    expect(groups).toHaveLength(1)
    expect(groups[0].series?.name).toBe('Series A')
    expect(groups[0].results).toHaveLength(2)
  })
})
```

### Integration Tests

**Backend Integration Tests (if added):**
- Test full request lifecycle (handler → service → repository → database)
- Use real SQLite database (not mocked)
- Test external service integration (qBittorrent, MAM) with test servers

**Frontend Integration Tests (if added):**
- Test component + store + API client interaction
- Mock API responses at fetch level
- Test user workflows (search → select → download)

## Mocking/Stubbing

### Backend Mocking Approach

**Prefer interfaces over concrete types:**
```go
// Good - service depends on interface
type DownloadService struct {
    repo persistence.DownloadRepository  // Interface
}

// In tests, create mock implementation
type mockDownloadRepo struct {
    // Mock fields
}

func (m *mockDownloadRepo) CreateDownload(ctx context.Context, d *models.Download) error {
    // Mock implementation
}
```

**External Services:**
- qBittorrent: Mock HTTP server with httptest.NewServer or mock client
- MAM API: Mock HTTP responses or mock provider interface

**Database:**
- Use in-memory SQLite (`:memory:`) for repository tests
- Run actual migrations for realistic testing

### Frontend Mocking Approach

**API Client Mocking:**
```typescript
// Mock fetch at global level
global.fetch = vi.fn(() =>
  Promise.resolve({
    ok: true,
    status: 200,
    json: async () => ({ downloads: [] }),
  })
)
```

**Component Props:**
- Pass mock functions as props
- Verify function calls with `expect(mockFn).toHaveBeenCalledWith(...)`

## Running Tests

### Backend Commands

```bash
# Run all tests
make test
# or
go test ./...

# Run tests with coverage
make test-coverage
# or
go test -cover ./...

# Generate HTML coverage report
make test-coverage
# Opens coverage.html in browser

# Run tests with race detector
make test-race
# or
go test -race ./...

# Run specific test
go test ./internal/server -run TestCreateDownload

# Verbose output
go test -v ./...
```

### Frontend Commands

```bash
# Run all tests (watch mode)
npm test

# Run tests once (CI mode)
npm run test:run

# Run with coverage
npm run test:coverage

# Run specific test file
npm test -- useDownloadStore.test.ts

# Update snapshots (if using)
npm test -- -u
```

## CI/CD Testing

**Current Setup:** Not explicitly configured (can be added)

**Recommended CI Workflow:**

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      - run: cd backend && go test -race -cover ./...

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '24'
      - run: cd frontend && npm ci
      - run: cd frontend && npm run test:run
      - run: cd frontend && npm run test:coverage
```

**Pre-commit Hooks:** Not currently configured (can be added with husky)

## Test Data & Fixtures

### Backend Test Fixtures

**Location:** backend/internal/testutil/

**Common Fixtures:**
- Sample Download entities
- Sample SearchResult data
- Mock qBittorrent responses
- Mock MAM API responses

**Example:**
```go
// testutil/fixtures.go
func NewTestDownload() *models.Download {
    return &models.Download{
        ID:         "test-id-123",
        Title:      "Test Book",
        Author:     "Test Author",
        Status:     models.StatusQueued,
        CreatedAt:  time.Now(),
    }
}
```

### Frontend Test Fixtures

**Location:** frontend/src/test/

**Common Fixtures:**
- Sample Download objects
- Sample SearchResult data
- Mock API responses
- Test utilities

**Example:**
```typescript
// test/fixtures.ts
export const mockDownload: Download = {
  id: 'test-id-123',
  title: 'Test Book',
  author: 'Test Author',
  status: 'queued',
  progress: 0,
  createdAt: new Date().toISOString(),
}
```

## Coverage Goals

**Backend:**
- Current: ~20% (estimated)
- Goal: Maintain pragmatic coverage on critical paths
- Focus areas: Handlers, services, repositories, utilities

**Frontend:**
- Current: 60% threshold enforced
- Goal: Maintain 60%+ coverage
- Focus areas: Stores, API client, utility functions

**Don't aim for 100% coverage** - prioritize tests that provide value over high numbers.
