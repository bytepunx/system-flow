---
id: T-2658
type: task
nature: improvement
title: flai accept reads git status --porcelain once per state of the checkout it needs, not six times
status: done
parent: S-0609
owner: alex
created: 2026-11-07T23:49:31Z
updated: 2026-11-08T02:56:48Z
transitions:
  - to: ready
    at: 2026-11-08T02:45:59Z
    by: agent-S-0609
  - to: in-progress
    at: 2026-11-08T02:46:00Z
    by: agent-S-0609
  - to: done
    at: 2026-11-08T02:56:48Z
    by: agent-S-0609
stream: S-0609
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_markdown.go, flai/internal/preview/accept.go, flai/cmd/accept_uncommitted_test.go]
after: [T-2649]
usage:
  source: log
  seconds: 648
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 25655
      cache_read: 10803167
      cache_write: 162683
      cost: 3.7202
---
# T-2658 flai accept reads git status --porcelain once per state of the checkout it needs, not six times

## Work

- An acceptance runs `git status --porcelain` six times today: the preview's `storygit.DirtyFiles` and its `uncommitted` per path (`flai/internal/preview/accept.go`), and `verify.ChangedPaths` in each `fixUncommitted` (`flai/cmd/accept_markdown.go`, called from `flai/cmd/accept.go`).
- Read the main checkout's status once before the acceptance changes anything, and pass the paths the preview already holds (`res.Uncommitted`) to the fix and the commits instead of reading them again. Read it again only after a step that changes the checkout, such as the merge or the archive, and only where a later step needs it.
- Answer the preview's per-path `uncommitted` from that one read where it runs in the same acceptance.
- Extend `flai/cmd/accept_uncommitted_test.go` with a counting runner that shows the reads, and that the uncommitted files the acceptance leaves (ADR-0248) are still left.
- Waits for T-2649: both change `fixUncommitted` and its callers in `flai/cmd/accept.go` and `flai/cmd/accept_markdown.go`.

## Done when

- A counting test shows one `git status --porcelain` of the main checkout per state of it that the acceptance reads, with each read after the first following a step that changed the checkout.
- The uncommitted-file tests of `flai accept` pass unchanged.
- `flai test flai/cmd/accept.go flai/cmd/accept_markdown.go flai/internal/preview/accept.go flai/cmd/accept_uncommitted_test.go` passes.

## Notes
