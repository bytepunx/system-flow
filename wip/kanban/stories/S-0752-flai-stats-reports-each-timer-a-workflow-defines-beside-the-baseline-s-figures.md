---
id: S-0752
type: story
nature: feature
title: flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:01:24Z
updated: 2026-11-12T16:01:24Z
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
draft: true
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

## Notes

- Until the project's own workflow file is read, a test gives the timers through a workflow other than the baseline, as S-0737's test hook allows.
