---
id: S-0791
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-15T07:18:49Z
updated: 2026-11-15T21:16:29Z
transitions:
  - to: ready
    at: 2026-11-15T07:18:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-15T20:59:23Z
    by: agent-S-0791
  - to: review
    at: 2026-11-15T21:01:55Z
    by: agent-S-0791
  - to: done
    at: 2026-11-15T21:16:29Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0791-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 162
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 2
      work: 25
  models:
    - model: claude-opus-5-5
      input: 58
      output: 8566
      cache_read: 4002407
      cache_write: 165733
      cost: 2.2979
  strategic:
    - kind: orchestrator
      seconds: 99
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 190
          cache_read: 1135219
          cache_write: 92760
          cost: 0.3699
formula: publish
---
# S-0791 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3373 Fetch the remote's branch and tags
- T-3374 Evaluate the release policy, stopping unless it is met
- T-3375 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-15: flai 1.97.0 → 1.98.0 (minor), tag `flai/v1.98.0`, with S-0743, S-0744, and S-0749. The tag and the branch were pushed.

### Accepted by the orchestrator

- Verified: d65c6bd7ea7cf344c0bcccf6353dcfe5dda07284
- Batch: 20261115-211624
- Enqueued verified: 93f1951aaa67a2396ed5a78ba5d540a1c27afd86
- At: 2026-11-15T21:16:29Z

Verdict: pass at 93f1951aaa67a2396ed5a78ba5d540a1c27afd86
- 1: T-3373 done
- 2: T-3374 done
- 3: T-3375 done
