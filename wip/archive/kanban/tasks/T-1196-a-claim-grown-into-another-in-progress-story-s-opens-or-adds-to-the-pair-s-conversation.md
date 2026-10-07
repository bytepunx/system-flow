---
id: T-1196
type: task
nature: improvement
title: A claim grown into another in-progress story's opens or adds to the pair's conversation
status: done
parent: S-0332
owner: alex
created: 2026-10-07T20:15:07Z
updated: 2026-10-07T23:03:37Z
transitions:
  - to: ready
    at: 2026-10-07T22:50:26Z
    by: agent-S-0332
  - to: in-progress
    at: 2026-10-07T22:50:27Z
    by: agent-S-0332
  - to: done
    at: 2026-10-07T23:03:37Z
    by: agent-S-0332
stream: S-0332
tags: [flai]
touches: [flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go, flai/cmd/touches.go, flai/internal/inbox/inbox.go, flai/internal/mcpserver/items_write_test.go]
after: [T-1194]
usage:
  source: log
  seconds: 790
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 27
      output: 10511
      cache_read: 1608586
      cache_write: 51000
      cost: 0.8379
---
# T-1196 A claim grown into another in-progress story's opens or adds to the pair's conversation

## Work

Extend `ClaimWatch.Grown` in `flai/internal/itemedit/claim.go`. It waits for T-1194, whose ADR settles the behaviour. It shares no path with T-1195 or T-1197 and runs beside them.

- For each other story in progress the grown claim reaches, open or add to the pair's conversation, `about` the paths gained, saying the claims grew to overlap and asking who changes them first.
- Keep the `overlapped` notice in `overlaps.jsonl` as it is.
- A conversation that cannot be written is logged, and the write that grew the claim stands, as `Untold` does today.

## Done when

- Tests cover a grown claim opening a conversation, a second growth adding an entry, a shared path telling nothing (S-0295), and a write failure leaving the edit in place.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
