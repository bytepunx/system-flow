---
id: S-0329
type: story
nature: improvement
title: Display all time stamps across the site in local time
status: done
owner: alex
created: 2026-10-07T19:44:24Z
updated: 2026-10-07T22:25:33Z
transitions:
  - to: ready
    at: 2026-10-07T19:44:25Z
    by: alex
  - to: in-progress
    at: 2026-10-07T20:24:21Z
    by: agent-S-0329
  - to: review
    at: 2026-10-07T20:43:10Z
    by: agent-S-0329
  - to: done
    at: 2026-10-07T22:25:33Z
    by: alex
tags: [dashboard]
touches: [flaiover/src, flaiover/vite.config.ts, flaiover/src/lib/localtime.ts, flaiover/src/lib/localtime.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/Threads.svelte.test.ts, flaiover/src/routes/threads/threads.svelte.test.ts, flaiover/src/lib/components/InboxView.svelte, flaiover/src/lib/components/Inbox.svelte.test.ts, flaiover/src/lib/components/OpenQuestions.svelte, flaiover/src/lib/components/OpenQuestions.svelte.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/ActivityView.svelte, flaiover/src/lib/components/ActivityView.svelte.test.ts, flaiover/src/routes/activity/+page.svelte, flaiover/src/routes/activity/activity.svelte.test.ts, flaiover/src/lib/components/AgentStopConfirm.svelte, flaiover/src/lib/components/HostFlai.svelte, flaiover/src/lib/components/HostFlai.svelte.test.ts, flaiover/src/lib/components/HostProcesses.svelte, flaiover/src/lib/components/HostProcesses.svelte.test.ts, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts, flaiover/src/lib/components/PlannerPanel.svelte.test.ts, flaiover/src/lib/components/AnalyzerPanel.svelte.test.ts, flaiover/src/lib/components/PlanAction.svelte, flaiover/src/lib/components/PlanAction.svelte.test.ts, flaiover/src/lib/components/HostAgentNotice.svelte, flaiover/src/lib/components/HostAgentNotice.svelte.test.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", flaiover/src/lib/components/SpendTable.svelte, flaiover/src/lib/components/SpendTable.svelte.test.ts, flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/WaitTable.svelte.test.ts, flaiover/src/lib/components/ForecastTable.svelte, flaiover/src/lib/components/ForecastTable.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1140
  turns:
    - day: 2026-10-07
      ceremony: 4
      test_runs: 1
      hand_edits: 3
      work: 51
  models:
    - model: claude-opus-5-5
      input: 420
      output: 152984
      cache_read: 26035100
      cache_write: 669419
      cost: 12.3712
  strategic:
    - kind: planner
      seconds: 392
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 362
          output: 10788
          cache_read: 2085625
          cache_write: 69282
          cost: 0.3495
        - model: claude-opus-5-5
          input: 124
          output: 29877
          cache_read: 9005877
          cache_write: 134159
          cost: 5.4833
    - kind: orchestrator
      seconds: 1602
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 178
          output: 2835
          cache_read: 17778673
          cache_write: 57870
          cost: 4.3951
        - model: claude-sonnet-5-5
          input: 14
          output: 56
          cache_read: 213782
          cache_write: 57248
          cost: 0.2336
cost_of_delay:
  inputs:
    penalty_per_week: 0
    by: alex
    at: 2026-10-07T19:44:24Z
  value: 0
  by: planner-S-0329
  at: 2026-10-07T19:48:50Z
forecast:
  duration: 1h7m
  delivery: 2026-10-07T23:08:00Z
  basis: "Median 78 s per unit of size over 44 done improvement stories on claude-opus-5-5 in the large band, times size 51 (2 criteria, 49 touches); delivery moved from 21:02 to after S-0328's forecast delivery at 22:01, since S-0328 in progress claims docs/users/flaiover.md and holds this story until it leaves in-progress."
  by: planner-S-0329
  at: 2026-10-07T19:48:50Z
---
# S-0329 Display all time stamps across the site in local time

## Goal

Every view in the site currently displays time stamps in UTC. It's important that we _record_ timestamps in UTC, but they should be displayed in the user's local time zone.

## Acceptance criteria
- [x] No times display in the dashboard in UTC
- [x] All times display in the dashboard in the local time zone

## Tasks
- T-1180 A shared formatter shows a recorded UTC time in the browser's local zone, and the tests run in a fixed non-UTC zone
- T-1181 Item, thread, inbox, review, activity, and host views show their times in local time
- T-1182 Strategic agent panels, the story agent, plan actions, host agent notices, and settings show their times in local time
- T-1183 Charts and the chart tables show their times and dates in local time, over flai's UTC buckets
- T-1184 The user guide and the dashboard's design say times are shown in local time, and a sweep finds no UTC display left

## Notes

### Planning

Planned by planner-S-0329 on 2026-10-07. The tasks are T-1180 to T-1184, in three layers: T-1180 alone; then T-1181, T-1182, and T-1183 together; then T-1184.

Touches:

| Touch | From |
|-------|------|
| `flaiover/src` | declared; kept as declared. It is a folder touch: the tasks name the files, which narrow it in the story's claim (ADR-0096), and the sweep in T-1184 may find a file no task names yet |
| `flaiover/src/lib/localtime.ts`, `localtime.test.ts` | layout: a new shared helper beside `age.ts`, since no formatter for an absolute time exists |
| `flaiover/vite.config.ts` | layout: the vitest projects pin a non-UTC zone so tests prove the conversion |
| the 43 other component, route, chart, and test files under `flaiover/src` the tasks name | layout and code: each formats or prints one of flai's UTC times today, or tests that it does, found by searching `flaiover/src` |
| `docs/users/flaiover.md` | co-change (21% of the commits that changed `flaiover/src`) and documentation.md: user-facing behaviour changes |
| `design/system/flaiover-dashboard.md` | co-change (27%) and design: the living design records the display rule |

`flai touches suggest` also listed `flai/internal/hostapi/*`, `flai/cmd/serve_actions.go`, `flai/internal/serve/*`, and `docs/users/flai.md`; none is predicted, since flai and the API keep recording and returning UTC.

`docs/users/flaiover.md` matters for the pull: S-0328, in progress, claims it, so this story is held in ready until S-0328 leaves in-progress.

Forecast: 1h7m, delivery 2026-10-07T23:08:00Z. The duration is `flai forecast`'s, 78 s per unit of size over 44 done improvement stories in the large band, times size 51 (2 criteria, 49 touches); it stands, since the work is many small, alike edits plus their tests. Before the touches were written it gave 5m on size 3, which counted only the one declared folder. The delivery is moved from the 21:02 `flai forecast` gave to S-0328's forecast delivery, 22:01, plus the duration, because of the hold above.

Cost of delay: value 0 USD a week, from `flai cod`. The operator's inputs give a penalty of 0 a week and no revenue or time lost, so the value stands at 0.

Assumptions, also on the plan thread:

- Times are shown as `YYYY-MM-DD HH:MM` with the zone's short name, the shape shown today with the zone in place of `UTC`; a date alone as `YYYY-MM-DD`.
- Chart day and week buckets stay UTC, as `flai stats` counts them (`design/system/metrics.md`, ADR-0056); only their labels, axes, and tooltips change. Rebucketing by local day needs a metrics ADR and a story of its own.
- A schedule's cron expression stays described as UTC, since flai reads it so; its next run shows in local time.
- A document's raw front matter panel and times written inside markdown bodies are the file's text, and stay as written.
