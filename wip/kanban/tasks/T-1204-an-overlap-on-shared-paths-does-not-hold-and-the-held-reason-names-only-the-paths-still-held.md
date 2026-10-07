---
id: T-1204
type: task
nature: feature
title: An overlap on shared paths does not hold, and the held reason names only the paths still held
status: backlog
parent: S-0334
owner: alex
created: 2026-10-07T20:16:13Z
updated: 2026-10-07T20:16:13Z
transitions: []
stream: S-0334
tags: [flai]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/boardview.go]
after: [T-1203]
---
# T-1204 An overlap on shared paths does not hold, and the held reason names only the paths still held

## Work

Make the hold honour shares. It waits for T-1203, whose `Shares` it reads. It shares no path with T-1205 and runs beside it.

- `Holds` takes the shares in force, as a list of pairs of stories with paths, beside `claims.shared`; an overlapping pair of entries does not hold when the narrower lies wholly inside a path the holding story shares with the held one.
- The `held` reason names only the paths still held, and every story that still holds.
- Wherever `Holds` is built (`boardview.go`, and each caller T-1205 or a test finds, such as `serve/start.go` and `mcpserver/work.go`), pass the shares in, so `flai board`, `inbox`, `wait_for_work`, and the launcher agree. Widen this task's touches to each caller you change.

## Done when

- Tests cover a share that clears a hold, one that clears part of it, a share with a third story that clears nothing, and a share that ended.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
