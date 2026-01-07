# Roadmap: Organizr

## Overview

Build complete audiobook automation from torrent discovery to organized files. Start with qBittorrent integration for torrent submission, add monitoring for download tracking, implement configurable folder organization, and connect frontend for real-time status updates. Each phase delivers a testable capability, culminating in end-to-end workflow validation.

## Domain Expertise

None

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
- [ ] **Phase 6: End-to-End Testing** - Verify complete workflow from search to organized files

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
**Research**: Unlikely (testing established features)
**Plans**: TBD

Plans:
- [ ] TBD during planning

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 1.1 → 2 → 3 → 4 → 5 → 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. qBittorrent Integration | 2/2 | Complete | 2026-01-06 |
| 1.1 Connection Testing (INSERTED) | 1/1 | Complete | 2026-01-06 |
| 2. Download Monitoring | 1/1 | Complete | 2026-01-06 |
| 3. Configuration System | 1/1 | Complete | 2026-01-07 |
| 4. File Organization Engine | 1/1 | Complete | 2026-01-07 |
| 5. Frontend Integration | 1/1 | Complete | 2026-01-07 |
| 6. End-to-End Testing | 0/TBD | Not started | - |
