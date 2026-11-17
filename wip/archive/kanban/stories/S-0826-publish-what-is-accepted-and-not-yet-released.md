---
id: S-0826
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-17T08:16:46Z
updated: 2026-11-17T08:22:23Z
transitions:
  - to: ready
    at: 2026-11-17T08:16:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:17:31Z
    by: agent-S-0826
  - to: review
    at: 2026-11-17T08:20:16Z
    by: agent-S-0826
  - to: done
    at: 2026-11-17T08:22:23Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0826-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 201
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 3
      hand_edits: 1
      work: 25
  models:
    - model: claude-opus-5-5
      input: 62
      output: 8521
      cache_read: 4455187
      cache_write: 174268
      cost: 2.4558
  strategic:
    - kind: orchestrator
      seconds: 56
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 53
          cache_read: 480872
          cache_write: 43513
          cost: 0.1624
        - model: claude-sonnet-5-5
          input: 1
          output: 12
          cache_read: 7012
          cache_write: 8411
          cost: 0.0124
formula: publish
---
# S-0826 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3517 Fetch the remote's branch and tags
- T-3518 Evaluate the release policy, stopping unless it is met
- T-3519 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-17 by step 3, with S-0750 and S-0805, and pushed:

| Component | From | To | Tag or version |
|-----------|------|----|----------------|
| flai | 1.103.0 | 1.104.0 | `flai/v1.104.0` |
| template | 1.1.44 | 1.1.45 | template repository v1.1.45 (315ea50) |

### Accepted by the orchestrator

- Verified: d399d4365fb3f7e51f171fbd04b0ab89945a3227
- At: 2026-11-17T08:22:23Z

Verdict: pass. Every step of the publish instance is done, and flai/v1.104.0 is tagged and pushed.
- 1: T-3517 done
- 2: T-3518 done
- 3: T-3519 done
