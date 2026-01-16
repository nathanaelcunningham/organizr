# Coding Conventions

## File Organization

### Backend (Go)
```
backend/
├── cmd/api/main.go              # Application entry point
├── internal/                    # Private application code
│   ├── server/                  # HTTP layer (handlers, routes, DTOs)
│   ├── models/                  # Domain entities (Download, Config, etc.)
│   ├── persistence/             # Repository interfaces and implementations
│   │   ├── interfaces.go        # Repository contracts
│   │   └── sqlite/              # SQLite implementations
│   ├── downloads/               # Download domain logic
│   ├── qbittorrent/             # qBittorrent client
│   ├── search/                  # Search service and providers
│   ├── config/                  # Configuration service
│   ├── fileutil/                # File utilities (templates, sanitization)
│   └── testutil/                # Test helpers and fixtures
├── assets/migrations/            # SQL migration files
└── docs/                        # API documentation
```

**Principles:**
- `cmd/` - Application entry points (only main.go for API server)
- `internal/` - Private packages (cannot be imported by external projects)
- Domain logic organized by feature (downloads, search, config)
- Shared utilities in dedicated packages (fileutil, testutil)

### Frontend (React)
```
frontend/
├── src/
│   ├── main.tsx                 # React entry point
│   ├── App.tsx                  # Route configuration
│   ├── pages/                   # Route components (SearchPage, DownloadsPage, etc.)
│   ├── components/              # Reusable components
│   │   ├── layout/              # Layout and navigation
│   │   ├── search/              # Search-specific components
│   │   ├── downloads/           # Download-specific components
│   │   ├── config/              # Config-specific components
│   │   └── common/              # Shared components (Button, Modal, etc.)
│   ├── stores/                  # Zustand state management
│   ├── api/                     # API client layer
│   ├── types/                   # TypeScript type definitions
│   ├── hooks/                   # Custom React hooks
│   ├── utils/                   # Utility functions
│   ├── test/                    # Test utilities and fixtures
│   └── index.css                # Global styles and Tailwind imports
├── public/                      # Static assets
└── node_modules/                # Dependencies (not committed)
```

**Principles:**
- Pages define routes and compose components
- Components organized by feature (search, downloads, config)
- Common/shared components in dedicated directory
- Zustand stores per domain (useDownloadStore, useSearchStore)
- API client wraps HTTP requests with error handling
- Types mirror backend models

## Naming Conventions

### Backend (Go)
**Files:**
- Lowercase with underscores: `downloads.go`, `mam_test.go`
- Test files: `*_test.go`
- Main packages: `main.go`

**Variables/Functions:**
- Exported: PascalCase (`CreateDownload`, `DownloadService`)
- Unexported: camelCase (`parseTemplate`, `sanitizePath`)
- Constants: PascalCase for exported, camelCase for unexported
- Acronyms: `ID`, `URL`, `API` (uppercase when part of exported name)

**Interfaces:**
- Descriptive nouns: `DownloadRepository`, `ConfigRepository`
- Avoid "I" prefix (Go convention)

**Structs:**
- PascalCase: `Download`, `SearchResult`, `QBittorrentClient`
- Fields: PascalCase for exported, camelCase for unexported

**Methods:**
- PascalCase: `GetDownload`, `UpdateStatus`
- Receiver names: Single letter or short abbreviation (`d *Download`, `s *DownloadService`)

### Frontend (TypeScript/React)
**Files:**
- PascalCase for components: `SearchPage.tsx`, `DownloadTable.tsx`
- camelCase for utilities: `groupSeries.ts`, `formatters.ts`
- Stores: `useDownloadStore.ts` (hook naming convention)
- Test files: `*.test.ts` or `*.test.tsx`

**Components:**
- PascalCase: `SearchPage`, `DownloadRow`
- File name matches component name

**Functions:**
- camelCase: `fetchDownloads`, `formatDate`
- Hooks: `use` prefix (`usePolling`, `useDownloadStore`)

**Variables:**
- camelCase: `downloadId`, `searchQuery`
- Constants: SCREAMING_SNAKE_CASE (`API_BASE_URL`, `MAX_RETRIES`)

**Types/Interfaces:**
- PascalCase: `Download`, `SearchResult`, `APIError`
- Props interfaces: `SearchFormProps`, `DownloadRowProps`

## Code Patterns

### Error Handling (Backend)
**Preferred: Wrapped errors with context**

```go
// Good
if err != nil {
    return nil, fmt.Errorf("failed to create download: %w", err)
}

// Also good - early returns
download, err := s.repo.GetDownload(id)
if err != nil {
    return nil, fmt.Errorf("get download: %w", err)
}
```

**HTTP Error Responses:**
```go
// Use helpers from server/errors.go
RespondError(w, http.StatusBadRequest, "Invalid download ID")
RespondValidationError(w, "title", "Title is required")
```

**Error Chain:**
- Repository layer: Returns raw errors or wrapped with context
- Service layer: Wraps with business context
- Handler layer: Logs and converts to HTTP responses

### Error Handling (Frontend)
**APIClientError pattern:**

```typescript
try {
  const downloads = await fetchDownloads()
  // Handle success
} catch (error) {
  if (error instanceof APIClientError) {
    // Check error.statusCode, error.apiError
    if (error.statusCode === 401) {
      // Handle auth error
    }
  }
  // Show user-friendly message
}
```

### Logging (Backend)
**Standard library log package:**

```go
log.Printf("Starting download monitor (interval: %v)", interval)
log.Printf("Download %s completed, triggering organization", d.ID)
log.Printf("ERROR: Failed to organize download %s: %v", d.ID, err)
```

**Conventions:**
- Info: `log.Printf("message")`
- Error: `log.Printf("ERROR: message: %v", err)`
- Debug: Development only, remove before commit

### Logging (Frontend)
**Development console logging:**

```typescript
// Only log in development
if (env.IS_DEV) {
  console.log('API request:', url, options)
}
```

**Conventions:**
- Avoid production logging
- Use browser dev tools for debugging
- Remove console.log before commit

### Comments & Documentation

**Backend (Go):**
- Godoc comments for exported symbols
- Handler comments with Swag annotations for API docs
- Inline comments for complex logic only

```go
// CreateDownload submits a torrent to qBittorrent and creates a download record
//
// @Summary Create download
// @Description Submit torrent to qBittorrent and track download
// @Tags downloads
// @Accept json
// @Produce json
// @Param download body CreateDownloadRequest true "Download details"
// @Success 201 {object} Download
// @Failure 400 {object} ErrorResponse
// @Router /api/downloads [post]
func (s *Server) CreateDownload(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

**Frontend (TypeScript):**
- JSDoc comments for complex functions
- Inline comments for non-obvious logic only
- Component prop descriptions via TypeScript types

```typescript
/**
 * Groups search results by series, handling both single books and series.
 * Results without series are returned in their own group.
 */
export function groupSeries(results: SearchResult[]): SeriesGroup[] {
  // Implementation
}
```

**When to Comment:**
- Complex algorithms or business logic
- Non-obvious workarounds or hacks
- API contracts (Swagger annotations)
- Public functions/methods (Godoc)
- **When NOT to comment:** Self-explanatory code, simple getters/setters, obvious logic

## Code Style

### Backend (Go)
- **Formatting:** `gofmt` (automatic via editor)
- **Linting:** golangci-lint with configuration in .golangci.yml
- **Line length:** No hard limit, but prefer 120 chars for readability
- **Imports:** Standard library first, third-party, then local packages
- **Error handling:** Always check errors, early returns preferred

### Frontend (TypeScript/React)
- **Formatting:** Prettier (recommended)
- **Linting:** ESLint (configured in package.json)
- **Line length:** Prettier default (80 chars)
- **Quotes:** Single quotes for strings, backticks for templates
- **Semicolons:** Required (ESLint rule)
- **Component style:** Functional components with hooks only

## Common Patterns

### Repository Pattern (Backend)
```go
// Define interface in persistence/interfaces.go
type DownloadRepository interface {
    CreateDownload(ctx context.Context, download *models.Download) error
    GetDownload(ctx context.Context, id string) (*models.Download, error)
    // ...
}

// Implement in persistence/sqlite/downloads.go
type SQLiteDownloadRepository struct {
    db *sql.DB
}

func (r *SQLiteDownloadRepository) CreateDownload(ctx context.Context, download *models.Download) error {
    // Implementation with SQL
}
```

### Service Pattern (Backend)
```go
// Service depends on repository interface
type DownloadService struct {
    repo          persistence.DownloadRepository
    qbitClient    *qbittorrent.Client
    configService *config.Service
}

func NewDownloadService(repo persistence.DownloadRepository, qbit *qbittorrent.Client, cfg *config.Service) *DownloadService {
    return &DownloadService{
        repo:          repo,
        qbitClient:    qbit,
        configService: cfg,
    }
}
```

### Zustand Store Pattern (Frontend)
```typescript
interface DownloadStore {
  downloads: Download[]
  loading: boolean
  error: string | null
  fetchDownloads: () => Promise<void>
  createDownload: (request: CreateDownloadRequest) => Promise<void>
}

export const useDownloadStore = create<DownloadStore>((set, get) => ({
  downloads: [],
  loading: false,
  error: null,

  fetchDownloads: async () => {
    set({ loading: true, error: null })
    try {
      const downloads = await fetchDownloads()
      set({ downloads, loading: false })
    } catch (error) {
      set({ error: error.message, loading: false })
    }
  },

  // More actions...
}))
```

### React Component Pattern (Frontend)
```typescript
interface SearchFormProps {
  onSubmit: (query: string) => void
  loading: boolean
}

export function SearchForm({ onSubmit, loading }: SearchFormProps) {
  const [query, setQuery] = useState('')

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    onSubmit(query)
  }

  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields */}
    </form>
  )
}
```

## Anti-Patterns

### Backend (Go)
- **Avoid:** Global state and singletons
  - **Use:** Dependency injection via constructors
- **Avoid:** Panic for expected errors
  - **Use:** Return errors explicitly
- **Avoid:** Ignoring errors (`_ = someFunc()`)
  - **Use:** Check and handle or wrap errors
- **Avoid:** Generic "utils" packages
  - **Use:** Domain-specific packages (fileutil, testutil)

### Frontend (React)
- **Avoid:** Class components
  - **Use:** Functional components with hooks
- **Avoid:** Prop drilling (passing props through many levels)
  - **Use:** Zustand stores for global state
- **Avoid:** Inline anonymous functions in render (causes re-renders)
  - **Use:** useCallback for event handlers
- **Avoid:** Mutating state directly
  - **Use:** Zustand's set function or useState updaters
- **Avoid:** Large components with mixed concerns
  - **Use:** Break into smaller, focused components
