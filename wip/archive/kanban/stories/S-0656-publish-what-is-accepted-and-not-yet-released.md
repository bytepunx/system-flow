---
id: S-0656
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-10T06:45:10Z
updated: 2026-11-10T07:10:30Z
transitions:
  - to: ready
    at: 2026-11-10T06:45:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T06:58:07Z
    by: agent-S-0656
  - to: review
    at: 2026-11-10T07:03:21Z
    by: system-flow
  - to: done
    at: 2026-11-10T07:10:30Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0656-publish-what-is-accepted-and-not-yet-released.md, design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 323
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 3
      work: 26
  models:
    - model: claude-opus-5-5
      input: 62
      output: 7513
      cache_read: 4295104
      cache_write: 163621
      cost: 2.3185
  strategic:
    - kind: orchestrator
      seconds: 16
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 33
          cache_read: 127642
          cache_write: 16888
          cost: 0.0388
formula: publish
---
# S-0656 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2831 Fetch the remote's branch and tags
- T-2832 Evaluate the release policy, stopping unless it is met
- T-2833 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: ea79e17c0cf7e047d2d5c296adac3b19bb8b6450
- At: 2026-11-10T07:10:30Z

Verdict: S-0656 meets all three criteria; the branch changes only paths in its touches and the last flai verify passed at ea79e17.
- 1: T-2831 done
- 2: T-2832 done
- 3: T-2833 done
