---
id: S-0773
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T15:20:30Z
updated: 2026-11-13T15:25:41Z
transitions:
  - to: ready
    at: 2026-11-13T15:20:32Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T15:21:24Z
    by: agent-S-0773
  - to: review
    at: 2026-11-13T15:24:15Z
    by: agent-S-0773
  - to: done
    at: 2026-11-13T15:25:41Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0773-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 189
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 2
      work: 27
  models:
    - model: claude-opus-5-5
      input: 62
      output: 8619
      cache_read: 4229001
      cache_write: 165564
      cost: 2.3429
  strategic:
    - kind: orchestrator
      seconds: 110
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 27
          output: 151
          cache_read: 1039666
          cache_write: 82186
          cost: 0.3265
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 7155
          cache_write: 5741
          cost: 0.01
formula: publish
---
# S-0773 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3307 Fetch the remote's branch and tags
- T-3308 Evaluate the release policy, stopping unless it is met
- T-3309 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: a34dfb428f53edd128cdb0518643fa25e5d0a502
- At: 2026-11-13T15:25:41Z

Verdict: pass at a34dfb428f53edd128cdb0518643fa25e5d0a502
- 1: T-3307 done
- 2: T-3308 done
- 3: T-3309 done
