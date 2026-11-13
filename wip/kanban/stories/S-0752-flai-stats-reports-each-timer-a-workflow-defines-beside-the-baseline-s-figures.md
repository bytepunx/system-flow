---
id: S-0752
type: story
nature: feature
title: flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:01:24Z
updated: 2026-11-13T08:22:07Z
transitions: []
tags: [cli]
topics: [cli, analysis]
touches: [flai/internal/metrics/timers.go, flai/internal/metrics/timers_test.go, flai/internal/metrics/metrics.go, flai/internal/metrics/window.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, design/system/metrics.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0750]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 13
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 15
          cache_read: 121381
          cache_write: 20507
          cost: 0.0405
draft: true
forecast:
  duration: 33m
  delivery: 2026-11-13T19:18:00Z
  basis: "Its own forecast of 33m; 26th in the pull order with an in-progress limit of 3, behind S-0660, S-0735, S-0760, S-0761, S-0745, S-0759, S-0751, S-0757, S-0740, S-0661, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749 and S-0750."
  by: flai
  at: 2026-11-13T08:22:07Z
---
# S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included

## Goal

An operator adds a timer to the workflow, such as the time from the first move to review to the acceptance, and `flai stats` reports it beside lead and cycle time with no change to flai. Since a timer is a reading of each item's recorded transitions, it applies to past items as well as new ones.

## Acceptance criteria

- [ ] `flai stats --json` gives, under a section the timers' ADR names, each timer the workflow defines: per item its seconds, and over the window its count, mean, p50, and p85, as the baseline's durations are given
- [ ] `flai stats` without `--json` prints the workflow's own timers after the baseline's
- [ ] A timer added to the workflow reads the transitions items already carry, and reports the items completed before it was added
- [ ] `stats.get` on the host carries the same section, so the dashboard reads the timers from flai
- [ ] The baseline's own figures and their fixtures are unchanged
- [ ] Tests give a timer from the first move to `review` to the acceptance, and a sum of the time in two states, on fixture items with known transitions
- [ ] `design/system/metrics.md` and `docs/users/flai.md` say how a timer is defined and where `flai stats` reports it

## Tasks

- T-3243 The metrics report each timer the workflow defines, per item and over the window, past items included
- T-3247 flai stats prints the workflow's own timers, and stats.get on the host carries them
- T-3248 metrics.md and the user guide say how a timer is defined and where flai stats reports it

## Notes

- Until the project's own workflow file is read, a test gives the timers through a workflow other than the baseline, as S-0737's test hook allows.

### Planning

- Touches, all files:
  - `flai/internal/metrics/timers.go`, `timers_test.go`, `metrics.go`, `window.go`: layout. The section is new; `window.go` holds the window's aggregates.
  - `flai/cmd/stats.go`, `stats_test.go`, `flai/internal/hostapi/reads.go`, `reads_test.go`: layout. `stats.get` is in `reads.go`.
  - `design/system/metrics.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: design (criterion 7); the reference is regenerated only if help text changes.
  - Left out from `flai touches suggest`: `flai-cli.md`, `docs/operators/`, and `flaiover-dashboard.md` (14% to 54%). The dashboard's side is S-0754's.
- `after`: S-0750. This story reads the timer evaluation S-0750 adds to `flai/internal/workflow` and the section name and definitions S-0750's ADR gives.
- Tasks: T-3243 (the section) first; then T-3247 (the command and `stats.get`) and T-3248 (the documents) together.
- Forecast 33m, as `flai forecast` gives it for size 18. The evaluation is S-0750's; this story reports it.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
