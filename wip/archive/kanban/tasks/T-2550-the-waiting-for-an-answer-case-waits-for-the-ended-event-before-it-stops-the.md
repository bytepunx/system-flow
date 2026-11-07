---
id: T-2550
type: task
nature: remediation
title: The waiting-for-an-answer case waits for the ended event before it stops the agent
status: done
parent: S-0561
owner: alex
created: 2026-11-07T08:18:53Z
updated: 2026-11-07T08:25:20Z
transitions:
  - to: ready
    at: 2026-11-07T08:23:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:23:49Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:25:20Z
    by: agent-S-0561
stream: S-0561
tags: [cli]
touches: [flai/internal/serve/stop_test.go]
usage:
  source: log
  seconds: 91
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 648
      cache_read: 950313
      cache_write: 12771
      cost: 0.3052
---
# T-2550 The waiting-for-an-answer case waits for the ended event before it stops the agent

## Work

In `TestTheOperatorsStopIsOneStoppedAgentEvent`'s case "an agent waiting for an answer", in `flai/internal/serve/stop_test.go`, wait for the run's ended event before calling `stop`.

The case waits for `lab.run(id).Outcome == OutcomeAsked` and then stops the agent. But `endRun` in `flai/internal/serve/agents.go` records the outcome with `l.ended(run)` first and calls `appendAgentEnd` after it. On a loaded host `Stop` appends its stopped event in that gap, so stopped comes before ended, which is what the three instances of I-0248 report at line 321.

Add, after the wait for `OutcomeAsked` and before `stop`, a `waitFor` that `len(lab.storyEvents(id))` has reached 2, the started event and the ended one. Say in a comment that the launcher appends the end just after it records it, as `eventsAre`'s comment does. Change no code outside the test: the event log is news, not state (ADR-0211 §4), so a stop landing in that gap costs a reader only the order of two events it still gets.

To reproduce the failure, add a `time.Sleep` of 200ms between `l.ended(run)` and `appendAgentEnd` in `endRun` for one local run, without committing it.

This task waits for nothing.

## Done when

- With the temporary sleep in `endRun`, the case fails at the event order without the new wait and passes with it, and the narrative's `## Decisions` records both runs.
- The sleep is gone from the diff: `flai/internal/serve/agents.go` is unchanged on the branch.
- `go test -race -count=20 -run TestTheOperatorsStopIsOneStoppedAgentEvent ./internal/serve` passes from `flai/`.

## Notes
