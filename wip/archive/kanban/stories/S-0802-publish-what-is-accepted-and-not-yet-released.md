---
id: S-0802
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-16T04:12:37Z
updated: 2026-11-16T04:19:55Z
transitions:
  - to: ready
    at: 2026-11-16T04:12:39Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T04:12:46Z
    by: agent-S-0802
  - to: review
    at: 2026-11-16T04:14:53Z
    by: agent-S-0802
  - to: done
    at: 2026-11-16T04:19:55Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0802-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 139
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 2
      work: 23
  models:
    - model: claude-opus-5-5
      input: 56
      output: 7259
      cache_read: 3783311
      cache_write: 163669
      cost: 2.2114
  strategic:
    - kind: orchestrator
      seconds: 59
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 179
          cache_read: 693280
          cache_write: 92102
          cost: 0.24
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 0
          cache_write: 13342
          cost: 0.0106
formula: publish
---
# S-0802 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3423 Fetch the remote's branch and tags
- T-3424 Evaluate the release policy, stopping unless it is met
- T-3425 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-16 from the main checkout, pushed to the remote: S-0799 as patches, `flai/v1.101.2` (flai 1.101.1 → 1.101.2) and `flaiover/v0.61.2` (flaiover 0.61.1 → 0.61.2).

### Accepted by the orchestrator

- Verified: 56ca61e99019665bc78109ed576efe55d81f6e8f
- Batch: 20261116-041948
- Enqueued verified: 56ca61e99019665bc78109ed576efe55d81f6e8f
- At: 2026-11-16T04:19:55Z

Verdict: pass, every step of the publish formula done and the close-out passed at 56ca61e, the branch's head
- 1: T-3423 done
- 2: T-3424 done
- 3: T-3425 done
