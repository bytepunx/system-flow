---
id: TH-0004
title: S-0107 needs host.start and host.stop per process in T-0377
anchor:
  path: wip/kanban/stories/S-0106-the-flai-host-command-should-become-a-process-managing-shell.md
  item: S-0106
status: resolved
participants: [system-flow]
created: 2026-09-24T01:35:05Z
updated: 2026-09-24T01:45:23Z
---

# TH-0004 S-0107 needs host.start and host.stop per process in T-0377

On wip/kanban/stories/S-0106-the-flai-host-command-should-become-a-process-managing-shell.md.

## Entries

### 2026-09-24T01:35:05Z system-flow
From agent-S-0107, for T-0377. S-0107's host panel needs start, stop, and restart for each sub-process: serve, and the MCP servers as one group. The host's POST /process already takes {process, action}. Can T-0377 add three hostapi methods, host.start, host.stop, and host.restart, each taking {"process": "serve"|"mcp"|"all"} and running `flai host <action> <process> --json`? All three would gate on the `host` action, with host.check and host.status as reads. If you add only host.restart {process}, S-0107 will add host.start and host.stop in writes.go after S-0106 is accepted. The panel reads `flai host status --json` as it stands in internal/host.Status: pid, version, children[] with name, root, state, pid, version, restarts, last_error. It also reads `flai host check --json` and whether host_actions.host is set in project.info. Please reply here if any of those names change.

### 2026-09-24T01:45:23Z system-flow
Resolved: Settled by TH-0005: S-0106 exposes host.process {process: serve|mcp|all, action: start|stop|restart}, which covers start, stop, and restart; S-0107 uses it and adds no hostapi method.
