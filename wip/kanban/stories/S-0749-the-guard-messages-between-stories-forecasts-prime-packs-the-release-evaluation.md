---
id: S-0749
type: story
nature: feature
title: The guard, messages between stories, forecasts, prime packs, the release evaluation, the digest, and turn counts read a state's kind instead of its name
status: review
parent: E-0034
owner: alex
created: 2026-11-12T16:00:56Z
updated: 2026-11-15T20:14:37Z
transitions:
  - to: ready
    at: 2026-11-15T02:00:15Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-15T19:44:10Z
    by: agent-S-0749
  - to: review
    at: 2026-11-15T20:14:37Z
    by: agent-S-0749
tags: [cli]
topics: [cli, planning]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/guard.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/call.go, flai/internal/guard/guard_test.go, flai/internal/messages/messages.go, flai/internal/messages/share.go, flai/internal/messages/messages_test.go, flai/internal/messages/share_test.go, flai/internal/planning/forecast.go, flai/internal/planning/forecast_test.go, flai/internal/context/reasons.go, flai/internal/context/reasons_figures.go, flai/internal/context/story.go, flai/internal/release/evaluate.go, flai/internal/digest/sources.go, flai/internal/usage/turns.go, flai/internal/statsread/statsread.go, flai/cmd/forecast.go, flai/internal/context/context_workflow_test.go, flai/internal/release/evaluate_workflow_test.go, flai/internal/statsread/statsread_workflow_test.go, design/issues/I-0379-callers-of-the-workflow-aware-guard-forecast-replay-and-share-hold-still-pass.md, design/issues/summary.md]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1851
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 3
      work: 47
  models:
    - model: claude-opus-5-5
      input: 364
      output: 116637
      cache_read: 19603297
      cache_write: 1019893
      cost: 12.1698
  strategic:
    - kind: orchestrator
      seconds: 84
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 117
          cache_read: 583665
          cache_write: 111182
          cost: 0.207
        - model: claude-sonnet-5-5
          input: 8
          output: 64
          cache_read: 86225
          cache_write: 28718
          cost: 0.0907
cost_of_delay:
  value: 74.82
  by: system-flow
  at: 2026-11-14T17:44:10Z
forecast:
  duration: 42m
  delivery: 2026-11-15T20:33:00Z
  basis: "Its own forecast of 42m; 1st in the pull order with an in-progress limit of 3, behind S-0741 and S-0743."
  by: flai
  at: 2026-11-15T19:43:53Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:56Z
---
# S-0749 The guard, messages between stories, forecasts, prime packs, the release evaluation, the digest, and turn counts read a state's kind instead of its name

## Goal

The rest of flai's packages that name a state follow the workflow. The guard tells a story's agent's moves apart by who may make each transition, a conversation between two stories closes when either leaves the open kinds, a forecast plays out the pull order through the workflow's active-work states and their WIP limit, a prime pack and the orchestrator's reasons pick their items by kind, and the release evaluation, the digest's sources, and the turn counts tell accepted from cancelled by kind.

## Acceptance criteria

- [x] `flai/internal/guard`, `messages`, `planning/forecast.go`, `context`, `release/evaluate.go`, `digest/sources.go`, `usage/turns.go`, and `statsread/statsread.go` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [x] `flai guard` refuses an agent a move by who may make its transition, with today's refusals word for word with the baseline
- [x] A conversation closes when either story reaches an accepted or a cancelled state, and `flai forecast` plays out the pull order through the active-work states with their WIP limit
- [x] Tests run the guard's move refusals, a conversation's close, and a forecast against a workflow with a state renamed and one added between `in-progress` and `review`
- [x] With the baseline, every existing test of these packages passes unchanged

## Tasks

- T-3232 flai guard refuses an agent a move by who may make its transition, with today's refusals on the baseline
- T-3233 A conversation between two stories closes when either reaches an accepted or cancelled state, and a share lasts while both are open
- T-3234 flai forecast plays out the pull order through the workflow's active-work states and their WIP limit
- T-3235 Prime packs, the release evaluation, the digest's sources, turn counts, and statsread read a state's kind

## Notes

- `messages` holds 14 references, `planning/forecast.go` 9, `guard` 6, `context` 3; the others one each. The role conditions in `guard/roles.go` are the role condition story's.

### Planning

- Touches, all files:
  - Layout: the files of these packages that name a state, with the tests of the guard, the messages, and the forecast.
  - `guard/roles.go` is left out: its role conditions are S-0748's.
  - Left out from `flai touches suggest`: `cmd/guard.go` and its test (17% to 27%), `mcpserver`, and `metrics/metrics.go` (S-0750's), none of which this story's change reaches.
- `after`: S-0737. The guard reads who may make each transition through S-0737's mover check, and the others read kinds and WIP limits from the workflow S-0737 gives.
- Tasks: T-3232 (guard), T-3233 (messages), T-3234 (forecast), and T-3235 (the rest), one layer: no two share a path.
- Forecast 42m, `flai forecast`'s rate for size 23, up from its 40m for the first draft's 22.
- Cost of delay: 74.82 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 42m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
