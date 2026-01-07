# Phase 2: Download Monitoring - Context

**Gathered:** 2026-01-06
**Status:** Ready for planning

<vision>
## How This Should Work

A persistent background service that constantly polls qBittorrent to track all active downloads. The monitor runs continuously, checking download states at regular intervals. When it detects a completed download, it immediately triggers the organization phase (Phase 4) without waiting or requiring manual intervention.

This is the glue between torrent submission (Phase 1) and file organization (Phase 4) — ensuring downloaded audiobooks automatically move through the pipeline.

</vision>

<essential>
## What Must Be Nailed

- **Reliable completion detection** - Never miss when a download finishes, even if qBittorrent restarts or network hiccups occur. This is non-negotiable — if we don't reliably detect completion, the entire automation chain breaks.

</essential>

<boundaries>
## What's Out of Scope

- **File organization** - Not touching files, that's Phase 4. This phase only monitors and detects completion.
- **Frontend updates** - No WebSocket or live UI updates yet, that's Phase 5. This is backend-only monitoring.
- **Error recovery and retry logic** - If a download fails in qBittorrent, we report it but don't try to fix it or retry automatically.

</boundaries>

<specifics>
## Specific Ideas

- **Track state transitions** - Monitor when downloads move between qBittorrent states (queued → downloading → seeding), not detailed metrics like speed or peer counts.
- **Graceful degradation** - If qBittorrent becomes unavailable, the monitor should keep trying without crashing. The service needs to be resilient to qBittorrent being down temporarily.
- **Automatic phase chaining** - On completion detection, immediately trigger file organization without waiting for external triggers.

</specifics>

<notes>
## Additional Context

The focus is on lifecycle state tracking rather than detailed download metrics. We care about "what state is this download in" not "how fast is it downloading."

The monitor is designed to be a reliable bridge in the automation chain — torrents flow from Phase 1 (submission) through Phase 2 (monitoring) to Phase 4 (organization) with minimal human intervention.

</notes>

---

*Phase: 02-download-monitoring*
*Context gathered: 2026-01-06*
