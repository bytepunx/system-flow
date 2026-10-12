---
id: S-0371
type: story
nature: remediation
title: flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's
status: backlog
owner: alex
created: 2026-10-12T16:52:47Z
updated: 2026-10-12T20:03:53Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/cmd/branch.go, flai/cmd/accept.go, flai/cmd/accept_undone_test.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 366
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 50
          output: 888
          cache_read: 15133318
          cache_write: 12571
          cost: 3.8181
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-12T16:52:47Z
  value: 25
  by: planner-S-0371
  at: 2026-10-12T17:16:59Z
forecast:
  duration: 45m
  delivery: 2026-10-13T00:59:00Z
  basis: "Its own forecast of 45m; 16th in the pull order with an in-progress limit of 5, behind S-0238, S-0360, S-0370, S-0356, S-0354, S-0357, S-0351, S-0352, S-0353, S-0239, S-0241, S-0304, S-0305, S-0306, S-0355 and S-0358."
  by: flai
  at: 2026-10-12T20:03:53Z
---
# S-0371 flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's

## Goal

This story remediates [I-0124](../../../design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md), "flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0124 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0124 is closed with `flai issue close I-0124 --reason` saying what fixed it

## Tasks
- T-1476 ADR: a rebase of a story branch that removes lines main gained since its old branch point is reported by flai stream sync and refused by flai accept
- T-1477 storygit.Rebase reports each path where the rebased branch removes lines main gained since its old branch point, with a test reproducing I-0124
- T-1478 flai stream sync and flai task done warn on each path where the rebased branch undoes main's change, and say how to restore it
- T-1479 flai accept refuses, merging nothing, a story whose rebase undoes main's change to a path, naming the path and the restore
- T-1480 The design and user docs say flai stream sync warns and flai accept refuses when a rebased branch undoes main's change
- T-1481 Close I-0124 as fixed by the rebase check that reports and refuses a branch undoing main's change

## Notes

Cost of delay inputs set by flai from I-0124. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T09:00:31Z, 4.3 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from the instance and the issue's suggestion. `storygit.Rebase` reads the branch's merge base with main before `git rebase`. After a clean rebase, it reports each path where the branch's net change removes lines main added since that old base, as `undoes_main`. `flai stream sync` and `flai task done` warn on each such path, inside the claim or outside it, and say how to restore it. `flai accept` refuses such a branch, merging nothing. T-1476 records this first as an ADR refining ADR-0069 and ADR-0098.

- Cause: git drops a commit that becomes empty on a rebase (S-0232's change, which S-0340 had made on main), then replays the revert after it cleanly. So no conflict shows, and the branch's net change before the rebase is empty. The check must therefore run after the rebase, on the old base read before it.
- Why lines, not the outside-touches note: `outsideClaim` in `flai/internal/storygit/sync.go` leaves out every path in the claim, so a note there would miss the same revert of a claimed path.
- Why accept refuses: it runs `storygit.Rebase` alone (`flai/cmd/branch.go`, `syncStoryBranch`), with no checks. A warning there would merge main's loss as I-0124's sync did.

Touches, each a file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/storygit/undone.go` | layout | the new check, beside `diff.go` |
| `flai/internal/storygit/undone_test.go` | layout | its unit tests |
| `flai/internal/storygit/sync.go` | layout, co-change | `Rebase` reads the old base and calls the check; `SyncResult` gains `undoes_main` |
| `flai/internal/storygit/sync_test.go` | co-change | reproduces I-0124 through `Sync` |
| `flai/cmd/stream_sync.go` | layout | `printSyncChecks`, shared by `flai stream sync` and `flai task done` |
| `flai/cmd/stream_sync_test.go` | co-change | the warning's wording |
| `flai/cmd/stream.go` | co-change | the sync's help |
| `flai/cmd/branch.go` | layout | `syncStoryBranch` and `mergeStoryBranch`, accept's rebase |
| `flai/cmd/accept.go` | co-change | accept's help |
| `flai/cmd/accept_undone_test.go` | layout | the refusal, after `accept_conflict_test.go` |
| `docs/users/flai-reference.md` | co-change | generated from the helps |
| `design/adrs/README.md` | design | the ADR index; `flai adr new` names the ADR's file, and `flai task done` widens the touches with it |
| `design/system/flai-cli.md` | co-change, design | the sync, task done, and accept rows |
| `design/system/workflow.md` | co-change, design | § Branches and collisions (ADR-0019) |
| `docs/users/flai.md` | co-change, design | § Story branches, § Closing a task, and acceptance |
| `design/issues/I-0124-…` | design | `flai issue close` writes it (criterion 2) |
| `design/issues/summary.md` | co-change | `flai issue close` regenerates it |

The story declared no touches, so `flai touches suggest` had nothing to start from. Seeded with `sync.go`, `stream_sync.go`, and `branch.go`, it listed `stream_sync_test.go` (50%), `stream.go` (46%), `flai-cli.md`, `docs/users/flai.md`, `accept.go`, and `sync_test.go` (23%), and `flai-reference.md` and `accept_conflict_test.go` (12%). All are predicted except `accept_conflict_test.go`: the refusal's test goes in a new file. The rest, at 8%, are not predicted: `flai/cmd/dashboard.go`, `check.go`, `release.go`, the flaiover docs, and `docs/operators/index.md`. No dashboard view, setting, or release step changes. `design/conventions/work-management.md` is not predicted either: the warning and the refusal name what to do, so no rule changes.

Overlaps: S-0367, in progress, claims `flai/cmd` and `flai/internal/storygit/sync.go`. This story is held while S-0367 is in progress; that is expected, since both change `sync.go`. `docs/users/flai.md`, `flai-reference.md`, `flai-cli.md`, `design/adrs`, and `design/issues` are in the manifest's shared paths and hold nothing.

Figures:

- Forecast 45m, delivery 2026-10-13T02:21Z. `flai forecast` gave 26m: the median of 82 s per unit of size over 23 done large-band remediation stories on claude-opus-5-5, times size 19 (2 criteria, 17 touches). It is raised because the work is a new line-level git check, two callers that behave differently, and git-fixture tests in two packages. S-0370's single filter forecasts 21m, and S-0368's and S-0369's planners raised theirs to 45m for comparable scope. The delivery moves by the 19m added to flai's 02:02Z.
- Cost of delay 25.00 USD a week: `flai cod` from flai's input, 10m lost per 168h cycle at 150 USD an hour. It is not adjusted. The issue has one instance, but an undetected loss of main's change can reach main through an acceptance, so it is not lowered either.
