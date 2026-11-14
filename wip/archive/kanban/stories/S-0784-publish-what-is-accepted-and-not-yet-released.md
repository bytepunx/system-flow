---
id: S-0784
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-14T07:05:35Z
updated: 2026-11-14T07:11:09Z
transitions:
  - to: ready
    at: 2026-11-14T07:05:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T07:05:55Z
    by: agent-S-0784
  - to: review
    at: 2026-11-14T07:09:01Z
    by: agent-S-0784
  - to: done
    at: 2026-11-14T07:11:09Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0784-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 207
  priced_by: harness
  turns:
    - day: 2026-11-14
      ceremony: 4
      hand_edits: 1
      work: 26
  models:
    - model: claude-opus-5-5
      input: 68
      output: 7281
      cache_read: 4804701
      cache_write: 164802
      cost: 2.4252
  strategic:
    - kind: orchestrator
      seconds: 119
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 30
          output: 247
          cache_read: 1031992
          cache_write: 120976
          cost: 0.3411
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 0
          cache_write: 12965
          cost: 0.0101
formula: publish
---
# S-0784 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3351 Fetch the remote's branch and tags
- T-3352 Evaluate the release policy, stopping unless it is met
- T-3353 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: 0ce7e92b0b23f43411c41d1bd29de5fa7bdfe7b9
- At: 2026-11-14T07:11:09Z

Verdict: S-0784 meets all three acceptance criteria; its close-out passed every step at 0ce7e92.
- 1: T-3351 done
- 2: T-3352 done
- 3: T-3353 done
