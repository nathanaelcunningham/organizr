# Testing Patterns

**Analysis Date:** 2026-01-07 (Updated after Phase 6)

## Test Framework

**Current State: Comprehensive Test Coverage Implemented**

**Backend (Go):**
- Standard Go testing framework (`testing` package)
- Test files: `handlers_test.go`, `monitor_test.go`, `organization_test.go`, `template_test.go`
- Race detection enabled (`go test -race`)
- Coverage: 12.6% (server), 25.4% (downloads)

**Frontend (TypeScript/React):**
- Vitest configured with @testing-library/react
- Test files: `useDownloadStore.test.ts` (22 tests)
- All tests passing in 347ms
- Coverage: Comprehensive store state management

## Test File Organization

**Current Pattern:**
- Backend: Co-located tests (`handlers.go` + `handlers_test.go`)
- Frontend: Co-located tests (`useDownloadStore.ts` + `useDownloadStore.test.ts`)

**Files:**
- `backend/internal/server/handlers_test.go` - HTTP handler tests
- `backend/internal/downloads/monitor_test.go` - Monitor service concurrency tests
- `backend/internal/downloads/organization_test.go` - Organization service tests
- `backend/internal/downloads/template_test.go` - Template validation tests
- `frontend/src/stores/useDownloadStore.test.ts` - Frontend store tests

**Naming:**
- Backend: `*_test.go` suffix
- Frontend: `*.test.ts` or `*.test.tsx` suffix

## Test Structure

**Backend (Go) - Recommended:**
```go
func TestDownloadService_CreateDownload(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        // arrange
        // act
        // assert
    })

    t.Run("error case", func(t *testing.T) {
        // test code
    })
}
```

**Frontend (TypeScript) - Recommended:**
```typescript
describe('SearchBar', () => {
  describe('handleSearch', () => {
    it('should debounce search input', () => {
      // arrange
      // act
      // assert
    });
  });
});
```

## Mocking

**Backend - Recommended:**
- Use interface-based mocking (interfaces already present in `backend/internal/persistence/interfaces.go`)
- Mock external clients (qBittorrent, MAM API)
- Pattern: Create mock implementations of repository interfaces

**Frontend - Recommended:**
- Mock API clients from `frontend/src/api/`
- Mock fetch/HTTP calls
- Use Vitest mocking: `vi.mock()`

**What to Mock:**
- Backend: Database repositories, external HTTP clients, file system operations
- Frontend: API calls, Zustand stores (in isolation), timers (for polling)

**What NOT to Mock:**
- Pure utility functions (`frontend/src/utils/formatters.ts`)
- Domain models
- TypeScript types

## Fixtures and Factories

**Not Currently Implemented**

**Recommended:**
```typescript
// Frontend factory pattern
function createTestDownload(overrides?: Partial<Download>): Download {
  return {
    id: 'test-id',
    title: 'Test Download',
    status: 'downloading',
    progress: 50,
    ...overrides
  };
}
```

```go
// Backend factory pattern
func NewTestDownload(overrides ...func(*models.Download)) *models.Download {
    dl := &models.Download{
        ID:    uuid.New().String(),
        Title: "Test Download",
        Status: models.StatusQueued,
    }
    for _, override := range overrides {
        override(dl)
    }
    return dl
}
```

## Coverage

**Current:**
- No coverage tracking

**Recommended:**
- Backend: `go test -cover`
- Frontend: Vitest coverage via c8 (built-in)
- Target: 70%+ for critical paths (services, API handlers, business logic)

## Test Types

**Unit Tests - Not Implemented:**
- Should test: Services, utilities, formatters, validation logic
- Mock: All external dependencies

**Integration Tests - Not Implemented:**
- Should test: HTTP handlers with real service layer, database operations
- Mock: Only external APIs (qBittorrent, MAM)

**E2E Tests - Not Implemented:**
- Should test: Full user flows (search → download → monitor)
- Tools: Playwright or Cypress (if needed)

## Common Patterns

**Not Currently Implemented**

**Backend - Recommended:**
```go
func TestSearchService_Search(t *testing.T) {
    // Mock repository
    mockRepo := &mockConfigRepo{
        config: map[string]string{
            "mam.baseurl": "https://test.example.com",
            "mam.secret": "test-secret",
        },
    }

    service := search.NewMAMService(mockRepo)

    results, err := service.Search(context.Background(), "test query")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(results) == 0 {
        t.Error("expected results, got none")
    }
}
```

**Frontend - Recommended:**
```typescript
describe('useDownloadStore', () => {
  it('should fetch downloads', async () => {
    // Mock API
    vi.mock('@/api/downloads', () => ({
      list: vi.fn().mockResolvedValue({ downloads: [] })
    }));

    const store = useDownloadStore.getState();
    await store.fetchDownloads();

    expect(store.downloads).toEqual([]);
  });
});
```

## Critical Test Coverage Gaps

**Backend:**
1. **Search Service** (`backend/internal/search/search_service.go`)
   - No tests for MAM API integration
   - No validation testing
   - Risk: Search failures could break silently

2. **Download Service** (`backend/internal/downloads/service.go`)
   - No tests for download creation, status updates, deletion
   - Risk: Core functionality untested

3. **Download Monitor** (`backend/internal/downloads/monitor.go`)
   - No tests for background polling logic
   - No tests for auto-organization
   - Risk: Context misuse bug (line 96) undetected

4. **HTTP Handlers** (`backend/internal/server/handlers.go`)
   - No tests for request validation, error handling, response format
   - Risk: API contract changes could break frontend

5. **File Organization** (`backend/internal/downloads/organization.go`)
   - No tests for path template processing
   - No tests for path traversal prevention
   - Risk: Security vulnerability (line 67) undetected

6. **qBittorrent Client** (`backend/internal/qbittorrent/client.go`)
   - No tests for authentication, torrent operations
   - Risk: Ignored errors (lines 57, 144) not caught

**Frontend:**
1. **API Clients** (`frontend/src/api/*.ts`)
   - No tests for error handling, timeouts, retries
   - Risk: Network failures mishandled

2. **Zustand Stores** (`frontend/src/stores/*.ts`)
   - No tests for state updates, polling logic
   - Risk: State inconsistencies

3. **Search Components** (`frontend/src/components/search/*.tsx`)
   - No tests for debouncing, filtering, result display
   - Risk: UI bugs

4. **Download Polling** (`frontend/src/stores/useDownloadStore.ts`)
   - No tests for 3-second polling interval
   - No tests for start/stop logic
   - Risk: Memory leaks, excessive polling

## Run Commands

**Backend - When Implemented:**
```bash
go test ./...                          # Run all tests
go test -v ./internal/downloads        # Run specific package
go test -cover ./...                   # With coverage
go test -race ./...                    # Race detection
```

**Frontend - When Implemented:**
```bash
npm test                              # Run all tests
npm test -- --watch                   # Watch mode
npm test -- SearchBar.test.tsx        # Single file
npm run test:coverage                 # Coverage report
```

## Priority Testing Recommendations

1. **High Priority:**
   - Backend: Search service, download service, HTTP handlers
   - Frontend: API clients, download store polling

2. **Medium Priority:**
   - Backend: File organization, monitor logic
   - Frontend: Component rendering, form validation

3. **Low Priority:**
   - Utility functions (already simple)
   - Type definitions (TypeScript provides safety)

---

*Testing analysis: 2026-01-06*
*Update when test patterns change*
