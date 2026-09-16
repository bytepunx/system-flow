---
id: T-058
type: task
nature: feature
title: Replace stderr messages and the error boundary with slog events
status: done
parent: S-031
owner: alex
created: 2026-09-16T23:45:51Z
updated: 2026-09-16T23:48:50Z
transitions:
  - to: ready
    at: 2026-09-16T23:48:49Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:48:49Z
    by: agent
  - to: done
    at: 2026-09-16T23:48:50Z
    by: agent
stream: S-031
tags: [cli, logging]
---

# T-058 Replace stderr messages and the error boundary with slog events

## Work
Root command builds the logger once per run, passes it via app, adds --verbose; created-config, fetching-template, WIP-limit and other warnings become slog events with component and fields; run() logs the final error once as fatal with err and exits 1; stdout output unchanged.

## Done when
In-process CLI tests see structured lines on stderr and unchanged stdout.

## Notes
