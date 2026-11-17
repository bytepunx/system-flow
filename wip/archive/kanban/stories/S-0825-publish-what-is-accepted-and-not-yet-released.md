---
id: S-0825
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-17T07:54:16Z
updated: 2026-11-17T08:20:01Z
transitions:
  - to: ready
    at: 2026-11-17T07:54:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:14:03Z
    by: agent-S-0825
  - to: review
    at: 2026-11-17T08:17:03Z
    by: agent-S-0825
  - to: done
    at: 2026-11-17T08:20:01Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0825-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 192
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 2
      work: 23
  models:
    - model: claude-opus-5-5
      input: 54
      output: 7690
      cache_read: 3842183
      cache_write: 179117
      cost: 2.3554
  strategic:
    - kind: orchestrator
      seconds: 74
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 19
          output: 127
          cache_read: 784560
          cache_write: 43270
          cost: 0.2564
        - model: claude-sonnet-5-5
          input: 2
          output: 18
          cache_read: 10734
          cache_write: 4191
          cost: 0.012
formula: publish
---
# S-0825 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3513 Fetch the remote's branch and tags
- T-3514 Evaluate the release policy, stopping unless it is met
- T-3515 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-17 by step 3, `flai release --pending`, which pushed the branch and these tags:

| Component | From | To | Tag | Items |
|-----------|------|----|-----|-------|
| flai | 1.102.0 | 1.103.0 | `flai/v1.103.0` | S-0742, S-0753, S-0817, S-0819 |
| flaiover | 0.61.3 | 0.62.0 | `flaiover/v0.62.0` | S-0753 |

### Accepted by the orchestrator

- Verified: e2e59593c729dbacc223bb942d5b5b7e5fcceb97
- At: 2026-11-17T08:20:01Z

Verdict: pass, all three steps done and the close-out passed at e2e5959
- 1: T-3513 done
- 2: T-3514 done
- 3: T-3515 done
