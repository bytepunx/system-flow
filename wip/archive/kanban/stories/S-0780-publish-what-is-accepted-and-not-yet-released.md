---
id: S-0780
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T20:51:49Z
updated: 2026-11-13T20:58:01Z
transitions:
  - to: ready
    at: 2026-11-13T20:51:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:52:08Z
    by: agent-S-0780
  - to: review
    at: 2026-11-13T20:55:47Z
    by: agent-S-0780
  - to: done
    at: 2026-11-13T20:58:01Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0780-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 246
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 2
      work: 27
  models:
    - model: claude-opus-5-5
      input: 62
      output: 8955
      cache_read: 4531058
      cache_write: 167554
      cost: 2.426
  strategic:
    - kind: orchestrator
      seconds: 76
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 129
          cache_read: 763600
          cache_write: 110916
          cost: 0.2576
formula: publish
---
# S-0780 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3334 Fetch the remote's branch and tags
- T-3335 Evaluate the release policy, stopping unless it is met
- T-3336 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-13: flai 1.95.1 → 1.95.2, tag `flai/v1.95.2`, pushed, releasing S-0777. The threshold policy was met with 1 story pending at a threshold of 1.

### Accepted by the orchestrator

- Verified: 5ff06e676f0c91df7bf6536e8af1349907495c5f
- At: 2026-11-13T20:58:01Z

Verdict: pass at 5ff06e676f0c91df7bf6536e8af1349907495c5f
- 1: T-3334 done
- 2: T-3335 done
- 3: T-3336 done
