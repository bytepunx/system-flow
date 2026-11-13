---
id: S-0749
type: story
nature: feature
title: The guard, messages between stories, forecasts, prime packs, the release evaluation, the digest, and turn counts read a state's kind instead of its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:56Z
updated: 2026-11-13T09:17:06Z
transitions: []
tags: [cli]
topics: [cli, planning]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/guard.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/call.go, flai/internal/guard/guard_test.go, flai/internal/messages/messages.go, flai/internal/messages/share.go, flai/internal/messages/messages_test.go, flai/internal/messages/share_test.go, flai/internal/planning/forecast.go, flai/internal/planning/forecast_test.go, flai/internal/context/reasons.go, flai/internal/context/reasons_figures.go, flai/internal/context/story.go, flai/internal/release/evaluate.go, flai/internal/digest/sources.go, flai/internal/usage/turns.go, flai/internal/statsread/statsread.go]
after: [S-0737]
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
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 108313
          cache_write: 18018
          cost: 0.0361
draft: true
forecast:
  duration: 42m
  delivery: 2026-11-13T18:21:00Z
  basis: "Its own forecast of 42m; 23rd in the pull order with an in-progress limit of 3, behind S-0763, S-0760, S-0745, S-0759, S-0751, S-0757, S-0740, S-0661, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747 and S-0748."
  by: flai
  at: 2026-11-13T09:17:06Z
---
# S-0749 The guard, messages between stories, forecasts, prime packs, the release evaluation, the digest, and turn counts read a state's kind instead of its name

## Goal

The rest of flai's packages that name a state follow the workflow. The guard tells a story's agent's moves apart by who may make each transition, a conversation between two stories closes when either leaves the open kinds, a forecast plays out the pull order through the workflow's active-work states and their WIP limit, a prime pack and the orchestrator's reasons pick their items by kind, and the release evaluation, the digest's sources, and the turn counts tell accepted from cancelled by kind.

## Acceptance criteria

- [ ] `flai/internal/guard`, `messages`, `planning/forecast.go`, `context`, `release/evaluate.go`, `digest/sources.go`, `usage/turns.go`, and `statsread/statsread.go` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [ ] `flai guard` refuses an agent a move by who may make its transition, with today's refusals word for word with the baseline
- [ ] A conversation closes when either story reaches an accepted or a cancelled state, and `flai forecast` plays out the pull order through the active-work states with their WIP limit
- [ ] Tests run the guard's move refusals, a conversation's close, and a forecast against a workflow with a state renamed and one added between `in-progress` and `review`
- [ ] With the baseline, every existing test of these packages passes unchanged

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
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
