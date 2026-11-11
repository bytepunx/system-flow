---
id: S-0704
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: review
owner: alex
created: 2026-11-10T08:56:47Z
updated: 2026-11-11T01:44:01Z
transitions:
  - to: ready
    at: 2026-11-10T08:56:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T01:41:18Z
    by: agent-S-0704
  - to: review
    at: 2026-11-11T01:44:01Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0704-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 175
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 1
      work: 32
  models:
    - model: claude-opus-5-5
      input: 72
      output: 10331
      cache_read: 5195203
      cache_write: 175755
      cost: 2.652
  strategic:
    - kind: orchestrator
      seconds: 231
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 47
          output: 256
          cache_read: 1743140
          cache_write: 159532
          cost: 0.5198
        - model: claude-sonnet-5-5
          input: 7
          output: 57
          cache_read: 75831
          cache_write: 23201
          cost: 0.0743
formula: publish
---
# S-0704 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2952 Fetch the remote's branch and tags
- T-2953 Evaluate the release policy, stopping unless it is met
- T-2954 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published 2026-11-11: flai 1.86.13 → 1.86.14, tag `flai/v1.86.14`, pushed with the branch. It bundles S-0665, S-0666, S-0675, S-0677, and S-0679. The policy was `threshold`, met at 5 stories worth 1087.58 USD/week. Each step's bare `flai` ran as the main checkout's `scripts/flai.sh`, since the installed flai (1.86.12) was behind the project's newest release (S-0681).
