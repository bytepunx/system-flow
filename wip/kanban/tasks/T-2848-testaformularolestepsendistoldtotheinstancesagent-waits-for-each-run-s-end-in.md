---
id: T-2848
type: task
nature: remediation
title: TestAFormulaRoleStepsEndIsToldToTheInstancesAgent waits for each run's end in the inbox instead of reading it once
status: backlog
parent: S-0653
owner: alex
created: 2026-11-10T07:01:18Z
updated: 2026-11-10T07:01:18Z
transitions: []
stream: S-0653
tags: [flai, test]
touches: [flai/internal/serve/role_test.go]
---
# T-2848 TestAFormulaRoleStepsEndIsToldToTheInstancesAgent waits for each run's end in the inbox instead of reading it once

## Work

First layer: it waits for no task, and shares no path with the task that fixes the activity measure.

- In `flai/internal/serve/role_test.go`, the test reads each run's end with `inbox.CatchUp` once, right after the run's `Ended` is set. `roleEnded` appends the agent event after it sets `Ended`, so under load the event may not be written yet. Poll for it with `waitFor`, reading the inbox with a fresh reader each time, until the one end of that run is there, then check its `by`, `to`, and `why` as now.
- Check `TestADefinedRoleRunIsJudgedByItsQuestionAndItsExit` and the other tests in the file that read an event or the activity once after `Ended`, and make each wait the same way.
- Leave the 30-second deadline of `waitFor` in `agents_test.go` as it is: the deadline is for a loaded host, not for speed.

## Done when

- No test in `role_test.go` reads an end event or the activity once after `Ended` without waiting for it.
- `go test ./internal/serve -run 'Role' -count 10` passes in `flai/`.
- `flai test flai/internal/serve/role_test.go` passes.
