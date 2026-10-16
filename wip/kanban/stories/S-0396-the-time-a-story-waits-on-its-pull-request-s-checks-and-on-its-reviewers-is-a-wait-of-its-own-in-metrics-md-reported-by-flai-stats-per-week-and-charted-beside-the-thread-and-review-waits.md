---
id: S-0396
type: story
nature: feature
title: The time a story waits on its pull request's checks and on its reviewers is a wait of its own in metrics.md, reported by flai stats per week and charted beside the thread and review waits
status: backlog
parent: E-0021
owner: alex
created: 2026-10-16T04:38:48Z
updated: 2026-10-16T05:24:55Z
transitions: []
tags: [cli, dashboard]
topics: [analysis, dashboard]
touches: [flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/WaitTable.svelte.test.ts, design/adrs, design/system/metrics.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [S-0394]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 42m
  delivery: 2026-10-16T09:55:00Z
  basis: "Its own forecast of 42m; 6th in the pull order with an in-progress limit of 5, behind S-0241, S-0392, S-0393, S-0394 and S-0395."
  by: flai
  at: 2026-10-16T05:24:55Z
---
# S-0396 The time a story waits on its pull request's checks and on its reviewers is a wait of its own in metrics.md, reported by flai stats per week and charted beside the thread and review waits

## Goal

A slow CI or a slow reviewer shows up where the operator and the analyzer already look for waiting. Under pull-request acceptance, the time a story in review spends waiting for its required checks, and then for its reviewers, is measured apart from the review wait. `flai stats` reports it per item, per week, and among the longest waits. The dashboard's agent-waiting chart and its table show it beside the thread and review waits.

## Acceptance criteria

- [ ] An ADR defines two waits, and `design/system/metrics.md` § Waiting carries them:
  - **The checks wait** runs from each push of a story's head to the time every required check on that head has finished.
  - **The reviewers wait** runs from then to the first approval that meets the main branch's rule, or to the first requested change.
  - Both are read from the events S-0394 keeps, and both fall inside the story's `review` intervals.
- [ ] `flai stats --json` carries `items[].wait_checks_seconds` and `items[].wait_reviewers_seconds`, absent when the story had no pull request. It also carries `checks` and `reviewers` in each of `waiting.weeks[]`, with `total_seconds` and `mean_seconds`. The longest waits gain the kinds `checks` and `reviewers`. `flai stats` prints one line for each, under the review line.
- [ ] The dashboard's agent-waiting chart draws the two waits as series of their own beside the thread and review waits, and its longest-waits table names their kinds.
- [ ] A project in `local` mode, or one with no kept events, reports the two as zero or absent, as each field's definition says, and its other figures do not change.
- [ ] `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, and `docs/users/flaiover.md` describe the new waits.

## Tasks
- T-1628 An ADR defines the checks wait and the reviewers wait, and metrics.md carries their definitions and fields
- T-1629 The metrics package computes the checks and reviewers waits per item, per week, and among the longest waits, from the kept forge events statsread loads
- T-1630 flai stats prints the checks and reviewers waits under the review line
- T-1631 The dashboard's agent-waiting chart draws the checks and reviewers waits as series of their own, and its table names their kinds
- T-1632 The dashboard and CLI design and the user docs describe the checks and reviewers waits

## Notes

- Drafted by the planner for E-0021, after S-0394, whose kept events it reads. It can run beside S-0395 and S-0397: it touches none of their paths.
- `CLAUDE.md` says a definition in `metrics.md` changes only with an ADR, hence the ADR first.
- The events are kept under `.flai-cache`, as the conflict log of ADR-0138 is, so they are not in git. A clone without them reports no forge waits.

### Planning

- **Touches.** None were declared. Every touch is the planner's.
- **Forecast and cost of delay.** Written after `flai forecast` and `flai cod`.
