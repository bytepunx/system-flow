---
id: S-0166
type: story
nature: remediation
title: Charts ignore the window drop-down
status: done
parent: E-0013
owner: alex
created: 2026-09-29T22:45:43Z
updated: 2026-09-29T23:14:09Z
transitions:
  - to: ready
    at: 2026-09-29T22:45:46Z
    by: alex
  - to: in-progress
    at: 2026-09-29T22:46:03Z
    by: agent-S-0166
  - to: review
    at: 2026-09-29T22:57:32Z
    by: agent-S-0166
  - to: done
    at: 2026-09-29T23:14:09Z
    by: alex
tags: [dashboard, cli]
topics: [server-side, client-side]
touches: [flaiover/src, flai/cmd, flai/internal/metrics, design/system/metrics.md, design/system/flaiover-dashboard.md, design/adrs, docs/users, design/issues/I-0027-work-items-written-in-the-main-checkout-reach-ci-without-a-markdown-lint.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 712
  models:
    - model: claude-opus-5-5
      input: 178
      output: 52731
      cache_read: 13035460
      cache_write: 205923
      cost: 5.3098
---
# S-0166 Charts ignore the window drop-down

## Goal

Neither the time axis or data points are correctly re-rendered when the time window drop down is set. From what I can tell each chart seems to have a pre-determined time window it renders with and no changes in the UI have any affect.

## Acceptance criteria
- [x] The axis matches the selected time window
- [x] The data points are rendered correctly for the selected time window
- [x] Changing the window changes the graph's time axis and plotted data to match

## Tasks
- T-0589 flai stats lays out burn-up, cumulative flow, and throughput over the window's days and weeks
- T-0590 The charts plot what the window holds, on a time axis that spans it
- T-0591 The design, the ADR, and the user guide say each chart spans the window

## Notes

The cause was on both sides. `flai stats` answered every item and the whole burn-up and cumulative flow history whatever the window, and throughput had only the weeks with something done. The dashboard fitted each time axis to the points it drew. ADR-0054 records the fix: every chart spans the window.

How the criteria were verified:

- Axis: `charts.test.ts` checks each time axis's ends against the report's window (cycle time, burn-up, cumulative flow, completion over time, the charts of spend over time by day and by week), and that a narrower window moves them. `TestSeriesCoverTheWindow` checks flai's burn-up, cumulative flow, and throughput ranges for 7, 30, and 365 days.
- Data points: the same tests check that a narrower window drops the items completed before it from cycle time, time in state, estimates, and cost by item. On this repository, `flai stats` built from the branch gives 2 days of cumulative flow for 1d, 8 for 7d, and 15 (all history) for 365d; before, it gave 15 for each.
- Changing the window: `charts.svelte.test.ts` changes the window on the page, sees flai asked for it, and sees the drawn axis and points follow. It also checks that a late answer to an earlier window does not replace the later one.

Not verified in a browser against the operator's running dashboard: that needs a `flai serve` built from this branch.
