---
id: S-0818
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-17T02:27:29Z
updated: 2026-11-17T02:32:36Z
transitions:
  - to: ready
    at: 2026-11-17T02:27:30Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T02:27:53Z
    by: agent-S-0818
  - to: review
    at: 2026-11-17T02:31:45Z
    by: agent-S-0818
  - to: done
    at: 2026-11-17T02:32:36Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0818-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 259
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 4
      work: 28
  models:
    - model: claude-opus-5-5
      input: 70
      output: 8269
      cache_read: 4864062
      cache_write: 165397
      cost: 2.4616
  strategic:
    - kind: orchestrator
      seconds: 103
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 27
          output: 181
          cache_read: 999988
          cache_write: 166346
          cost: 0.3586
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 0
          cache_write: 13391
          cost: 0.0107
formula: publish
---
# S-0818 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3492 Fetch the remote's branch and tags
- T-3493 Evaluate the release policy, stopping unless it is met
- T-3494 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published 2026-11-17 by step 3 (`flai release --pending`, exit 0):

| Component | From | To | Tag | Items |
|-----------|------|----|-----|-------|
| flai | 1.101.3 | 1.101.4 | `flai/v1.101.4` | S-0797 (patch) |

The branch and the tag were pushed (`pushed: true`).

### Accepted by the orchestrator

- Verified: 0385a7c04d9bf34808c7cf864bd8516ded64f8d4
- At: 2026-11-17T02:32:36Z

Verdict: PASS. All three criteria met at 0385a7c, the head of story/S-0818; the close-out passed every step.
- 1: T-3492 done
- 2: T-3493 done
- 3: T-3494 done
