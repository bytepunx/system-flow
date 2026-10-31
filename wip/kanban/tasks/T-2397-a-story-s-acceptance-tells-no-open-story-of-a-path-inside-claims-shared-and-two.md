---
id: T-2397
type: task
nature: improvement
title: A story's acceptance tells no open story of a path inside claims.shared, and two ready stories that both edit docs/users/flai.md start together
status: in-progress
parent: S-0550
owner: alex
created: 2026-10-31T07:43:30Z
updated: 2026-10-31T08:17:29Z
transitions:
  - to: ready
    at: 2026-10-31T08:17:28Z
    by: agent-S-0550
  - to: in-progress
    at: 2026-10-31T08:17:29Z
    by: agent-S-0550
stream: S-0550
tags: [cli]
touches: [flai/cmd/accept_overlap.go, flai/cmd/accept_overlap_test.go, flai/internal/serve/hold_test.go]
---
# T-2397 A story's acceptance tells no open story of a path inside claims.shared, and two ready stories that both edit docs/users/flai.md start together

## Work

The hold already skips the shared paths. `workitem.Overlaps` (`flai/internal/workitem/hold.go`) returns false when the narrower entry lies inside a `claims.shared` pattern (ADR-0096 §2), and the grown-claim notice (`itemedit.Grown`) goes through it too. The acceptance notice does not. `tellAccepted` in `flai/cmd/accept.go` calls `tellOverlaps` in `flai/cmd/accept_overlap.go`, which calls `overlapsOf`, which calls `itemedit.Covering`. Its comment says "The shared paths are no exception". The notice is written to `.flai-cache/overlaps.jsonl` and reaches the open story as an `overlapped` change. The merge queue lands each story through `flai accept`, so it takes the same path.

- In `overlapsOf`, leave out each changed path that lies inside a `claims.shared` pattern (`manifest.Shared`, `Match` in `flai/internal/manifest/shared.go`) before `Covering` is asked. A story whose changed paths all lie inside shared patterns tells nobody. Leave `itemedit.Covering` and `flai task done`'s message (`flai/internal/taskdone/taskdone.go`) as they are, as the story's ADR decides.
- In `flai/cmd/accept_overlap_test.go`, beside `TestAcceptTellsOverlappingOpenStories`, add a manifest with `claims.shared`. An acceptance that changes `docs/users/flai.md` and one file outside the shared paths tells the open story only of that file. An acceptance that changes only `docs/users/flai.md` tells nobody.
- In `flai/internal/serve/hold_test.go`, beside `TestOneLookStartsNoTwoThatOverlap`, write `claims.shared: [docs/users/flai.md]` into the lab's `system-flow.yaml`, as `orchestrate_test.go` does. Give two ready stories touches of `docs/users/flai.md` and one file each of their own, and show that one look of the launcher starts both. `newAgentLab` writes no shared paths, so the test writes them itself.

First layer: it waits for nothing, and no other task of the story changes these files.

## Done when

The two new tests pass, and `TestOverlapsOf`, `TestAcceptTellsOverlappingOpenStories`, and `TestCovering` still pass. `go test -race ./cmd ./internal/serve ./internal/itemedit` passes. Story criterion 3.

## Notes

Drafted by the planner of S-0550 on 2026-10-31. If `TestAnOverlapAtAcceptanceReachesTheAgent` in `flai/internal/mcpserver/server_test.go` relies on a shared path, change it too, and `flai task done` widens the touches.
