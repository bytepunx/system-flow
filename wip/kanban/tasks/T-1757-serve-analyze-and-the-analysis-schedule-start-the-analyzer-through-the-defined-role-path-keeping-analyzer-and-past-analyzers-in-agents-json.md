---
id: T-1757
type: task
nature: improvement
title: serve.Analyze and the analysis schedule start the analyzer through the defined-role path, keeping analyzer and past_analyzers in agents.json
status: cancelled
parent: S-0419
owner: alex
created: 2026-10-16T05:46:19Z
updated: 2026-10-30T13:03:50Z
transitions:
  - to: cancelled
    at: 2026-10-30T13:03:50Z
    by: agent-S-0419
stream: S-0419
tags: [cli]
touches: [flai/internal/serve/analyze.go, flai/internal/serve/analysis_schedule.go]
after: [T-1756]
---
# T-1757 serve.Analyze and the analysis schedule start the analyzer through the defined-role path, keeping analyzer and past_analyzers in agents.json

## Work

- Make `serve.Analyze`, `analyzeCheck`, `analyzeStart`, and `analyzeEnded` call the defined-role path for `analyzer` with the focus as the run's input.
- Keep the report detection (`newestReport`), the issues an entry names, the record under `analyzer`, and the refusal texts.
- The analysis schedule starts its runs through the same path.

Waits for the planner's task, which shapes the path for the built-in roles. Runs beside the orchestrator's task, since they share no path.

## Done when

`analyze_test.go`, `analysis_schedule_test.go`, and the `flai/cmd` analyze tests pass unchanged.

## Notes
- 2026-10-30T13:03:50Z: moved to cancelled: T-1756's startOf gives serve.Analyze and the analysis schedule the analyzer's role file at the start every run goes through; nothing is left for this task (TH-0760)
