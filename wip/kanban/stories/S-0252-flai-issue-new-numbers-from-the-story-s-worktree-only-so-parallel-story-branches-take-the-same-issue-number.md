---
id: S-0252
type: story
nature: remediation
title: flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number
status: ready
owner: alex
created: 2026-10-03T18:33:16Z
updated: 2026-10-04T23:56:22Z
transitions:
  - to: ready
    at: 2026-10-04T21:42:51Z
    by: alex
tags: []
topics: [cli]
touches: [flai/internal/storygit, flai/internal/issues, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/mcpserver/issues_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/issues/I-0065-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: planner-S-0252
    at: 2026-10-04T23:56:01Z
  value: 25
  by: planner-S-0252
  at: 2026-10-04T23:56:18Z
forecast:
  duration: 40m
  delivery: 2026-10-05T00:42:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 3, with nothing ahead of it."
  by: flai
  at: 2026-10-04T23:56:22Z
---
# S-0252 flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number

## Goal

This story remediates [I-0065](../../../design/issues/I-0065-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md), "flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0065 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0065 is closed with `flai issue close I-0065 --reason` saying what fixed it

## Tasks
- T-0831 storygit lists the file names under a design folder on main, every story worktree, and every story branch
- T-0832 flai issue new numbers past every issue on main, every story worktree, and every story branch
- T-0833 The design and the user guide say how issues are numbered, and I-0065 is closed

## Notes

### Planning

Proposed solution, from I-0065's three instances: `issues.NextID` (`flai/internal/issues/issues.go`) globs only the current checkout's `design/issues`. It should number one past the highest I-number found in four places: the main checkout, every story worktree (which catches issues not yet committed), every local `story/*` branch (read with `git ls-tree`), and the main branch. The lookup goes in a `storygit` helper that takes a folder, so that S-0245 (I-0063, the same defect for ADR numbers) can reuse it for `design/adrs`.

Touches, none declared before:

- `flai/internal/storygit`: layout. It holds the story-branch helpers, and the new helper goes there.
- `flai/internal/issues`: layout. `NextID` and `New` live there.
- `flai/cmd/issue.go`, `flai/cmd/issue_test.go`: co-change (6 of 20 commits for the test) and layout. These hold the help text and the reproduction test.
- `flai/internal/mcpserver/issues_test.go`: co-change (3 of 20). It calls `issues.New`, so it changes if the signature does.
- `docs/users/flai-reference.md`: co-change (4 of 20). It is regenerated from the help text.
- `docs/users/flai.md`: co-change (9 of 20). Its issues section describes `flai issue new`.
- `design/system/flai-cli.md`, `design/system/continuous-improvement.md`: co-change (6 and 2 of 20) and design. These describe how issues are written and where they live.
- `design/issues/I-0065-…md`, `design/issues/summary.md`: the second criterion closes I-0065, and closing regenerates the summary.

Forecast: flai gave 12m (55 s per unit of size times size 13). I raised it to 40m. The work adds a git helper with a repository fixture of branches and worktrees, a reproduction test, three documents, and a reference regeneration. S-0181, S-0203 and S-0179, comparable flai remediations with tests and docs, took 23, 26 and 35 minutes. The delivery of 2026-10-05T00:30Z allows for S-0252 being held behind S-0248, which also touches `design/system/flai-cli.md` and `design/issues/summary.md`.

Cost of delay: 25.00 USD a week, as `flai cod` works it out from the operator's input of `time_lost_per_cycle: 10m` (TH-0110): 10m per 168h cycle at 150 USD an hour. It stands unadjusted. It matches S-0258 and S-0262, sibling remediations whose issues cost the same, so the pull order compares them fairly. The issue's instances came closer together than one a week, but the operator chose the figure knowing the instances.

Topics: `cli` is added, since every task changes flai or the CLI's design.
