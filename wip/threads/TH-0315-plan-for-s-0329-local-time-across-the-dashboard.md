---
id: TH-0315
title: "Plan for S-0329: local time across the dashboard"
anchor:
  path: wip/kanban/stories/S-0329-display-all-time-stamps-across-the-site-in-local-time.md
  item: S-0329
status: answered
participants: [planner-S-0329, agent-S-0329]
created: 2026-10-07T19:50:54Z
updated: 2026-10-07T20:43:07Z
---

# TH-0315 Plan for S-0329: local time across the dashboard

On wip/kanban/stories/S-0329-display-all-time-stamps-across-the-site-in-local-time.md.

## Entries

### 2026-10-07T19:50:54Z planner-S-0329
S-0329 is planned as five tasks in three layers. Its touches, forecast (1h7m, delivery 2026-10-07T23:08:00Z), and cost of delay value (0 USD a week) are written, with reasons under Notes › Planning.

## Tasks and layers

| Layer | Task | What it does | Waits for |
|-------|------|--------------|-----------|
| 1 | T-1180 | `flaiover/src/lib/localtime.ts`: `localTime` and `localDate`, with tests; vitest pinned to `America/New_York` in `flaiover/vite.config.ts` | none |
| 2 | T-1181 | Item page, threads, inbox, open questions, review, activity, agent stop, host flai, host updates | T-1180 |
| 2 | T-1182 | Planner, orchestrator, and analyzer panels, plan action, host agent notice, story agent, settings next run; the four local `at` helpers go | T-1180 |
| 2 | T-1183 | Charts: `useUTC` dropped, hour labels local, day and week bars placed at local midnight; the chart tables' dates local | T-1180 |
| 3 | T-1184 | `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`; a sweep of `flaiover/src` for any UTC display left | T-1181, T-1182, T-1183 |

The three layer-2 tasks share no file and run together.

## Assumptions

1. A time shows as `2026-10-07 15:35 EDT`: today's shape with the browser's zone in place of `UTC`. A date alone shows as `YYYY-MM-DD` in local time.
2. Chart day and week buckets stay UTC days and weeks, as `flai stats` counts them (`design/system/metrics.md`, ADR-0056). Only axes, tooltips, labels, and bar positions change. Bucketing by your local day would change the metrics contract, which needs an ADR and its own story.
3. The planner's and analyzer's schedules stay described as UTC cron expressions, since flai reads them that way. Their next run shows in local time.
4. A document's raw front matter panel and times written inside markdown bodies, such as a narrative's log headings, are the file's own text and stay as written.
5. Relative ages such as `3h ago` stay as they are.

## Pull

S-0328, in progress, claims `docs/users/flaiover.md`. That holds S-0329 in ready until S-0328 leaves in-progress. The delivery date allows for it.

## What I need from you

Nothing, unless an assumption is wrong. If you want local-day buckets in the charts (assumption 2), or the front matter and markdown text converted (assumption 4), say so here. I would then add a follow-up story rather than widen this one.

### 2026-10-07T20:43:07Z agent-S-0329
S-0329 was worked on these five assumptions, with one addition: flai's own messages shown in a page's notice, such as a refusal that names when a planner started, now show those times in local time too (`localTimes`, applied in `DismissibleNotice.svelte`). Chart buckets stay UTC days and weeks. Each is drawn on its own date and named by it, and hours are labelled in local time. The close-out passed, and the story is going to review.
