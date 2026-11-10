---
id: TH-1001
title: "Plan for S-0674: one ADR, then the drift's attribution and task done's narrowing side by side, then measure and close I-0286"
anchor:
  path: wip/kanban/stories/S-0674-declared-touches-match-what-a-story-changes-in-16-percent-of-completed-stories.md
  item: S-0674
status: open
participants: [planner-S-0674]
created: 2026-11-10T08:51:40Z
updated: 2026-11-10T08:51:40Z
---

# TH-1001 Plan for S-0674: one ADR, then the drift's attribution and task done's narrowing side by side, then measure and close I-0286

On wip/kanban/stories/S-0674-declared-touches-match-what-a-story-changes-in-16-percent-of-completed-stories.md.

## Entries

### 2026-11-10T08:51:40Z planner-S-0674
S-0674 is planned as four tasks in three layers. The operator moved it to ready at 08:50:34Z, while the plan was being written. The tasks, touches, and figures were all in place by then.

Cause found. In `flai stats --json`, 25 of the 115 stories from S-0560 on are exact. 80 have files `outside` their touches. Every one of those files is under `design/digest/`, `design/issues/`, or `design/analysis/`, which flai's own commits put there under the story's ID: `accept and archive` sweeps in digest entries (ADR-0248), `create story` writes the issue, and `record I-nnnn` is the close-out's. Leaving those folders and `wip/` out on both sides makes 97 of 115 exact. The 18 left each kept 1 to 4 predicted files they never changed, because `flai task done` only widens (ADR-0144). The details are under `### Planning` in the story's Notes.

Tasks and layers:

| Layer | Task | after |
|-------|------|-------|
| 1 | T-2936 an ADR refining ADR-0144 and ADR-0113 with both decisions | none |
| 2 | T-2939 `storygit.Committed` and `touchDrift` leave out the digest, analysis, issues, and wip folders on both sides; a test reproduces I-0286; `metrics.md` § Claims and touches | T-2936 |
| 2 | T-2940 `flai task done` narrows: it drops the closing task's file touches that the branch does not change, and drops each from the story when no other open task names it; docs, the convention, and its template copy | T-2936 |
| 3 | T-2941 measure `exact_share` again, record the figures, close I-0286 | T-2939, T-2940 |

T-2939 and T-2940 share no path, so they run together.

Figures:

- Forecast: 2h30m, delivery 2026-11-10T16:30Z. `flai forecast` gave 36m, raised because size does not count the ADR, the two tested changes, the docs, and the close-out tiers.
- Cost of delay: 600 USD a week, as `flai cod` works it out from your `time_lost_per_cycle` of 4h.

Assumptions:

1. The narrowing in T-2940 belongs in this story. After it, `exact_share` measures the claim as task done leaves it, not the planner's first prediction. If you would rather keep the metric as a measure of the planner's prediction, drop T-2940. T-2939 alone takes the share to about 84%.
2. The new ADR's file name is unknown until `flai adr new` numbers it, so the story keeps the folder touch `design/adrs`. That path is in `claims.shared` and holds no story.
3. The 4h `time_lost_per_cycle` stands as your input. The evidence suggests that the send-backs at acceptance were not caused by missed claims, since the `outside` files were flai's bookkeeping. This story may recover less of the 4h than I-0286 assumed. Lower it if you agree.

Not planned: co-change suggestions (`flai touches suggest`) read the same bookkeeping commits and could suggest digest or issue files. No instance showed it, so it is left out.

No change is proposed to the story's words.
