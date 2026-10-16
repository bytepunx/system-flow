---
id: T-1744
type: task
nature: feature
title: serve starts a role whose trigger is a schedule each time it comes round, through the cadence the analysis scheduler uses
status: backlog
parent: S-0418
owner: alex
created: 2026-10-16T05:45:31Z
updated: 2026-10-16T05:45:44Z
transitions: []
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/serve/analysis_schedule.go]
after: [T-1741]
---
# T-1744 serve starts a role whose trigger is a schedule each time it comes round, through the cadence the analysis scheduler uses

## Work

Add a scheduler per role with a schedule, built on the `cadence` that `analysis_schedule.go` and `replan.go` share. Move `cadence` out of `analysis_schedule.go` into a shared helper if it lives there, without changing the analysis scheduler's behaviour.

- Each time the schedule comes round with the action on, start the role through `RunRole` with the trigger `schedule <spec>`.
- Log a refused start as `scheduled role not started` and do not retry it before the next time.
- Times missed between two looks come round once.
- Test it with a fake clock.

Waits for T-1741. Runs beside T-1743 and T-1746.

## Done when

`go test ./internal/serve/...` passes. `analysis_schedule_test.go` is unchanged.

## Notes
