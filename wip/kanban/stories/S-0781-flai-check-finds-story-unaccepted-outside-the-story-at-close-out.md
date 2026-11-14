---
id: S-0781
type: story
nature: improvement
title: "flai check finds `story.unaccepted` outside the story at close-out"
status: in-progress
owner: alex
created: 2026-11-13T22:46:13Z
updated: 2026-11-13T23:45:09Z
transitions:
  - to: ready
    at: 2026-11-13T22:58:27Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-13T22:58:54Z
    by: agent-S-0781
blocked:
  - from: 2026-11-13T23:45:09Z
    reason: awaiting the operator on TH-1190
tags: [flai, check]
topics: [cli, conventions, template, testing]
touches: [flai/cmd/cmd_test.go, flai/internal/mcpserver/server_test.go, flai/cmd/story_start_test.go, design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/issues/I-0371-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md, design/issues/summary.md, flai/cmd/verify_test.go, flai/internal/verify/story_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2754
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 2
      test_runs: 1
      work: 35
  models:
    - model: claude-opus-5-5
      input: 262
      output: 79201
      cache_read: 13200979
      cache_write: 589139
      cost: 7.9584
  strategic:
    - kind: planner
      seconds: 280
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 73
          output: 24030
          cache_read: 5884489
          cache_write: 149572
          cost: 2.8544
    - kind: orchestrator
      seconds: 223
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 67
          output: 601
          cache_read: 2585163
          cache_write: 384148
          cost: 0.877
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: orchestrator
    at: 2026-11-13T22:56:05Z
  value: 2.5
  by: planner-S-0781
  at: 2026-11-13T22:57:46Z
forecast:
  duration: 1h
  delivery: 2026-11-14T06:24:00Z
  basis: "flai forecast's 36m at size 20 (2 criteria, 18 touches), raised to 1h: S-0761, the same scoped-check fix with a one-command repair, used 67m of agent time against its 33m forecast, and this story adds a test-harness fix with a real-git reproduction; delivery moved by the 24m added."
  by: planner-S-0781
  at: 2026-11-13T22:54:53Z
finalized:
  by: orchestrator
  at: 2026-11-13T22:58:07Z
---
# S-0781 flai check finds `story.unaccepted` outside the story at close-out

## Goal

This story remediates [I-0371](../../../design/issues/I-0371-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md), "flai check finds `story.unaccepted` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0371 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0371 is closed with `flai issue close I-0371 --reason` saying what fixed it

## Tasks
- T-3343 Stop a test's git from walking up out of its temp folder into the repository that holds it
- T-3345 Remove the four branches and worktrees S-0777's test run left in the main repository, once the operator confirms
- T-3347 Record in an ADR whether a close-out leaves out story.unaccepted outside the story
- T-3348 Leave story.unaccepted outside the story out of a check scoped to a story, with a test that reproduces I-0371
- T-3349 Say in the design, the user guide, and the work-management convention that a close-out leaves out story.unaccepted
- T-3350 Close I-0371 saying what fixed it

## Notes

### Planning

Cause, from I-0371's five instances (close-outs of S-0778, S-0730, S-0777, S-0780, and S-0779 on 2026-11-13): each names the same four branches, `story/S-0001`, `story/S-0002`, `story/S-0003`, and `story/S-0005`. They are not those stories' branches. All four point at 5c3b47d0c, S-0777's commit "test: [S-0777] flai/cmd's kluster fixtures resolve the project root". Their worktrees, now prunable, were under `.flai-cache/worktrees/S-0777/bin/S-0777/real/<Test…>/002/.flai-cache/worktrees/`. S-0777 ran the Go tests by hand with `TMPDIR` inside its worktree (`wip/archive/agents/S-0777.md`). `tempProject` (`flai/cmd/workitems_test.go`) runs no `git init`, so `storygit.InWorkTree` walked up into S-0777's worktree, and `storygit.Open` added real `story/S-nnnn` branches to the main repository. The same branches make `make flai-test` fail on this host (I-0370).

Proposed fix, in three parts:

1. Stop the leak: the `TestMain` of `flai/cmd` and `flai/internal/mcpserver` sets `GIT_CEILING_DIRECTORIES` to the temp root, with a test that reproduces the leak (T-3343).
2. Remove the four branches and prune their worktrees, once the operator confirms, since deleting a ref is hard to reverse (T-3345).
3. Leave `story.unaccepted` outside the story out of a check scoped to a story, as S-0761 (ADR-0296) did for `item.blocked-closed`, decided in an ADR that weighs the root fix alone (T-3347, T-3348, T-3349). No story branch clears either of the rule's cases: only `flai accept`, or the operator's merge or deletion of the branch, in the main checkout.

Tasks and layers:

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-3343 The git ceiling in the tests, with a reproduction | nothing |
| 1 | T-3345 Remove the four branches, on the operator's word | nothing |
| 1 | T-3347 The ADR | nothing |
| 2 | T-3348 The scoped check leaves the finding out, with tests | T-3347, for its decision and number |
| 2 | T-3349 The design, guide, convention, and changelog say so | T-3347, for its decision and number |
| 3 | T-3350 Close I-0371 | T-3343, T-3345, T-3348, T-3349: the reason names each fix |

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/cmd/cmd_test.go`, `flai/internal/mcpserver/server_test.go` | Layout: the explorer found the two `TestMain`s of the packages whose tests leaked (`TestInTheSessionOfAgentA`, `TestTheOrchestratorSetsCostOfDelayInputsByItsPermissions`); `flai/internal/guard`'s copy of the latter runs no git |
| `flai/cmd/story_start_test.go` | Layout: the story-start tests that assume a temp project is "not a git work tree", where the reproduction fits |
| `flai/internal/check/scope.go`, `flai/internal/check/scope_test.go`, `flai/cmd/check.go`, `flai/cmd/check_test.go`, `docs/users/flai-reference.md` | Design and layout: the files S-0761's fix commit changed for `item.blocked-closed` |
| `design/adrs`, `design/adrs/README.md` | Design: S-0366, S-0555, and S-0761 each needed an ADR for the same kind of fix |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md`, `design/system/flai-cli.md`, `design/system/continuous-improvement.md`, `docs/users/flai.md` | Design: each lists what a scoped check leaves out, as S-0761 changed them; `docs/users/flai.md` also shows a `story.unaccepted` as its example of a recorded outside finding |
| `design/issues/I-0371-…md`, `design/issues/summary.md` | Criterion 2: `flai issue close I-0371` |

The folder touch `design/adrs` stays: the ADR's file name comes from the title T-3347 decides, so no task can name it yet. `flai touches suggest` found no co-change above 8% on these paths (`docs/operators/index.md` was the highest), so none was added. `flai/internal/check/check_items.go`, the rule itself, does not change. T-3345 touches no file: it removes git refs in the main checkout.

No `after`: the story reads no output of an open story. S-0777 and S-0761, whose work it builds on, are done.

Forecast 1h, adjusted from `flai forecast`'s 36m at size 20 (2 criteria, 18 touches). S-0761, the same scoped-check fix with a one-command repair task, used 67m of agent time against its 33m forecast. This story adds a test-harness fix with a real-git reproduction. Delivery moves by the 24m added, to 2026-11-14T06:24:00Z.

Cost of delay 2.50 USD a week stands as `flai cod` gives it. The input, `time_lost_per_cycle: 1m`, is the one recommended on TH-1188 and set by the orchestrator, as TH-1137 set it for S-0761 on the same kind of finding. At 150 USD an hour and one 168h cycle a week, that is 2.50 USD.
