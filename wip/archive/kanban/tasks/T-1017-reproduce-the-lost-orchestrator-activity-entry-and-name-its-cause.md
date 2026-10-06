---
id: T-1017
type: task
nature: remediation
title: Reproduce the lost orchestrator activity entry and name its cause
status: cancelled
parent: S-0292
owner: alex
created: 2026-10-06T11:34:16Z
updated: 2026-10-06T11:35:51Z
transitions:
  - to: cancelled
    at: 2026-10-06T11:35:51Z
    by: agent-S-0292
stream: S-0292
tags: [flai, serve, test]
touches: [flai/internal/serve/orchestrate_test.go]
---
# T-1017 Reproduce the lost orchestrator activity entry and name its cause

## Work

Make `TestTheOrchestratorIsStoppedWhenTheActionIsTurnedOff/a_run_this_flai_serve_waits_for` fail as it did in S-0285's close-out. Run it under load: `go test -race -count=N -run TestTheOrchestratorIsStopped ./internal/serve` alongside the rest of the module, or with `-cpu 1` and a busy machine. Then find which of these left the activity document without its entry within 5 s:

- the `await` goroutine reaping the process late, so the deadline in `waitFor` passed first
- `logRunEnd` returning nil without logging: no newest run, a span outside the run, or no usage read from the stream
- `stop` and the `await` goroutine both settling the run, so that the one that logs reads the run's state before the other has written it

This task waits for nothing: it is the first layer, and the fix depends on what it finds.

## Done when

- The cause is named in the narrative's `## Decisions`, with the evidence: the failing run's output, and the code path that dropped the entry or ran late.
- A test fails on the current code because of that cause, reliably, not one time in many. If no test can, `## Decisions` says why.
- The proposed fix is recorded in `## Decisions` before it is built, as the story's goal asks.

## Notes

Planned by planner-S-0292. The three candidate causes come from reading `orchestrate.go`, `activity.go`, and `agents.go`; none is confirmed.
- 2026-10-06T11:35:51Z: moved to cancelled: T-1016, written before the planner's draft arrived, already reproduced, fixed, and closes I-0090 in one task
