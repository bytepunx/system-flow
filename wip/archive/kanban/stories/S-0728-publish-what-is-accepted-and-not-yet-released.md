---
id: S-0728
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-11T23:40:03Z
updated: 2026-11-11T23:44:17Z
transitions:
  - to: ready
    at: 2026-11-11T23:40:04Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T23:40:13Z
    by: agent-S-0728
  - to: review
    at: 2026-11-11T23:42:25Z
    by: agent-S-0728
  - to: done
    at: 2026-11-11T23:44:17Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0728-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 146
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 2
      work: 25
  models:
    - model: claude-opus-5-5
      input: 58
      output: 6868
      cache_read: 3824748
      cache_write: 154205
      cost: 2.1362
  strategic:
    - kind: orchestrator
      seconds: 102
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 31
          output: 185
          cache_read: 1072204
          cache_write: 80951
          cost: 0.3265
        - model: claude-sonnet-5-5
          input: 4
          output: 32
          cache_read: 19996
          cache_write: 23382
          cost: 0.0332
formula: publish
---
# S-0728 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3127 Fetch the remote's branch and tags
- T-3128 Evaluate the release policy, stopping unless it is met
- T-3129 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: 45d9d08430a9babaddc357ac9d69dfc3c4ff5dbb
- At: 2026-11-11T23:44:17Z

Verdict: pass, every step of the publish formula's instance done with exit 0 at 45d9d08
- 1: T-3127 done
- 2: T-3128 done
- 3: T-3129 done
