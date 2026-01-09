# Roadmap: Organizr

## Milestones

- ✅ [v1.0 MVP](milestones/v1.0-MVP.md) — Phases 1-6 (plus 1.1 inserted) — SHIPPED 2026-01-07
- ✅ [v1.1 Enhancements](milestones/v1.1-Enhancements.md) — Phases 7-9 (plus 7.1 inserted) — SHIPPED 2026-01-08
- 🚧 **v1.2 Developer Experience** — Phases 10-15 — IN PROGRESS

## Completed Milestones

<details>
<summary>✅ v1.1 Enhancements (Phases 7-9) — SHIPPED 2026-01-08</summary>

Enhanced search and download experience with MAM series detection, batch operations for multiple downloads, and series number organization templates.

**Stats:** 7 plans, 46 files, +3,636/-264 lines, 2 days

See [milestones/v1.1-Enhancements.md](milestones/v1.1-Enhancements.md) for full details.

</details>

<details>
<summary>✅ v1.0 MVP (Phases 1-6) — SHIPPED 2026-01-07</summary>

Complete audiobook automation from torrent submission through organized files. Full qBittorrent integration, background monitoring, configurable folder templates, auto-organization, and comprehensive testing.

**Stats:** 8 plans, 50 files, 8,968 LOC, 6 days

See [milestones/v1.0-MVP.md](milestones/v1.0-MVP.md) for full details.

</details>

## Current Milestone

### 🚧 v1.2 Developer Experience (In Progress)

**Milestone Goal:** Improve code organization, documentation, testing infrastructure, and developer workflows to make the codebase more maintainable and contributor-friendly.

#### Phase 10: Code Organization & Architecture Review
**Goal**: Audit current structure, identify patterns, document architectural decisions, create contribution guidelines
**Depends on**: Previous milestone complete (v1.1)
**Research**: Unlikely (internal code review and documentation)
**Plans**: 1 complete
**Status**: ✅ Complete (2026-01-08)

Plans:
- [x] 10-01: Architecture Decision Record, contribution guidelines, codebase documentation updates

#### Phase 11: API Layer Cleanup
**Goal**: Standardize error handling, improve request/response patterns, add OpenAPI/Swagger documentation
**Depends on**: Phase 10
**Research**: Unlikely (established patterns, standard API documentation tools)
**Plans**: 3 complete
**Status**: ✅ Complete (2026-01-08)

Plans:
- [x] 11-01: Error handling standardization with typed helpers
- [x] 11-02: OpenAPI/Swagger documentation
- [x] 11-03: Response pattern standardization

#### Phase 12: Testing Infrastructure
**Goal**: Add integration test helpers, improve test coverage reporting, create test data fixtures, document testing patterns
**Depends on**: Phase 11
**Research**: Unlikely (extending existing test patterns)
**Plans**: 1 complete
**Status**: ✅ Complete (2026-01-09)

Plans:
- [x] 12-01: Test helpers, fixtures, and coverage reporting

#### Phase 13: Developer Documentation
**Goal**: Comprehensive README, architecture diagrams, API documentation, setup guides, troubleshooting docs
**Depends on**: Phase 12
**Research**: Unlikely (documentation of existing systems)
**Plans**: TBD

Plans:
- [ ] 13-01: TBD

#### Phase 14: Code Quality Tools
**Goal**: Add linting (golangci-lint, ESLint configs), pre-commit hooks, CI pipeline improvements, code formatting standards
**Depends on**: Phase 13
**Research**: Likely (tool selection and configuration best practices)
**Research topics**: golangci-lint configuration, ESLint + Prettier setup for TypeScript, pre-commit hooks, GitHub Actions workflows
**Plans**: TBD

Plans:
- [ ] 14-01: TBD

#### Phase 15: Refactoring Opportunities
**Goal**: Address technical debt, improve naming consistency, extract reusable utilities, optimize hot paths
**Depends on**: Phase 14
**Research**: Unlikely (internal refactoring based on established patterns)
**Plans**: TBD

Plans:
- [ ] 15-01: TBD

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: qBittorrent Integration** - Connect to qBittorrent Web API with authentication and torrent submission
- [x] **Phase 1.1: qBittorrent Connection Testing (INSERTED)** - Add connection test endpoint and frontend UI for diagnostics
- [x] **Phase 2: Download Monitoring** - Poll qBittorrent for status updates and detect completion
- [x] **Phase 3: Configuration System** - Build settings for folder templates and destination paths
- [x] **Phase 4: File Organization Engine** - Create folder structures and copy files on completion
- [x] **Phase 5: Frontend Integration** - Connect UI to backend for real-time download status
- [x] **Phase 6: End-to-End Testing** - Verify complete workflow from search to organized files
- [x] **Phase 7: MAM Series Detection** - Parse and display series information from MAM search results
- [x] **Phase 7.1: Fix Series Download Field (INSERTED)** - Fix download requests to send series name only
- [x] **Phase 8: Batch Operations** - Support adding multiple torrents at once
- [x] **Phase 9: Series Number Organization** - Add series_number template variable for folder organization
- [x] **Phase 10: Code Organization & Architecture Review** - Audit structure, document decisions, create contribution guidelines
- [ ] **Phase 11: API Layer Cleanup** - Standardize error handling, add OpenAPI/Swagger documentation (1+ plans, in progress)
- [ ] **Phase 12: Testing Infrastructure** - Integration test helpers, coverage reporting, test fixtures
- [ ] **Phase 13: Developer Documentation** - Comprehensive README, architecture diagrams, setup guides
- [ ] **Phase 14: Code Quality Tools** - Linting, pre-commit hooks, CI improvements, formatting standards
- [ ] **Phase 15: Refactoring Opportunities** - Address technical debt, improve consistency, optimize hot paths

## Phase Details

### Phase 1: qBittorrent Integration
**Goal**: Authenticate with qBittorrent Web API and submit torrents from MAM search results
**Depends on**: Nothing (first phase)
**Research**: Skipped (client code exists, Level 0 discovery)
**Research topics**: N/A - existing implementation
**Plans**: 2 plans

Plans:
- [x] 01-01: Torrent file upload and category support
- [x] 01-02: Integration testing and error handling

### Phase 1.1: qBittorrent Connection Testing (INSERTED)
**Goal**: Add backend test endpoint and frontend UI to verify qBittorrent connectivity and authentication
**Depends on**: Phase 1
**Addresses**: ISS-001 (discovered during Phase 1 integration testing)
**Plans**: 1 plan

Plans:
- [x] 01.1-01: Diagnostic endpoint and connection test UI

**Details:**
Urgent insertion to address ISS-001. During Phase 1 testing, qBittorrent authentication issues were encountered. A frontend connection test button helps users diagnose connectivity and authentication problems before attempting downloads. This phase added a dedicated test endpoint and UI component for troubleshooting qBittorrent configuration.

### Phase 2: Download Monitoring
**Goal**: Background monitor that polls qBittorrent for download progress and detects completion
**Depends on**: Phase 1
**Research**: Skipped (existing monitor code, refinement only)
**Plans**: 1 plan

Plans:
- [x] 02-01: Monitor refinement with resilience and remote path support

### Phase 3: Configuration System
**Goal**: User-configurable folder structure templates and destination paths
**Depends on**: Phase 2
**Research**: Skipped (standard validation patterns)
**Plans**: 1 plan

Plans:
- [x] 03-01: Template validation and path preview

### Phase 4: File Organization Engine
**Goal**: Automated folder creation and file copying based on configuration templates
**Depends on**: Phase 3
**Research**: Skipped (validated existing implementation)
**Plans**: 1 plan

Plans:
- [x] 04-01: Organization testing and validation

### Phase 5: Frontend Integration
**Goal**: Real-time UI updates showing download status and organization progress
**Depends on**: Phase 4
**Research**: Skipped (existing integration verified and enhanced)
**Plans**: 1 plan

Plans:
- [x] 05-01: Frontend integration verification and UX enhancements

### Phase 6: End-to-End Testing
**Goal**: Complete workflow verification from MAM search through organized files
**Depends on**: Phase 5
**Research**: Skipped (testing established features using existing patterns)
**Plans**: 1 plan

Plans:
- [x] 06-01: Backend handler and monitor tests with E2E manual test scenarios

### Phase 7: MAM Series Detection
**Goal**: Parse and extract series information from MAM search results and display in UI
**Depends on**: Previous milestone complete
**Research**: Completed (MAM JSON structure, series metadata format)
**Research topics**: MAM API series_info format, structured data parsing
**Plans**: 2 plans

Plans:
- [x] 07-01: Backend series parsing - return structured series data
- [x] 07-02: Frontend series grouping and display

### Phase 7.1: Fix Series Download Field (INSERTED)
**Goal**: Fix download requests to send series name only, not name with number
**Depends on**: Phase 7
**Plans**: 1 plan

Plans:
- [x] 07.1-01: Fix series field in download handlers

**Details:**
Urgent fix discovered during Phase 7 testing. When clicking download, the series field being sent includes the series number (e.g., "Discworld #1" instead of "Discworld"), causing the organization function to create incorrect folder structures like `Discworld #1/book.m4b` instead of `Discworld/book.m4b`. Must extract series name only before sending to backend.

### Phase 8: Batch Operations
**Goal**: Support adding multiple torrents simultaneously from search results
**Depends on**: Phase 7
**Research**: Skipped (internal patterns, existing qBittorrent integration)
**Plans**: 2 plans

Plans:
- [x] 08-01: Backend batch endpoint with partial success handling
- [x] 08-02: Frontend multi-select and batch UI

### Phase 9: Series Number Organization
**Goal**: Add {series_number} template variable for folder organization
**Depends on**: Phase 8
**Research**: Skipped (extending existing template system)
**Plans**: 3 plans

Plans:
- [x] 09-01: Backend data model - add SeriesNumber field to Download
- [x] 09-02: Template and organization - support {series_number} variable
- [x] 09-03: Frontend integration - series_number input and display

## Progress

| Phase | Milestone | Plans Complete | Status | Completed |
|-------|-----------|----------------|--------|-----------|
| 1. qBittorrent Integration | v1.0 | 2/2 | Complete | 2026-01-06 |
| 1.1 Connection Testing (INSERTED) | v1.0 | 1/1 | Complete | 2026-01-06 |
| 2. Download Monitoring | v1.0 | 1/1 | Complete | 2026-01-06 |
| 3. Configuration System | v1.0 | 1/1 | Complete | 2026-01-07 |
| 4. File Organization Engine | v1.0 | 1/1 | Complete | 2026-01-07 |
| 5. Frontend Integration | v1.0 | 1/1 | Complete | 2026-01-07 |
| 6. End-to-End Testing | v1.0 | 1/1 | Complete | 2026-01-07 |
| 7. MAM Series Detection | v1.1 | 2/2 | ✅ Shipped | 2026-01-07 |
| 7.1 Fix Series Download Field (INSERTED) | v1.1 | 1/1 | ✅ Shipped | 2026-01-08 |
| 8. Batch Operations | v1.1 | 2/2 | ✅ Shipped | 2026-01-08 |
| 9. Series Number Organization | v1.1 | 3/3 | ✅ Shipped | 2026-01-08 |
| 10. Code Organization & Architecture Review | v1.2 | 1/1 | ✅ Complete | 2026-01-08 |
| 11. API Layer Cleanup | v1.2 | 3/3 | ✅ Complete | 2026-01-08 |
| 12. Testing Infrastructure | v1.2 | 1/1 | ✅ Complete | 2026-01-09 |
| 13. Developer Documentation | v1.2 | 0/? | Not started | - |
| 14. Code Quality Tools | v1.2 | 0/? | Not started | - |
| 15. Refactoring Opportunities | v1.2 | 0/? | Not started | - |
