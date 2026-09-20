---
id: T-0277
type: task
nature: feature
title: "flai mcp http, start, stop, status, and token: the server as a process of its own"
status: done
parent: S-0076
owner: alex
created: 2026-09-20T12:16:18Z
updated: 2026-09-20T12:25:24Z
transitions:
  - to: ready
    at: 2026-09-20T12:25:24Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:25:24Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:25:24Z
    by: system-flow
stream: S-0076
tags: []
---
# T-0277 flai mcp http, start, stop, status, and token: the server as a process of its own

## Work
flai mcp stays stdio. flai mcp http runs the server in the foreground on 127.0.0.1:4243 unless --addr says otherwise; start detaches it, stop ends it by its recorded PID, status says where it listens and how to connect; token prints or rotates .flai-cache/mcp.token (0600). State in .flai-cache/mcp.json. A non-loopback address is allowed and warned about: TLS is not flai's job.

## Done when
- Command tests: start, status, a tool call over HTTP, stop; a second start says it is running; a stale state file is cleaned
- Works from a worktree and with FLAI_CACHE_DIR set
- make flai-test passes

## Notes
