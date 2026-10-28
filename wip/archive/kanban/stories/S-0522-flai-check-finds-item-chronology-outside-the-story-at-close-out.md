---
id: S-0522
type: story
nature: improvement
title: "flai check finds `item.chronology` outside the story at close-out"
status: done
owner: alex
created: 2026-10-28T12:53:36Z
updated: 2026-10-28T15:18:40Z
transitions:
  - to: ready
    at: 2026-10-28T14:24:16Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T14:24:22Z
    by: agent-S-0522
  - to: review
    at: 2026-10-28T14:41:18Z
    by: agent-S-0522
  - to: done
    at: 2026-10-28T15:18:40Z
    by: orchestrator
tags: [cli]
touches: [design/adrs, design/adrs/README.md, design/system/work-hierarchy.md, flai/internal/workitem/repo.go, flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/widen_test.go, design/issues/I-0183-flai-check-finds-item-chronology-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1024
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 1
      hand_edits: 2
      work: 28
  models:
    - model: claude-opus-5-5
      input: 182
      output: 73810
      cache_read: 7873193
      cache_write: 419329
      cost: 5.7776
  strategic:
    - kind: orchestrator
      seconds: 1397
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 113
          output: 1921
          cache_read: 37596557
          cache_write: 34285
          cost: 8.2612
        - model: claude-sonnet-5-5
          input: 12
          output: 82
          cache_read: 166540
          cache_write: 43780
          cost: 0.1438
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: orchestrator
    at: 2026-10-28T14:16:06Z
  value: 37.5
  by: planner-S-0522
  at: 2026-10-28T14:22:36Z
forecast:
  duration: 45m
  delivery: 2026-10-29T03:57:00Z
  basis: "flai forecast S-0522 gives 23m from size 15 (2 criteria, 13 touches; median 89 s per unit over 98 done improvement stories on claude-opus-5-5, large band), raised to 45m for an ADR, seven tests with controlled clocks across three packages, serve's among them, and a close-out of a flai/** change that runs the integration and smoke tiers in full; delivery is flai's play-out start, 2026-10-29T03:35Z less 23m, plus 45m."
  by: planner-S-0522
  at: 2026-10-28T14:23:28Z
finalized:
  by: alex
  at: 2026-10-28T13:01:26Z
---
# S-0522 flai check finds `item.chronology` outside the story at close-out

## Goal

This story remediates [I-0183](../../../design/issues/I-0183-flai-check-finds-item-chronology-outside-the-story-at-close-out.md), "flai check finds `item.chronology` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0183 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0183 is closed with `flai issue close I-0183 --reason` saying what fixed it

## Tasks
- T-2237 An ADR refining ADR-0184 records that an item write never leaves updated before its last transition, and work-hierarchy.md states the rule
- T-2238 Repo.Save keeps updated at or after the last transition and created, and tests reproduce I-0183 with a stale clock after a move and with a same-second tie in the merge
- T-2239 The replanner's forecast write, the merge queue's send-back, and the touches widening take their clock after reading the item they stamp
- T-2240 The main checkout's items pass item.chronology, and I-0183 is closed with the fix named

## Notes

### Planning

Planned by planner-S-0522 on 2026-10-28 (the plan's thread names the tasks, their layers, and the assumptions; TH-0692 settled the cost of delay input).

The cause, from the instance and the code. `item.chronology` is `history` in `flai/internal/check/check.go`: an error when `updated` is before the last transition's `at` or before `created`, to the second. A move stamps the transition and `updated` from one clock (`workitem/rules.go`), but `Repo.Save` (`workitem/repo.go`, `itemlock.go`) keeps no order between them, and three writers take their clock before they read the item: the replanner's forecast writes (`serve/replan.go`, one `now` for a loop that runs a `flai check` per edit and spans minutes), the merge queue's send-back (`mergequeue/bisect.go`), and the touches widening (`itemedit/widen.go`). A write that reads the item after another process moved it, and stamps `updated` from a clock read before the move, leaves `updated` behind the transition. The instance: S-0482 was moved to `ready` at 2026-10-21T12:14:19Z by the orchestrator; S-0473's close-out at 12:18:49Z read the main checkout's file and found `updated` before that; the writes between were not committed, so git does not name the writer, and S-0482's forecast kept its planner's stamp, so the replanner is not shown to be it. ADR-0184's merge landed at 12:40Z, after the occurrence; it keeps the order no better, since `updated` takes the later value only when both sides changed it. The check scoped to a story keeps the finding as an outside note and the close-out records it in an issue, which is how I-0183 came to be.

The remediation proposed: `Repo.Save` raises `updated` to the later of the last transition's `at` and `created` when the stamp it is about to write is earlier (T-2238, the invariant), the three writers read their clock after the item (T-2239, the honest stamp), an ADR refining ADR-0184 and a rule in `work-hierarchy.md` record it (T-2237), and I-0183 is closed once the main checkout's items pass (T-2240). Not proposed: leaving `item.chronology` outside the story out of a scoped check, as ADR-0122 and its kin do for findings only the main checkout clears; the finding is a symptom of a write flai made wrong, and with the order kept at save it cannot arise from flai's writes, so a hand edit that makes it should still be seen. The operator may decide otherwise on the plan's thread.

Touches. `flai touches suggest S-0522` names no path on its own, since the story declared none; seeded with `itemlock.go`, `repo.go`, `check.go`, `scope.go`, and `replan.go` (82 of 2483 commits) it lists `check_test.go` (29%), `flai-cli.md` (26%), `docs/users/flai.md` (23%), `cmd/check.go` (17%), `scope_test.go` (15%), `flai-reference.md` (13%), `cmd/root.go`, `manifest.go`, `docs/operators/settings.md`, `cmd/check_test.go`, `design/issues/summary.md` (11%), `workitem/item.go` (11%), and `work-hierarchy.md` (9%). Kept from it (co-change): `design/issues/summary.md` (rewritten by `flai issue close`) and `design/system/work-hierarchy.md` (the schema's rules on `updated`). Left out: the check's files and their tests (the check does not change; its fixture already covers the finding), `flai-cli.md`, `docs/users/flai.md`, `flai-reference.md`, `docs/operators/settings.md` (no command, flag, or setting changes), `cmd/root.go`, `manifest.go`, `workitem/item.go` (the stamp is raised at save, not in the item's marshal). The rest are the tasks' (layout):

- `flai/internal/workitem/repo.go`, `itemlock.go`, `itemlock_test.go` (layout): `Save`, `mergeSaved`, and the tests beside `TestSaveKeepsAMoveSavedSinceTheCopyWasRead`. T-2238.
- `flai/internal/serve/replan.go`, `replan_test.go`, `flai/internal/mergequeue/bisect.go`, `bisect_test.go`, `flai/internal/itemedit/widen.go`, `widen_test.go` (layout): the three writers and their tests. T-2239.
- `design/adrs/README.md` (layout): the index row `flai adr new` writes. `design/system/work-hierarchy.md` (design): § Front matter schema, Rules. T-2237.
- `design/issues/I-0183-….md`, `design/issues/summary.md` (criterion 2): `flai issue close` writes both; `design/issues` is in the manifest's shared paths and holds nothing (ADR-0096). T-2240.
- The folder `design/adrs` (layout, folder kept): `flai adr new` in T-2237 names the ADR's file by its number and title, which no task can name; the folder is kept for that file alone, is in the manifest's shared paths, so it holds no ready story, and `flai task done` widens the touches with the committed path (ADR-0144).

Overlaps on 2026-10-28: `flai/internal/serve/replan.go` is touched by S-0419 (ready, E-0024) and `design/system/work-hierarchy.md` by S-0421 (ready, E-0024), so whichever of the two starts first holds the other, and the pull order settles which; no story in progress touches any of the paths. No `after`: the story needs no other story's change on main.

Forecast. `flai forecast S-0522` gave 3m with no touches (size 2) and gives 23m with them (size 15: 2 criteria, 13 touches; median 89 s per unit over 98 done improvement stories on claude-opus-5-5, large band), delivery 2026-10-29T03:35Z, 34th in the pull order. Raised to 45m: the story writes an ADR, seven tests with controlled clocks across three packages, serve's among them, and its close-out of a flai/** change runs the integration and smoke tiers in full. The delivery written is flai's play-out start, 03:35Z less 23m, plus 45m: 2026-10-29T03:57Z.

Cost of delay. `flai cod S-0522` gives 37.50 USD a week: 15m lost per 168h cycle at 150 USD an hour, from the input the orchestrator set on TH-0692 at my recommendation. Written as flai gives it: the issue has one instance, and the input is the operator's to raise if a stale stamp ever blocks an acceptance.

### Accepted by the orchestrator

- Verified: 79975cf6f9455b0b7d6df31a31e15fbd50f6e46b
- Batch: 20261028-150510
- Enqueued verified: eaa9d4f5086a1e69318d10a3d862a37b5391341e
- At: 2026-10-28T15:18:40Z

Verdict: pass at eaa9d4f5086a1e69318d10a3d862a37b5391341e
- 1: flai/internal/workitem/itemlock.go, flai/internal/workitem/repo.go, flai/internal/workitem/itemlock_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/widen_test.go, design/system/work-hierarchy.md
- 2: design/issues/I-0183-flai-check-finds-item-chronology-outside-the-story-at-close-out.md, design/issues/summary.md
