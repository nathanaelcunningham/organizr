# Deployment Guide

This guide covers deploying Organizr to production. Organizr consists of a Go backend binary, React frontend static files, and a SQLite database, with qBittorrent running as a separate service.

---

## Overview

**Organizr Architecture:**
- **Backend** - Go binary serving REST API on port 8080
- **Frontend** - React single-page application (static HTML/JS/CSS files)
- **Database** - SQLite database file (created automatically on first run)
- **qBittorrent** - Torrent client with Web UI enabled (can be same server or remote)

**Deployment Patterns:**
1. **All-in-one** - Backend, frontend (served by backend or nginx), and qBittorrent on same server
2. **Distributed** - Backend on one server, qBittorrent on another (e.g., dedicated seedbox)

---

## Prerequisites

Before deploying Organizr:

- **Linux server** with systemd (Ubuntu 20.04+, Debian 11+, or similar)
- **qBittorrent instance** with Web UI enabled
  - Can be local (same server) or remote
  - Web UI must be accessible from backend server
  - Credentials required (username/password)
- **Storage** for audiobook files (local disk or network mount)
- **Go 1.23+** (for building backend) - only needed on build machine, not production server
- **Node.js 20+** (for building frontend) - only needed on build machine

---

## Backend Deployment

### 1. Build Backend Binary

On your build machine (can be different from production server):

```bash
# Clone repository
git clone https://github.com/yourusername/organizr.git
cd organizr/backend

# Build binary (statically linked with SQLite)
make build
# Creates binary at backend/bin/organizr
```

**Build output:**
- Binary: `backend/bin/organizr` (~15-20 MB)
- No additional files needed (embedded SQLite)

### 2. Copy Binary to Server

```bash
# Create deployment directory on server
ssh user@server "sudo mkdir -p /opt/organizr && sudo chown $USER:$USER /opt/organizr"

# Copy binary to server
scp backend/bin/organizr user@server:/opt/organizr/

# Make binary executable
ssh user@server "chmod +x /opt/organizr/organizr"
```

### 3. Create systemd Service

Create `/etc/systemd/system/organizr.service`:

```ini
[Unit]
Description=Organizr Audiobook Automation
After=network.target

[Service]
Type=simple
User=organizr
Group=organizr
WorkingDirectory=/opt/organizr
ExecStart=/opt/organizr/organizr
Restart=on-failure
RestartSec=10s

# Resource limits (adjust based on your needs)
MemoryLimit=512M
CPUQuota=100%

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=organizr

[Install]
WantedBy=multi-user.target
```

**Create dedicated user (recommended for security):**

```bash
# Create organizr user (no login shell)
sudo useradd -r -s /bin/false organizr

# Grant ownership of deployment directory
sudo chown -R organizr:organizr /opt/organizr
```

**Enable and start service:**

```bash
# Reload systemd to recognize new service
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable organizr

# Start service
sudo systemctl start organizr

# Check status
sudo systemctl status organizr
```

**View logs:**

```bash
# Follow logs in real-time
sudo journalctl -u organizr -f

# View recent logs
sudo journalctl -u organizr -n 100
```

---

## Frontend Deployment

You have two options for deploying the frontend:

### Option A: Static File Server (nginx) - Recommended

This approach uses nginx to serve frontend files and proxy API requests to the backend.

**1. Build Frontend:**

```bash
cd frontend

# Install dependencies
npm install

# Build for production
npm run build
# Creates frontend/dist/ directory with static files
```

**2. Copy Files to Server:**

```bash
# Create web directory
ssh user@server "sudo mkdir -p /var/www/organizr && sudo chown www-data:www-data /var/www/organizr"

# Copy built files to server
scp -r frontend/dist/* user@server:/var/www/organizr/
```

**3. Configure nginx:**

Create `/etc/nginx/sites-available/organizr`:

```nginx
server {
    listen 80;
    server_name organizr.example.com;

    # Frontend files
    root /var/www/organizr;
    index index.html;

    # Serve frontend files
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Proxy API requests to Go backend
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (for future real-time updates)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Timeout settings (adjust for large file operations)
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

**4. Enable nginx site:**

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/organizr /etc/nginx/sites-enabled/

# Test nginx configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

**5. (Optional) Enable HTTPS with Let's Encrypt:**

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d organizr.example.com

# Certbot automatically updates nginx config and sets up auto-renewal
```

### Option B: Backend Static Serving - Future Enhancement

**Note:** Backend does not currently serve static files. This feature could be added in a future version by:

1. Adding static file middleware to Chi router
2. Embedding frontend dist/ directory into backend binary
3. Serving frontend from root path and API from /api/

For now, use Option A (nginx) for production deployments.

---

## Database Setup

SQLite database is created automatically on first run.

**Database Location:**
- Default: `./organizr.db` in backend working directory
- For systemd: `/opt/organizr/organizr.db`

**Initialization:**
- Database file created on first startup
- Migrations run automatically
- No manual setup required

**Backup Procedure:**

```bash
# Stop backend service
sudo systemctl stop organizr

# Backup database file
sudo cp /opt/organizr/organizr.db /opt/organizr/backup/organizr-$(date +%Y%m%d).db

# Alternative: SQLite backup command (works while running)
sqlite3 /opt/organizr/organizr.db ".backup /opt/organizr/backup/organizr-$(date +%Y%m%d).db"

# Restart service
sudo systemctl start organizr
```

**Automated Backups:**

Add cron job to backup database daily:

```bash
# Edit crontab
sudo crontab -e

# Add daily backup at 3 AM
0 3 * * * sqlite3 /opt/organizr/organizr.db ".backup /opt/organizr/backup/organizr-$(date +\%Y\%m\%d).db" && find /opt/organizr/backup -name "organizr-*.db" -mtime +30 -delete
```

**Database Maintenance:**

```bash
# Vacuum database (reclaim space, optimize)
sqlite3 /opt/organizr/organizr.db "VACUUM;"

# Analyze database (update query planner statistics)
sqlite3 /opt/organizr/organizr.db "ANALYZE;"
```

---

## Configuration

After deployment, configure Organizr through the web interface:

**1. Access UI:**
- Open browser to `http://organizr.example.com` (or `http://server-ip:8080` if no nginx)

**2. Configure qBittorrent Connection:**
- Navigate to **Config** page
- Enter qBittorrent Web UI URL (e.g., `http://localhost:8080` or `http://seedbox.example.com:8080`)
- Enter qBittorrent username and password
- Click "Test Connection" to verify
- Save configuration

**3. Configure Audiobook Organization:**
- Set **Destination Path** - where organized audiobooks should be placed
  - Example: `/media/audiobooks` or `/mnt/nas/audiobooks`
  - Must be writable by organizr user
- Configure **Path Template** - folder structure for organized files
  - Default: `{author}/{series}/{title}`
  - Variables: `{author}`, `{series}`, `{title}`, `{series_number}`
  - Example: `{author}/{series}/{series_number} - {title}` → "Terry Pratchett/Discworld/01 - The Colour of Magic"
- Preview template with sample data
- Save configuration

**4. (Optional) Path Mapping for Remote qBittorrent:**

If qBittorrent runs on a different server with different paths:
- Set **Path Mapping** - maps qBittorrent's path to local path
  - Example: qBittorrent downloads to `/downloads` but files appear at `/mnt/seedbox/downloads` on organizr server
  - Path Mapping: `/mnt/seedbox` (prepended to qBittorrent paths)

See `docs/CONFIGURATION.md` for detailed configuration options.

---

## Environment Variables

**Frontend Build Time:**
- `VITE_API_URL` - Backend API URL (default: `/api` for same-origin deployment)
  - Set in `frontend/.env.production` before build
  - Example: `VITE_API_URL=https://api.organizr.example.com` for separate backend server

**Backend Runtime:**
- Currently no environment variables (uses database configuration)
- Future: May add `PORT`, `DATABASE_PATH`, `LOG_LEVEL` env vars

---

## Health Checks

**Backend Health:**

```bash
# HTTP health check (if endpoint added in future)
curl http://localhost:8080/api/health

# Check if backend is listening
ss -tlnp | grep 8080

# Check systemd status
sudo systemctl status organizr
```

**Frontend Health:**

```bash
# Check if nginx is serving frontend
curl -I http://organizr.example.com

# Should return 200 OK with index.html
```

**qBittorrent Connection:**
- Use "Test Connection" button in Organizr Config page
- Verifies connectivity, authentication, and API access

**Overall System Health:**

```bash
# Backend logs (look for errors)
sudo journalctl -u organizr -n 50

# nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# qBittorrent logs
# (location depends on qBittorrent setup)
```

---

## Troubleshooting

Common deployment issues and solutions:

**Backend won't start:**
- Check logs: `sudo journalctl -u organizr -n 100`
- Verify port 8080 not in use: `ss -tlnp | grep 8080`
- Check file permissions: `ls -l /opt/organizr/organizr`
- Verify binary is executable: `chmod +x /opt/organizr/organizr`

**Frontend can't reach backend (CORS errors):**
- If using nginx: verify `proxy_pass` configuration is correct
- Check nginx error logs: `sudo tail -f /var/log/nginx/error.log`
- Verify backend is running: `curl http://localhost:8080/api/downloads`

**qBittorrent connection fails:**
- Verify qBittorrent Web UI is accessible: `curl http://qbittorrent-host:port/api/v2/app/version`
- Check firewall rules if qBittorrent is remote
- Verify credentials are correct in Organizr config
- Ensure qBittorrent Web UI is enabled (Settings → Web UI → Enable)

**Database errors:**
- Check database file permissions: `ls -l /opt/organizr/organizr.db`
- Verify organizr user has write access to directory
- Check disk space: `df -h /opt/organizr`

**File organization fails:**
- Verify destination path exists and is writable
- Check disk space on destination: `df -h /destination/path`
- Review logs for specific error messages
- Verify path template is valid (use preview in UI)

See `docs/TROUBLESHOOTING.md` for comprehensive troubleshooting guide.

---

## Upgrading

**Backend Upgrade:**

```bash
# Build new version
cd backend
git pull
make build

# Stop service
sudo systemctl stop organizr

# Backup current binary and database
sudo cp /opt/organizr/organizr /opt/organizr/organizr.backup
sudo cp /opt/organizr/organizr.db /opt/organizr/organizr.db.backup

# Copy new binary
sudo cp bin/organizr /opt/organizr/

# Start service (migrations run automatically)
sudo systemctl start organizr

# Check logs for migration success
sudo journalctl -u organizr -n 50
```

**Frontend Upgrade:**

```bash
# Build new version
cd frontend
git pull
npm install
npm run build

# Copy new files
scp -r dist/* user@server:/var/www/organizr/

# No service restart needed (static files)
```

---

## Docker Deployment

Organizr runs in Docker containers with support for both development (named volumes) and production (host path mounts) deployments.

### Prerequisites

- Docker and Docker Compose installed
- qBittorrent with Web UI enabled (running on host or in container)
- MyAnonamouse account with API secret
- Storage paths for downloads and audiobooks

### Unraid Deployment

Unraid is a popular NAS operating system with Docker support. This guide covers deploying Organizr to Unraid.

#### Step 1: Prepare Directories

Typical Unraid storage paths:
- `/mnt/user/downloads` - qBittorrent download location
- `/mnt/user/audiobooks` - Audiobookshelf library
- `/mnt/user/appdata/organizr` - Organizr database

Ensure these directories exist and have proper permissions:

```bash
# Create directories if needed
mkdir -p /mnt/user/downloads
mkdir -p /mnt/user/audiobooks
mkdir -p /mnt/user/appdata/organizr

# Set ownership (uid/gid 1001 matches container user)
chown -R 1001:1001 /mnt/user/appdata/organizr
```

#### Step 2: Clone Repository

```bash
# Navigate to appdata
cd /mnt/user/appdata

# Clone repository
git clone https://github.com/yourusername/organizr.git
cd organizr
```

#### Step 3: Configure docker-compose.yml

Replace named volumes with host paths. Edit `docker-compose.yml`:

```yaml
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    env_file:
      - .env
    environment:
      - ORGANIZR_DB_PATH=${ORGANIZR_DB_PATH:-/data/organizr.db}
      - QBITTORRENT_URL=${QBITTORRENT_URL}
      - QBITTORRENT_USERNAME=${QBITTORRENT_USERNAME}
      - QBITTORRENT_PASSWORD=${QBITTORRENT_PASSWORD}
      - PATHS_DESTINATION=${PATHS_DESTINATION:-/audiobooks}
      - PATHS_TEMPLATE=${PATHS_TEMPLATE:-{author}/{series}/{title}}
      - PATHS_NO_SERIES_TEMPLATE=${PATHS_NO_SERIES_TEMPLATE:-{author}/{title}}
      - PATHS_OPERATION=${PATHS_OPERATION:-copy}
      - PATHS_LOCAL_MOUNT=${PATHS_LOCAL_MOUNT:-/downloads}
      - MONITOR_INTERVAL_SECONDS=${MONITOR_INTERVAL_SECONDS:-30}
      - MONITOR_AUTO_ORGANIZE=${MONITOR_AUTO_ORGANIZE:-true}
      - MAM_BASEURL=${MAM_BASEURL:-https://www.myanonamouse.net}
      - MAM_SECRET=${MAM_SECRET}
    volumes:
      - /mnt/user/downloads:/downloads              # qBittorrent downloads
      - /mnt/user/audiobooks:/audiobooks            # Organized audiobooks
      - /mnt/user/appdata/organizr:/data            # Database
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 5s
    networks:
      - organizr-network

  frontend:
    build: ./frontend
    ports:
      - "8081:8080"
    environment:
      - VITE_API_URL=http://backend:8080
    depends_on:
      backend:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 5s
    networks:
      - organizr-network

networks:
  organizr-network:
    driver: bridge
```

**Note:** Remove the `volumes:` section at the bottom (no longer needed with host path mounts).

#### Step 4: Configure Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
nano .env  # or your preferred editor
```

**Required configuration:**

```bash
# qBittorrent Connection
# If qBittorrent in container on same network: http://container-name:8080
# If qBittorrent on Unraid host: http://192.168.1.100:8080 (use Unraid IP)
QBITTORRENT_URL=http://192.168.1.100:8080

# qBittorrent credentials
QBITTORRENT_USERNAME=admin
QBITTORRENT_PASSWORD=your_password

# MyAnonamouse API secret
MAM_SECRET=your_mam_api_secret_here

# Path configuration (container paths - match volume mounts)
PATHS_LOCAL_MOUNT=/downloads
PATHS_DESTINATION=/audiobooks

# Path templates (customize folder structure)
PATHS_TEMPLATE={author}/{series}/{title}
PATHS_NO_SERIES_TEMPLATE={author}/{title}

# Operation mode
PATHS_OPERATION=copy
```

**Path Configuration Notes:**
- `PATHS_LOCAL_MOUNT=/downloads` - Container path to qBittorrent downloads (matches volume mount)
- `PATHS_DESTINATION=/audiobooks` - Container path for organized files (matches volume mount)
- Host paths are configured in docker-compose.yml volumes section

**qBittorrent URL Options:**
- **qBittorrent in Docker:** Use container name (e.g., `http://qbittorrent:8080`)
- **qBittorrent on Unraid host:** Use Unraid IP (e.g., `http://192.168.1.100:8080`)
- **qBittorrent on different server:** Use that server's IP

#### Step 5: Start Services

```bash
# Build and start containers
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

**Expected output:**
```
Creating organizr_backend_1  ... done
Creating organizr_frontend_1 ... done
```

#### Step 6: Verify Configuration

**Access UI:**
- Navigate to `http://unraid-ip:8081` in your browser
- You should see the Organizr interface

**Test qBittorrent Connection:**
1. Click **Config** in sidebar
2. Verify qBittorrent connection settings
3. Click **Test Connection** button
4. Should see "Connection successful" message

**Test Download and Organization:**
1. Search for an audiobook
2. Download a torrent
3. Wait for download to complete
4. Verify files organized to `/mnt/user/audiobooks`
5. Check folder structure matches your template

#### Step 7: Configure Audiobookshelf (Optional)

If using Audiobookshelf on Unraid:

1. Install Audiobookshelf from Community Applications
2. Map library to `/mnt/user/audiobooks`
3. Scan library to detect organized audiobooks
4. Verify audiobooks appear with correct metadata

### Docker Compose Deployment (Non-Unraid)

For standard Docker Compose deployments:

**1. Clone repository:**
```bash
git clone https://github.com/yourusername/organizr.git
cd organizr
```

**2. Configure docker-compose.yml:**

Adjust host paths in volumes section to match your system:

```yaml
volumes:
  - /path/to/downloads:/downloads          # Your qBittorrent download path
  - /path/to/audiobooks:/audiobooks        # Your audiobook library path
  - /path/to/data:/data                    # Database persistence
```

**3. Configure .env:**

```bash
cp .env.example .env
# Edit .env with your settings
```

**4. Start services:**

```bash
docker-compose up -d
```

**5. Access UI:**
- Frontend: `http://localhost:8081`
- Backend API: `http://localhost:8080`

### Troubleshooting Docker Deployment

#### qBittorrent Connection Failed

**Symptom:** "Connection failed" when testing qBittorrent in Config

**Solutions:**
- **Check QBITTORRENT_URL:** Must be accessible from container
  - Test from container: `docker-compose exec backend wget -O- http://your-qbittorrent-url`
- **Verify Web UI enabled:** qBittorrent Settings → Web UI → Enable Web UI
- **Check credentials:** Ensure username/password correct in .env
- **Network issues:** If qBittorrent on host, use host IP not `localhost`
  - Linux: Use host IP (e.g., `192.168.1.100:8080`)
  - Mac/Windows: Use `host.docker.internal:8080`

#### Files Not Organizing

**Symptom:** Downloads complete but files don't move to audiobooks directory

**Solutions:**
- **Verify PATHS_LOCAL_MOUNT:** Must match qBittorrent download volume mount
  - Check docker-compose.yml: If mounted as `/downloads`, set `PATHS_LOCAL_MOUNT=/downloads`
- **Check logs:** `docker-compose logs backend | grep -i error`
- **Verify download completion:** qBittorrent shows "Completed" status, not paused/seeding
- **Check permissions:** Host directories must be readable/writable by uid 1001
  ```bash
  ls -la /mnt/user/downloads
  ls -la /mnt/user/audiobooks
  ```

#### Permission Errors

**Symptom:** "Permission denied" errors in logs

**Solutions:**
- **Container user:** Both containers run as uid/gid 1001 (non-root)
- **Host permissions:** Ensure directories writable by uid 1001
  ```bash
  chown -R 1001:1001 /mnt/user/appdata/organizr
  # Or grant world write (less secure)
  chmod 777 /mnt/user/downloads /mnt/user/audiobooks
  ```
- **Unraid:** Set "Nobody" user permissions on shares (Shares → Edit Share → Security)

#### Database Issues

**Symptom:** Backend won't start, database errors in logs

**Solutions:**
- **Database location:** Verify mounted at `/data` in container
  ```bash
  docker-compose exec backend ls -la /data
  ```
- **Permissions:** Database directory must be writable
  ```bash
  ls -la /mnt/user/appdata/organizr
  chown -R 1001:1001 /mnt/user/appdata/organizr
  ```
- **Backup and reset:** Stop containers, backup .db files, delete database, restart
  ```bash
  docker-compose down
  cp /mnt/user/appdata/organizr/organizr.db /mnt/user/appdata/organizr/organizr.db.backup
  rm /mnt/user/appdata/organizr/organizr.db*
  docker-compose up -d
  ```

#### Container Won't Start

**Symptom:** Container exits immediately or fails health check

**Solutions:**
- **Check logs:** `docker-compose logs backend` or `docker-compose logs frontend`
- **Port conflicts:** Ensure ports 8080/8081 not in use
  ```bash
  netstat -tlnp | grep -E '8080|8081'
  ```
- **Missing .env:** Ensure .env file exists with required variables
- **Build issues:** Rebuild containers
  ```bash
  docker-compose down
  docker-compose build --no-cache
  docker-compose up -d
  ```

### Upgrading Docker Deployment

**1. Stop containers:**
```bash
docker-compose down
```

**2. Backup database:**
```bash
cp /mnt/user/appdata/organizr/organizr.db /mnt/user/appdata/organizr/organizr.db.backup
```

**3. Pull latest code:**
```bash
git pull origin main
```

**4. Rebuild and start:**
```bash
docker-compose build --no-cache
docker-compose up -d
```

**5. Verify:**
```bash
docker-compose logs -f
# Check for migration messages and startup success
```

### Alternative: Unraid Community Applications Template

**Note:** Unraid Community Applications template is a future enhancement.

**Planned features:**
- UI-based installation from Community Applications
- Pre-configured volume mappings
- Web UI configuration fields
- One-click updates

For now, use docker-compose deployment method above.

---

## Security Considerations

**Important:** Organizr v1.2 does not include authentication. Deploy only in trusted environments.

**Recommended Security Practices:**

1. **Network Security:**
   - Deploy behind VPN or only expose to trusted network
   - Use firewall rules to restrict access
   - Consider nginx basic auth as temporary measure

2. **File Permissions:**
   - Run backend as dedicated user (not root)
   - Restrict database file permissions: `chmod 600 organizr.db`
   - Limit destination directory write access

3. **HTTPS:**
   - Always use HTTPS in production (Let's Encrypt recommended)
   - Prevents credential interception for qBittorrent auth

4. **Future Work:**
   - Authentication/authorization planned for future release
   - Multi-user support requires authentication layer

---

## Performance Tuning

**Backend:**
- Adjust systemd resource limits based on workload
- Monitor memory usage with `systemctl status organizr`
- Consider increasing Go's `GOMAXPROCS` for multi-core systems

**Database:**
- SQLite WAL mode enabled by default (good concurrency)
- Run `VACUUM` periodically to reclaim space
- Consider moving to PostgreSQL if scaling beyond single-user

**Frontend:**
- nginx gzip compression enabled by default
- Consider CDN for static assets if serving many users
- Browser caching configured in nginx

---

## Monitoring

**Systemd Journal:**
```bash
# Real-time logs
journalctl -u organizr -f

# Last 100 lines
journalctl -u organizr -n 100

# Logs since specific time
journalctl -u organizr --since "1 hour ago"
```

**Process Monitoring:**
```bash
# CPU/memory usage
top -p $(pgrep organizr)

# Detailed stats
systemctl status organizr
```

**Disk Usage:**
```bash
# Database size
du -h /opt/organizr/organizr.db

# Destination directory size
du -sh /destination/path
```

---

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/yourusername/organizr/issues
- Documentation: See `docs/` directory
- Troubleshooting: `docs/TROUBLESHOOTING.md`
- Configuration: `docs/CONFIGURATION.md`

---

*Last updated: 2026-01-09*
*Version: 1.2*
