---
id: T-0872
type: task
nature: remediation
title: flai accept --dry-run lists a blocker naming each conflict marker the story branch carries
status: backlog
parent: S-0276
owner: alex
created: 2026-10-05T04:06:44Z
updated: 2026-10-05T04:06:44Z
transitions: []
stream: S-0276
tags: [flai]
touches: [flai/internal/preview, flai/cmd/accept_conflict_test.go]
after: [T-0871]
---
# T-0872 flai accept --dry-run lists a blocker naming each conflict marker the story branch carries

## Work

In `preview.Accept` (`flai/internal/preview/accept.go`), when the item is a story with git and `storygit.Branch(it.ID)` exists, call the `conflictmark` branch scan from T-0871. If it finds markers, append one blocker that names each one as `path:line`, shortened as the refusal shortens them, and says how to resolve them as `refuseConflictMarkers` does: in the worktree, or after `flai stream open` when the story has none. The dashboard's acceptance dialog, `GET /api/items/:id/acceptance`, shows `blockers[]` as they come, so nothing in `flaiover/` changes. An error from the scan is a blocker too, not a failed preview.

Add tests to `flai/cmd/accept_conflict_test.go`, next to S-0253's acceptance tests:

- `flai accept --dry-run --json` on a story branch that adds a marker lists the blocker with the path and line.
- On a clean branch it lists no such blocker.

This task waits for T-0871, whose function it calls.

## Done when

- `flai accept --dry-run` on a story branch with a conflict marker lists a blocker naming each path and line, and accepting the same story is refused for the same markers.
- On a clean branch the preview lists no conflict marker blocker.
- The new tests and `scripts/flai-test.sh` pass.

## Notes
