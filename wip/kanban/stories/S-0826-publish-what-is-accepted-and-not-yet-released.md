---
id: S-0826
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-17T08:16:46Z
updated: 2026-11-17T08:19:56Z
transitions:
  - to: ready
    at: 2026-11-17T08:16:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:17:31Z
    by: agent-S-0826
tags: []
touches: [wip/kanban/stories/S-0826-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 149
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-17
      ceremony: 2
      hand_edits: 1
      work: 17
  models:
    - model: claude-opus-5-5
      input: 42
      output: 445
      cache_read: 2601429
      cache_write: 164005
      cost: 1.0538
  strategic:
    - kind: orchestrator
      seconds: 39
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 46
          cache_read: 370539
          cache_write: 41784
          cost: 0.1277
        - model: claude-sonnet-5-5
          input: 1
          output: 12
          cache_read: 7012
          cache_write: 8411
          cost: 0.0124
formula: publish
---
# S-0826 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3517 Fetch the remote's branch and tags
- T-3518 Evaluate the release policy, stopping unless it is met
- T-3519 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-17 by step 3, with S-0750 and S-0805, and pushed:

| Component | From | To | Tag or version |
|-----------|------|----|----------------|
| flai | 1.103.0 | 1.104.0 | `flai/v1.104.0` |
| template | 1.1.44 | 1.1.45 | template repository v1.1.45 (315ea50) |
