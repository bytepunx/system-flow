---
id: T-3480
type: task
nature: feature
title: resolve.status gives the dashboard a job a state's runs started, with its state, its current step, its skipped cluster steps and why, and how it ended
status: in-progress
parent: S-0814
owner: alex
created: 2026-11-16T23:51:03Z
updated: 2026-11-28T00:43:59Z
transitions:
  - to: ready
    at: 2026-11-28T00:43:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-28T00:43:59Z
    by: system-flow
stream: S-0814
tags: [flai]
touches: [flai/internal/hostapi/resolve_status.go, flai/internal/hostapi/resolve_status_test.go]
---
# T-3480 resolve.status gives the dashboard a job a state's runs started, with its state, its current step, its skipped cluster steps and why, and how it ended

## Work

Criterion 2, the host API's part. S-0804's `resolve.status` gives the dashboard a Rebase and resolve job; S-0810 makes the job a state's `runs:` starts, with its steps recorded as they start and end and each skipped `kluster:` step's reason. Make `resolve.status` answer that job too: the state it runs for, the step it is on, each step's start, end, and outcome, each skipped step with the reason S-0810 records (no deploy configured, the kluster host action off, no agent namespace available), the log's path for `log.read`, and how it ended, passed, sent back by the state's `send-back` transition with its thread, or interrupted.

If S-0810 already answers all of this through `resolve.status`, write only the tests that show it, and say so in this task's notes.

Waits for no task: it reads S-0810's job record, on main before this story starts.

## Done when

- `resolve.status` answers a running `runs:` job, a passed one, one sent back, one interrupted, and one whose cluster steps were skipped, each with the fields above.
- A Rebase and resolve job answers as S-0804 made it.
- `resolve_status_test.go` covers these, and `flai test` on the changed paths passes.

## Notes

- Drafted by planner-S-0814. S-0804's T-3430 creates `resolve_status.go`; if S-0804 or S-0810 names the method or the file otherwise, follow theirs and widen the touches.
