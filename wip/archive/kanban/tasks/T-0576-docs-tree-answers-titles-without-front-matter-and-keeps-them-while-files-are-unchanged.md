---
id: T-0576
type: task
nature: improvement
title: docs.tree answers titles without front matter and keeps them while files are unchanged
status: done
parent: S-0162
owner: alex
created: 2026-09-29T20:15:55Z
updated: 2026-09-29T20:17:52Z
transitions:
  - to: ready
    at: 2026-09-29T20:16:04Z
    by: agent-S-0162
  - to: in-progress
    at: 2026-09-29T20:16:04Z
    by: agent-S-0162
  - to: done
    at: 2026-09-29T20:17:52Z
    by: agent-S-0162
stream: S-0162
tags: []
touches: [flai/internal/hostapi]
usage:
  source: log
  seconds: 108
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 26
      output: 6964
      cache_read: 1960598
      cache_write: 26109
      cost: 0.7404
---
# T-0576 docs.tree answers titles without front matter and keeps them while files are unchanged

## Work

- `docs.tree` answers each file's name, path, kind, and title, without its front matter, which no page shows.
- `flai serve` keeps each file's title while the file has the same identity, modification time, and size, as `internal/workitem`'s store keeps items (S-0156), so a warm answer reads no file.
- Behaviour tests for the shape, for a title that changes, and for a file that is added or removed.

## Done when

- `flai hostapi --timing docs.tree` repeated in one process answers warm in under 10 ms on this repository, measured and recorded.
- `go test -race ./internal/hostapi/...` and lint pass.

## Notes
