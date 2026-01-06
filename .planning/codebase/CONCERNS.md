# Codebase Concerns

**Analysis Date:** 2026-01-06

## Tech Debt

**Database-driven configuration with hardcoded defaults:**
- Issue: Default credentials (`admin`/`adminpass`) hardcoded in SQL migration and main.go
- Files:
  - `backend/assets/migrations/001_init.up.sql` (lines 34-37)
  - `backend/cmd/api/main.go` (lines 51-62)
- Why: Quick initial setup for development
- Impact: Exposed credentials in version control, insecure defaults could reach production
- Fix approach: Remove hardcoded defaults from migration, require explicit configuration, fail startup if critical config missing

**CORS wildcard configuration:**
- Issue: `AllowedOrigins: []string{"*"}` allows requests from any origin
- File: `backend/internal/server/server.go` (line 44)
- Why: Development convenience
- Impact: Security vulnerability - any website can make requests to API
- Fix approach: Use environment variable for allowed origins, restrict to specific domains in production

**Duplicate filter logic:**
- Issue: Search filtering implemented in both store and component
- Files:
  - `frontend/src/stores/useSearchStore.ts` (lines 80-106)
  - `frontend/src/pages/SearchPage.tsx` (lines 18-41)
- Why: Incremental development without refactoring
- Impact: Maintenance burden, risk of inconsistency between implementations
- Fix approach: Consolidate filtering logic in store only, remove from component

**Manual database migration strategy:**
- Issue: Only one migration file, adding new migrations requires code changes
- File: `backend/cmd/api/main.go` (lines 115-168)
- Why: Initial implementation, single migration sufficient so far
- Impact: No versioning, no rollback mechanism, manual tracking required
- Fix approach: Use migration library like `golang-migrate` or `goose` for proper migration management

## Known Bugs

**Ignored I/O errors in HTTP response handling:**
- Symptoms: Errors reading HTTP response bodies are silently discarded
- Files:
  - `backend/internal/qbittorrent/client.go` (lines 57, 144) - `body, _ := io.ReadAll()`
  - `backend/internal/search/providers/mam.go` (line 140) - `io.ReadAll()` error ignored in error handler
- Trigger: Any HTTP response reading failure (network issues, malformed responses)
- Root cause: Error values ignored with `_`
- Fix: Check and log/return all I/O errors

**Context misuse in background organization:**
- Symptoms: Organization goroutines don't respect monitor cancellation
- File: `backend/internal/downloads/monitor.go` (line 96)
- Trigger: Application shutdown or monitor stop
- Root cause: `context.Background()` used instead of parent context
- Impact: Goroutines may continue running after monitor stops, delaying graceful shutdown
- Fix: Pass `monitorCtx` to `organizeDownload()` instead of `context.Background()`

**Cookie jar initialization error ignored:**
- Symptoms: qBittorrent client may fail to maintain sessions if cookie jar creation fails
- File: `backend/internal/qbittorrent/client.go` (line 23)
- Trigger: Very rare (cookie jar initialization failure)
- Root cause: `cookiejar.New(nil)` error return value ignored
- Fix: Check and handle error from `cookiejar.New()`

## Security Considerations

**Hardcoded default credentials:**
- Risk: Default admin credentials in version control and database
- Files:
  - `backend/assets/migrations/001_init.up.sql` (lines 34-37) - `admin`/`adminpass` in SQL
  - `backend/cmd/api/main.go` (lines 51-62) - Fallback to insecure defaults
- Current mitigation: None
- Recommendations:
  - Remove defaults from migration
  - Require explicit configuration via environment or CLI
  - Fail startup if qBittorrent credentials not configured

**CORS allows all origins:**
- Risk: Cross-origin requests from any domain accepted
- File: `backend/internal/server/server.go` (line 44)
- Current mitigation: None
- Recommendations:
  - Add `ALLOWED_ORIGINS` environment variable
  - Restrict to specific trusted domains in production

**No authentication layer:**
- Risk: API endpoints publicly accessible, no user authentication
- Files: All handlers in `backend/internal/server/handlers.go`
- Current mitigation: None
- Recommendations: Add authentication middleware (JWT, session-based, or basic auth)

**Potential path traversal in file organization:**
- Risk: User-derived data used in file path construction
- File: `backend/internal/downloads/organization.go` (line 67)
- Code: `filepath.Join(destBase, path)` combines base with user-derived path
- Current mitigation: `fileutil.SanitizePath()` replaces invalid characters
- Recommendations: Add explicit path traversal check after `filepath.Join()`, ensure path stays within `destBase`

**Missing .env.example files:**
- Risk: Developers don't know what configuration is required
- Files: No `.env.example` in backend or frontend root
- Current mitigation: README files partially document config
- Recommendations: Create `.env.example` files with all required variables and safe defaults

## Performance Bottlenecks

**Frontend polling inefficiency:**
- Problem: Hard-coded 3-second polling for all downloads regardless of state
- File: `frontend/src/stores/useDownloadStore.ts` (line 135)
- Measurement: Continuous polling even when downloads are idle
- Cause: No smart polling or exponential backoff
- Improvement path:
  - Only poll when downloads are in active states (downloading, queuing)
  - Increase interval for stable downloads
  - Use WebSocket for real-time updates instead of polling

**No query result caching:**
- Problem: Every search or download list request hits backend
- Files: `frontend/src/api/search.ts`, `frontend/src/api/downloads.ts`
- Cause: No caching layer
- Improvement path: Add short-term caching in Zustand stores, cache search results for 30-60 seconds

## Fragile Areas

**Monitor goroutine context handling:**
- File: `backend/internal/downloads/monitor.go`
- Why fragile: Context misuse could cause hanging goroutines
- Common failures: Shutdown delays, goroutine leaks
- Safe modification: Always use parent context for spawned goroutines, test shutdown paths
- Test coverage: None

**File organization template system:**
- File: `backend/internal/downloads/organization.go`, `backend/internal/fileutil/template.go`
- Why fragile: Complex path manipulation with multiple steps (template → sanitize → join)
- Common failures: Path traversal, invalid characters in paths
- Safe modification: Always validate final paths, unit test edge cases
- Test coverage: None

**qBittorrent client authentication:**
- File: `backend/internal/qbittorrent/client.go`
- Why fragile: Cookie-based session management, errors ignored
- Common failures: Session expiration, authentication failures silently ignored
- Safe modification: Add session validation, retry logic, proper error handling
- Test coverage: None

## Scaling Limits

**SQLite database:**
- Current capacity: Suitable for single-user deployment
- Limit: Limited concurrency (even with WAL mode), not suitable for high-traffic multi-user
- Symptoms at limit: Database locked errors under heavy concurrent writes
- Scaling path: Migrate to PostgreSQL or MySQL for multi-user production deployment

**No connection pooling for external APIs:**
- Current capacity: Direct HTTP calls without pooling
- Limit: May hit rate limits with concurrent requests
- Symptoms at limit: 429 rate limit errors from MAM
- Scaling path: Add connection pooling, request queuing, backoff/retry logic

**In-memory state in monitor:**
- Current capacity: Single-instance only
- Limit: Cannot horizontally scale with multiple API instances
- Symptoms at limit: Duplicate monitoring if multiple instances run
- Scaling path: Use distributed lock (Redis) or job queue for monitor coordination

## Dependencies at Risk

**No dependency version constraints:**
- Risk: Go dependencies have `// indirect` comments but no explicit version pins for some packages
- File: `backend/go.mod`
- Impact: Potential for breaking changes on dependency updates
- Migration plan: Use `go mod tidy` and explicit version constraints

**React 19 is very recent:**
- Risk: React 19.2.0 released recently, potential for ecosystem compatibility issues
- File: `frontend/package.json`
- Impact: Some third-party libraries may not be compatible yet
- Migration plan: Monitor React ecosystem, be prepared to downgrade to React 18 if needed

## Missing Critical Features

**No authentication/authorization:**
- Problem: API endpoints are publicly accessible
- Current workaround: None (assumes private deployment)
- Blocks: Multi-user deployment, production readiness
- Implementation complexity: Medium (add auth middleware, session/token management)

**No error notifications in frontend:**
- Problem: API errors logged to console but not shown to user
- Current workaround: Users don't see error messages
- Blocks: User experience, error recovery
- Implementation complexity: Low (useNotificationStore exists but not fully wired)

**No download retry mechanism:**
- Problem: Failed downloads stay in error state permanently
- Current workaround: Manual deletion and re-download
- Blocks: Reliability for flaky network connections
- Implementation complexity: Medium (add retry logic with backoff)

**No download queue limit:**
- Problem: Users can queue unlimited downloads
- Current workaround: None
- Blocks: Resource management, qBittorrent capacity
- Implementation complexity: Low (add queue limit check in service)

## Test Coverage Gaps

**Entire codebase untested:**
- What's not tested: Everything (no test files present)
- Risk: Bugs in critical paths (search, download, file organization, API handlers)
- Priority: High
- Difficulty to test: Low to Medium (interfaces already present for mocking)

**Critical areas needing tests first:**
1. Download service (`backend/internal/downloads/service.go`)
2. Search service (`backend/internal/search/search_service.go`)
3. HTTP handlers (`backend/internal/server/handlers.go`)
4. File organization (`backend/internal/downloads/organization.go`)
5. API clients (`frontend/src/api/*.ts`)
6. Zustand stores (`frontend/src/stores/*.ts`)

---

*Concerns audit: 2026-01-06*
*Update as issues are fixed or new ones discovered*
