---
id: S-0749
type: story
nature: feature
title: The guard, messages between stories, forecasts, prime packs, the release evaluation, the digest, and turn counts read a state's kind instead of its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:56Z
updated: 2026-11-12T16:00:56Z
transitions: []
tags: [cli]
topics: [cli, planning]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/guard.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/call.go, flai/internal/guard/guard_test.go, flai/internal/messages/messages.go, flai/internal/messages/share.go, flai/internal/messages/messages_test.go, flai/internal/planning/forecast.go, flai/internal/planning/forecast_test.go, flai/internal/context/reasons.go, flai/internal/context/reasons_figures.go, flai/internal/context/story.go, flai/internal/release/evaluate.go, flai/internal/digest/sources.go, flai/internal/usage/turns.go, flai/internal/statsread/statsread.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

## Notes

- `messages` holds 14 references, `planning/forecast.go` 9, `guard` 6, `context` 3; the others one each. The role conditions in `guard/roles.go` are the role condition story's.
