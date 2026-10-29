---
id: T-1741
type: task
nature: feature
title: "A role file's trigger is read and checked: events, a schedule, or continuous, with roles.file findings for one flai cannot read"
status: cancelled
parent: S-0418
owner: alex
created: 2026-10-16T05:45:19Z
updated: 2026-10-29T16:36:17Z
transitions:
  - to: cancelled
    at: 2026-10-29T16:36:17Z
    by: agent-S-0418
stream: S-0418
tags: [cli]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go]
---
# T-1741 A role file's trigger is read and checked: events, a schedule, or continuous, with roles.file findings for one flai cannot read

## Work

Parse a role file's `trigger` into its kinds. The events are `accepted`, `cancelled`, `ready`, `review`, `edited`, and `thread-opened`. The schedule is parsed with `flai/internal/cron`. `continuous` stands alone. `asked` is always allowed. Report `roles.file` for an unknown event, a schedule that does not parse or never comes round, and `continuous` combined with another trigger. Add tests for each.

First layer: the three triggers' tasks read it.

## Done when

`go test ./internal/roles/... ./internal/check/...` passes.

## Notes
- 2026-10-29T16:36:17Z: moved to cancelled: TH-0749 option A: S-0414 and S-0536 already read and check triggers
