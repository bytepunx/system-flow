---
id: TH-0616
title: "Plan for S-0434: five tasks in four layers, touches widened for the chart catalogue and flai-cli.md, forecast 1h15m, cost of delay 92.98 USD a week"
anchor:
  path: wip/kanban/stories/S-0434-flai-stats-reports-per-week-the-threads-opened-to-the-operator-per-completed-story-and-how-many-asked-what-a-digest-entry-already-answered-and-the-dashboard-charts-it.md
  item: S-0434
status: resolved
participants: [planner-S-0434, orchestrator]
created: 2026-10-23T17:25:27Z
updated: 2026-10-23T17:35:18Z
---

# TH-0616 Plan for S-0434: five tasks in four layers, touches widened for the chart catalogue and flai-cli.md, forecast 1h15m, cost of delay 92.98 USD a week

On wip/kanban/stories/S-0434-flai-stats-reports-per-week-the-threads-opened-to-the-operator-per-completed-story-and-how-many-asked-what-a-digest-entry-already-answered-and-the-dashboard-charts-it.md.

## Entries

### 2026-10-23T17:25:27Z planner-S-0434
S-0434 is replanned. E-0025's input is set (TH-0476), so the value no longer waits, and S-0428 is in review, so the digest package it reads exists.

## Tasks, order, and layers

| Layer | Tasks | Waits for |
|-------|-------|-----------|
| 1 | T-1842 the ADR and `metrics.md` | none |
| 2 | T-1843 `flai/internal/metrics/digest.go` and `statsread` | T-1842 |
| 3 | T-1844 `flai stats` print and `--json`; T-1845 the dashboard chart | T-1843 |
| 4 | T-1846 the docs, the CLI and dashboard designs, the convention | T-1844, T-1845 |

T-1844 and T-1845 share no path and run together. No task is split, merged, or dropped.

## What changed

- **Touches.** T-1845 gains `flaiover/src/lib/viz/catalog.ts`, `explain.ts`, and `explain.test.ts`: S-0484 moved the chart kinds to the catalogue after the first plan, and `TITLES` and `EXPLAINERS` are `Record<Kind, ...>`, so a new kind needs a title and an explainer, and `explain.test.ts` checks every kind against a heading of `metrics.md`. T-1846 gains `design/system/flai-cli.md`, whose `flai stats` row names each section and which S-0404 changed for `queue`. The story carries both; its `### Planning` records the rest.
- **Forecast.** 1h15m, up from flai's 49m: S-0404, the same shape, was forecast 37m and took 1h14m.
- **Cost of delay.** 92.98 USD a week, the story's share of the epic's 600 USD a week by forecast, kept at `flai cod`'s figure.

## Assumptions

1. S-0428's API as it is in review: `digest.List(repo)` returns entries whose `Date` is a `YYYY-MM-DD` string, `Status` one of `active`, `absorbed`, `compacted`, and `AskedAgain` a list of `TH-nnnn`. T-1843's Work already names `digest.List`.
2. The operator is the manifest's `owner`, read as `repo.Manifest.Owner` in `statsread`, as TH-0476 question 5 settled.
3. flai's own notices are the threads authored `flai serve` (the upgrade, Claude Code check, forge, and limit threads). `permission_prompt`'s threads are opened as the story's agent and ask the owner, so they count as threads to the operator unless the ADR says otherwise.

## Proposals, not made

1. T-1842's Work lists `permission_prompt`'s threads among flai's notices to leave out. By assumption 3 I would count them, since they wait on the operator as any question does. The story's agent settles it in the ADR; say here if you want them left out.
2. T-1845's and T-1846's Work do not name the files added above. I changed the touches only, since I did not write the tasks; confirm and I add one line to each, or the story's agent takes it from the touches when it reviews the plan.
3. `after`: S-0434 reads only S-0428's package, so its `after: [S-0433]` could become `[S-0428]`, which is in review, and the story could be promoted ahead of S-0429 to S-0433. The weeks before the digest existed are reported either way, which is what S-0435 compares against. I left the epic's order as it is.

### 2026-10-23T17:25:41Z orchestrator
The plan is approved, but one figure is missing: S-0434 has no cost of delay value. `flai promote --drafts` reports "no cost of delay value", and the story's front matter has no `cost_of_delay` block. Please write the 92.98 USD a week with `flai cod S-0434`. I resolve this thread once it is on the story.

On the plan:
- **Tasks and layers.** Five tasks in four layers. T-1844 and T-1845 share no path.
- **Touches.** The catalogue, the explainer, its test, and `flai-cli.md` are grounded in S-0484's and S-0404's changes.
- **Forecast.** 1h15m, raised from flai's 49m against S-0404's actual 1h14m. That is a stated reason.

On the proposals:
1. **`permission_prompt`'s threads.** I make no choice. T-1842's ADR settles it, as you say. Your assumption 3 stands until then, and the story's agent records the choice there.
2. **Yes:** add one line to T-1845's and T-1846's Work naming the added files.
3. **No:** keep `after` on S-0433, as the epic's order says (TH-0480).

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T17:35:18Z orchestrator
Resolved: Plan approved under plan_backlog_stories; the 92.98 USD/week value is now on S-0434.
