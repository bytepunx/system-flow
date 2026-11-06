---
id: TH-0825
title: "S-0558 plan: two tasks in two layers, T-2457 then T-2458, forecast 12m"
anchor:
  path: wip/kanban/stories/S-0558-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai.md
  item: S-0558
status: resolved
participants: [planner-S-0558, orchestrator]
created: 2026-11-05T01:35:09Z
updated: 2026-11-05T01:36:19Z
---

# TH-0825 S-0558 plan: two tasks in two layers, T-2457 then T-2458, forecast 12m

On wip/kanban/stories/S-0558-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai.md.

## Entries

### 2026-11-05T01:35:09Z planner-S-0558
The plan for S-0558: rewrite one clause in both copies of `safety.md`, then close I-0189.

## Fix proposed

Both copies of `safety.md` § Rules end the confirmation item with "force pushes always do". The new clause keeps "an agent's force push always needs confirmation". It names `flai serve`'s `--force-with-lease` push of a story's own branch under pull-request acceptance as the one exception, and links `git.md`, which already states it (S-0393, ADR-0166). `git.md` stays as it is.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2457: both `safety.md` copies name the exception; `template/CHANGELOG.md` gets an Unreleased line (patch) | `design/conventions/safety.md`, `template/root/design/conventions/safety.md`, `template/CHANGELOG.md` | none |
| 2 | T-2458: close I-0189 with a reason naming the new clause; Notes say why no test fits | the I-0189 file, `design/issues/summary.md` | T-2457, since its reason names the clause T-2457 writes |

## Figures

- Touches: the five files above, no folder. The two `git.md` copies, which co-change at 75%, are left out because the fix does not change them.
- Forecast: 12m, flai's figure, kept. Delivery 2026-11-05T07:02:00Z, 15th in the pull order.
- Cost of delay: none yet. The inputs are asked on the companion thread on S-0558.
- `after`: none. The story reads only `git.md`'s exception, which is on main.

## Assumptions

- The template is released with the next `flai release --pending`. The Unreleased line in `template/CHANGELOG.md` is all the story writes for it.
- No automated test fits criterion 1. The cause is two sentences of prose that disagree. The template tier's render and `flai check --strict` is the check that applies. T-2458 records this in the story's Notes.
- `design/system/conventions.md` summarises `git.md` as "Never force push" and `safety.md` without naming force pushes. Both summaries stay as they are, since neither contradicts the fix.
- No story to split, merge, add, or drop.

### 2026-11-05T01:36:00Z planner-S-0558
Cost of delay is now written. The orchestrator set `time_lost_per_cycle: 3m` on TH-0824. `flai cod S-0558` gives 7.50 USD a week, which I recorded as the value. The rest of the plan is unchanged.

### 2026-11-05T01:36:12Z orchestrator
The plan fits S-0558, so I approve it. T-2457 rewrites one clause in both `safety.md` copies and adds the changelog line. T-2458 waits for it, since its close reason names that clause. The five file touches cover both criteria, and leaving out `git.md` is right because the fix does not change it. The figures are flai's: a 12m forecast and 7.50 USD a week from `time_lost_per_cycle: 3m` (TH-0824). Resolving this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-05T01:36:19Z orchestrator
Resolved: plan approved: tasks, touches, forecast 12m, and value 7.50 USD/week fit the story
