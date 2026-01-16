# Technology Stack

## Languages

### Backend
- **Go 1.25.3** - Primary backend language
  - Compiled binary with static linking
  - CGO enabled for SQLite support
  - Usage: REST API server, background monitoring, file operations

### Frontend
- **TypeScript 5.9.3** - Primary frontend language
  - Strict mode enabled
  - Full type safety across React components
  - Usage: All frontend code (components, stores, API client)

## Frameworks

### Backend
- **Chi v5.2.3** - HTTP router and middleware
  - Lightweight, idiomatic Go router
  - Middleware: CORS, request logging, panic recovery
  - Route grouping and sub-routers
  - File: backend/internal/server/routes.go

### Frontend
- **React 19.2.0** - UI framework
  - Features: Suspense, new hooks (useTransition, useDeferredValue)
  - Functional components only
  - File: frontend/src/main.tsx, frontend/src/App.tsx

- **React Router v7.11.0** - Client-side routing
  - BrowserRouter with route configuration
  - Nested routes for layout
  - File: frontend/src/App.tsx

## Libraries

### Backend Core Dependencies
- **mattn/go-sqlite3 v1.14.28** - SQLite driver with CGO bindings
  - WAL mode for concurrent access
  - File: backend/internal/persistence/sqlite/db.go

- **google/uuid v1.6.0** - UUID generation for download IDs
  - v4 (random) UUIDs
  - File: backend/internal/models/download.go

- **swaggo/swag v1.16.6** - Swagger/OpenAPI documentation generation
  - Annotations in handler comments
  - Generate: `swag init -g cmd/api/main.go`

- **go-chi/cors v1.2.2** - CORS middleware
  - Configurable allowed origins
  - File: backend/internal/server/server.go

- **joho/godotenv v1.5.1** - .env file loading
  - Development configuration
  - File: backend/cmd/api/main.go

### Frontend Core Dependencies
- **zustand v5.0.9** - State management
  - Lightweight Redux alternative
  - Hook-based API
  - Files: frontend/src/stores/*

- **tailwindcss v4.1.0** - Utility-first CSS framework
  - Vite plugin integration
  - Files: frontend/src/index.css, frontend/vite.config.ts

- **react-hook-form v7.70.0** - Form management
  - Validation and error handling
  - Files: frontend/src/components/config/, frontend/src/components/search/

### Frontend Development Dependencies
- **vite v7.2.0** - Build tool and dev server
  - Fast HMR, optimized builds
  - React plugin, TailwindCSS plugin
  - File: frontend/vite.config.ts

- **vitest v4.0.0** - Testing framework
  - 60% coverage threshold (lines, functions, branches, statements)
  - happy-dom environment for lightweight DOM
  - File: frontend/vitest.config.ts

- **@testing-library/react v17.0.1** - React component testing
  - User-centric queries
  - File: frontend/src/stores/useDownloadStore.test.ts

## Databases & Storage

### SQLite 3
- **Version:** 3.x (via mattn/go-sqlite3 v1.14.28)
- **Mode:** WAL (Write-Ahead Logging) for concurrent reads
- **Connection Pool:** 25 max open, 25 max idle connections
- **Location:** Configurable via ORGANIZR_DB_PATH (default: /data/organizr.db in Docker)
- **Tables:**
  - `downloads` - Download entities with progress and status
  - `configs` - Key-value configuration storage
  - `schema_migrations` - Migration version tracking
- **Indexes:** status, qbit_hash (unique), created_at

### Migrations
- **Tool:** Custom Go migration runner in main.go
- **Location:** backend/assets/migrations/
- **Format:** SQL files with .up.sql extension (001, 002, 003, 004)

## Build & Deploy

### Backend Build Tool
- **Go toolchain 1.25.3**
  - Build command: `CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o organizr ./cmd/api`
  - Make targets: build, run, test, test-coverage, test-race, lint
  - File: backend/Makefile

### Frontend Build Tool
- **Vite 7.2.0**
  - Dev server: `npm run dev` (HMR on localhost:5173)
  - Production build: `npm run build` (outputs to dist/)
  - Preview: `npm run preview`
  - File: frontend/package.json

### Package Managers
- **Go:** go mod (v1.25)
  - Files: backend/go.mod, backend/go.sum
- **Node:** npm (comes with Node 24)
  - Files: frontend/package.json, frontend/package-lock.json

### Docker
- **Multi-stage builds** - Separate build and runtime stages for smaller images

**Backend Dockerfile:**
- Build stage: golang:1.25-alpine with gcc, musl-dev, sqlite-dev
- Runtime stage: alpine:3.21 with ca-certificates, sqlite-libs
- User: app (uid 1001, gid 1001)
- Port: 8080
- Command: /app/organizr

**Frontend Dockerfile:**
- Build stage: node:24-alpine with npm ci and Vite build
- Runtime stage: nginx:alpine with custom nginx.conf
- User: app (uid 1001, gid 1001)
- Port: 8080 (Nginx)
- Command: nginx -g "daemon off;"

### Docker Compose
- **Version:** 3.8
- **Network:** organizr-network (bridge)
- **Services:**
  - backend: Port 8080, health check via /api/health
  - frontend: Port 8081, depends on backend health, proxies /api requests
- **Volumes:** Downloads, audiobooks, data directories (configurable paths)
- **File:** docker-compose.yml

### Deployment Platforms
1. **Docker Compose (Recommended)** - Full stack orchestration
2. **Docker Individually** - Separate container management
3. **Bare Metal** - Standalone binaries with systemd services
4. **Unraid** - Docker Compose with NAS volume mappings
5. **Synology** - Similar to Unraid with DSM paths

## Development Environment

### Required Tools
- **Go 1.25+** - Backend development and compilation
- **Node 24+** - Frontend development and builds
- **npm** - Frontend package management
- **Docker** - Container builds and testing
- **Docker Compose** - Multi-service orchestration
- **Make** - Backend build automation (optional)
- **SQLite 3** - Database CLI for manual inspection (optional)

### Setup Steps
1. Clone repository
2. Install Go 1.25+ and Node 24+
3. Backend setup:
   - `cd backend`
   - `go mod download`
   - Copy .env.example to .env and configure
   - `make run` or `go run cmd/api/main.go`
4. Frontend setup:
   - `cd frontend`
   - `npm install`
   - `npm run dev`
5. Docker setup (alternative):
   - `docker-compose up --build`

### IDE Recommendations
- **Backend:** VSCode with Go extension, GoLand
- **Frontend:** VSCode with ESLint and Prettier extensions, WebStorm

### Code Quality Tools
- **golangci-lint** - Go linter with configuration in .golangci.yml
- **ESLint** - Frontend linting (configured in frontend package.json)
- **Prettier** - Code formatting (frontend)
