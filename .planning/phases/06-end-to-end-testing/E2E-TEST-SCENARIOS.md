# End-to-End Test Scenarios

This document provides comprehensive manual testing scenarios for validating the complete audiobook download automation workflow. These tests should be executed before production deployment to ensure all components work together correctly.

## Prerequisites

Before executing these test scenarios, ensure the following components are properly configured and running:

### 1. qBittorrent Setup
- qBittorrent installed and running
- Web UI enabled (Tools > Options > Web UI)
- Authentication configured with username and password
- Port accessible (default: 8080)
- Verify access: Navigate to `http://localhost:8080` and confirm login works

### 2. Backend Service
```bash
cd backend
go run cmd/api/main.go
```
- Server should start on port 8080 (or configured port)
- Database initialized (SQLite at `backend/organizr.db`)
- Configuration loaded successfully

### 3. Frontend Application
```bash
cd frontend
npm run dev
```
- Development server running (typically on port 5173)
- Can access frontend in browser

### 4. Test Environment
- **Test destination folder**: Create a writable directory for organized files (e.g., `/tmp/audiobooks-test`)
- **Test magnet links**: Have 1-2 small, legal torrent magnet links ready (avoid large files for faster testing)
- **MAM credentials**: If testing MAM downloads, ensure valid MyAnonamouse credentials configured

### 5. Configuration Verification
Navigate to the configuration page and verify:
- qBittorrent URL is correct (e.g., `http://localhost:8080`)
- qBittorrent credentials are set
- Destination path points to test folder
- Organization template is configured (e.g., `{{Author}}/{{Series}}/{{Book}}`)
- Auto-organization is enabled (default)

---

## Scenario 1: Happy Path - Magnet Link Download

**Objective**: Verify the complete end-to-end workflow from magnet link submission to organized files.

### Steps

1. **Navigate to frontend**
   - Open browser to frontend URL (e.g., `http://localhost:5173`)
   - Verify page loads without errors

2. **Submit magnet link**
   - Paste a magnet link into the "Download URL" input field
   - Click "Add Download" button
   - **Expected**: Success notification appears
   - **Expected**: Download appears in the downloads list immediately

3. **Verify initial state**
   - **Expected**: Download shows "queued" status
   - **Expected**: Download card displays:
     - Title/name of the download
     - Progress bar at 0%
     - Status indicator
     - Action buttons (Cancel, etc.)

4. **Monitor status transition to downloading**
   - Wait approximately 3-10 seconds (monitor polling interval)
   - **Expected**: Status automatically transitions to "downloading"
   - **Expected**: No page refresh required (real-time update via polling)

5. **Monitor progress updates**
   - Observe the download card while content downloads
   - **Expected**: Progress bar updates incrementally (e.g., 0% → 25% → 50% → 75%)
   - **Expected**: Progress percentage displayed numerically
   - **Expected**: Updates occur every ~3 seconds during active download

6. **Monitor completion**
   - Wait for download to reach 100%
   - **Expected**: Status transitions to "completed"
   - **Expected**: Progress bar shows 100%

7. **Verify auto-organization triggers**
   - Within 3-10 seconds after completion:
   - **Expected**: Status automatically transitions to "organizing"
   - **Expected**: Progress indicator shows organization in progress

8. **Verify organized state**
   - After organization completes (typically 5-30 seconds depending on file size):
   - **Expected**: Status transitions to "organized"
   - **Expected**: `organized_path` field displays the final destination path
   - **Expected**: Path follows the configured template structure

9. **Verify files in destination**
   - Navigate to the destination folder in file explorer
   - **Expected**: Files exist at the path shown in `organized_path`
   - **Expected**: Folder structure matches template (e.g., `Author/Series/Book/`)
   - **Expected**: All files from the torrent are present
   - **Expected**: File permissions are correct (readable)

10. **Test clipboard functionality**
    - Click "Copy Path" button on the organized download card
    - **Expected**: Success notification appears
    - **Expected**: Clipboard contains the full path to organized files
    - Paste into terminal/file explorer to verify

### Common Failure Modes

- **Download stuck in "queued"**: qBittorrent may be unreachable or credentials incorrect
- **No status updates**: Polling may have stopped; check browser console for errors
- **Organization fails**: Check destination path permissions, disk space, or template syntax
- **Files missing after organization**: Check qBittorrent still has the torrent; files may have been moved/deleted

---

## Scenario 2: MAM Torrent File Download

**Objective**: Verify that MyAnonamouse torrent URLs are properly downloaded as .torrent files before submission to qBittorrent.

### Steps

1. **Navigate to frontend**
   - Ensure MAM credentials are configured in settings

2. **Submit MAM torrent URL**
   - Paste a MAM torrent download URL: `https://www.myanonamouse.net/tor/download.php?tid=XXXXXX`
   - Click "Add Download"
   - **Expected**: Success notification appears

3. **Verify backend downloads .torrent file**
   - Check backend logs for torrent download activity
   - **Expected**: Log shows torrent file downloaded before submission
   - **Expected**: Download appears in list with correct title

4. **Verify normal workflow**
   - **Expected**: Download proceeds through states: queued → downloading → completed
   - **Expected**: Auto-organization triggers as normal
   - **Expected**: Files organized correctly

5. **Verify in qBittorrent**
   - Open qBittorrent Web UI
   - **Expected**: Torrent appears in list
   - **Expected**: Torrent name matches the MAM entry

### Common Failure Modes

- **401 Authentication error**: MAM credentials invalid or expired
- **404 Not found**: Invalid torrent ID or torrent removed from MAM
- **Rate limiting**: MAM may rate-limit requests; wait and retry
- **Cookie issues**: MAM session cookies may need refresh

---

## Scenario 3: Manual Organization (Auto-Org Disabled)

**Objective**: Verify that disabling auto-organization allows manual control of organization timing.

### Steps

1. **Disable auto-organization**
   - Navigate to Configuration page
   - Uncheck "Automatically organize completed downloads"
   - Click "Save Configuration"
   - **Expected**: Success notification appears
   - **Expected**: Configuration persists (refresh page to verify)

2. **Add download**
   - Submit a magnet link or torrent URL
   - Wait for download to complete

3. **Verify download stops at "completed"**
   - **Expected**: Status transitions from downloading → completed
   - **Expected**: Status remains "completed" (does NOT auto-transition to organizing)
   - **Expected**: "Organize Now" button appears on the download card

4. **Manually trigger organization**
   - Click "Organize Now" button
   - **Expected**: Status immediately transitions to "organizing"
   - **Expected**: Organization proceeds normally

5. **Verify organized state**
   - **Expected**: Status transitions to "organized"
   - **Expected**: Files appear in destination with correct structure
   - **Expected**: `organized_path` displayed correctly

6. **Re-enable auto-organization (cleanup)**
   - Navigate to Configuration page
   - Check "Automatically organize completed downloads"
   - Save configuration

### Common Failure Modes

- **Auto-organization still triggers**: Config change may not have persisted; verify config in database or backend logs
- **"Organize Now" button missing**: Frontend may not have refreshed; reload page
- **Organization fails**: Same failure modes as Scenario 1

---

## Scenario 4: Organization Retry After Failure

**Objective**: Verify that failed organizations can be retried after fixing the underlying issue.

### Steps

1. **Create organization failure condition**
   - Configure an invalid destination path (e.g., `/nonexistent/path/audiobooks`)
   - OR remove write permissions on destination folder: `chmod 000 /path/to/destination`
   - Save configuration

2. **Add and complete download**
   - Submit magnet link
   - Wait for download to complete
   - Wait for auto-organization to trigger

3. **Verify organization failure**
   - **Expected**: Status transitions to "organizing" briefly
   - **Expected**: Status returns to "completed" or shows "organization_failed"
   - **Expected**: Error message displayed on download card or in notifications
   - **Expected**: Error indicates the specific problem (e.g., "permission denied", "directory does not exist")

4. **Fix the underlying issue**
   - Update configuration with valid, writable destination path
   - OR restore permissions: `chmod 755 /path/to/destination`
   - Save configuration

5. **Retry organization**
   - Click "Retry Organization" or "Organize Now" button
   - **Expected**: Status transitions to "organizing"
   - **Expected**: Organization proceeds without error

6. **Verify successful organization**
   - **Expected**: Status transitions to "organized"
   - **Expected**: Files appear in destination
   - **Expected**: No error messages
   - **Expected**: `organized_path` populated correctly

### Common Failure Modes

- **Same error repeats**: Fix may not have been applied; verify config saved and monitor logs
- **Partial organization**: Some files may have been copied before failure; check destination folder
- **Permission errors persist**: Verify user running backend has write access to destination

---

## Scenario 5: Download Cancellation

**Objective**: Verify that downloads can be cancelled at any stage and properly cleaned up.

### Steps

1. **Add download**
   - Submit magnet link
   - **Expected**: Download appears in list

2. **Cancel while downloading**
   - While download is in "downloading" state (not yet completed):
   - Click "Cancel" button on the download card
   - **Expected**: Confirmation dialog may appear (if implemented)
   - Confirm cancellation

3. **Verify removal from frontend**
   - **Expected**: Download immediately disappears from the downloads list
   - **Expected**: Success notification: "Download cancelled"

4. **Verify removal from qBittorrent**
   - Open qBittorrent Web UI
   - **Expected**: Torrent no longer appears in the list
   - **Expected**: Downloaded files removed from qBittorrent's download directory

5. **Verify database cleanup**
   - Refresh the frontend page
   - **Expected**: Cancelled download does not reappear
   - **Expected**: Download record removed from database

6. **Test cancellation at different stages**
   - Repeat with download in "queued" state
   - Repeat with download in "completed" state (before organization)
   - **Expected**: Cancellation works at all stages

### Common Failure Modes

- **Download reappears after refresh**: Backend may not have removed from database
- **Torrent remains in qBittorrent**: qBittorrent client call may have failed; check logs
- **Error notification**: qBittorrent may return error if torrent already removed
- **Organized files remain**: Cancellation after organization does not remove organized files (by design)

---

## Scenario 6: Connection Testing

**Objective**: Verify that the qBittorrent connection test provides accurate feedback.

### Steps

1. **Test successful connection**
   - Navigate to Configuration page
   - Ensure qBittorrent is running and configured correctly
   - Click "Test Connection" button
   - **Expected**: Success message appears: "Successfully connected to qBittorrent"
   - **Expected**: Message includes qBittorrent version info (if available)

2. **Test failed connection - qBittorrent stopped**
   - Stop qBittorrent service
   - Click "Test Connection" button
   - **Expected**: Error message appears
   - **Expected**: Error indicates connection refused or timeout
   - **Expected**: Message is user-friendly (not raw error dump)

3. **Test failed connection - wrong credentials**
   - Update qBittorrent credentials in config to incorrect values
   - Click "Test Connection" button
   - **Expected**: Error message appears
   - **Expected**: Error indicates authentication failure (e.g., "401 Unauthorized")

4. **Test failed connection - wrong URL**
   - Update qBittorrent URL to invalid address (e.g., `http://localhost:9999`)
   - Click "Test Connection" button
   - **Expected**: Error message appears
   - **Expected**: Error indicates connection failure

5. **Verify connection restoration**
   - Restore correct qBittorrent configuration
   - Restart qBittorrent if stopped
   - Click "Test Connection" button
   - **Expected**: Success message appears

### Common Failure Modes

- **False positive**: Connection test may succeed even if qBittorrent is not fully functional
- **Timeout too short**: Connection test may fail on slow networks
- **Certificate errors**: HTTPS connections may fail with self-signed certificates

---

## Scenario 7: Polling Behavior

**Objective**: Verify that polling starts/stops automatically based on download activity to conserve resources.

### Steps

1. **Verify initial state (no downloads)**
   - Open frontend with no active downloads
   - Open browser DevTools > Network tab
   - Filter for requests to `/api/downloads`
   - **Expected**: No polling requests occur (polling is stopped)

2. **Add first download**
   - Submit magnet link
   - Observe Network tab
   - **Expected**: Polling starts immediately after download creation
   - **Expected**: GET requests to `/api/downloads` appear every 3 seconds

3. **Add multiple downloads**
   - Submit 2-3 additional downloads while first is still active
   - **Expected**: Polling continues at 3-second interval
   - **Expected**: All downloads appear in list and update together

4. **Monitor during completion**
   - Wait for all downloads to complete and organize
   - Continue observing Network tab
   - **Expected**: Polling continues while any download is in "downloading" or "organizing" state

5. **Verify polling stops automatically**
   - Wait for all downloads to reach "organized" state
   - Continue observing Network tab for 15-30 seconds
   - **Expected**: After all downloads inactive, polling stops automatically
   - **Expected**: No more GET requests to `/api/downloads`

6. **Verify polling resumes on new download**
   - With polling stopped, add a new download
   - **Expected**: Polling resumes immediately
   - **Expected**: 3-second interval restored

7. **Test manual refresh**
   - With polling stopped, refresh the browser page
   - **Expected**: Downloads list loads correctly
   - **Expected**: Polling remains stopped (no active downloads)
   - Add new download
   - **Expected**: Polling starts

### Common Failure Modes

- **Polling never stops**: Logic may not detect all inactive states; check for downloads stuck in transitional states
- **Polling never starts**: Download creation may not trigger startPolling(); check store implementation
- **Polling interval wrong**: May be using different interval value; verify 3-second interval in logs/network
- **Memory leak**: Long-running frontend with many start/stop cycles may accumulate intervals; monitor browser memory

---

## Additional Verification Points

### Browser Console
- No JavaScript errors during normal operation
- No warnings about failed API calls (except during intentional failure tests)
- WebSocket or polling connections managed properly

### Backend Logs
- Minimal log spam during normal operation
- State changes logged clearly (e.g., "Download XYZ transitioned from downloading to completed")
- Errors include actionable information
- No repeated error messages indicating infinite retry loops

### qBittorrent Integration
- Torrents appear in qBittorrent with correct names
- Categories applied correctly (if configured)
- Download locations correct
- Torrents removed when cancelled

### Database Integrity
- Download records persist across backend restarts
- Status updates reflected in database
- Organized paths saved correctly
- No orphaned records after cancellation

### Performance
- Frontend remains responsive with 10+ downloads
- Polling does not cause UI jank
- Organization does not block the monitor from tracking other downloads
- Large files (1GB+) organize without timeout

---

## Troubleshooting Guide

### Download Stuck in "Queued"
- **Check**: qBittorrent running and accessible
- **Check**: qBittorrent credentials correct
- **Check**: Backend logs for connection errors
- **Action**: Test connection via configuration page

### Organization Fails Repeatedly
- **Check**: Destination path exists and is writable
- **Check**: Sufficient disk space (requires 110% of download size)
- **Check**: Template syntax is valid
- **Check**: No special characters in paths causing issues
- **Action**: Try manual organization with simplified template

### Frontend Not Updating
- **Check**: Browser console for errors
- **Check**: Polling is running (Network tab)
- **Check**: Backend is responding to API requests
- **Action**: Hard refresh browser (Ctrl+Shift+R)

### Missing Files After Organization
- **Check**: Files still exist in qBittorrent download folder
- **Check**: Organization completed successfully (status "organized")
- **Check**: organized_path is correct
- **Action**: Re-run organization if needed

---

## Acceptance Criteria

Before declaring the application production-ready, ALL of the following must be verified:

- [ ] Scenario 1 (Happy Path) completes successfully with real torrent
- [ ] Scenario 2 (MAM) downloads and organizes correctly (if MAM integration used)
- [ ] Scenario 3 (Manual Organization) works when auto-org disabled
- [ ] Scenario 4 (Retry) successfully recovers from organization failure
- [ ] Scenario 5 (Cancellation) removes download at all stages
- [ ] Scenario 6 (Connection Testing) provides accurate feedback
- [ ] Scenario 7 (Polling) starts/stops automatically as expected
- [ ] No errors in browser console during normal operation
- [ ] No repeated errors in backend logs
- [ ] All organized files accessible and complete
- [ ] Frontend remains responsive with multiple concurrent downloads
- [ ] Backend handles qBittorrent temporary unavailability gracefully

---

## Notes

- **Test Data**: Use small, legal torrents (e.g., public domain audiobooks) to minimize test time
- **Cleanup**: After testing, clear test data from qBittorrent and destination folder
- **Automation**: These manual tests can be automated using Playwright/Cypress in future phases
- **Regression**: Re-run these scenarios after any significant code changes to backend or frontend
