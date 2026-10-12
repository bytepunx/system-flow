---
id: T-1435
type: task
nature: remediation
title: Confirm S-0348's scoped-check test covers I-0122's instance and name I-0122 in it
status: done
parent: S-0361
owner: alex
created: 2026-10-12T14:14:11Z
updated: 2026-10-12T14:20:51Z
transitions:
  - to: ready
    at: 2026-10-12T14:20:38Z
    by: agent-S-0361
  - to: in-progress
    at: 2026-10-12T14:20:38Z
    by: agent-S-0361
  - to: done
    at: 2026-10-12T14:20:51Z
    by: agent-S-0361
stream: S-0361
tags: [flai, check]
touches: [flai/internal/check/scope_test.go]
usage:
  source: log
  seconds: 13
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 881
      cache_read: 325556
      cache_write: 19563
      cost: 0.2393
---
# T-1435 Confirm S-0348's scoped-check test covers I-0122's instance and name I-0122 in it

## Work

I-0122 has one instance: S-0291's close-out at 2026-10-08T08:30:23Z found `wip/kanban/board.md`: 6 stories in in-progress, limit 5. It has the same title and cause as I-0123, which S-0348 fixed at fd5c2970 (2026-10-08T09:35:05Z, released in flai 1.39.15) under [ADR-0133](../../../design/adrs/0133-a-close-out-records-no-board-wip-limit-and-a-check-scoped-to-a-story-leaves-out.md). The instance predates that fix.

Confirm that `TestScopedCheckLeavesOutABoardWIPLimitOutsideTheStory` in `flai/internal/check/scope_test.go`, case "in progress over its limit", reproduces I-0122's instance: in progress over its limit on the main checkout's board, left out of a check scoped to a story. Run it with `flai test flai/internal/check/scope_test.go`. Add I-0122 beside I-0123 in the test's comment, so the test names both issues it reproduces.

If the test does not cover the instance, add the case that does, and say why in the narrative's `## Decisions`. Waits for nothing: it is the first layer.

## Done when

- The test's comment names I-0122 and I-0123.
- `flai test flai/internal/check/scope_test.go` passes.
- The narrative's `## Decisions` records that S-0348's fix covers I-0122's instance, or the case added and why.

## Notes
