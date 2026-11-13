---
id: S-0774
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T15:43:49Z
updated: 2026-11-13T15:47:57Z
transitions:
  - to: ready
    at: 2026-11-13T15:43:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T15:44:00Z
    by: agent-S-0774
  - to: review
    at: 2026-11-13T15:46:08Z
    by: agent-S-0774
  - to: done
    at: 2026-11-13T15:47:57Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0774-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 140
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 2
      work: 25
  models:
    - model: claude-opus-5-5
      input: 58
      output: 7032
      cache_read: 3975331
      cache_write: 163462
      cost: 2.2436
  strategic:
    - kind: orchestrator
      seconds: 84
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 233
          cache_read: 876326
          cache_write: 94206
          cost: 0.2829
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 7155
          cache_write: 6082
          cost: 0.0103
formula: publish
---
# S-0774 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3310 Fetch the remote's branch and tags
- T-3311 Evaluate the release policy, stopping unless it is met
- T-3312 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: 855eccae0ac86b09d4c90c6e4f7476c357c73e3b
- At: 2026-11-13T15:47:57Z

Verdict: S-0774 meets all three criteria; close-out passed at 855eccae0ac86b09d4c90c6e4f7476c357c73e3b.
- 1: T-3310 done
- 2: T-3311 done
- 3: T-3312 done
