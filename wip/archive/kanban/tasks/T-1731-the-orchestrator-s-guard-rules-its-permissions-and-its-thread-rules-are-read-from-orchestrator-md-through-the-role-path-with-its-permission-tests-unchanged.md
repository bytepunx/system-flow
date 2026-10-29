---
id: T-1731
type: task
nature: feature
title: The orchestrator's guard rules, its permissions and its thread rules, are read from orchestrator.md through the role path, with its permission tests unchanged
status: done
parent: S-0416
owner: alex
created: 2026-10-16T05:44:15Z
updated: 2026-10-28T22:27:34Z
transitions:
  - to: ready
    at: 2026-10-28T21:41:02Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T21:41:03Z
    by: system-flow
  - to: done
    at: 2026-10-28T22:05:27Z
    by: agent-S-0416
stream: S-0416
tags: [cli, template]
touches: [flai/internal/guard/guard.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/guard/call.go, flai/cmd/guard.go, flai/internal/roles/baseline/orchestrator.md]
after: [T-1730]
usage:
  source: log
  seconds: 1464
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 72
      output: 28519
      cache_read: 6023084
      cache_write: 277845
      cost: 3.2929
---
# T-1731 The orchestrator's guard rules, its permissions and its thread rules, are read from orchestrator.md through the role path, with its permission tests unchanged

## Work

Route `orchestrate` through the role path. The orchestrator's file lists each permission with the calls it allows, as `design/system/strategic-agents.md` § Its permissions and the guard tabulates them. Keep in Go, as named conditions, the rules a list cannot state:

- `finalizesDraft`, and cost of delay inputs alone;
- the thread rules by opener and `answer_threads`;
- `--by orchestrator` on `accept`;
- the refusal of a placement by hand.

Keep `AppendRefusal`'s entry and the refusal texts word for word.

Waits for T-1730, which changes `guard.go` first.

## Done when

`flai/internal/guard/orchestrate_test.go`, `thread_test.go`, `flai/cmd/guard_permissions_change_test.go`, and `flai/cmd/orchestrate_permissions_test.go` pass unchanged.

## Notes

- Planner (S-0416, 2026-10-21): the role file is `design/workflows/roles/orchestrator.md`, as TH-0471's answer 2 moved it in S-0414, and the touches say so.
- The touches add `flai/internal/roles/conditions.go` and its test, in case a rule here needs a named condition that S-0414's set lacks. Leave them alone if the set has every one.
- `flai/cmd/guard.go` loads the orchestrator's permissions and thread openers by role name. That load stays, so this task does not touch the file.
