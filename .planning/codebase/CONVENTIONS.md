# Coding Conventions

**Analysis Date:** 2026-01-06

## Naming Patterns

**Files:**
- Backend (Go): snake_case for all files (`search_service.go`, `request_types.go`, `main.go`)
- Frontend components: PascalCase (`SearchBar.tsx`, `DownloadCard.tsx`, `Button.tsx`)
- Frontend utilities: camelCase (`useDebounce.ts`, `formatters.ts`, `constants.ts`)
- Frontend API: kebab-case (`client.ts`, `downloads.ts`, `search.ts`)

**Functions:**
- Backend (Go): PascalCase for exported (`NewService`, `CreateDownload`), camelCase for unexported (`initProvider`)
- Frontend: camelCase for all functions (`formatFileSize`, `handleSearch`, `createDownload`)
- Frontend handlers: `handle<EventName>` pattern (`handleClick`, `handleSubmit`)
- Hooks: Prefix with `use` (`useDebounce`, `useSearchStore`, `useDownloadStore`)

**Variables:**
- Backend (Go): camelCase for local variables, PascalCase for exported
- Frontend: camelCase for all variables (`searchQuery`, `downloadStatus`, `isLoading`)
- Constants: UPPER_SNAKE_CASE (`MIN_SEARCH_LENGTH`, `DOWNLOAD_POLL_INTERVAL`, `SEARCH_DEBOUNCE_DELAY`)

**Types:**
- Backend (Go): PascalCase for all types (`Download`, `SearchResult`, `DownloadStatus`)
- Frontend: PascalCase for interfaces and types (`Download`, `SearchResult`, `ButtonProps`)
- Go interfaces: PascalCase with suffix (`DownloadRepository`, `ConfigRepository`)
- No `I` prefix for interfaces (Go or TypeScript)

## Code Style

**Formatting:**

Backend (Go):
- Tab indentation (Go standard)
- No semicolons (Go standard)
- Comments use `//` for single-line
- Error wrapping: `fmt.Errorf("context: %w", err)`
- Deferred cleanup: `defer db.Close()`

Frontend (TypeScript):
- 2-space indentation
- Single quotes for string literals in code
- Double quotes in JSX attributes: `<Input type="text" />`
- Semicolons used consistently
- Template literals for multiline strings (Tailwind classes)

**Linting:**

Backend:
- No linter configured (uses `go fmt`)

Frontend:
- ESLint with flat config format (`eslint.config.js`)
- Extends: `@eslint/js`, `typescript-eslint`, `eslint-plugin-react-hooks`, `eslint-plugin-react-refresh`
- Targets: `**/*.{ts,tsx}`
- Browser globals enabled
- Run: `npm run lint`

**TypeScript:**
- Strict mode enabled (`tsconfig.app.json`)
- Target: ES2022, Module: ESNext
- `noUnusedLocals: true`, `noUnusedParameters: true`, `noFallthroughCasesInSwitch: true`

## Import Organization

**Backend (Go):**
1. Standard library imports
2. External packages
3. Internal packages
Blank line between groups, alphabetical within groups

**Frontend (TypeScript):**
1. External packages (`react`, `zustand`, etc.)
2. Internal modules (relative imports from `./` or `@/`)
3. Type imports (`import type { }`)

**Order:**
- Blank line between groups
- Destructured imports: `import { create } from 'zustand'`

**Path Aliases:**
- No path aliases configured (uses relative imports)

## Error Handling

**Backend (Go):**
- Throw errors, catch at handler boundaries
- Error wrapping: `fmt.Errorf("failed to do X: %w", err)`
- Return errors to caller, handle at HTTP handlers
- Deferred cleanup: `defer db.Close()`
- **Concern**: Multiple ignored errors with `_` (should be logged)

**Frontend (TypeScript):**
- Custom `APIClientError` class extends `Error` (`frontend/src/api/client.ts`)
- try/catch in async operations
- Error state in Zustand stores
- Error handling at API client layer with timeout support

**Error Types:**
- Backend: Return errors from all functions, wrap with context
- Frontend: Catch at API boundary, display via notification store
- Logging: Backend uses standard log package, frontend uses console in dev mode

## Testing Conventions

**Backend (Go):**
- **Table-driven tests with subtests** - Test multiple cases efficiently using `t.Run()`
- **Interface mocking for dependencies** - Repository and client interfaces enable mock implementations
- **Race detection required** - All concurrent code must pass `go test -race`
- **Test naming:** `Test<FunctionName>` (e.g., `TestDownloadService_CreateDownload`)
- **File location:** `*_test.go` alongside source files

**Frontend (TypeScript):**
- **Vitest with jsdom** - DOM testing environment for React components
- **Test stores and utilities** - Focus on logic, avoid implementation details
- **Test naming:** `describe/it` blocks (e.g., `describe('formatFileSize', () => { it('formats bytes correctly', ...) })`)
- **File location:** `*.test.ts` files or `test/` directory

## Logging

**Backend:**
- Standard Go `log` package
- Output: stdout/stderr
- No structured logging

**Frontend:**
- `console.log` in development (`frontend/src/api/client.ts`)
- No production logging framework

**Patterns:**
- Log at service boundaries
- Log errors before returning
- Development-only logging in frontend

## Comments

**When to Comment:**
- Backend: Single-line comments for exported functions, inline for complex logic
- Frontend: JSDoc-style for utility functions, inline for complex logic
- Explain why, not what
- Document business logic and edge cases

**JSDoc/TSDoc:**
- Frontend: Used for utility functions (`frontend/src/utils/formatters.ts`, `frontend/src/hooks/useDebounce.ts`)
- Format: `@param`, `@returns` tags
- Example:
  ```typescript
  /**
   * Format bytes to human-readable file size
   */
  export function formatFileSize(bytes: number | string): string {
  ```

**Backend (Go):**
- Single-line comments above functions: `// MAMService handles torrent search`
- No formal JSDoc equivalent

**TODO Comments:**
- None detected in codebase

## Commit Message Format

**Convention:** Conventional Commits with phase-plan scoping

**Format:** `type(scope): description`

**Types:**
- `feat` - New feature, endpoint, component
- `fix` - Bug fix, error correction
- `test` - Test-only changes
- `refactor` - Code cleanup, no behavior change
- `docs` - Documentation changes
- `chore` - Config, tooling, dependencies

**Scope:** Phase-plan format for milestone work (e.g., `09-02`, `10-01`) or general scope for standalone changes

**Examples from recent commits:**
```bash
feat(09-02): add series_number template variable
fix(07.1-01): extract series name only for downloads
test(08-01): add batch handler integration tests
docs(10-01): create Architecture Decision Record
refactor(06-01): extract organization logic to service
chore: update frontend dependencies
```

## Function Design

**Size:**
- Keep functions focused and readable
- Backend: No explicit limit, but functions are generally concise
- Frontend: Similar approach

**Parameters:**
- Backend: Struct-based for multiple parameters
- Frontend: Destructured object parameters for 4+ params
- Example: `function create(options: CreateOptions)`

**Return Values:**
- Backend: Multiple return values `(result, error)`, always check errors
- Frontend: Explicit returns, no implicit undefined
- Return early for guard clauses

## Module Design

**Exports:**
- Backend: PascalCase for exported, camelCase for unexported
- Frontend: Named exports preferred
- Default exports: Not commonly used

**Barrel Files:**
- Not used in this codebase
- Direct imports from source files

**Backend Packages:**
- Interface-based for dependency inversion (`backend/internal/persistence/interfaces.go`)
- Constructor pattern: `New<Type>()` functions
- Dependency injection via constructors

**Frontend Modules:**
- Domain-specific API clients (`frontend/src/api/downloads.ts`, etc.)
- Zustand stores for state management
- Type definitions separated in `frontend/src/types/`

---

*Convention analysis: 2026-01-08 (v1.1 complete)*
*Update when patterns change*
