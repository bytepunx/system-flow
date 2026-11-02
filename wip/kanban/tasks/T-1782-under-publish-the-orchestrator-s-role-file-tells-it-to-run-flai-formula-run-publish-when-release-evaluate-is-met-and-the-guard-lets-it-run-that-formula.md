---
id: T-1782
type: task
nature: improvement
title: Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula
status: done
parent: S-0423
owner: alex
created: 2026-10-16T05:48:37Z
updated: 2026-11-02T04:45:47Z
transitions:
  - to: ready
    at: 2026-11-02T04:27:39Z
    by: agent-S-0423
  - to: in-progress
    at: 2026-11-02T04:27:40Z
    by: agent-S-0423
  - to: done
    at: 2026-11-02T04:45:47Z
    by: agent-S-0423
stream: S-0423
tags: [cli, template]
touches: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/cmd/orchestrate_permissions_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/roles/baseline/orchestrator.md]
after: [T-1780, T-2013]
usage:
  source: log
  seconds: 1087
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 111
      output: 671
      cache_read: 6568007
      cache_write: 155646
      cost: 2.3662
---
# T-1782 Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula

## Work

In `design/workflows/roles/orchestrator.md`, template and here:

- list `flai formula run publish` under the `orchestration.permissions.publish` permission's commands;
- tell the orchestrator, under `publish`, to run `flai formula run publish --ready --autocommit` when `release_evaluate` is met and the project has a `publish` formula, and to log the instance's ID with `activity_log`;
- keep `release_publish` for a project with no `publish` formula.

In `flai/cmd/formula.go`, flai holds the run itself, as it holds the orchestrator's other writes: `flai formula run` under the orchestrator's role (`FLAI_ROLE=orchestrate`) is refused unless `orchestration.permissions.publish` is on and the formula is `publish`.

Add cases to `orchestrate_permissions_test.go` (or `formula_test.go`): the run passes with `publish` on and is refused with it off, and `formula run` of any other formula is refused. Change `flai/internal/guard/roles.go` only if the guard does not already allow a permission's `formula run` line (S-0422's T-2012 allows one in `commands`).

Waits for T-1780, whose formula file the test runs, and T-2013, which changes `flai/cmd/formula.go` first.

## Done when

`flai test` on the changed paths passes. The existing orchestrator permission cases are unchanged.

## Notes

Replanned 2026-10-21 by planner-S-0423 on TH-0471's answer 2: the role file is `design/workflows/roles/orchestrator.md`, in both copies. S-0419 reads the orchestrator's prompt from it, so no `.claude/agents` file changes. 2026-11-02, agent-S-0423: added the flai-side hold in `flai/cmd/formula.go` that the planner's note asked for.
