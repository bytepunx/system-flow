---
id: T-2765
type: task
nature: remediation
title: Close I-0284 on the manifest change that turned set_cost_of_delay on, naming the two regression tests
status: backlog
parent: S-0622
owner: alex
created: 2026-11-09T18:05:47Z
updated: 2026-11-09T18:05:47Z
transitions: []
stream: S-0622
tags: [cli]
touches: [design/issues/I-0284-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md, design/issues/summary.md]
after: [T-2763, T-2764]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 8
          cache_read: 26246
          cache_write: 880
          cost: 0.0072
---
# T-2765 Close I-0284 on the manifest change that turned set_cost_of_delay on, naming the two regression tests

## Work

Close I-0284 from the story's worktree with `flai issue close I-0284 --reason "<reason>" --commit`. The reason says what fixed the cause: commit c76b23a2d turned `orchestration.permissions.set_cost_of_delay` on in `system-flow.yaml` at 2026-11-07T09:17:55Z. It names the tests that now pin the path: T-2763's in `flai/internal/serve/orchestrate_confirm_test.go` and T-2764's in `flai/cmd/orchestrate_permissions_test.go`, by their test names. If T-2763 had to fix code, the reason names that fix too.

Before closing, check that `system-flow.yaml` on main still has `set_cost_of_delay: true`. If it does not, do not close: ask the operator on a thread on the story.

This task waits for T-2763 and T-2764, because the reason names their tests and any fix they made.

## Done when

- I-0284's status is `closed`, and its body ends with the reason.
- `design/issues/summary.md` no longer lists I-0284 as open.
- `flai check --strict` is clean for the story.

## Notes
