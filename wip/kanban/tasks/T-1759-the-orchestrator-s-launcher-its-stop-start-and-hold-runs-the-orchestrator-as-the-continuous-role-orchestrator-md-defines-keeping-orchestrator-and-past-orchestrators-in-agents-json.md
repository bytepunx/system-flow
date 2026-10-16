---
id: T-1759
type: task
nature: improvement
title: The orchestrator's launcher, its stop, start, and hold, runs the orchestrator as the continuous role orchestrator.md defines, keeping orchestrator and past_orchestrators in agents.json
status: backlog
parent: S-0419
owner: alex
created: 2026-10-16T05:46:24Z
updated: 2026-10-16T05:46:24Z
transitions: []
stream: S-0419
tags: [cli]
touches: [flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/role_continuous.go]
after: [T-1756]
---
# T-1759 The orchestrator's launcher, its stop, start, and hold, runs the orchestrator as the continuous role orchestrator.md defines, keeping orchestrator and past_orchestrators in agents.json

## Work

Make the orchestrator's launcher the continuous-role launcher from S-0418, over `orchestrator.md`.

- Keep `flai serve orchestrate stop` and `start`, `orchestrate.stop` and `orchestrate.start`, the hold kept on the run, the minute's wait after a failure, the refusal recorded once, the stop summaries, the journal entries under `orchestrate`, and the record under `orchestrator` and `past_orchestrators`.
- `flai role stop orchestrator` and `start orchestrator` act the same.

Waits for the planner's task. Runs beside the analyzer's task.

## Done when

`orchestrate_test.go`, `orchestrate_hold_test.go`, and the `flai/cmd` orchestrate tests pass unchanged.

## Notes
