---
id: S-0107
type: story
nature: improvement
title: Host panel in the dashboard allows for view and management of serve and MCP processes
status: done
parent: E-0003
owner: alex
created: 2026-09-24T01:26:46Z
updated: 2026-09-24T04:31:28Z
transitions:
  - to: ready
    at: 2026-09-24T01:28:15Z
    by: alex
  - to: in-progress
    at: 2026-09-24T01:35:00Z
    by: system-flow
  - to: review
    at: 2026-09-24T04:09:05Z
    by: system-flow
  - to: done
    at: 2026-09-24T04:31:28Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0107 Host panel in the dashboard allows for view and management of serve and MCP processes

## Goal

The host panel should show the status and version of the serve and MCP sub-processes.

## Acceptance criteria
- [x] the host panel should show an areas (like the one for dashboard) with the status and version of the serve and MCP sub-processes managed by the host process
- [x] each sub-process should have controls for start/stop/restart
- [x] a single `check for upgrade` button should check to see if a new flai cli version is available
- [x] a single `upgrade` button should download, install, and then restart the serve and MCP sub processes via a command passed back to the host

## Tasks
- T-0379 The dashboard's /api/host route reads the host and asks it to start, stop, restart, check, and upgrade
- T-0380 The host panel shows serve and the MCP servers with start, stop, and restart, and one check and one upgrade for flai
- T-0381 The panel works against a real flai host, with the design and docs saying so

## Notes

- Verified live on 2026-09-24 against a scratch `flai host` built from `story/S-0106`, with its own config, address, and project, and a scratch flaiover dev server from `story/S-0107`, driven in headless Chromium. It showed the host, serve, and the MCP server with state and version; MCP stop, start, and restart; serve restart, reconnecting to the new serve; serve stop with its confirmation; the check (1.14.0 available while running 1.13.0, and "the latest" when running 1.14.0); an upgrade that found nothing newer; and the 403 while the `host` action is off.
- Criterion 4 was verified live on 2026-09-24, with no published release that has `flai host` (v1.15.0's release failed; S-0108). `FLAI_RELEASES_API` points self-upgrade at a stand-in for GitHub's releases API, which served a build of this branch as flai 1.99.0, with its archive and `checksums.txt`. A scratch host built as 1.15.0, from main after S-0106, was used from the panel:
  - Check for upgrade said "flai 1.99.0 is available (running 1.15.0)".
  - Upgrade downloaded the archive and checksums, verified them, installed 1.99.0 over the host's binary, and restarted the host on it with the same PID. Serve and the MCP server came back on 1.99.0.
  - The panel reconnected and said "Upgraded flai 1.15.0 to 1.99.0; serve and MCP restarted".
  - Check then said "flai 1.99.0 is the latest".
- Found live and fixed in S-0107: a serve restart ran twice, because flaiover retried the write after serve went down. Recorded on S-0106 as TH-0006.
