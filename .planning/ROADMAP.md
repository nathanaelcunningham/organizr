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

- [ ] **Phase 1: qBittorrent Integration** - Connect to qBittorrent Web API with authentication and torrent submission
- [ ] **Phase 2: Download Monitoring** - Poll qBittorrent for status updates and detect completion
- [ ] **Phase 3: Configuration System** - Build settings for folder templates and destination paths
- [ ] **Phase 4: File Organization Engine** - Create folder structures and copy files on completion
- [ ] **Phase 5: Frontend Integration** - Connect UI to backend for real-time download status
- [ ] **Phase 6: End-to-End Testing** - Verify complete workflow from search to organized files

## Phase Details

### Phase 1: qBittorrent Integration
**Goal**: Authenticate with qBittorrent Web API and submit torrents from MAM search results
**Depends on**: Nothing (first phase)
**Research**: Skipped (client code exists, Level 0 discovery)
**Research topics**: N/A - existing implementation
**Plans**: 2 plans

Plans:
- [ ] 01-01: Torrent file upload and category support
- [ ] 01-02: Integration testing and error handling

### Phase 2: Download Monitoring
**Goal**: Background monitor that polls qBittorrent for download progress and detects completion
**Depends on**: Phase 1
**Research**: Unlikely (established polling patterns exist in codebase)
**Plans**: TBD

Plans:
- [ ] TBD during planning

### Phase 3: Configuration System
**Goal**: User-configurable folder structure templates and destination paths
**Depends on**: Phase 2
**Research**: Unlikely (standard Go configuration patterns)
**Plans**: TBD

Plans:
- [ ] TBD during planning

### Phase 4: File Organization Engine
**Goal**: Automated folder creation and file copying based on configuration templates
**Depends on**: Phase 3
**Research**: Unlikely (file system operations, Go standard library)
**Plans**: TBD

Plans:
- [ ] TBD during planning

### Phase 5: Frontend Integration
**Goal**: Real-time UI updates showing download status and organization progress
**Depends on**: Phase 4
**Research**: Unlikely (existing React patterns, Zustand store integration)
**Plans**: TBD

Plans:
- [ ] TBD during planning

### Phase 6: End-to-End Testing
**Goal**: Complete workflow verification from MAM search through organized files
**Depends on**: Phase 5
**Research**: Unlikely (testing established features)
**Plans**: TBD

Plans:
- [ ] TBD during planning

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. qBittorrent Integration | 0/2 | Not started | - |
| 2. Download Monitoring | 0/TBD | Not started | - |
| 3. Configuration System | 0/TBD | Not started | - |
| 4. File Organization Engine | 0/TBD | Not started | - |
| 5. Frontend Integration | 0/TBD | Not started | - |
| 6. End-to-End Testing | 0/TBD | Not started | - |
