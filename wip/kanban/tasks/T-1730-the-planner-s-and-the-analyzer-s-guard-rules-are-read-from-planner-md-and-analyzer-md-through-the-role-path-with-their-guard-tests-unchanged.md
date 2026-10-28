---
id: T-1730
type: task
nature: feature
title: The planner's and the analyzer's guard rules are read from planner.md and analyzer.md through the role path, with their guard tests unchanged
status: done
parent: S-0416
owner: alex
created: 2026-10-16T05:44:09Z
updated: 2026-10-28T22:27:34Z
transitions:
  - to: ready
    at: 2026-10-28T21:05:56Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T21:05:57Z
    by: system-flow
  - to: done
    at: 2026-10-28T21:40:33Z
    by: agent-S-0416
stream: S-0416
tags: [cli, template]
touches: [flai/internal/guard/guard.go, design/workflows/roles/planner.md, design/workflows/roles/analyzer.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/analyzer.md, flai/internal/guard/call.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/cmd/guard.go, docs/operators/settings.md]
after: [T-1729]
usage:
  source: log
  seconds: 2076
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 68
      output: 26973
      cache_read: 5696638
      cache_write: 262786
      cost: 3.1145
---
# T-1730 The planner's and the analyzer's guard rules are read from planner.md and analyzer.md through the role path, with their guard tests unchanged

## Work

Route `plan` and `analyze` through the role path from T-1729. Remove from `guard.go` the fixed lists the role files now state: `cliPlans`, `MCPAnalyzes`, `cliAnalyzes`, `cliAuthors`, and the planner's and analyzer's tool lists. Keep the named conditions in Go: draft-only `item_new`, moves to backlog, and `inReports`. Correct the two role files, here and in the template, wherever the tests show they differ from today's rules. The refusal texts stay word for word.

Waits for T-1729. T-1731 waits for this task, since both change `guard.go`.

## Done when

`guard_test.go`, `call_test.go`, `writes_test.go`, `authors_test.go`, and `flai/cmd/guard_test.go` pass unchanged.

## Notes

- Planner (S-0416, 2026-10-21): the role files are under `design/workflows/roles/`, as TH-0471's answer 2 moved them in S-0414, and the touches say so.
- The touches add `flai/internal/roles/conditions.go` and its test, in case a rule here needs a named condition that S-0414's set lacks. Leave them alone if the set has every one.
- `flai/cmd/guard.go` is not touched here. Its project load for `analyze` keys on the role name and keeps working. T-1734 rewrites its help.
