---
id: S-0765
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T09:17:34Z
updated: 2026-11-13T09:21:33Z
transitions:
  - to: ready
    at: 2026-11-13T09:17:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T09:17:43Z
    by: agent-S-0765
  - to: review
    at: 2026-11-13T09:19:33Z
    by: system-flow
  - to: done
    at: 2026-11-13T09:21:33Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0765-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 119
  priced_by: harness
  turns:
    - day: 2026-11-13
      work: 21
  models:
    - model: claude-opus-5-5
      input: 46
      output: 6439
      cache_read: 3098331
      cache_write: 167597
      cost: 2.0894
  strategic:
    - kind: orchestrator
      seconds: 33
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 40
          cache_read: 411060
          cache_write: 4448
          cost: 0.1199
formula: publish
---
# S-0765 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3281 Fetch the remote's branch and tags
- T-3282 Evaluate the release policy, stopping unless it is met
- T-3283 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: ed72e7710ec169379dbe9576ccb2cfc138b971b5
- At: 2026-11-13T09:21:33Z

Verdict: pass, every step of the publish formula done with exit 0; close-out passed at ed72e77
- 1: T-3281 done
- 2: T-3282 done
- 3: T-3283 done
