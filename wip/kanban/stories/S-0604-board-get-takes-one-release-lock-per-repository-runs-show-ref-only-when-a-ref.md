---
id: S-0604
type: story
nature: improvement
title: board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:40:34Z
updated: 2026-11-07T22:40:34Z
transitions: []
tags: [cli]
touches: [flai/internal/release/history.go, flai/internal/mergequeue/eligible.go, flai/internal/hostapi/hostapi.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0604 board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing

## Goal

Cause 14 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `board.get` has 16,641 slow calls in 4 days at p50 800 ms, mostly `release.pending` (438 ms) and `merge.queue` (279 ms). `release`'s `kept` mutex is one for every repository (`flai/internal/release/history.go`) and is held across a `git show-ref`, so the 13 projects' boards queue on it: 23 to 88 ms alone, 250 ms with 8 at once. `merge.queue` reads every thread even with no story in review (`flai/internal/mergequeue/eligible.go`), 185 to 630 ms. Each story's undone paths are read one after another (`flai/internal/hostapi/hostapi.go`).

## Acceptance criteria

- [ ] `release`'s kept history takes one lock per repository root, and holds none across a git process another repository's call would wait for.
- [ ] While flai serve's main-refs watcher has seen no change, `release.pending` starts no git process.
- [ ] `merge.queue` reads no thread while no story is in review, and the undone paths of the stories are read at once.
- [ ] 13 `board.get` calls at once on the operator's host, one per served project, each answer under 300 ms, and server-performance.md records the before and after.

## Tasks

## Notes

- `publish.preview` holds the same lock and gains too.
