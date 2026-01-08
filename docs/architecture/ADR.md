# Architecture Decision Record (ADR)

**Project:** Organizr - Audiobook Torrent Automation
**Version Coverage:** v1.0 MVP, v1.1 Enhancements
**Last Updated:** 2026-01-08

This document captures key technical decisions made during the initial development of Organizr, documenting the context, rationale, and tradeoffs for major architectural choices.

---

## 1. Technology Stack Choices

### 1.1 Go for Backend

**Context:** Need a backend language for REST API, background monitoring, file operations, and qBittorrent integration. Requirements include good concurrency support, static typing, and reliable tooling.

**Decision:** Use Go 1.23 for backend implementation.

**Rationale:**
- **Concurrency:** Goroutines and channels provide lightweight concurrency perfect for background monitoring without complex thread management
- **Static typing:** Compile-time type safety catches errors early, critical for file operations and API contracts
- **Standard library:** Excellent stdlib for HTTP servers, JSON handling, file I/O - minimal external dependencies needed
- **Single binary:** Compiles to standalone executable, simplifies deployment (no runtime dependencies)
- **Performance:** Fast enough for file operations and API responses, low memory footprint

**Consequences:**
- ✅ Goroutine-based monitor is clean and efficient
- ✅ Interfaces enable excellent testability
- ✅ Single binary deployment is simple
- ⚠️ Smaller ecosystem than Node.js (but stdlib sufficient for our needs)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 1.2 React 19 + TypeScript + Vite for Frontend

**Context:** Need modern frontend framework for search interface, real-time download tracking, and configuration management. Must support type safety and fast development iteration.

**Decision:** Use React 19 with TypeScript and Vite build tooling.

**Rationale:**
- **React 19:** Latest stable version with improved performance, familiar component model
- **TypeScript:** Type safety catches bugs at compile time, excellent IDE support, self-documenting code
- **Vite:** Lightning-fast HMR (hot module replacement) for rapid development, modern ESM-based build
- **Ecosystem:** Massive library ecosystem (React Hook Form, Zustand, TailwindCSS)

**Consequences:**
- ✅ Type safety across frontend codebase
- ✅ Sub-second hot reload speeds up development
- ✅ Rich component library ecosystem
- ⚠️ React 19 is very recent - some third-party libraries may lag in compatibility

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 1.3 SQLite with WAL Mode

**Context:** Need persistent storage for downloads, configuration, and state. Single-user deployment, embedded database preferred over network database.

**Decision:** Use SQLite with Write-Ahead Logging (WAL) mode.

**Rationale:**
- **Single-user fit:** Designed for embedded use cases, perfect for single-user application
- **Zero configuration:** No database server to install/manage, database is just a file
- **WAL mode:** Enables concurrent reads while writing (critical for monitor + API queries)
- **Sufficient performance:** More than adequate for expected data volumes (hundreds of downloads)
- **Simple backup:** Copy a single .db file to backup entire database

**Consequences:**
- ✅ Zero operational overhead (no DB server)
- ✅ Simple deployment and backup
- ✅ WAL mode eliminates "database locked" errors during monitoring
- ⚠️ Not suitable for high-concurrency multi-user deployment (would need PostgreSQL migration)
- ⚠️ No built-in replication (acceptable for single-user case)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 1.4 Chi Router

**Context:** Need HTTP router for Go backend with middleware support, route parameters, and standard net/http compatibility.

**Decision:** Use Chi router instead of alternatives (Gin, Echo, Gorilla Mux).

**Rationale:**
- **Lightweight:** Minimal abstraction over net/http, doesn't hide stdlib patterns
- **Standard middleware:** Compatible with standard net/http middleware (easy to add CORS, logging, etc.)
- **Context-aware:** Uses context.Context for request-scoped values (idiomatic Go)
- **Route patterns:** Supports route parameters and wildcards without complex DSL
- **No magic:** Explicit routing, no hidden reflection or code generation

**Consequences:**
- ✅ Easy to understand for Go developers familiar with net/http
- ✅ Standard middleware ecosystem works out-of-box
- ✅ Fast enough for our use case (not a bottleneck)
- ⚠️ Less "batteries included" than Gin (no built-in validation, binding)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 1.5 Zustand for State Management

**Context:** Need frontend state management for downloads, search results, configuration, and notifications. Options: Redux, Context API, Zustand, Jotai.

**Decision:** Use Zustand for global state management.

**Rationale:**
- **Simplicity:** Minimal boilerplate compared to Redux (no actions, reducers, thunks)
- **Performance:** Selector-based subscriptions prevent unnecessary re-renders
- **TypeScript support:** Excellent TypeScript inference, type-safe stores
- **Small bundle size:** ~1KB gzipped vs Redux's larger footprint
- **Hooks-first:** Natural API for React hooks (`useDownloadStore`, `useSearchStore`)

**Consequences:**
- ✅ Easy to add new stores without boilerplate
- ✅ Fast implementation velocity
- ✅ Good performance with polling (no render thrashing)
- ⚠️ Less structured than Redux (no enforced action patterns) - acceptable tradeoff for smaller team

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

## 2. Architecture Patterns

### 2.1 Repository Pattern with Interface-Based Dependency Injection

**Context:** Need clean separation between business logic and data access. Want to enable testing without database.

**Decision:** Use repository pattern with Go interfaces for all data access.

**Rationale:**
- **Testability:** Services depend on interfaces, tests can use mock repositories
- **Flexibility:** Can swap SQLite for PostgreSQL by implementing same interface
- **Clear boundaries:** Service layer doesn't know about SQL, persistence layer doesn't know about business rules
- **Interface-based DI:** Constructor injection makes dependencies explicit

**Consequences:**
- ✅ Comprehensive handler tests achieved using mock repositories
- ✅ Zero race conditions due to clear concurrency boundaries
- ✅ Easy to understand data flow (handler → service → repository)
- ⚠️ Slightly more code (interface + implementation) - worthwhile for testability

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 2.2 Service Layer Separation

**Context:** HTTP handlers were getting large, mixing request parsing, validation, business logic, and response formatting.

**Decision:** Extract business logic into dedicated service layer (ConfigService, DownloadService, MAMService).

**Rationale:**
- **Single Responsibility:** Handlers focus on HTTP concerns (parsing, validation, response codes)
- **Reusability:** Service methods can be called from multiple handlers or background tasks
- **Testability:** Services can be tested independently of HTTP layer
- **Business logic isolation:** Core logic doesn't depend on HTTP framework

**Consequences:**
- ✅ Handlers are thin, easy to read
- ✅ Business logic is reusable (e.g., organization logic used by both manual and automatic triggers)
- ✅ Testing is focused (handler tests mock services, service tests mock repositories)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 2.3 Background Monitoring via Goroutines

**Context:** Need continuous polling of qBittorrent for download status without blocking API requests.

**Decision:** Use goroutine-based monitor with context cancellation for graceful shutdown.

**Rationale:**
- **Go idiom:** Goroutines are designed for concurrent background tasks
- **Lightweight:** Minimal overhead compared to thread-based approaches
- **Context control:** context.Context enables graceful shutdown
- **Channel communication:** Safe cross-goroutine communication without locks

**Consequences:**
- ✅ Monitor runs independently without affecting API latency
- ✅ Graceful shutdown works reliably (all goroutines exit cleanly)
- ✅ Zero race conditions achieved with proper context usage
- ⚠️ Context misuse found in v1.0 (organization goroutines using context.Background instead of parent context) - requires fix

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 2.4 Polling vs WebSockets for Real-Time Updates

**Context:** Frontend needs real-time download status updates. Options: WebSockets, Server-Sent Events (SSE), or polling.

**Decision:** Use HTTP polling (3-second intervals) instead of WebSockets or SSE.

**Rationale:**
- **Simplicity:** No connection management, no reconnection logic, standard HTTP requests
- **Stateless backend:** No need to track active connections or handle connection drops
- **Good enough performance:** 3-second updates sufficient for user experience (not millisecond-critical)
- **Auto-stop optimization:** Polling automatically stops when no downloads active (smart polling implementation)
- **Easier debugging:** Standard HTTP requests, visible in network tab

**Consequences:**
- ✅ Backend remains stateless (easier to scale horizontally in future)
- ✅ No WebSocket complexity (reconnection, heartbeats, etc.)
- ✅ Auto-stop optimization prevents unnecessary requests
- ⚠️ Slightly higher latency than WebSockets (3s vs real-time) - acceptable tradeoff
- ⚠️ More HTTP requests than WebSockets (mitigated by auto-stop)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 2.5 Template-Based Folder Organization

**Context:** Need flexible folder structure for organized audiobooks. Different users prefer different patterns (Author/Series/Book vs Author/Book vs Series/Book).

**Decision:** Use template system with variables ({author}, {series}, {title}, {series_number}) instead of hardcoded structure.

**Rationale:**
- **User control:** Users can define their own folder structure
- **Flexibility:** Supports any pattern without code changes
- **Validation:** Real-time preview validates template before use
- **Sanitization:** Path sanitization prevents security issues (no path traversal)

**Consequences:**
- ✅ Works for any folder preference (Audiobookshelf, Plex, manual organization)
- ✅ Template validation prevents errors before downloads complete
- ✅ Real-time preview provides immediate feedback
- ⚠️ More complex than hardcoded paths (worthwhile for flexibility)

**Status:** Accepted (v1.0), Enhanced (v1.1 with {series_number})
**Decided:** 2026-01-06, Extended: 2026-01-08

---

## 3. Data Model Decisions

### 3.1 SeriesInfo as Structured Array

**Context:** MAM provides series information as structured data with ID, name, and number. Need to decide between structured representation vs concatenated string.

**Decision:** Store series as `[]SeriesInfo` array with ID, Name, and Number fields instead of concatenated string like "Discworld #1".

**Rationale:**
- **Data fidelity:** Preserves structured information from MAM API
- **Enables grouping:** Frontend can group search results by series ID
- **Enables sorting:** Can sort by series number (1, 2, 3) instead of string comparison
- **Multiple series support:** Books can belong to multiple series (e.g., "Discworld" and "Discworld - Death")

**Consequences:**
- ✅ Frontend can group and sort series properly
- ✅ Series ID enables precise grouping (handles name variations)
- ✅ Supports books in multiple series
- ⚠️ Slightly more complex than string (worthwhile for features enabled)

**Status:** Accepted (v1.1)
**Decided:** 2026-01-07

---

### 3.2 Series Numbers as Strings

**Context:** Series numbers come in various formats: "1", "1.5", "Book 1", "1-2", "Omnibus 1-3".

**Decision:** Store SeriesNumber field as string (TEXT in SQLite) instead of integer or float.

**Rationale:**
- **Format flexibility:** Accommodates any format MAM provides
- **No data loss:** Doesn't force numeric conversion that could lose information
- **Frontend parsing:** Frontend can parse and display as needed (extract leading digits for sorting)
- **Simple implementation:** No complex parsing logic in backend

**Consequences:**
- ✅ Handles all series number formats without special cases
- ✅ Frontend has flexibility in display (can extract numbers if needed)
- ⚠️ Sorting requires frontend logic (acceptable - frontend already does sorting)

**Status:** Accepted (v1.1)
**Decided:** 2026-01-08

---

### 3.3 Empty series_number as Empty String

**Context:** Not all books in series have numbers (standalone books in series, unnumbered entries). Need to decide between empty string, null, or omitted field.

**Decision:** Use empty string ("") instead of null for missing series numbers.

**Rationale:**
- **Template simplicity:** Template system replaces {series_number} with empty string cleanly
- **No null handling:** Avoids null checks throughout codebase
- **Database compatibility:** SQLite TEXT columns default to empty string easily
- **JSON serialization:** Empty string serializes consistently (null requires omitempty handling)

**Consequences:**
- ✅ Templates work without special null handling
- ✅ Consistent API responses (always string, never null)
- ✅ Simple database migrations (DEFAULT '')
- ⚠️ Empty string vs null distinction lost (acceptable - functionally equivalent for our use case)

**Status:** Accepted (v1.1)
**Decided:** 2026-01-08

---

### 3.4 First Series as Primary for Organization

**Context:** Books can belong to multiple series (e.g., main series + subseries). Need to decide which series to use for folder organization.

**Decision:** Use first series in SeriesInfo array for {series} template variable.

**Rationale:**
- **Simple rule:** No complex logic or configuration needed
- **MAM ordering:** MAM typically lists main series first, then subseries
- **Consistency:** Predictable behavior for users
- **Good enough:** Main series is usually correct choice

**Consequences:**
- ✅ Simple implementation (just use `series[0]`)
- ✅ Usually correct (main series first in MAM data)
- ⚠️ Edge cases possible (if MAM orders differently) - rare, acceptable tradeoff

**Status:** Accepted (v1.1)
**Decided:** 2026-01-08

---

## 4. API Design Decisions

### 4.1 Partial Success Pattern for Batch Operations

**Context:** Batch downloads can partially fail (some torrents valid, others invalid). Need to decide between all-or-nothing vs partial success.

**Decision:** Return HTTP 200 with separate `successful` and `failed` arrays in batch response.

**Rationale:**
- **User feedback:** Users see exactly which downloads succeeded and which failed
- **Graceful degradation:** Some progress is better than none (don't fail entire batch for one error)
- **Detailed errors:** Each failure includes index, request, and error message
- **Resumable:** Users can retry only the failed items

**Consequences:**
- ✅ Better UX than all-or-nothing (5 succeed, 2 fail is better than 0 succeed)
- ✅ Clear error reporting (users know exactly what failed and why)
- ✅ Frontend can show granular notifications ("3 downloads started, 2 failed")
- ⚠️ HTTP 200 with errors might confuse some clients (acceptable - response structure is clear)

**Status:** Accepted (v1.1)
**Decided:** 2026-01-08

---

### 4.2 Sequential Batch Processing

**Context:** Batch downloads can add multiple torrents to qBittorrent. Need to decide between sequential vs concurrent processing.

**Decision:** Process batch downloads sequentially (one at a time) instead of concurrently.

**Rationale:**
- **Prevent overwhelming qBittorrent:** Concurrent requests could overload qBittorrent Web API
- **Predictable behavior:** Sequential processing is easier to reason about
- **Error isolation:** Failures don't cause cascading issues
- **Good enough performance:** Even 50 sequential requests complete in <10 seconds

**Consequences:**
- ✅ Stable qBittorrent integration (no rate limiting or connection issues)
- ✅ Predictable order (downloads added in request order)
- ⚠️ Slightly slower than concurrent (acceptable - not time-critical)

**Status:** Accepted (v1.1)
**Decided:** 2026-01-08

---

### 4.3 50-Item Batch Limit

**Context:** Need to prevent abuse and maintain system stability for batch operations.

**Decision:** Enforce maximum of 50 items per batch request, return 400 error if exceeded.

**Rationale:**
- **Prevent abuse:** Stops malicious or accidental huge batches
- **System stability:** Limits memory usage and processing time
- **Reasonable limit:** 50 items more than sufficient for legitimate use cases
- **Clear error:** 400 Bad Request with "batch too large" message

**Consequences:**
- ✅ System protected from resource exhaustion
- ✅ Predictable worst-case performance
- ⚠️ Users wanting >50 must split into multiple batches (rare case, acceptable)

**Status:** Accepted (v1.1)
**Decided:** 2026-01-08

---

### 4.4 RESTful Endpoints with JSON

**Context:** Need API design pattern for HTTP endpoints and payload format.

**Decision:** Use RESTful HTTP endpoints with JSON request/response bodies.

**Rationale:**
- **Standard:** Widely understood, large ecosystem
- **Tooling support:** Every HTTP client supports JSON
- **Type safety:** TypeScript interfaces map directly to JSON
- **Debuggability:** JSON is human-readable in network tools

**Consequences:**
- ✅ Standard patterns (POST for create, GET for list, etc.)
- ✅ Excellent tooling support (curl, Postman, browser devtools)
- ✅ Easy to document (OpenAPI/Swagger compatibility)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

## 5. Error Handling Strategy

### 5.1 User-Friendly Error Messages

**Context:** Need to decide what information to include in API error responses.

**Decision:** Return user-friendly error messages without exposing internal implementation details.

**Rationale:**
- **Security:** Don't expose stack traces, database schema, file paths
- **UX:** Users need actionable messages ("Invalid torrent URL") not technical details
- **Debugging:** Log full errors server-side for debugging, show friendly message to user

**Consequences:**
- ✅ No information leakage to potential attackers
- ✅ Better user experience (clear, actionable messages)
- ⚠️ Debugging requires server logs (acceptable - proper logging implemented)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 5.2 5-Minute Timeout for File Organization

**Context:** File organization can take long time for large audiobooks (multi-GB files). Need to balance patience with preventing hanging.

**Decision:** Use 5-minute context timeout for organization operations.

**Rationale:**
- **Large file support:** Multi-GB audiobooks can take minutes to copy
- **Prevents hanging:** If something goes wrong, don't wait forever
- **Graceful shutdown:** Timeout enables clean application shutdown

**Consequences:**
- ✅ Large files (up to ~10GB) complete successfully
- ✅ Application can shut down cleanly (doesn't wait forever)
- ⚠️ Extremely large files (>10GB) might timeout (rare, users can increase timeout if needed)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 5.3 Monitor Resilience When qBittorrent Unavailable

**Context:** qBittorrent might be temporarily unavailable (restarting, network issues). Need to decide whether monitor should stop or continue.

**Decision:** Continue monitoring when qBittorrent queries fail, log warnings but don't stop monitor.

**Rationale:**
- **Resilience:** Temporary qBittorrent unavailability shouldn't require app restart
- **Self-healing:** Monitor automatically resumes when qBittorrent comes back
- **Operational simplicity:** Users don't need to restart Organizr when qBittorrent restarts

**Consequences:**
- ✅ Resilient to transient failures
- ✅ Monitor resumes automatically
- ⚠️ Warning logs during downtime (acceptable - users can see in logs)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 5.4 All-or-Nothing Copy Operations with Cleanup

**Context:** File organization involves copying multiple files. If copy fails halfway, need to decide whether to leave partial copies or clean up.

**Decision:** Delete all copied files if any copy fails (all-or-nothing pattern).

**Rationale:**
- **Consistency:** Never leave incomplete audiobook folders
- **User clarity:** Either complete audiobook exists or nothing (no partial state)
- **Retry-safe:** User can safely retry without cleanup

**Consequences:**
- ✅ No partial/incomplete audiobook folders
- ✅ Clean state for retries
- ⚠️ Disk space usage during copy (need free space for full audiobook temporarily)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-07

---

## 6. Testing Philosophy

### 6.1 Interface-Based Testing with Mock Dependencies

**Context:** Need testing strategy for services that depend on database, external APIs, and file system.

**Decision:** Use interface-based dependency injection with mock implementations for tests.

**Rationale:**
- **Fast tests:** Mock dependencies eliminate database/network I/O
- **Isolated tests:** Each test verifies one component without side effects
- **Deterministic:** No flaky tests due to external factors
- **Comprehensive:** Can test error cases that are hard to reproduce with real dependencies

**Consequences:**
- ✅ Comprehensive handler tests achieved (100+ test cases)
- ✅ Fast test suite (seconds, not minutes)
- ✅ Tests are deterministic (no flakes)
- ✅ Easy to test edge cases and error conditions

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 6.2 Race Detection Requirement

**Context:** Application uses goroutines for background monitoring. Concurrent code can have race conditions that only appear under load.

**Decision:** All concurrent code must pass `go test -race` before commit.

**Rationale:**
- **Catch races early:** Race detector finds data races at test time, not production
- **Concurrent safety:** Zero tolerance for race conditions
- **Go tooling:** Race detector is built into Go toolchain

**Consequences:**
- ✅ Zero race conditions found in v1.0 and v1.1
- ✅ Confidence in concurrent code correctness
- ⚠️ Race detector adds overhead to test runs (acceptable - only run when testing concurrent code)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 6.3 Comprehensive Handler Tests

**Context:** HTTP handlers are API contract - need to verify all request/response scenarios.

**Decision:** Write comprehensive tests for all handlers covering success, validation errors, and service errors.

**Rationale:**
- **API contract verification:** Tests document expected behavior
- **Regression prevention:** Changes that break API are caught immediately
- **Confidence in refactoring:** Can refactor internals knowing API contract is verified

**Consequences:**
- ✅ All handlers tested with multiple scenarios
- ✅ API contract is documented by tests
- ✅ Refactoring is safer

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 6.4 Zero Race Conditions Achieved

**Context:** Goal was to achieve zero race conditions in concurrent code.

**Decision:** Declare "zero race conditions" as a quality bar, verify with race detector.

**Rationale:**
- **Quality signal:** Zero races demonstrates mature concurrent code
- **User trust:** No mysterious crashes or corruption due to races
- **Maintainability:** Race-free code is easier to understand and modify

**Consequences:**
- ✅ No race conditions in v1.0 or v1.1
- ✅ Concurrent code is correct and safe
- ✅ Builds confidence in Go concurrency patterns

**Status:** Accepted (v1.0), Verified (v1.1)
**Decided:** 2026-01-06

---

## 7. Security Posture

### 7.1 Single-User Deployment Assumption

**Context:** Need to decide whether to implement authentication/authorization in v1.

**Decision:** Ship v1 without authentication, assume single-user private deployment.

**Rationale:**
- **Primary use case:** Self-hosted single-user tool
- **Scope management:** Authentication is complex, defer to future milestone
- **Deployment model:** Private network deployment (localhost or LAN)

**Consequences:**
- ✅ Faster initial development (auth is significant effort)
- ✅ Simpler codebase for v1
- ⚠️ Not suitable for public deployment without authentication layer
- ⚠️ Must add authentication before multi-user or internet-facing deployment

**Status:** Accepted (v1.0), Deferred (authentication will be added in future milestone)
**Decided:** 2026-01-06

---

### 7.2 Private Tracker Authentication via Torrent File Downloads

**Context:** MyAnonamouse requires authentication to download torrent files. Generic torrent URLs need direct upload, MAM URLs need torrent file download first.

**Decision:** Detect MAM URLs and download torrent file before uploading to qBittorrent.

**Rationale:**
- **MAM requirement:** MAM embeds authentication in torrent file download
- **User cookies:** Use user's browser cookies to authenticate MAM requests
- **Transparent:** Works seamlessly for users

**Consequences:**
- ✅ MAM downloads work correctly with authentication
- ✅ Users don't need to manually download torrent files
- ⚠️ Requires user's MAM cookies (users must be logged into MAM in browser)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

### 7.3 Path Sanitization for File Organization

**Context:** User data (author, series, title) is used to construct file paths. Need to prevent path traversal attacks.

**Decision:** Sanitize all user-derived data before using in file paths.

**Rationale:**
- **Security:** Prevent path traversal (e.g., author: "../../etc")
- **Cross-platform:** Handle Windows, Mac, Linux path restrictions
- **Invalid characters:** Replace characters that filesystems reject

**Consequences:**
- ✅ No path traversal vulnerability
- ✅ Works across platforms
- ⚠️ Special characters in names get replaced (acceptable - files still organized correctly)

**Status:** Accepted (v1.0)
**Decided:** 2026-01-06

---

## Summary of Status

**Accepted and Stable (v1.0):**
- All technology stack choices
- Core architecture patterns
- Error handling strategy
- Testing philosophy
- Security baseline

**Extended in v1.1:**
- SeriesInfo structured data model
- Series number support
- Batch operations with partial success
- Template system extended with {series_number}

**Future Work:**
- Authentication/authorization (deferred from v1)
- PostgreSQL migration for multi-user scaling (if needed)
- WebSocket real-time updates (if polling insufficient)

---

*Last updated: 2026-01-08*
*Covers: v1.0 MVP (Phases 1-6), v1.1 Enhancements (Phases 7-9)*
