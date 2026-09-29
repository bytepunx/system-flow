---
id: T-0577
type: task
nature: improvement
title: items.count answers how many items are active and archived
status: done
parent: S-0162
owner: alex
created: 2026-09-29T20:15:55Z
updated: 2026-09-29T20:19:16Z
transitions:
  - to: ready
    at: 2026-09-29T20:16:04Z
    by: agent-S-0162
  - to: in-progress
    at: 2026-09-29T20:17:52Z
    by: agent-S-0162
  - to: done
    at: 2026-09-29T20:19:16Z
    by: agent-S-0162
stream: S-0162
tags: []
touches: [flai/internal/hostapi]
usage:
  source: log
  seconds: 84
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 4327
      cache_read: 1218034
      cache_write: 16220
      cost: 0.46
---
# T-0577 items.count answers how many items are active and archived

## Work

- A read method `items.count` answering `{ active, archived }`, so the overview can say how many items are archived without asking for the archive.
- A valid call and its refusals in the host API's tests; the method in `flai-cli.md`.

## Done when

- `go test -race ./internal/hostapi/... ./cmd/...` and lint pass.

## Notes
