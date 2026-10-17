---
id: T-1880
type: task
nature: feature
title: The components methods and flai components serve take a story, flai serve stops a story's dev server when it leaves review, and components.max_servers caps them
status: done
parent: S-0442
owner: alex
created: 2026-10-16T21:37:24Z
updated: 2026-10-17T03:40:27Z
transitions:
  - to: ready
    at: 2026-10-17T02:56:02Z
    by: agent-S-0442
  - to: in-progress
    at: 2026-10-17T02:56:03Z
    by: agent-S-0442
  - to: done
    at: 2026-10-17T03:40:27Z
    by: agent-S-0442
stream: S-0442
tags: [cli, dashboard]
touches: [flai/internal/serve/components.go, flai/internal/serve/components_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flai/cmd/serve_actions.go, docs/users/flai-reference.md, flai/internal/config/config.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md]
after: [T-1879]
usage:
  source: log
  seconds: 2664
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 196
      output: 78985
      cache_read: 14031065
      cache_write: 369086
      cost: 6.5496
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
