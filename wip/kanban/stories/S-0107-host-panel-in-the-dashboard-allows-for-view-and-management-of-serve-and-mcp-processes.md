---
id: S-0107
type: story
nature: improvement
title: Host panel in the dashboard allows for view and management of serve and MCP processes
status: backlog
parent: E-0003
owner: alex
created: 2026-09-24T01:26:46Z
updated: 2026-09-24T01:26:46Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0107 Host panel in the dashboard allows for view and management of serve and MCP processes

## Goal

The host panel should show the status and version of the serve and MCP sub-processes.

## Acceptance criteria
- [ ] the host panel should show an areas (like the one for dashboard) with the status and version of the serve and MCP sub-processes managed by the host process
- [ ] each sub-process should have controls for start/stop/restart
- [ ] a single `check for upgrade` button should check to see if a new flai cli version is available
- [ ] a single `upgrade` button should download, install, and then restart the serve and MCP sub processes via a command passed back to the host

## Tasks

## Notes
