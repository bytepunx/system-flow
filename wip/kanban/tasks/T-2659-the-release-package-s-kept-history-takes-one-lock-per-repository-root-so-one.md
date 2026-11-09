---
id: T-2659
type: task
nature: improvement
title: The release package's kept history takes one lock per repository root, so one root's git processes never hold another root's call
status: done
parent: S-0604
owner: alex
created: 2026-11-07T23:49:37Z
updated: 2026-11-08T02:19:10Z
transitions:
  - to: ready
    at: 2026-11-08T02:07:51Z
    by: agent-S-0604
  - to: in-progress
    at: 2026-11-08T02:07:52Z
    by: agent-S-0604
  - to: done
    at: 2026-11-08T02:19:10Z
    by: agent-S-0604
stream: S-0604
tags: [cli]
touches: [flai/internal/release/history.go, flai/internal/release/release.go, flai/internal/release/history_test.go]
usage:
  source: log
  seconds: 678
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 22
      output: 5071
      cache_read: 2974989
      cache_write: 100884
      cost: 1.2696
---
# T-2659 The release package's kept history takes one lock per repository root, so one root's git processes never hold another root's call

## Work

`kept` in `flai/internal/release/history.go` is one mutex over a map of every root's history. `PendingBatch` in `flai/internal/release/release.go` takes it and holds it across every git process it runs: `show-ref`, `git log`, `git log -S`, and `diff-tree --stdin`. So 13 served projects' boards queue on one lock.

Give each root its own entry with its own mutex. Hold the outer lock only to find or create that entry, never across a git process. Keep the kept history's behaviour as it is: one `show-ref` while nothing changed, and the same answers as the per-process reference that `history_test.go` holds it to.

Waits for no task.

## Done when

- [ ] A test in `flai/internal/release/history_test.go` holds one root's call inside a git process and shows that a second root's call answers meanwhile.
- [ ] The existing kept-history tests pass unchanged, `TestPendingStartsOneGitProcessWhileNothingChanged` among them.
- [ ] `go test -race ./internal/release/...` passes.

## Notes
