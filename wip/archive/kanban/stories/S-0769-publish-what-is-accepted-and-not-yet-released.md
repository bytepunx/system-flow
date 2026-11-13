---
id: S-0769
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T12:40:07Z
updated: 2026-11-13T12:56:30Z
transitions:
  - to: ready
    at: 2026-11-13T12:40:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T12:52:42Z
    by: agent-S-0769
  - to: review
    at: 2026-11-13T12:55:06Z
    by: agent-S-0769
  - to: done
    at: 2026-11-13T12:56:30Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0769-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 156
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      hand_edits: 4
      work: 20
  models:
    - model: claude-opus-5-5
      input: 56
      output: 7560
      cache_read: 3916272
      cache_write: 180499
      cost: 2.3787
  strategic:
    - kind: orchestrator
      seconds: 42
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 15
          output: 76
          cache_read: 470284
          cache_write: 47425
          cost: 0.15
formula: publish
---
# S-0769 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3294 Fetch the remote's branch and tags
- T-3295 Evaluate the release policy, stopping unless it is met
- T-3296 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published 2026-11-13: `flai/v1.90.0` (minor, from 1.89.1), for S-0763 (minor) and S-0767 (patch); branch and tag pushed. No other component released.

### Accepted by the orchestrator

- Verified: 79002dbb0e4ce8e6b9a6e116338406321782154c
- At: 2026-11-13T12:56:30Z

Verdict: pass
- 1: T-3294 done
- 2: T-3295 done
- 3: T-3296 done
