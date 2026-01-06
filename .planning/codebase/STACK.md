# Technology Stack

**Analysis Date:** 2026-01-06

## Languages

**Primary:**
- TypeScript 5.9.3 - Frontend application code (`frontend/package.json`, `frontend/tsconfig.json`)
- Go 1.25.3 - Backend application code (`backend/go.mod`)

**Secondary:**
- SQL - Database migrations (`backend/assets/migrations/`)
- JavaScript - Build configuration (`frontend/eslint.config.js`)

## Runtime

**Environment:**
- Node.js - Frontend runtime (module type: "module" in `frontend/package.json`)
- Go - Backend runtime with native compilation

**Package Manager:**
- npm - Frontend package manager (`frontend/package-lock.json`)
- Go modules - Backend dependency management (`backend/go.mod`, `backend/go.sum`)

## Frameworks

**Core:**
- React 19.2.0 - Frontend UI framework (`frontend/package.json`, `frontend/src/App.tsx`)
- Chi 5.2.3 - Backend HTTP router/middleware (`backend/go.mod`, `backend/internal/server/server.go`)

**Testing:**
- None currently configured

**Build/Dev:**
- Vite 7.2.4 - Frontend bundler and dev server (`frontend/vite.config.ts`)
- TypeScript compiler - TypeScript compilation (`frontend/tsconfig.json`, `frontend/tsconfig.app.json`, `frontend/tsconfig.node.json`)
- Air - Backend hot-reload for development (`backend/.air.toml`)
- ESLint 9.39.1 - Frontend linting (`frontend/package.json`, `frontend/eslint.config.js`)

## Key Dependencies

**Critical:**
- React Router DOM 7.11.0 - Frontend routing (`frontend/package.json`, `frontend/src/App.tsx`)
- Zustand 5.0.9 - Frontend state management (`frontend/package.json`, `frontend/src/stores/`)
- TailwindCSS 4.1.18 - Utility-first CSS (`frontend/package.json`, `frontend/vite.config.ts`)
- Chi CORS 1.2.2 - CORS middleware (`backend/go.mod`, `backend/internal/server/server.go`)
- SQLite3 1.14.32 - Embedded database driver (`backend/go.mod`, `backend/internal/persistence/sqlite/db.go`)

**Infrastructure:**
- React Hook Form 7.70.0 - Form state management (`frontend/package.json`)
- google/uuid 1.6.0 - UUID generation (`backend/go.mod`)
- @vitejs/plugin-react 5.1.1 - Vite React plugin (`frontend/package.json`)
- @tailwindcss/vite 4.1.18 - Tailwind Vite plugin (`frontend/package.json`, `frontend/vite.config.ts`)
- Go net/http - Standard HTTP client library (`backend/internal/qbittorrent/client.go`, `backend/internal/search/providers/mam.go`)

## Configuration

**Environment:**
- Frontend: Vite environment variables in `frontend/.env.development` (VITE_API_URL)
- Backend: Database-driven configuration stored in SQLite (`backend/internal/config/service.go`)
- No `.env.example` files present (documentation gap)

**Build:**
- `frontend/vite.config.ts` - Vite configuration with React and Tailwind plugins
- `frontend/tsconfig.json` - TypeScript base configuration
- `frontend/tsconfig.app.json` - Application TypeScript config (strict mode enabled)
- `frontend/eslint.config.js` - ESLint flat config with TypeScript and React rules
- `backend/.air.toml` - Hot reload configuration for development

## Platform Requirements

**Development:**
- Any platform with Node.js and Go installed
- No Docker configuration present
- SQLite database at `backend/organizr.db` (WAL mode enabled)

**Production:**
- Frontend: Static build output from Vite (SPA)
- Backend: Standalone Go binary
- Database: SQLite with WAL mode for concurrency
- No deployment configuration found

---

*Stack analysis: 2026-01-06*
*Update after major dependency changes*
