---
id: S-0323
type: story
nature: improvement
title: "flai check finds `narrative.state` outside the story at close-out"
status: in-progress
owner: alex
created: 2026-10-07T18:59:54Z
updated: 2026-10-08T06:12:33Z
transitions:
  - to: ready
    at: 2026-10-08T05:54:04Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T06:11:40Z
    by: agent-S-0323
tags: [flai, template]
topics: [cli, conventions, template]
touches: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0109-flai-check-finds-narrative-state-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0111-flai-check-finds-narrative-state-outside-the-story-at-close-out.md]
after: [S-0318]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 96
  estimated: true
  turns:
    - day: 2026-10-08
      ceremony: 1
      hand_edits: 1
      work: 15
  models:
    - model: claude-opus-5-5
      input: 34
      output: 276
      cache_read: 1749699
      cache_write: 124275
      cost: 0.8361
  strategic:
    - kind: orchestrator
      seconds: 314
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 88
          output: 1334
          cache_read: 31178287
          cache_write: 44197
          cost: 7.6926
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: orchestrator
    at: 2026-10-08T00:26:13Z
  value: 37.5
  by: planner-S-0323
  at: 2026-10-08T00:27:28Z
forecast:
  duration: 24m
  delivery: 2026-10-08T06:20:00Z
  basis: "Its own forecast of 24m; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0309 and S-0319."
  by: flai
  at: 2026-10-08T05:54:06Z
finalized:
  by: orchestrator
  at: 2026-10-08T00:27:54Z
---
# S-0323 flai check finds `narrative.state` outside the story at close-out

## Goal

This story remediates [I-0109](../../../design/issues/I-0109-flai-check-finds-narrative-state-outside-the-story-at-close-out.md), "flai check finds `narrative.state` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0109 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0109 is closed with `flai issue close I-0109 --reason` saying what fixed it
- [ ] I-0111 is closed with `flai issue close I-0109 --reason` saying what fixed it

## Tasks
- T-1313 An ADR refining ADR-0085 and S-0318's ADR records the remedy for I-0109, proposed from its instances
- T-1314 A check scoped to a story leaves out a narrative.state finding on another story's narrative, so a close-out records none
- T-1315 I-0109 is closed with flai issue close, saying that a close-out leaves out a narrative.state on another story's narrative
- T-1316 The convention, the design, and the users' guide say that a close-out leaves out a narrative.state on another story's narrative

## Notes

### Planning

The proposed remedy, which T-1313's ADR decides: a check scoped to a story leaves out a `narrative.state` finding on another story's narrative, as ADR-0122 does for `item.archive`. All three of I-0109's instances are another story in progress, just started, whose narrative still holds the template's placeholder. Only that story's agent can write it, and its own close-out stops on it in the `narrative` step. Every other finding on another story's narrative is still recorded.

It waits for S-0318, which changes the same `ScopeToStory` and adds the test for another open story's narrative that T-1314 reuses.

Where each touch came from. The story declared none, so `flai touches suggest` first listed nothing. Run again from the touches below, its co-change list held only files changed by many unrelated commits (13% at most), so none was added.

| Touch | Source |
|-------|--------|
| `design/adrs` | layout: the ADR T-1313 writes |
| `flai/internal/check/scope.go`, `scope_test.go` | design: `ScopeToStory` leaves out `item.archive` (ADR-0122) |
| `flai/cmd/check.go`, `check_test.go` | design: `recordOutside` and the command's help |
| `docs/users/flai-reference.md` | layout: generated from the help by `make flai-reference` |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | co-change: S-0280 changed all three for `item.archive` |
| `design/system/continuous-improvement.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | co-change: S-0280 changed these for the same rule |
| `design/issues/I-0109-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

One folder touch is kept: `design/adrs`. `flai adr new` allocates the ADR's number when it runs, so no file can be named now.

Figures:

- Forecast: 24m, delivery 2026-10-08T07:48:00Z. `flai forecast` gave 21m (78 s per unit of size times 16). Raised to 24m, the mean of S-0280's 22m and S-0279's 26m, the same shape of remedy.
- Cost of delay: 37.5 USD a week, from `flai cod` on the input `time_lost_per_cycle: 15m`, which the orchestrator set as recommended on TH-0350. Left as `flai cod` gives it.
