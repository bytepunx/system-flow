---
id: T-3503
type: task
nature: feature
title: "S-0742 settles its rebase onto S-0748 and S-0755: replays keeps main's lines, and the built-in triggers are read on the project's board by kind"
status: done
parent: S-0742
owner: alex
created: 2026-11-17T02:56:00Z
updated: 2026-11-17T02:56:16Z
transitions:
  - to: ready
    at: 2026-11-17T02:56:06Z
    by: agent-S-0742
  - to: in-progress
    at: 2026-11-17T02:56:06Z
    by: agent-S-0742
  - to: done
    at: 2026-11-17T02:56:16Z
    by: agent-S-0742
stream: S-0742
tags: []
touches: [flai/internal/serve/replan.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go]
after: [T-3372]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 2171
      cache_read: 879058
      cache_write: 20226
      cost: 0.3463
---
# T-3503 S-0742 settles its rebase onto S-0748 and S-0755: replays keeps main's lines, and the built-in triggers are read on the project's board by kind

## Work

- TH-1268: the acceptance's rebase stopped in `flai/internal/serve/replan.go`, where S-0755 gave `replays` the project's workflow (`planning.ReplayIn(flow, …)`) and S-0742 gave it the queue states. Keep main's doc line, its `flow` parameter, and its `ReplayIn` call, with S-0742's `replayStates(flow)`, so the rebased branch undoes none of main's lines.
- TH-1233: S-0748's new test calls `reasonsOf` without the workflow S-0742 added; pass it the test's own workflow.
- S-0748 makes a role's `to` condition name any state of the project's workflow (ADR-0304), so the occurrence carries the moved state's current name, and `conditionState`, S-0742's bridge to the baseline's names, goes. The orchestrator's built-in triggers, used when its role file is missing or says continuous, are read on the project's board by kind (`builtinTriggersIn`): the baseline's waiting state is each waiting state, and its initial state the project's initial state.
- Waits for T-3372, the last task before the send back.

## Done when

- `flai stream sync S-0742` is clean and `flai test` passes on `flai/internal/serve`.

## Notes
