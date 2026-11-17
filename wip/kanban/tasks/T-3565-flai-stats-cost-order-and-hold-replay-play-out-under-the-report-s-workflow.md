---
id: T-3565
type: task
nature: remediation
title: flai stats' cost order and hold replay play out under the report's workflow
status: backlog
parent: S-0843
owner: alex
created: 2026-11-17T12:29:16Z
updated: 2026-11-17T12:29:16Z
transitions: []
stream: S-0843
tags: [flai, workflow, metrics]
touches: [flai/internal/metrics/costorder.go, flai/internal/metrics/costorder_test.go, flai/internal/metrics/claims.go, flai/internal/metrics/claims_test.go]
after: [T-3561]
---
# T-3565 flai stats' cost order and hold replay play out under the report's workflow

## Work

Two callers in `flai/internal/metrics` still run the baseline although `Options.flow()` gives the report's workflow (I-0379):

- `costOrder` in `costorder.go` builds its puller with `newPuller = planning.NewPuller`, the baseline wrapper. Make `newPuller` take the workflow and default to `planning.NewPullerIn`, and call it with `opt.flow()`; update the test in `costorder_test.go` that replaces `newPuller` to the new signature.
- The hold replay in `claims.go` builds each judgement's Holds with `memo.NewHolds(in, opt.Projects)`. Call `memo.NewHoldsIn(flow, in, opt.Projects)` with the flow the function reads already.

Waits for T-3561, which adds `HoldsMemo.NewHoldsIn`.

## Done when

- Neither file calls `planning.NewPuller` or `HoldsMemo.NewHolds`.
- A test in `costorder_test.go` shows the projected pull times follow a workflow whose active state is not named `in-progress`: a story in that state holds a lane.
- A test in `claims_test.go` shows the replay counts no hold from a task closed under a workflow whose closed task state is not `done`.
- `flai test` on the four files passes.

## Notes
