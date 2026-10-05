---
id: T-0972
type: task
nature: feature
title: flai serve starts the analyzer when analysis.schedule comes round, while the analyze host action is on
status: backlog
parent: S-0223
owner: alex
created: 2026-10-05T05:47:36Z
updated: 2026-10-05T05:47:36Z
transitions: []
stream: S-0223
tags: [flai]
touches: [flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/serve/serve.go]
after: [T-0968]
---
# T-0972 flai serve starts the analyzer when analysis.schedule comes round, while the analyze host action is on

## Work

In `flai/internal/serve/analysis_schedule.go`, run a schedule for each project whose manifest sets `analysis.schedule`, as the replanner runs `planning.schedule` (`replan.go`, `scheduleDue`): at each look, when the schedule has come round since the last and the `analyze` host action is on, start an analyzer run for all three focuses with the trigger `schedule`, unless one is running. A schedule set or changed comes round first at its next time from then; times missed between two looks come round once. Start it with the serve loop (`serve.go`).

It waits for T-0968, whose run it starts.

## Done when

- a test with a fixed clock starts a run when the schedule comes round with the action on, and none with it off, none with a run going, and one only for times missed between two looks
- a changed schedule comes round first at its next time, in a test
- the run records its trigger as `schedule`
- `go test ./internal/serve/` passes

## Notes
