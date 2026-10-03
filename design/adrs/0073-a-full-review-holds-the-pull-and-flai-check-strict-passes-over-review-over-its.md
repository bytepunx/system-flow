---
id: ADR-0073
title: "A full review holds the pull, and flai check --strict passes over review over its limit"
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
refines: [ADR-0043]
---

# ADR-0073 A full review holds the pull, and flai check --strict passes over review over its limit

## Context

I-0007: the agent moves stories to review as they finish, and the operator accepts them in batches. Pulling looked only at the in-progress limit ([ADR-0043](0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md): "the in-progress limit is the only thing that holds back a ready story"), so `flai serve`, `inbox`, and `wait_for_work` kept starting stories while review grew without end. `flai check --strict` fails on every warning, and `board.wip-limit` warns when review is over its limit. Once it was, every other story's close-out (`scripts/close-out.sh`) stopped on a warning its agent could not clear, since only acceptance clears it. Raising the limit only moves the threshold. The designer chose this decision on TH-0078.

## Decision

A full review holds the pull. While review is at or over its limit, `flai serve` starts no agent for a ready story and says why ("review is full (5 of 5): accept or send back a story"). `inbox` and `wait_for_work` report `can_pull` false, with the reason in `pull_hold`, and `wait_for_work` waits with `waiting_for` `review`. Stories already in progress, and agents already started, go on. The operator's start-now goes past a full review with a warning, as it goes past a full in-progress limit. A queued retry waits for room.

`flai check` still warns `board.wip-limit` for review over its limit, but `--strict` passes over it: the result counts it in `advisory`, and the summary line says so. Ready and in-progress over their limits still fail `--strict`.

## Consequences

- Review stays near its limit while acceptance is batched: at most the stories in progress when it filled can join it.
- A story's close-out no longer stops on the review column's size, which is the operator's queue, not the story's.
- An idle agent waits for acceptance rather than pulling; the board shows review at its limit, and `waiting_for` says what it waits on.
- `flai check --json` gains `advisory`; `warnings` still counts every warning.

## Alternatives considered

- **Raise the review limit to match the operator's cadence.** The breach only comes later, and the close-out still stops when it does.
- **Only the pull hold.** A story in progress that finishes into a full review still takes it over the limit, and the next close-out stops.
- **Only the strict change.** Review grows without end while acceptance is batched.
