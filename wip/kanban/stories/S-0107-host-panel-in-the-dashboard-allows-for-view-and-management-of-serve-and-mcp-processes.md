---
id: S-0107
type: story
nature: improvement
title: Host panel in the dashboard allows for view and management of serve and MCP processes
status: in-progress
parent: E-0003
owner: alex
created: 2026-09-24T01:26:46Z
updated: 2026-09-24T01:35:33Z
transitions:
  - to: ready
    at: 2026-09-24T01:28:15Z
    by: alex
  - to: in-progress
    at: 2026-09-24T01:35:00Z
    by: system-flow
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
- [ ] a single `upgrade` button should download, install, and then restart the serve and MCP sub processes via a command passed back to the host

## Tasks
- T-0379 The dashboard's /api/host route reads the host and asks it to start, stop, restart, check, and upgrade
- T-0380 The host panel shows serve and the MCP servers with start, stop, and restart, and one check and one upgrade for flai
- T-0381 The panel works against a real flai host, with the design and docs saying so

## Notes

- Verified live on 2026-09-24 against a scratch `flai host` built from `story/S-0106`, with its own config, address, and project, and a scratch flaiover dev server from `story/S-0107`, driven in headless Chromium. It showed the host, serve, and the MCP server with state and version; MCP stop, start, and restart; serve restart, reconnecting to the new serve; serve stop with its confirmation; the check (1.14.0 available while running 1.13.0, and "the latest" when running 1.14.0); an upgrade that found nothing newer; and the 403 while the `host` action is off.
- Criterion 4 is not ticked. The button sends `host.upgrade` to the host, and its up-to-date answer was seen live. The download, install, and restart could not be exercised, because no published flai newer than 1.14.0 contains `flai host` yet: installing 1.14.0 would leave the host a binary without the command. Once S-0106 is released, a scratch host built from it as an older version can upgrade to that release for real.
- Found live and fixed in S-0107: a serve restart ran twice, because flaiover retried the write after serve went down. Recorded on S-0106 as TH-0006.
