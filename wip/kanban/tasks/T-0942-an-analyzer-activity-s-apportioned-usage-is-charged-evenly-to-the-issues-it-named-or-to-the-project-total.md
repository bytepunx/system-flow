---
id: T-0942
type: task
nature: improvement
title: An analyzer activity's apportioned usage is charged evenly to the issues it named, or to the project total
status: backlog
parent: S-0227
owner: alex
created: 2026-10-05T05:45:35Z
updated: 2026-10-06T21:31:12Z
transitions: []
stream: S-0227
tags: [flai]
touches: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, design/system/flai-cli.md, design/system/strategic-agents.md]
after: [T-0940]
---
# T-0942 An analyzer activity's apportioned usage is charged evenly to the issues it named, or to the project total

## Work

Extend `activityMeasure.charge` in `flai/internal/serve/activity.go`, which charges only a planner's activity today, to the analyzer. It waits for T-0940, which gives issues their usage and the function that charges one.

- Take the issue IDs (`I-nnnn`) among the items the analyzer's activity entry names, and give each an even share of the usage apportioned to the activity's span, seconds and each model's tokens and cost, under its `analyzer` entry. Save each issue and regenerate nothing else.
- An activity that names no issue is charged to the project strategic total S-0226 adds, as S-0226 charges an orchestrator activity that names no item.
- An issue ID that names no issue is left out and said in what the activity reports, as S-0225 reports what it charged on `Logged`.
- Say in `design/system/flai-cli.md`, where `activity_log` and flai serve's activity log are described, that an analyzer activity is charged to the issues it names.

## Done when

- Tests pin an activity naming two issues, each charged half; one naming none, charged to the project total; and one naming an unknown issue, which charges nothing to it.
- A planner's and an orchestrator's charge are unchanged: their tests still pass.
- `scripts/flai-test.sh` passes.

## Notes
