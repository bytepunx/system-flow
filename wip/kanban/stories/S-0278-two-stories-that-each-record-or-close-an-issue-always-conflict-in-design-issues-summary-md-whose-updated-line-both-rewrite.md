---
id: S-0278
type: story
nature: improvement
title: Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite
status: ready
owner: alex
created: 2026-10-05T04:40:47Z
updated: 2026-10-06T11:44:49Z
transitions:
  - to: ready
    at: 2026-10-06T11:29:58Z
    by: alex
tags: []
touches: [flai/cmd/branch.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, design/adrs, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-10-05T04:40:47Z
  value: 22.5
  by: planner-S-0278
  at: 2026-10-06T11:31:29Z
forecast:
  duration: 40m
  delivery: 2026-10-06T16:11:00Z
  basis: "Its own forecast of 40m; 8th in the pull order with an in-progress limit of 3, behind S-0221, S-0222, S-0224, S-0226, S-0223, S-0227, S-0229 and S-0284."
  by: flai
  at: 2026-10-06T11:44:49Z
finalized:
  by: alex
  at: 2026-10-06T11:29:55Z
---
# S-0278 Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite

## Goal

This story remediates [I-0074](../../../design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md), "Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0074 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0074 is closed with `flai issue close I-0074 --reason` saying what fixed it

## Tasks
- T-1011 Choose how flai keeps design/issues/summary.md from stopping a sync, from I-0074's instances, and record it in an ADR
- T-1012 flai stream sync and flai accept regenerate design/issues/summary.md when a rebase stops on it alone, with a test reproducing I-0074
- T-1013 flai stream sync's trial merge does not report design/issues/summary.md as a conflict between two open story branches
- T-1014 The design, the user guide, and flai stream sync's help say that sync and acceptance regenerate design/issues/summary.md instead of stopping on it
- T-1015 Close I-0074 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0074. time_lost_per_cycle 9m: 3m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-05T03:22:09Z, 0.1 days before this story; under one cycle counts as one).

### Planning

The story declared no touches, so `flai touches suggest` started from the three files the fix lives in (22 of 948 commits changed them). Where each touch came from:

- `flai/cmd/branch.go`: layout. `syncStoryBranch` runs the rebase that stops on `summary.md`, and `mergeStoryBranch` (`flai accept`) syncs through it.
- `flai/cmd/stream_sync.go`: layout. `trialMerge` and `reportConflicts` opened TH-0126 and TH-0130 in the first and third instances.
- `flai/cmd/stream_sync_test.go`: co-change, 18%. The reproduction tests use real git here.
- `flai/cmd/stream.go`: co-change, 27%. It holds `flai stream sync`'s help.
- `flai/internal/storygit/sync.go` and `sync_test.go`: layout. They hold `Conflicts` and `RebaseInProgress`, where a rebase-continue helper belongs.
- `flai/internal/issues/issues.go` and `issues_test.go`: layout and co-change (`issues_test.go` 27%). `Summary` and `WriteSummary` render the file and its `updated` line.
- `design/adrs`: design. The goal asks for a solution proposed before it is built, recorded as a decision.
- `design/system/flai-cli.md`: co-change, 36%, and design. It holds the `flai stream sync` row.
- `design/system/continuous-improvement.md`: design. Its `## Summary` section describes the file.
- `docs/users/flai.md`: co-change, 36%. It holds the sync conflict walkthrough.
- `docs/users/flai-reference.md`: layout. It is generated from the command help.
- I-0074's file and `design/issues/summary.md`: criterion 2. `flai issue close` writes both.

`flai/cmd/accept.go` (co-change, 23%) is left out: acceptance reaches the sync through `mergeStoryBranch` in `branch.go`.

Forecast: 40m, raised from `flai forecast`'s 26m. flai's figure is 89 s per unit of size over 21 large improvement stories, times size 17. This story adds a design decision and an integration test that builds two story branches with real git. S-0253 (conflict markers on acceptance) took 35m, and S-0197 (stream sync conflicts) took 31m. Delivery moves by the same 14m, to 16:05Z, from 8th in the pull order.

Cost of delay: 22.50 USD a week, `flai cod`'s figure from flai's input of 9m a cycle, unchanged. I-0089 has the same cause, 3m a cycle more. It is not added, because the story does not close it. See the plan's thread.
