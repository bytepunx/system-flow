---
id: S-0776
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T17:20:42Z
updated: 2026-11-13T17:24:08Z
transitions:
  - to: ready
    at: 2026-11-13T17:20:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T17:20:53Z
    by: agent-S-0776
  - to: review
    at: 2026-11-13T17:22:59Z
    by: agent-S-0776
  - to: done
    at: 2026-11-13T17:24:08Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0776-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 135
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 24
  models:
    - model: claude-opus-5-5
      input: 54
      output: 7195
      cache_read: 3670837
      cache_write: 157591
      cost: 2.139
  strategic:
    - kind: orchestrator
      seconds: 45
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 91
          cache_read: 508448
          cache_write: 36276
          cost: 0.1591
formula: publish
---
# S-0776 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3316 Fetch the remote's branch and tags
- T-3317 Evaluate the release policy, stopping unless it is met
- T-3318 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: 4f9055659e445cd01d17972211f96e71a3df479e
- At: 2026-11-13T17:24:08Z

Verdict: pass — publish formula instance; close-out passed at 4f90556, every step done
- 1: T-3316 done
- 2: T-3317 done
- 3: T-3318 done
