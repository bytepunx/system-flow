---
id: S-0371
type: story
nature: remediation
title: flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's
status: done
owner: alex
created: 2026-10-12T16:52:47Z
updated: 2026-10-14T23:09:12Z
transitions:
  - to: ready
    at: 2026-10-14T20:55:53Z
    by: alex
  - to: in-progress
    at: 2026-10-14T21:33:36Z
    by: agent-S-0371
  - to: review
    at: 2026-10-14T23:02:14Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T23:09:12Z
    by: alex
tags: [cli]
topics: [cli]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_restore.go, flai/cmd/stream_restore_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/preview/accept.go, flai/cmd/branch.go, flai/cmd/accept.go, flai/cmd/accept_undone_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, docs/users/flai-reference.md, docs/operators/settings.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/summary.md, design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md, flai/cmd/move.go, flai/internal/preview/orchestrator.go, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5333
  priced_by: harness
  turns:
    - day: 2026-10-14
      ceremony: 2
      hand_edits: 2
      work: 87
  models:
    - model: claude-opus-5-5
      input: 578
      output: 215863
      cache_read: 34057181
      cache_write: 1596318
      cost: 19.8956
  strategic:
    - kind: planner
      seconds: 435
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 90
          output: 34187
          cache_read: 4794128
          cache_write: 168461
          cost: 0.3316
        - model: claude-opus-5-5
          input: 94
          output: 43697
          cache_read: 8056525
          cache_write: 220584
          cost: 4.2503
    - kind: orchestrator
      seconds: 3862
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 256
          output: 4031
          cache_read: 41817550
          cache_write: 87424
          cost: 10.5292
        - model: claude-sonnet-5-5
          input: 14
          output: 78
          cache_read: 295074
          cache_write: 77557
          cost: 0.305
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-12T16:52:47Z
  value: 25
  by: planner-S-0371
  at: 2026-10-12T17:16:59Z
forecast:
  duration: 1h15m
  delivery: 2026-10-14T22:35:00Z
  basis: "Its own forecast of 1h15m; 2nd in the pull order with an in-progress limit of 5, behind S-0304, S-0372 and S-0374."
  by: flai
  at: 2026-10-14T21:15:26Z
finalized:
  by: alex
  at: 2026-10-14T20:55:48Z
---
# S-0371 flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's

## Goal

This story remediates [I-0124](../../../design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md), "flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0124 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0124 is closed with `flai issue close I-0124 --reason` saying what fixed it

## Tasks
- T-1476 ADR: a rebased story branch that undoes main's change to a path is an acceptance blocker the operator resolves per path, by restoring it, keeping the removal with a reason, or sending the story back
- T-1477 storygit.Rebase reports and records each path where the rebased branch removes lines main gained since its old branch point, with a test reproducing I-0124
- T-1478 flai stream sync and flai task done warn on each path where the rebased branch undoes main's change, and name flai stream restore
- T-1479 flai accept's preview lists each path the branch undoes as a blocker with its lines as data, and accept refuses it unless --keep-removal names it with a reason, which the operator alone may give
- T-1480 The design and user docs say how flai finds a rebased branch that undoes main's change, and how the operator restores the path, keeps the removal, or sends the story back
- T-1481 Close I-0124 as fixed by the rebase check whose findings block acceptance until the operator restores the path or keeps the removal
- T-1574 flai stream restore S-nnnn --path commits main's version of a path the branch undoes onto the story branch, and clears it from the record
- T-1575 The host write stream.restore runs flai stream restore, and accept.run takes keep_removals, each a path and a reason, for the operator's acceptance only

## Notes

Cost of delay inputs set by flai from I-0124. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T09:00:31Z, 4.3 days before this story; under one cycle counts as one).

### Planning

Replanned on 2026-10-14 to the shape the operator confirmed on TH-0420 (2026-10-14T18:14:21Z): the operator resolves a branch that undoes main's change from the dashboard, path by path, and flai does the git work. The first plan's `git checkout` advice is gone.

The remedy:

- **Check.** `storygit.Rebase` reads the branch's old merge base before `git rebase`. After a clean rebase, it reports each path where the branch's net change removes lines main added since that base, as `undoes_main`.
- **Cause.** Git drops the commit that became empty (S-0232's change, which S-0340 had made on main) and replays the revert after it cleanly. No conflict shows, and the branch's net change before the rebase is empty. So the check runs after the rebase.
- **Record.** The next rebase starts from a base that already holds main's change and cannot find the path again. The close-out syncs before review, so the acceptance's rebase is usually a no-op. Each finding is therefore recorded per story under `.flai-cache/` when the sync makes it, after the precedent of `.flai-cache/verify/<story>.json`. A recorded path that the branch no longer undoes drops out.
- **Warning at sync.** `flai stream sync` and `flai task done` warn and name `flai stream restore`. Neither fails.
- **Blocker at acceptance.** `preview.AcceptWith` reads the record, so `flai accept --dry-run` and the host read `accept.preview` both give each path as a blocker and as data, `undoes_main`, with main's lines that the branch removes. `flai accept` refuses on it, merging nothing. With `--json` it exits 4 with the paths as data, as the markdown refusal does.
- **Ways out.** `flai stream restore S-nnnn --path <p>` commits main's version onto the branch. `flai accept --keep-removal <p> --reason` keeps the removal, records the reason in the Notes, and is refused with `--by orchestrator`. Sending back is the existing move to in-progress.
- **Host writes.** `stream.restore` and `accept.run`'s `keep_removals` expose the two ways out to the dashboard.
- **Dashboard.** The review page's cards and the board's badge are a separate draft story, after this one (see the plan thread).

Why lines, not the outside-touches note: `outsideClaim` in `flai/internal/storygit/sync.go` leaves out every path in the claim. A note there would miss the same revert of a claimed path.

Touches, each a file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/storygit/undone.go` | declared, layout | the check and the record, beside `diff.go` |
| `flai/internal/storygit/undone_test.go` | declared, layout | their unit tests |
| `flai/internal/storygit/sync.go` | declared, co-change | `Rebase` reads the old base, calls the check, and records; `SyncResult` gains `undoes_main` |
| `flai/internal/storygit/sync_test.go` | declared, co-change | reproduces I-0124 through `Sync` |
| `flai/cmd/stream_restore.go` | layout | the new `flai stream restore`, after `stream_diff.go` |
| `flai/cmd/stream_restore_test.go` | layout | its tests |
| `flai/cmd/stream_sync.go` | declared, layout | `printSyncChecks`, shared by `flai stream sync` and `flai task done` |
| `flai/cmd/stream_sync_test.go` | declared, co-change | the warning's wording |
| `flai/cmd/stream.go` | declared, co-change | registers `restore`; the sync's help |
| `flai/internal/preview/accept.go` | layout | `Acceptance.UndoesMain`, the record read in `AcceptWith`, and the blockers |
| `flai/cmd/branch.go` | declared, layout | `syncStoryBranch` and `mergeStoryBranch`, the acceptance's own rebase |
| `flai/cmd/accept.go` | declared, co-change | `--keep-removal` and `--reason`, the Notes, and the refusal's data |
| `flai/cmd/accept_undone_test.go` | declared, layout | the blocker, the refusal, and the keep, after `accept_conflict_test.go` |
| `flai/internal/hostapi/writes.go` | co-change, layout | `stream.restore`, and `accept.run`'s `keep_removals` |
| `flai/internal/hostapi/writes_test.go` | co-change | `TestEveryWriteIsCovered` needs both tables' entries |
| `flaiover/src/lib/server/agent.ts` | layout | `REQUIRED_METHODS`: `TestTheDashboardAsksForExactlyWhatFlaiOffers` fails without `stream.restore` |
| `docs/users/flai-reference.md` | declared, co-change | generated from the helps |
| `docs/operators/settings.md` | co-change | its flag index is generated with the reference |
| `design/adrs/README.md` | declared, design | the ADR index; `flai task done` widens the touches with the ADR's own file |
| `design/system/flai-cli.md` | declared, co-change, design | the sync, task done, restore, accept, and hostapi rows |
| `design/system/workflow.md` | declared, design | § Branches and collisions (ADR-0019) |
| `design/system/strategic-agents.md` | design | § Accepting a story: the orchestrator never keeps a removal |
| `docs/users/flai.md` | declared, co-change, design | § Story branches, § Closing a task, § Accept and release, and `flai hostapi` |
| `design/issues/I-0124-…` | declared, design | `flai issue close` writes it (criterion 2) |
| `design/issues/summary.md` | declared, co-change | `flai issue close` regenerates it |

All seventeen declared touches are kept. `flai touches suggest S-0371` over them listed `docs/operators/index.md` (11%), `design/system/flaiover-dashboard.md` (11%), `docs/operators/settings.md` (10%), and `docs/users/flaiover.md` (9%). `settings.md` is predicted, because `make flai-reference` writes its flag index. The flaiover documents belong to the dashboard story, and no setting is added, so `docs/operators/index.md` is not predicted. `flai/internal/hostapi/writes.go` and `writes_test.go` (3%) are predicted for the host writes. `design/conventions/work-management.md` (4%) is not predicted: the warning and the blocker name what to do, so no rule changes.

Overlaps: none in progress claims these paths today (`who_touches` on `sync.go` and `accept.go`). `docs/users/flai.md`, `flai-reference.md`, `flai-cli.md`, `design/adrs`, and `design/issues` are in the manifest's shared paths and hold nothing. `flaiover/src/lib/server/agent.ts` holds any flaiover story that touches it while this one is in progress.

Figures:

- **Forecast 1h15m**, delivery 2026-10-14T22:20Z. `flai forecast` gave 40m: the median of 88 s per unit of size over 28 done large-band remediation stories on claude-opus-5-5, times size 27 (2 criteria, 25 touches). I raised it because the scope is now eight tasks in six layers:
  - a new command and a per-story record
  - a preview field and a refusal with an override
  - two host writes
  - git-fixture tests in three packages

  The first plan's six tasks were forecast at 45m. S-0387's five tasks and S-0372's seven were raised to 1h15m and 1h. The delivery moves by the 35m added to flai's 21:45Z.
- **Cost of delay 25.00 USD a week**, from `flai cod` on flai's input: 10m lost per 168h cycle at 150 USD an hour. It is not adjusted. The issue has one instance, but an undetected loss of main's change can reach main through an acceptance, so it is not lowered either.
