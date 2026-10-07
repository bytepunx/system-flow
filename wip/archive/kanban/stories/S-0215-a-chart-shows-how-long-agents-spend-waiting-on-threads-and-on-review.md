---
id: S-0215
type: story
nature: feature
title: A chart shows how long agents spend waiting on threads and on review
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-07T14:26:00Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:56Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:41Z
    by: alex
  - to: ready
    at: 2026-10-07T03:29:20Z
    by: alex
  - to: in-progress
    at: 2026-10-07T08:54:06Z
    by: agent-S-0215
  - to: review
    at: 2026-10-07T09:18:15Z
    by: agent-S-0215
  - to: done
    at: 2026-10-07T14:26:00Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts, flai/internal/metrics, design/system/metrics.md, design/adrs, flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/WaitTable.svelte.test.ts]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1683
  models:
    - model: claude-opus-5-5
      input: 358
      output: 113032
      cache_read: 15311872
      cache_write: 547156
      cost: 8.7553
  strategic:
    - kind: orchestrator
      seconds: 1233
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 96
          output: 1535
          cache_read: 11458030
          cache_write: 23423
          cost: 2.9957
        - model: claude-sonnet-5-5
          input: 8
          output: 48
          cache_read: 70344
          cache_write: 36453
          cost: 0.0716
cost_of_delay:
  value: 76.01
  by: planner-S-0215
  at: 2026-10-05T05:44:37Z
forecast:
  duration: 1h15m
  delivery: 2026-10-07T10:32:00Z
  basis: "Its own forecast of 1h15m; 1st in the pull order with an in-progress limit of 3, behind S-0214, S-0265 and S-0275."
  by: flai
  at: 2026-10-07T08:48:14Z
---
# S-0215 A chart shows how long agents spend waiting on threads and on review

## Goal

Agents wait for the operator: on threads they open and on stories in review. That time is the operator's bottleneck and the orchestrator's reason to exist. The designer asked for a chart of it on 2026-10-02.

## Acceptance criteria
- [x] `/charts/agent-waiting`: stacked bar per week of the window of hours agents waited, split into thread waits (opened to answered, from thread timestamps while the story was in progress) and review waits (review to done, from transitions), with the mean wait per story as a line
- [x] A table under it lists the longest waits in the window with the story, the thread, and who was awaited
- [x] Spans the window, reads `/api/stats`, matches `flai stats --json`; listed under Flow; design and user guide describe it; tests cover the mapping

## Tasks
- T-0919 metrics.md defines each wait of the window, with its story, thread, and who was awaited, and an ADR records it
- T-0923 flai stats reports the longest waits of the window with their story, thread, and who was awaited
- T-0927 The agent-waiting chart maps the weeks of waiting to stacked thread and review bars with the mean wait per story as a line
- T-0932 A table under the agent-waiting chart lists the longest waits with the story, the thread, and who was awaited
- T-0934 The dashboard design and the user guide describe the agent-waiting chart and its table

## Notes

### Planning

Touches, and where each came from:

- Declared, all kept: `flaiover/src/routes/charts`, `flaiover/src/lib/charts`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `flaiover/src/lib/viz`, `flaiover/src/lib/sitemenu.ts`, `flai/internal/metrics`, `design/system/metrics.md`. `flaiover/src/lib/charts` does not exist: the chart builders are in `flaiover/src/lib/viz/charts.ts`. The Flow group is `FLOW_KINDS` there, so `sitemenu.ts` should not need to change.
- Design: `design/adrs`. `flai stats --json` gives waiting only as totals per story and per week (`flai/internal/metrics/waiting.go`), with no thread ID and no one awaited. The table needs `waiting.longest[]`, and `metrics.md` changes only with an ADR.
- Layout: `flaiover/src/lib/components/WaitTable.svelte` and its test, on the pattern of `SpendTable.svelte`, which the chart page already renders under the spend charts.
- Co-change: `flai touches suggest` ranks `design/system/flai-cli.md` (41%) and `docs/users/flai.md` (39%) highest. Neither is predicted: the work changes no command or flag, and `metrics.md` is where `flai stats --json` fields are defined.

Forecast: `flai forecast` gave 14m on the declared touches, and 31m once the touches above were added (132 s per unit of size over 14 large stories on claude-opus-5-5, size 14). I raised it to 1h15m. Size counts touches, not the four layers of work: an ADR, a Go metrics extension with second-exact fixtures, a chart, and a table component. 1h15m is in line with S-0212 (1h15m) and S-0214 (1h45m), sibling chart stories. Delivery is flai's 2026-10-05T15:00Z, 12th in the pull order behind the in-progress limit of 3, plus the 44m added: 2026-10-05T15:44Z.

Cost of delay: 76.01 USD a week, from `flai cod`, which had given 60.81 at the old 1h forecast. The story has no inputs of its own, so this is its share of E-0016's 1500.00 USD a week: 1h15m of the 24h40m forecast over the epic's 17 open stories without inputs. It replaces the 48.78 the epic's planner recorded on 2026-10-04. It stands as flai worked it out: no input was missing.
