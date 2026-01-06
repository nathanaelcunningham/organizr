# Codebase Structure

**Analysis Date:** 2026-01-06

## Directory Layout

```
organizr/
├── backend/               # Go API server
│   ├── cmd/api/          # Application entry point
│   ├── internal/         # Internal packages
│   ├── assets/           # Embedded resources (migrations)
│   ├── go.mod            # Go dependencies
│   ├── Makefile          # Build commands
│   └── .air.toml         # Hot reload config
├── frontend/             # React SPA
│   ├── src/              # Source code
│   ├── public/           # Static assets
│   ├── package.json      # npm dependencies
│   ├── vite.config.ts    # Vite configuration
│   └── eslint.config.js  # ESLint configuration
└── .planning/            # Project planning (newly created)
    └── codebase/         # Codebase documentation
```

## Directory Purposes

**backend/**
- Purpose: Go REST API server
- Contains: All backend application code
- Key files: `backend/cmd/api/main.go` (entry point), `backend/go.mod` (dependencies)
- Subdirectories: `cmd/` (entry points), `internal/` (application code), `assets/` (embedded resources)

**backend/cmd/api/**
- Purpose: Application entry point
- Contains: `main.go` - Database initialization, service setup, HTTP server start, graceful shutdown
- Key files: `main.go`
- Subdirectories: None

**backend/internal/**
- Purpose: Internal application packages (not importable externally)
- Contains: Business logic, services, models, persistence, HTTP handlers
- Key directories:
  - `config/` - Configuration service
  - `downloads/` - Download domain (service, monitor, organization)
  - `fileutil/` - File utilities (sanitizer, template)
  - `models/` - Domain models
  - `persistence/` - Repository interfaces and SQLite implementations
  - `qbittorrent/` - qBittorrent client integration
  - `search/` - Search service and providers
  - `server/` - HTTP server, routes, handlers, validation

**backend/assets/migrations/**
- Purpose: Database schema migrations
- Contains: SQL migration files
- Key files: `001_init.up.sql` - Initial schema with downloads and configs tables
- Subdirectories: None

**frontend/**
- Purpose: React single-page application
- Contains: All frontend code
- Key files: `package.json`, `vite.config.ts`, `eslint.config.js`
- Subdirectories: `src/` (source code), `public/` (static assets)

**frontend/src/**
- Purpose: Frontend source code
- Contains: Components, pages, stores, API clients, types, utilities
- Key files: `main.tsx` (React entry), `App.tsx` (router), `index.css` (global styles)
- Key directories:
  - `api/` - HTTP client and API endpoints
  - `components/` - React components (layout, common, feature-specific)
  - `pages/` - Page components
  - `stores/` - Zustand state stores
  - `types/` - TypeScript type definitions
  - `hooks/` - Custom React hooks
  - `utils/` - Utility functions

**frontend/src/components/**
- Purpose: Reusable UI components
- Contains:
  - `layout/` - Layout structure (Layout, Header, Sidebar, PageHeader)
  - `common/` - Reusable UI (Button, Input, Card, Modal, Spinner, Badge, ProgressBar, etc.)
  - `search/` - Search feature components (SearchBar, SearchResults, SearchResultCard)
  - `downloads/` - Download feature components (DownloadList, DownloadCard, DownloadFilters)
  - `config/` - Config feature components (ConfigForm, ConfigSection)

**frontend/src/stores/**
- Purpose: Zustand state management stores
- Contains: `useDownloadStore.ts`, `useSearchStore.ts`, `useConfigStore.ts`, `useNotificationStore.ts`
- Pattern: Each store manages domain-specific state with actions and computed getters

## Key File Locations

**Entry Points:**
- `backend/cmd/api/main.go` - Backend application entry
- `frontend/src/main.tsx` - Frontend React entry
- `frontend/src/App.tsx` - Frontend router configuration

**Configuration:**
- `backend/go.mod` - Go dependencies
- `frontend/package.json` - npm dependencies and scripts
- `frontend/vite.config.ts` - Vite build configuration
- `frontend/tsconfig.json` - TypeScript base config
- `frontend/tsconfig.app.json` - TypeScript app config (strict mode)
- `frontend/eslint.config.js` - ESLint flat config
- `backend/.air.toml` - Hot reload configuration
- `frontend/.env.development` - Development environment variables

**Core Logic:**
- `backend/internal/downloads/service.go` - Download business logic
- `backend/internal/search/search_service.go` - Search service
- `backend/internal/config/service.go` - Configuration management
- `backend/internal/server/handlers.go` - HTTP request handlers
- `backend/internal/persistence/sqlite/` - Database operations
- `frontend/src/stores/` - Frontend state management

**Testing:**
- None detected

**Documentation:**
- `backend/README.md` - Backend documentation
- `frontend/README.md` - Frontend documentation (Vite template)
- `.planning/codebase/` - Codebase analysis documents (newly created)

## Naming Conventions

**Files:**

Backend (Go):
- snake_case for Go files: `search_service.go`, `request_types.go`
- Entry point: `main.go`
- Test files: `*_test.go` (none present)

Frontend (TypeScript/React):
- PascalCase for component files: `SearchBar.tsx`, `DownloadCard.tsx`, `Button.tsx`
- camelCase for utilities: `useDebounce.ts`, `formatters.ts`, `constants.ts`
- kebab-case for API files: `client.ts`, `downloads.ts`, `search.ts`

**Directories:**
- Backend: lowercase, short names: `models`, `server`, `downloads`
- Frontend: camelCase or kebab-case: `components`, `api`, `stores`, `hooks`, `utils`

**Special Patterns:**
- Go: `main.go` for entry points
- React: `main.tsx` for React entry, `App.tsx` for router
- TypeScript: `index.css` for global styles
- Hooks: Prefix with `use`: `useDebounce.ts`, `useToast.ts`
- Stores: Prefix with `use`, suffix with `Store`: `useDownloadStore.ts`

## Where to Add New Code

**New Backend Feature:**
- Primary code: `backend/internal/<domain>/service.go`
- Models: `backend/internal/models/<domain>.go`
- Persistence: `backend/internal/persistence/interfaces.go` (interface), `backend/internal/persistence/sqlite/<domain>.go` (implementation)
- HTTP handlers: `backend/internal/server/handlers.go`, `backend/internal/server/routes.go`
- Tests: None present (should add `*_test.go` alongside source)

**New Frontend Component:**
- Implementation: `frontend/src/components/<category>/<ComponentName>.tsx`
- Types: `frontend/src/types/<domain>.ts`
- Styles: Inline with TailwindCSS classes

**New Frontend Page:**
- Implementation: `frontend/src/pages/<PageName>.tsx`
- Route: Add to `frontend/src/App.tsx` router
- Store (if needed): `frontend/src/stores/use<Domain>Store.ts`
- API client (if needed): `frontend/src/api/<domain>.ts`

**New API Endpoint:**
- Backend handler: `backend/internal/server/handlers.go`
- Backend route: `backend/internal/server/routes.go`
- Frontend client: `frontend/src/api/<domain>.ts`
- Types: `frontend/src/types/<domain>.ts`

**Utilities:**
- Backend: `backend/internal/<domain>/` (domain-specific) or new package for shared utilities
- Frontend: `frontend/src/utils/` (formatters, validators, constants)
- Type definitions: `frontend/src/types/`

## Special Directories

**backend/assets/**
- Purpose: Embedded resources (compiled into binary)
- Source: Manually created migration files
- Committed: Yes

**frontend/public/**
- Purpose: Static assets served directly
- Source: Images, fonts, favicon
- Committed: Yes

**frontend/dist/**
- Purpose: Build output (production bundle)
- Source: Generated by Vite build
- Committed: No (in .gitignore)

**backend/organizr.db** (implied)
- Purpose: SQLite database file
- Source: Created at runtime
- Committed: No (should be in .gitignore)

---

*Structure analysis: 2026-01-06*
*Update when directory structure changes*
