---
id: TH-1046
title: "Plan for S-0675: a verified batch lands over main's wip-only commits, and every drop is logged; five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0675-the-merge-queue-builds-one-story-batches-and-runs-three-verifications-per-story.md
  item: S-0675
status: open
participants: [planner-S-0675]
created: 2026-11-10T10:37:46Z
updated: 2026-11-10T10:38:37Z
---

# TH-1046 Plan for S-0675: a verified batch lands over main's wip-only commits, and every drop is logged; five tasks in three layers

On wip/kanban/stories/S-0675-the-merge-queue-builds-one-story-batches-and-runs-three-verifications-per-story.md.

## Entries

### 2026-11-10T10:37:46Z planner-S-0675
The plan for S-0675 is written: touches, a forecast of 1h30m, a cost of delay of 1050 USD a week, and five tasks in three layers. The evidence and the reasoning are under `### Planning` in the story's Notes.

## The cause

Most of I-0287's waste comes from dropped batches, not from batch size. 37 of the 94 batches that passed between 2026-10-21 and 2026-11-02 never landed. `mergequeue.DropIfMainMoved` drops a verified batch whenever main's head moved at all. In the drops sampled, flai's own `chore: replan forecasts after reordered` commits moved it, and those change only `wip/`. Each drop costs another run of `integration` and `smoke`.

## Proposed fix

- A verified batch survives main commits that change only `wip` paths none of its stories change. It lands on the new main without a second verification, as `flai verify`'s sync already allows (ADR-0135). Any other move of main still drops it.
- Every drop is logged as a `dropped` event with its reason, and `flai stats` counts it as `queue.weeks[].dropped`. The metrics change comes with an ADR.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-3073 ADR for both decisions | none | `design/adrs/README.md` (the story keeps `design/adrs` for the new file) |
| 2 | T-3074 pass-over in `DropIfMainMoved` and `Landable`, with the test that reproduces I-0287 | T-3073 | `flai/internal/mergequeue/land.go`, `land_test.go`, `flai/cmd/accept_queue_test.go` |
| 2 | T-3075 `dropped` event and its weekly count | T-3073 | `flai/internal/mergequeue/log.go`, `log_test.go`, `flai/internal/metrics/queue.go`, `queue_test.go`, `flai/cmd/stats.go` |
| 3 | T-3076 workflow.md, metrics.md, flai-cli.md, docs/users/flai.md, changelog | T-3074, T-3075 | those files |
| 3 | T-3077 close I-0287 | T-3074, T-3075 | the issue file, `design/issues/summary.md` |

## Assumptions

1. `merge_queue.coalesce` stays at 3 minutes. Stories reached review hours apart in W44, so a longer window would delay landings more than it would join stories. Say so if you want the window revisited in this story.
2. Passing over `wip`-only commits is safe because no tier result depends on `wip` files the batch's stories do not change. ADR-0135 rests on the same reasoning.
3. The story has no `after`. S-0698 also changes `flai/internal/mergequeue`, and the trial merge judges that overlap.
4. The forecast is raised from flai's 31m to 1h30m because of the ADR, the real-git tests, and the long `integration` and `smoke` tiers at close-out.

### 2026-11-10T10:38:37Z planner-S-0675
The explorer's read of the code and the main branch's reflog adds two findings that bear on the plan.

1. Some drops were not main moving. S-0541 was batched 16 times between 2026-10-29 14:21 and 15:12, and each batch passed. Only one of its drops coincides with a main commit, a replan at 14:44:56. The other 14 point to a refused landing: `Land` drops the batch when `flai accept --batch` fails (`flai/internal/mergequeue/land.go:351`). That cause cannot be confirmed, because the reason lived only in `state.json` and the serve journal, and neither kept it. T-3074 does not fix these drops. T-3075 makes the next ones visible, since each `dropped` line carries its reason. If that log shows refused landings recurring, they need a story of their own.
2. The 3-minute coalesce window has applied only since 2026-11-01. Before that, batches were built 4 to 13 s after the enqueue. Since then, the orchestrator's enqueues have come 30 to 60 minutes apart. One-story batches therefore follow from how stories arrive, which supports assumption 1.

T-3074 should also update `TestALandingKeptOnMainIsDroppedWhenMainMoved` (`land_test.go:867`) alongside the two tests it names. The plan otherwise stands.
