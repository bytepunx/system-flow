---
id: S-0775
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T16:13:36Z
updated: 2026-11-13T16:17:57Z
transitions:
  - to: ready
    at: 2026-11-13T16:13:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T16:13:47Z
    by: agent-S-0775
  - to: review
    at: 2026-11-13T16:16:09Z
    by: agent-S-0775
  - to: done
    at: 2026-11-13T16:17:57Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0775-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 153
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 27
  models:
    - model: claude-opus-5-5
      input: 60
      output: 7370
      cache_read: 4094330
      cache_write: 160402
      cost: 2.2497
  strategic:
    - kind: orchestrator
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 33
          cache_read: 193500
          cache_write: 29748
          cost: 0.0651
formula: publish
---
# S-0775 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3313 Fetch the remote's branch and tags
- T-3314 Evaluate the release policy, stopping unless it is met
- T-3315 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: 017d98a96adc8aa7001865827f41c76f8254bd4b
- At: 2026-11-13T16:17:57Z

Verdict: pass, all three steps done with exit 0; flai/v1.94.0 tagged and pushed for S-0751
- 1: T-3313 done
- 2: T-3314 done
- 3: T-3315 done
