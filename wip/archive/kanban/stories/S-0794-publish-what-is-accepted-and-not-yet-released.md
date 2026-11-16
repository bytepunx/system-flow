---
id: S-0794
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-15T21:47:08Z
updated: 2026-11-16T00:44:38Z
transitions:
  - to: ready
    at: 2026-11-15T21:47:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T00:36:45Z
    by: agent-S-0794
  - to: review
    at: 2026-11-16T00:39:49Z
    by: agent-S-0794
  - to: done
    at: 2026-11-16T00:44:38Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0794-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 210
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 2
      work: 28
  models:
    - model: claude-opus-5-5
      input: 64
      output: 7825
      cache_read: 4359087
      cache_write: 165677
      cost: 2.354
  strategic:
    - kind: orchestrator
      seconds: 85
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 190
          cache_read: 621414
          cache_write: 94582
          cost: 0.2171
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 0
          cache_write: 13437
          cost: 0.0107
formula: publish
---
# S-0794 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3387 Fetch the remote's branch and tags
- T-3388 Evaluate the release policy, stopping unless it is met
- T-3389 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-16 (T-3389). `flai release --pending` pushed the branch and the tags.

| Component | From | To | Tag or publish | Items |
|-----------|------|----|----------------|-------|
| flai | 1.100.0 | 1.101.0 | `flai/v1.101.0` | S-0755 (minor), S-0790 (patch) |
| template | 1.1.42 | 1.1.43 | `git@github.com:bytepunx/system-flow-template.git` v1.1.43 (9da14f1) | S-0755 (patch) |

### Accepted by the orchestrator

- Verified: 72621b47683b779ed045dc71ebf7deb9e61c296d
- Batch: 20261116-004429
- Enqueued verified: 36c00764409212170a4f99e2d78c1f73832510a4
- At: 2026-11-16T00:44:38Z

Verdict: pass. Formula instance; close-out passed at 36c0076, every step's task done.
- 1: T-3387 done
- 2: T-3388 done
- 3: T-3389 done
