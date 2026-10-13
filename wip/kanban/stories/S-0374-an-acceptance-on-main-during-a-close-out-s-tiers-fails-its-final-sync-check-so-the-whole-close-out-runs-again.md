---
id: S-0374
type: story
nature: improvement
title: An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again
status: backlog
owner: alex
created: 2026-10-12T16:52:54Z
updated: 2026-10-13T22:00:19Z
transitions: []
tags: [cli]
topics: [cli, conventions, template, testing]
touches: [design/adrs/README.md, flai/internal/verify/story.go, flai/internal/verify/paths.go, flai/internal/verify/story_test.go, flai/internal/verify/paths_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 322
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 32
          output: 18651
          cache_read: 964719
          cache_write: 90573
          cost: 0.0303
        - model: claude-opus-5-5
          input: 46
          output: 19274
          cache_read: 2230760
          cache_write: 116361
          cost: 1.7627
    - kind: orchestrator
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 20
          cache_read: 658068
          cache_write: 2445
          cost: 0.1665
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 18m
    by: flai
    at: 2026-10-12T16:52:54Z
  value: 45
  by: planner-S-0374
  at: 2026-10-12T17:35:05Z
forecast:
  duration: 1h
  delivery: 2026-10-13T23:11:00Z
  basis: "Its own forecast of 1h; 7th in the pull order with an in-progress limit of 5, behind S-0378, S-0380, S-0241, S-0304, S-0305, S-0306, S-0371 and S-0372."
  by: flai
  at: 2026-10-13T22:00:19Z
---
# S-0374 An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again

## Goal

This story remediates [I-0136](../../../design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md), "An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0136 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0136 is closed with `flai issue close I-0136 --reason` saying what fixed it

## Tasks
- T-1493 Record the ADR: a close-out that main outran syncs and runs again only the tiers main's new paths select
- T-1494 flai verify reuses the tiers a clean rebase onto a newer main cannot reach, and runs again those main's new paths select
- T-1495 The close-out syncs and verifies again, at most three rounds, when main moved during its tiers
- T-1496 Document the close-out's rounds and verify's reuse across a rebase in the design, the user docs, and the work-management convention
- T-1497 Close I-0136 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0136. time_lost_per_cycle 18m: 6m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-12T15:19:01Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, for the plan thread to settle: when the close-out's last check fails because main moved, the close-out runs `flai stream sync` and verifies again, at most three rounds, and `flai verify` reuses across that clean rebase every tier main's new paths do not select. Passing over main's non-wip commits outright, as I-0136's first instance suggests, would let a main change that breaks the branch without touching its files, such as a Go package it imports, reach review unverified.

Touches, every one a file; no folder touch is kept:

| Touch | From |
|-------|------|
| `flai/internal/verify/story.go`, `paths.go`, `story_test.go`, `paths_test.go` | layout: `sync()`, `reusable`, and `BaseChanges` live there (ADR-0135, S-0341) |
| `scripts/close-out.sh`, `template/root/scripts/close-out.sh` | layout: the last check, in both copies (ADR-0135 decision 4) |
| `template/CHANGELOG.md` | design: the template's close-out names the flai it needs (ADR-0135 consequences) |
| `design/adrs/README.md` | design: the new ADR refining ADR-0135 and ADR-0110; `flai task done` adds its file, not known yet |
| `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | design: where the close-out and verify are described |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md` | design: the baseline close-out rule names what stops it |
| `design/issues/I-0136-...md`, `design/issues/summary.md` | goal: criterion 2 closes the issue |

`flai touches suggest` found no co-change above 12% that the remedy reaches; its top paths (`docs/operators/index.md`, `design/system/flaiover-dashboard.md`) change with nearly everything, and `docs/users/flai-reference.md` changes only with a new flag, which the plan does not add.

Forecast 1h, raised from flai's 23m: five tasks in four layers, with an ADR, a change to verify's reuse tested in a real repository, a loop in two copies of the close-out, and docs in five files. Recent planned stories of this shape (S-0369, S-0371, S-0372) were raised to 45m to 1h for the same reasons. Delivery moves by the 37m added.

Cost of delay 45 USD a week stands, as `flai cod` works it from flai's 18m input at 150 USD an hour. It is likely low: the three instances came within 1.5 hours of one another, so the 168h cycle undercounts. The input is not the planner's to change.
