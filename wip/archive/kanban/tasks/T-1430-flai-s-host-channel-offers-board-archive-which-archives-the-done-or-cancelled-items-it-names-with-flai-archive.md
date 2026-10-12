---
id: T-1430
type: task
nature: feature
title: flai's host channel offers board.archive, which archives the done or cancelled items it names with flai archive
status: done
parent: S-0343
owner: alex
created: 2026-10-09T16:38:50Z
updated: 2026-10-09T16:45:14Z
transitions:
  - to: ready
    at: 2026-10-09T16:39:21Z
    by: agent-S-0343
  - to: in-progress
    at: 2026-10-09T16:39:21Z
    by: agent-S-0343
  - to: done
    at: 2026-10-09T16:45:14Z
    by: agent-S-0343
stream: S-0343
tags: []
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, design/system/flai-cli.md, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 353
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 20
      output: 5661
      cache_read: 1100754
      cache_write: 35384
      cost: 0.5728
---

# T-1430 flai's host channel offers board.archive, which archives the done or cancelled items it names with flai archive

## Work

Add `board.archive` to `itemSpecs()` in `flai/internal/hostapi/writes.go`: params `{ ids }`, one or more item IDs, each checked with `needID`; an empty list is refused, since `flai archive` without IDs archives every closed item. It runs `flai archive --json -- <ids...>` and answers flai's `{ archived, narratives, resolved_threads, closed_messages, dry_run }`. Describe it in `design/system/flai-cli.md` beside `flai archive` and the host channel's methods. Waits for nothing: the first layer.

## Done when

- `writes_test.go`'s argument table and refusal table cover `board.archive`.
- `flai test` on the changed paths passes.

## Notes
