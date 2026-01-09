# Organizr

> Automated audiobook torrent management with qBittorrent integration and smart folder organization

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue)](https://go.dev/)
[![Node Version](https://img.shields.io/badge/Node-20%2B-green)](https://nodejs.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

Organizr automates the complete workflow for audiobook torrent downloads: search MyAnonamouse for audiobooks, send torrents to qBittorrent, monitor download progress in real-time, and automatically organize completed files into Audiobookshelf-compatible folder structures. No more manual file management - just search, download, and listen.

Perfect for audiobook collectors who want their library organized consistently without the tedious copy/paste/rename workflow.

## Features

- **MyAnonamouse Integration** - Search torrents directly from the UI with series detection and metadata extraction
- **qBittorrent Automation** - One-click torrent submission with category support and connection testing
- **Background Monitoring** - Automatic progress tracking with real-time UI updates
- **Smart Organization** - Template-based folder creation (e.g., `{author}/{series}/{title}`) with automatic file copying
- **Series Detection** - Parse and display series information from MAM results, group by series in UI
- **Batch Operations** - Select and download multiple audiobooks at once with partial success handling
- **Template Variables** - Support for `{author}`, `{title}`, `{series}`, `{series_number}` with live preview
- **Path Mapping** - Remote qBittorrent support for Docker and network share deployments

## Quick Start

### Prerequisites

- **Go 1.23+** - [Download](https://go.dev/dl/)
- **Node 20+** - [Download](https://nodejs.org/)
- **qBittorrent with Web UI enabled** - [Download](https://www.qbittorrent.org/download.php)
  - Enable Web UI in qBittorrent: Tools → Options → Web UI
  - Note your port (default: 8080) and credentials

### Backend Setup

```bash
cd backend
make build    # Build the Go binary
make run      # Start the backend server (port 8080)
```

### Frontend Setup

```bash
cd frontend
npm install           # Install dependencies
npm run dev           # Start dev server (port 5173)
```

Visit `http://localhost:5173` to access the UI.

### Basic Configuration

1. **Configure qBittorrent connection** (via UI or API):
   - Go to Settings in the UI
   - Enter qBittorrent URL (e.g., `http://localhost:8080`)
   - Enter Web UI username and password
   - Click "Test Connection" to verify

2. **Set destination path**:
   - Configure where organized audiobooks should be copied (e.g., `/mnt/nas/audiobooks`)

3. **Customize folder templates** (optional):
   - Default: `{author}/{series}/{title}` for books in a series
   - Default: `{author}/{title}` for standalone books
   - Preview your template before saving

> **💡 Production Deployment:** The Quick Start above is for local development. For production deployment with Docker, see [Deployment](#deployment) below or jump directly to [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md).

## Deployment

Organizr can be deployed with Docker (recommended for production) or as standalone binaries.

**Docker Deployment (Recommended):**

```bash
# 1. Create .env file with your configuration
cp .env.example .env
# Edit .env with your qBittorrent URL, credentials, and paths

# 2. Start services with Docker Compose
docker compose up -d

# 3. Access Organizr at http://localhost:3000
```

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for comprehensive Docker deployment instructions, including environment variable configuration and health checks.

**Unraid Deployment:** Organizr supports Unraid with docker-compose.yml deployment. See the [Unraid section](docs/DEPLOYMENT.md#unraid-deployment) for volume mapping examples and path configuration.

**Bare Metal Deployment:** For systemd service deployment without Docker, see the [systemd section](docs/DEPLOYMENT.md#backend-deployment) for service configuration and binary installation.

**Troubleshooting:** See [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) for Docker-specific issues, connection problems, and common deployment errors.

## Screenshots

_Coming soon - screenshots of search interface, download monitoring, and configuration settings_

## Core Concepts

### Search

Search MyAnonamouse directly from the UI. Results include series information, metadata, and torrent details. Books in the same series are automatically grouped for easy browsing.

### Download

Click to send torrents to qBittorrent. Organizr uploads the torrent file (with authentication for private trackers) and starts tracking progress immediately. Batch select multiple audiobooks to download them all at once.

### Monitor

Background service polls qBittorrent every few seconds for progress updates. Real-time UI shows download status, speed, ETA, and completion percentage. Resilient to qBittorrent restarts.

### Organization

When downloads complete, Organizr automatically creates folder structures based on your template and copies files to the destination. Templates use metadata variables like `{author}`, `{series}`, and `{title}` to build paths that Audiobookshelf can parse correctly.

## Configuration

See [backend/docs/CONFIGURATION.md](backend/docs/CONFIGURATION.md) for complete configuration reference.

**Available template variables:**
- `{author}` - Author name
- `{title}` - Book title
- `{series}` - Series name (empty if book has no series)
- `{series_number}` - Series number (empty if book has no series)

**Path mapping for remote deployments:**
- Configure `paths.mount_point` to prepend a local path to qBittorrent's reported paths
- Useful for Docker containers or network shares where qBittorrent sees different paths than Organizr

## API Documentation

See [backend/docs/API.md](backend/docs/API.md) for complete REST API reference.

Key endpoints:
- `POST /api/downloads` - Create download
- `GET /api/downloads` - List all downloads
- `GET /api/search?q=<query>` - Search torrents
- `GET /api/config` - Get configuration
- `PUT /api/config/{key}` - Update configuration

## Development

### Running Tests

```bash
# Backend tests
cd backend
make test                    # Run all tests
make test-coverage           # Run with coverage report
go test -race ./...          # Run with race detection

# Frontend tests
cd frontend
npm test                     # Run Vitest tests
npm test -- --coverage       # Run with coverage report
```

**Coverage thresholds:** 60% for both backend and frontend (baseline quality metric)

### Architecture

See [docs/architecture/ADR.md](docs/architecture/ADR.md) for architectural decisions and technical rationale.

Key patterns:
- **Repository pattern** for data access
- **Background monitor** goroutine for download tracking
- **Typed error helpers** for consistent API responses
- **Zustand stores** for frontend state management

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development setup instructions
- Code organization and conventions
- Pull request guidelines
- Testing requirements

## Project Structure

```
organizr/
├── backend/               # Go REST API server
│   ├── cmd/api/          # Application entry point
│   ├── internal/         # Business logic, services, handlers
│   └── assets/           # Embedded resources (migrations)
├── frontend/             # React single-page application
│   ├── src/              # Components, pages, stores, API clients
│   └── public/           # Static assets
└── docs/                 # Documentation
    └── architecture/     # Architecture decisions and diagrams
```

## Tech Stack

**Backend:**
- Go 1.23 - REST API server
- Chi - HTTP router and middleware
- SQLite - Embedded database with WAL mode
- qBittorrent Web API - Torrent management

**Frontend:**
- React 19 - UI framework
- TypeScript 5.9 - Type safety
- Vite 7 - Build tooling and dev server
- Zustand - State management
- TailwindCSS - Utility-first styling

**Integrations:**
- qBittorrent Web API - Torrent automation
- MyAnonamouse - Audiobook torrent search

## Troubleshooting

See [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) for common issues and solutions.

Quick checks:
- **qBittorrent connection fails**: Verify Web UI enabled, check URL format, test credentials
- **Downloads not organizing**: Check destination path exists, verify template syntax, ensure qBittorrent finished download
- **Frontend can't reach backend**: Verify backend running on port 8080, check CORS settings

## License

MIT License - see [LICENSE](LICENSE) for details

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Key areas for contribution:
- Additional torrent site integrations
- Enhanced metadata extraction
- UI/UX improvements
- Documentation and examples
