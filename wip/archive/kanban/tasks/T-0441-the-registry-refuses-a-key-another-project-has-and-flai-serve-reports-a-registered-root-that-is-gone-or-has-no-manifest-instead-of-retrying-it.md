---
id: T-0441
type: task
nature: feature
title: The registry refuses a key another project has, and flai serve reports a registered root that is gone or has no manifest instead of retrying it
status: done
parent: S-0121
owner: alex
created: 2026-09-26T05:25:59Z
updated: 2026-09-26T05:28:34Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:18Z
    by: agent-S-0118
  - to: in-progress
    at: 2026-09-26T05:26:18Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T05:28:34Z
    by: agent-S-0118
stream: S-0121
tags: []
---

# T-0441 The registry refuses a key another project has, and flai serve reports a registered root that is gone or has no manifest instead of retrying it

## Work

`serve.Dir.Register` refuses an entry whose key another root already has, naming that root, so every caller (flai dashboard, the board's import, flai serve project add) gets the same refusal. `serve.Unavailable(root)` says why a registered root cannot be served (the folder is gone, it has no system-flow.yaml, its manifest does not load). `serve.Run` skips such an entry instead of starting a client, watcher, and launcher for it, logs once when it goes and once when it comes back, and records the reason in `Status.Unavailable` by root.

## Done when

- `go test ./internal/serve/` covers the key conflict, a missing root, and a root that lost its manifest: the missing one has no client, one log line, and a reason in the state
- `scripts/flai-test.sh` passes

## Notes
