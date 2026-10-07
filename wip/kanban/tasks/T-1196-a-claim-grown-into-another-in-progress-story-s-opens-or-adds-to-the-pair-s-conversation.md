---
id: T-1196
type: task
nature: improvement
title: A claim grown into another in-progress story's opens or adds to the pair's conversation
status: backlog
parent: S-0332
owner: alex
created: 2026-10-07T20:15:07Z
updated: 2026-10-07T20:15:07Z
transitions: []
stream: S-0332
tags: [flai]
touches: [flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go]
after: [T-1194]
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
