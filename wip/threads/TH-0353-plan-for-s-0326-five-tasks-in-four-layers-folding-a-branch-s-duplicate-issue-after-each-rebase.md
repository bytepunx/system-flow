---
id: TH-0353
title: "Plan for S-0326: five tasks in four layers, folding a branch's duplicate issue after each rebase"
anchor:
  path: wip/kanban/stories/S-0326-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md
  item: S-0326
status: open
participants: [planner-S-0326, orchestrator]
created: 2026-10-08T00:35:03Z
updated: 2026-10-08T00:35:17Z
---

# TH-0353 Plan for S-0326: five tasks in four layers, folding a branch's duplicate issue after each rebase

On wip/kanban/stories/S-0326-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md.

## Entries

### 2026-10-08T00:35:03Z planner-S-0326
S-0326 is planned: 15 touches, a forecast of 45m delivered by 2026-10-08T08:44:00Z, and a cost of delay of 25 USD a week. It stays a draft for you to finalize.

## Cause

`issues.RecordOnce` (`flai/internal/issues/record.go`) matches titles in the story's own checkout only. A close-out cannot see an issue another story branch opened under the same title. Once that branch is accepted and this one syncs, both files stand on the branch, and the next close-out bumps the lower ID, leaving the other a duplicate.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1319 ADR records the remedy, refining ADR-0085 and ADR-0098 | — | `design/adrs/`, `design/adrs/README.md` |
| 2 | T-1320 `issues.Fold`, with unit tests | T-1319 | `flai/internal/issues/fold.go`, `fold_test.go` |
| 2 | T-1322 design and users' guide | T-1319 | `design/system/continuous-improvement.md`, `design/system/flai-cli.md`, `docs/users/flai.md` |
| 3 | T-1321 fold after the rebase in sync, task done, and accept, with a test reproducing I-0112 | T-1320 | `flai/internal/storygit/sync.go`, `sync_test.go`, `flai/cmd/stream.go`, `flai/cmd/branch.go`, `flai/internal/taskdone/taskdone.go`, `flai/cmd/stream_sync_test.go` |
| 4 | T-1323 close I-0112 | T-1321, T-1322 | `design/issues/I-0112-…md`, `design/issues/summary.md` |

## Assumptions

- The remedy folds an open issue the branch added into the main branch's open issue of the same title once a rebase completes, and commits the fold on the branch. Fixing it at record time does not work: the close-out cannot bump a file its branch does not have. T-1319's ADR weighs this and may decide otherwise.
- The fold reaches all three rebase callers through `storygit.SyncOptions`, as `issues.Generated` does, so `storygit` stays free of the issues package.
- `design/adrs/` stays a folder touch because the ADR's number is unknown until `flai adr new` assigns it.
- The forecast is raised from flai's 20m to 45m for the ADR, the hook through three callers, and the two-branch test. S-0315, two tasks in the same area, took 17m.
- The cost of delay comes from flai's input (10m per cycle) unchanged.

## Related

TH-0346 proposes merging S-0323 and S-0325, two stories from the duplicate issues this story's cause produced. S-0297, issue front-matter conflicts on a bump, is in the same area but touches none of the same files. I propose no split, merge, or drop.

### 2026-10-08T00:35:17Z orchestrator (recommendation)
Recommendation: take the planner's remedy and approve the plan as it stands.

The remedy is that once a rebase completes, flai folds an open issue the branch added into the main branch's open issue of the same title, and commits the fold on the branch.

- **It removes I-0112's cause.** The record-time fix cannot work: a close-out cannot bump a file its branch does not have. Folding after the rebase is the first point where both files are visible.
- **It goes beyond ADR-0098.** ADR-0098 lets a sync regenerate one derived file, `summary.md`, and only when it is the sole conflict. Folding merges two source issue files and deletes one, after every completed rebase in sync, task done, and accept. That is a new automatic change to recorded data, so it refines ADR-0085 and ADR-0098 in a way only you can confirm. T-1319 records it.
- **The alternative is to leave folding to people.** flai would only warn on a duplicate title after a sync, and the story's agent or you would fold it by hand. That is safer, but it keeps the duplicates that I-0109/I-0111 and S-0323/S-0325 show.

The plan fits the story:

- T-1321 reproduces I-0112 across two branches (criterion 1).
- T-1323 closes I-0112 (criterion 2).
- T-1320 and T-1322 share no path.
- The 45m forecast and the 25 USD a week value stand.

S-0326 stays a draft until this thread settles.

Source: design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md § Decision
