---
id: T-2851
type: task
nature: remediation
title: TestAFormulaRoleStepsEndIsToldToTheInstancesAgent waits for each run's end in the inbox instead of reading it once
status: done
parent: S-0653
owner: alex
created: 2026-11-10T07:01:30Z
updated: 2026-11-10T07:35:01Z
transitions:
  - to: ready
    at: 2026-11-10T07:34:25Z
    by: agent-S-0653
  - to: in-progress
    at: 2026-11-10T07:34:26Z
    by: agent-S-0653
  - to: done
    at: 2026-11-10T07:35:01Z
    by: agent-S-0653
stream: S-0653
tags: [flai, test]
touches: [flai/internal/serve/role_test.go]
usage:
  source: log
  seconds: 35
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 932
      cache_read: 849664
      cache_write: 11071
      cost: 0.2772
---
# T-2851 TestAFormulaRoleStepsEndIsToldToTheInstancesAgent waits for each run's end in the inbox instead of reading it once

## Work

First layer: it waits for no task, and shares no path with the task that fixes the activity measure.

- In `flai/internal/serve/role_test.go`, the test reads each run's end with `inbox.CatchUp` once, right after the run's `Ended` is set. `roleEnded` appends the agent event after it sets `Ended`, so under load the event may not be written yet. Poll for it with `waitFor`, reading the inbox with a fresh reader each time, until the one end of that run is there, then check its `by`, `to`, and `why` as now.
- Check `TestADefinedRoleRunIsJudgedByItsQuestionAndItsExit` and the other tests in the file that read an event or the activity once after `Ended`, and make each wait the same way.
- Leave the 30-second deadline of `waitFor` in `agents_test.go` as it is: the deadline is for a loaded host, not for speed.

## Done when

- No test in `role_test.go` reads an end event or the activity once after `Ended` without waiting for it.
- `go test ./internal/serve -run 'Role' -count 10` passes in `flai/`.
- `flai test flai/internal/serve/role_test.go` passes.

## Notes
