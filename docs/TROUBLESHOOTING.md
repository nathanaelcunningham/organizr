# Troubleshooting Guide

Common issues and solutions for Organizr.

---

## qBittorrent Connection Issues

### Symptom

- "Failed to connect to qBittorrent" errors in the UI
- Connection test fails
- Downloads won't submit to qBittorrent

### Solutions

1. **Verify qBittorrent Web UI is enabled:**
   - Open qBittorrent
   - Go to Tools → Options → Web UI
   - Ensure "Enable the Web User Interface" is checked
   - Note the port (default: 8080)

2. **Check URL format:**
   - Correct: `http://192.168.1.100:8080`
   - Incorrect: `192.168.1.100:8080` (missing http://)
   - Incorrect: `https://192.168.1.100:8080` (qBittorrent Web UI uses http by default)

3. **Verify credentials:**
   - Use the same username/password configured in qBittorrent Web UI settings
   - Test credentials by opening qBittorrent Web UI in a browser

4. **Use connection test button:**
   - In Organizr Settings, click "Test Connection"
   - Review error message for specific issues (auth failure, network unreachable, etc.)

5. **Check firewall and network access:**
   - Ensure port 8080 (or your configured port) is not blocked
   - If qBittorrent is on a different machine, verify network connectivity: `ping <qbittorrent-host>`
   - Try accessing qBittorrent Web UI directly: `http://<qbittorrent-host>:8080`

6. **Verify qBittorrent is running:**
   - Check if qBittorrent application is running
   - Restart qBittorrent and try again

### Verification

- Connection test in Settings shows "Connected successfully"
- Can submit a test download and see it appear in qBittorrent

### Related Documentation

- [Configuration Guide](../backend/docs/CONFIGURATION.md) - qBittorrent connection setup
- [API Documentation](../backend/docs/API.md) - Connection test endpoint

---

## Download Not Starting

### Symptom

- Download stuck in "queued" status
- Torrent doesn't appear in qBittorrent
- No progress after several minutes

### Solutions

1. **Check qBittorrent logs:**
   - Open qBittorrent
   - View → Execution Log
   - Look for errors related to the torrent

2. **Verify torrent URL or magnet link is valid:**
   - Test the torrent link manually by pasting into qBittorrent
   - Ensure URL is accessible (not behind authentication you haven't configured)
   - For MAM torrents, Organizr automatically handles authentication

3. **Ensure qBittorrent has disk space:**
   - Check available space in qBittorrent's download directory
   - qBittorrent will queue torrents if disk space is insufficient

4. **Check category exists in qBittorrent:**
   - If using categories, ensure the category is created in qBittorrent
   - Go to qBittorrent → Categories → Right-click → Add category

5. **Verify qBittorrent queue settings:**
   - Check if qBittorrent has queue limits configured
   - Tools → Options → Downloads → Queue
   - Ensure max active downloads allows new torrents

6. **Check Organizr monitor is running:**
   - Look for "Monitor started" in backend logs
   - Restart backend if monitor isn't running

### Verification

- Download status changes to "downloading" or shows progress percentage
- Torrent appears in qBittorrent with active status

### Related Documentation

- [API Documentation](../backend/docs/API.md) - Download endpoints and status codes

---

## Organization Failures

### Symptom

- Files download to 100% but don't organize
- Download stuck in "completed" status without organizing
- Organized files not appearing in destination directory

### Solutions

1. **Check destination path exists and is writable:**
   ```bash
   # Test if path exists
   ls -la /your/destination/path

   # Test write permissions
   touch /your/destination/path/test.txt
   rm /your/destination/path/test.txt
   ```

2. **Verify path template is valid:**
   - Go to Settings in UI
   - Use the template preview to see the generated path
   - Ensure no invalid characters for your filesystem
   - Templates should use `{author}`, `{title}`, `{series}`, `{series_number}` variables

3. **Ensure qBittorrent finished downloading:**
   - Check qBittorrent shows 100% complete
   - Verify torrent status is "seeding" not "downloading"
   - Organization only triggers after qBittorrent reports completion

4. **Check monitor is running:**
   - Look for "Monitor started" in backend logs
   - Backend logs should show periodic progress checks
   - Restart backend if monitor stopped

5. **Verify path mapping for Docker/remote qBittorrent:**
   - If qBittorrent sees `/downloads/file.m4b` but Organizr sees `/mnt/qbittorrent/downloads/file.m4b`:
   - Configure `paths.mount_point` to `/mnt/qbittorrent` in Settings
   - This prepends the mount point to qBittorrent's reported paths
   - See [Configuration Guide](../backend/docs/CONFIGURATION.md) for details

6. **Check backend logs for organization errors:**
   - Look for copy/move errors in backend output
   - Common issues: permission denied, disk full, path too long

7. **Verify download metadata is complete:**
   - Check download details in UI shows author, title, series populated
   - Missing metadata means template variables will be empty
   - Re-create download with correct metadata if needed

### Verification

- Files appear in destination directory with correct folder structure
- Download status changes to "organized"
- Backend logs show "Organization complete" message

### Related Documentation

- [Configuration Guide](../backend/docs/CONFIGURATION.md) - Path templates and mount points
- [Architecture Decision Record](architecture/ADR.md) - Organization strategy and path mapping

---

## Template Variables Not Working

### Symptom

- Folders created with `{author}` literal instead of actual author name
- Template preview shows variable names instead of values
- Organized files have incomplete paths

### Solutions

1. **Verify metadata is populated:**
   - Check download details in UI
   - Ensure author, title, series fields are filled
   - Metadata comes from MAM search results when creating download

2. **Ensure template syntax is correct:**
   - Correct: `{author}/{series}/{title}`
   - Incorrect: `{ author }/{ series }` (spaces inside braces)
   - Incorrect: `{{author}}` (double braces)

3. **Check for empty metadata fields:**
   - If series is empty, `{series}` becomes empty string
   - Use "No Series Template" for books without series
   - Default no-series template: `{author}/{title}`

4. **Test template with preview:**
   - In Settings, use template preview to see actual output
   - Preview shows example path using sample metadata
   - Verify preview looks correct before saving

5. **Verify variable names are valid:**
   - Supported: `{author}`, `{title}`, `{series}`, `{series_number}`
   - Case-sensitive (must be lowercase)
   - No custom variables (only these four are supported)

### Verification

- Template preview shows real author/title/series values, not variable names
- Organized files have folders named with actual metadata

### Related Documentation

- [Configuration Guide](../backend/docs/CONFIGURATION.md) - Complete template variable reference

---

## Frontend Won't Connect to Backend

### Symptom

- Frontend shows connection errors
- "Network Error" messages in UI
- Search and download features don't work

### Solutions

1. **Verify backend is running:**
   ```bash
   # Check if backend process is running
   ps aux | grep organizr

   # Check if port 8080 is listening
   lsof -i :8080
   ```

2. **Check VITE_API_URL in frontend/.env.development:**
   ```bash
   cat frontend/.env.development
   # Should show: VITE_API_URL=http://localhost:8080
   ```
   - If backend is on different host, update to `http://<backend-host>:8080`
   - Restart frontend dev server after changing

3. **Verify CORS is not blocking requests:**
   - Backend enables CORS by default for all origins
   - Check browser console (F12) for CORS errors
   - Look for "Access-Control-Allow-Origin" in error messages

4. **Check browser console for specific errors:**
   - Open browser DevTools (F12)
   - Go to Console tab
   - Look for network errors with details
   - Common: "Failed to fetch", "ERR_CONNECTION_REFUSED", "net::ERR_EMPTY_RESPONSE"

5. **Test backend directly:**
   ```bash
   # Test health endpoint
   curl http://localhost:8080/api/health
   # Should return: {"status":"ok"}
   ```

6. **Restart both backend and frontend:**
   ```bash
   # Backend
   cd backend
   make run

   # Frontend (in new terminal)
   cd frontend
   npm run dev
   ```

### Verification

- Frontend loads without connection errors
- Search bar accepts input and returns results
- Settings page shows current configuration

### Related Documentation

- [CONTRIBUTING.md](../CONTRIBUTING.md) - Development setup
- [API Documentation](../backend/docs/API.md) - Backend endpoints

---

## MAM Search Not Working

### Symptom

- Search returns no results
- Search shows error message
- Search takes very long and times out

### Solutions

1. **Verify MAM API key is configured (if required in future):**
   - Currently Organizr searches MAM without authentication
   - Future versions may require API key configuration

2. **Check network connectivity to MyAnonamouse:**
   ```bash
   # Test connectivity
   curl -I https://www.myanonamouse.net
   # Should return 200 OK or similar
   ```

3. **Ensure search query format is valid:**
   - Avoid special characters that might break search
   - Try simpler queries if complex ones fail
   - Example valid queries: "stephen king", "discworld", "sanderson"

4. **Check backend logs for search errors:**
   - Backend logs show search provider errors
   - Look for HTTP errors, timeouts, or parsing failures

5. **Test search endpoint directly:**
   ```bash
   # Test search API
   curl "http://localhost:8080/api/search?q=test"
   # Should return JSON with results array
   ```

6. **Verify MAM website is accessible:**
   - Open https://www.myanonamouse.net in browser
   - Ensure site is not down or blocked in your region
   - Check MAM status on social media if site seems down

### Verification

- Search returns results within a few seconds
- Results include title, author, series metadata
- Can click results to create downloads

### Related Documentation

- [API Documentation](../backend/docs/API.md) - Search endpoints

---

## Tests Failing

### Symptom

- `make test` or `npm test` shows failures
- CI pipeline failing on tests
- Race conditions detected

### Solutions

1. **Run tests with race detection (Go):**
   ```bash
   cd backend
   go test -race ./...
   ```
   - Race detector finds concurrency issues
   - Fix any race conditions before proceeding

2. **Ensure database is not in use by running app:**
   - Stop backend server before running tests
   - Tests use in-memory database or separate test database
   - Conflicts occur if app locks database file

3. **Check test database has write permissions:**
   ```bash
   # Ensure test directory is writable
   ls -la backend/
   # Remove stale test databases
   rm backend/*.db-test 2>/dev/null
   ```

4. **Run frontend tests with verbose output:**
   ```bash
   cd frontend
   npm test -- --reporter=verbose
   ```
   - Shows detailed test output
   - Helps identify which specific test is failing

5. **Check for missing dependencies:**
   ```bash
   # Backend
   cd backend
   go mod tidy
   go mod verify

   # Frontend
   cd frontend
   npm install
   ```

6. **Run tests with coverage to see what's tested:**
   ```bash
   # Backend
   cd backend
   make test-coverage

   # Frontend
   cd frontend
   npm test -- --coverage
   ```

7. **Verify test fixtures and helpers:**
   - Tests use testutil package for helpers and fixtures
   - Ensure test data is valid and realistic
   - Check frontend test/fixtures.ts for correct types

### Verification

- All tests pass without errors
- No race conditions detected with `-race` flag
- Coverage meets 60% threshold (baseline)

### Related Documentation

- [CONTRIBUTING.md](../CONTRIBUTING.md) - Testing requirements
- [Architecture Decision Record](architecture/ADR.md) - Testing strategy

---

## Still Having Issues?

If you've tried these solutions and still experiencing problems:

1. **Check backend logs** for detailed error messages
2. **Check browser console** (F12 → Console) for frontend errors
3. **Review configuration** in Settings page - verify all paths and URLs
4. **Test components individually** - backend health endpoint, qBittorrent Web UI, frontend dev server
5. **Create an issue** with:
   - Steps to reproduce
   - Error messages from logs
   - Your environment (OS, Go version, Node version, qBittorrent version)
   - Configuration (without sensitive credentials)

---

**Quick Reference:**

| Problem | Quick Check |
|---------|-------------|
| qBittorrent connection | Test Web UI in browser: `http://host:8080` |
| Downloads not starting | Check qBittorrent logs for errors |
| Organization not working | Verify destination path exists and is writable |
| Template variables broken | Use template preview in Settings |
| Frontend errors | Check browser console (F12) |
| Search not working | Test backend health: `curl http://localhost:8080/api/health` |
| Tests failing | Run with race detection: `go test -race ./...` |
