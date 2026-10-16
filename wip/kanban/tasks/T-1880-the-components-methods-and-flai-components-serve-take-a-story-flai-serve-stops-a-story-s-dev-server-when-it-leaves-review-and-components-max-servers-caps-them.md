---
id: T-1880
type: task
nature: feature
title: The components methods and flai components serve take a story, flai serve stops a story's dev server when it leaves review, and components.max_servers caps them
status: backlog
parent: S-0442
owner: alex
created: 2026-10-16T21:37:24Z
updated: 2026-10-16T21:37:24Z
transitions: []
stream: S-0442
tags: [cli, dashboard]
touches: [flai/internal/serve/components.go, flai/internal/serve/components_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flai/cmd/serve_actions.go, docs/users/flai-reference.md]
after: [T-1879]
---
# T-1880 The components methods and flai components serve take a story, flai serve stops a story's dev server when it leaves review, and components.max_servers caps them

## Work

- `writes.go` and `hostapi/components.go`: `story` on `components.start`, `components.stop`, `components.status`, and `components.list`; the list reads the worktree's components.
- `serve/components.go`: key the servers by project, sub-project, and story. Stop a story's server when the watcher sees the story leave review, and refuse a start beyond `components.max_servers` with the servers running.
- `flai/cmd/components.go`: `serve --story`; `serve_actions.go`: `flai serve components max-servers <n>`. Regenerate the reference.

Waits for T-1879, whose worktree start it supervises.

## Done when

Tests cover a story's start and status beside main's, the stop on leaving review, the cap, and `serve --story`. `flai test flai/internal/serve flai/internal/hostapi flai/cmd` passes.

## Notes

None.
