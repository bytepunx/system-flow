---
id: T-2359
type: task
nature: feature
title: The analyzer and the curator start from their role files' triggers through the same launchers, in place of the analysis and digest schedulers' own rules
status: backlog
parent: S-0554
owner: alex
created: 2026-10-30T20:49:20Z
updated: 2026-10-30T20:49:20Z
transitions: []
stream: S-0554
tags: [cli]
touches: [flai/internal/serve/analysis_schedule.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/serve.go]
after: [T-2358]
---
# T-2359 The analyzer and the curator start from their role files' triggers through the same launchers, in place of the analysis and digest schedulers' own rules

## Work

Read the analyzer's `{schedule: analysis.schedule}` and the curator's `{event: accepted}` and `{schedule: digest.schedule}` from their role files, the project's or the baseline, through the launchers the planner's triggers task opened to built-in roles. Starts still go through `analyzeCheck` and the curator's check, with an empty focus and the window since its last run as today, so their records, coalescing, and journal entries stay as they are. Wire it in `flai/internal/serve/serve.go`. If S-0424 has landed, an `analysis.schedule` that names a formula still starts it.

Fourth layer: it waits for the planner's triggers task, which changes the launchers this task uses.

## Done when

`analysis_schedule_test.go` and `curate_schedule_test.go` pass unchanged. A new case shows a project's analyzer role file with no schedule trigger, which then starts no analysis on `analysis.schedule`.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. In S-0419's worktree that day, the schedulers are at `analysis_schedule.go:67` and `curate_schedule.go:110`. The new case goes in a test file of `flai/internal/serve`, which `flai task done` adds to the touches.
