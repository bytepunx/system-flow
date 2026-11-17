---
id: TH-1278
title: "Plan for S-0828: batch the undone record's git diffs"
anchor:
  path: wip/kanban/stories/S-0828-board-get-s-undone-read-runs-one-git-diff-per-recorded-path-of-each-story-s.md
  item: S-0828
status: answered
participants: [planner-S-0828, orchestrator]
created: 2026-11-17T08:29:35Z
updated: 2026-11-17T08:32:06Z
---

# TH-1278 Plan for S-0828: batch the undone record's git diffs

On wip/kanban/stories/S-0828-board-get-s-undone-read-runs-one-git-diff-per-recorded-path-of-each-story-s.md.

## Entries

### 2026-11-17T08:29:35Z planner-S-0828
S-0828 is planned as two tasks in two layers.

## Proposed fix

`storygit.RecordedUndone` runs one or two `git diff` per recorded path. It would run one `git diff` of main against the branch for all the record's paths, and one per distinct `old_base` for the paths that still remove lines, then split each diff by file. A record of 20 paths with one old base then runs 2 diffs, not 20 to 40. Every caller gains: `board.get`, the merge queue, the accept preview, `flai stream restore`, and the sync.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3525 batch the diffs in `RecordedUndone`, with a test that counts the git processes | `flai/internal/storygit/undone.go`, `flai/internal/storygit/undone_test.go` | none |
| 2 | T-3527 fix the `hostapi.go` comments, measure `undone.read` before and after in `design/system/server-performance.md`, close I-0349 | `flai/internal/hostapi/hostapi.go`, `design/system/server-performance.md`, the I-0349 file, `design/issues/summary.md` | T-3525 |

## Figures

- Forecast: 1h, up from flai's 2m. Three stories of the same shape took about an hour each: S-0698, S-0604, and S-0800.
- Delivery: 2026-11-17T16:00:00Z.
- Cost of delay: 5.00 USD a week, from the operator's input of 2m lost per cycle.

## Assumptions

- I left out the issue's second candidate: caching the answer while the two branch heads stay unchanged. Checking the heads costs a `git rev-parse` per call, about what the two batched diffs cost. T-3527 records whether a cache is still worth a story.
- The story's behaviour does not change, so no ADR is needed.
- If no story on the host has a large undone record, T-3527 measures on a scratch repository and says so.
- My start prompt named a thread `asked`, but no thread by that ID exists and S-0828 had no threads. This was its first plan.

I would split, merge, or drop nothing.

### 2026-11-17T08:32:06Z orchestrator
Approved. The plan fits S-0828:

- T-3525 batches the diffs with a test that counts the git processes, which meets criterion 1. T-3527 measures the change and closes I-0349, which meets criterion 2. T-3527 waits for T-3525 because it measures what T-3525 builds.
- Every touch is a file one of the two tasks names. There is no folder touch and no `after`.
- The 1h forecast rests on S-0698, S-0604, and S-0800, three stories of the same shape. The 5.00 USD a week value comes from `flai cod`.
- Leaving the cache out until T-3527's measurement shows it is needed keeps the story small.

I am finalizing S-0828 under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator
