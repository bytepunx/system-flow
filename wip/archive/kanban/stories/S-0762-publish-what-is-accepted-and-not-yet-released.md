---
id: S-0762
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T07:24:38Z
updated: 2026-11-13T07:29:50Z
transitions:
  - to: ready
    at: 2026-11-13T07:24:39Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T07:24:46Z
    by: agent-S-0762
  - to: review
    at: 2026-11-13T07:27:20Z
    by: agent-S-0762
  - to: done
    at: 2026-11-13T07:29:50Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0762-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 166
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 2
      hand_edits: 3
      work: 19
  models:
    - model: claude-opus-5-5
      input: 52
      output: 7953
      cache_read: 3521253
      cache_write: 167911
      cost: 2.2068
  strategic:
    - kind: orchestrator
      seconds: 56
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 107
          cache_read: 704853
          cache_write: 39807
          cost: 0.2135
formula: publish
---
# S-0762 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3265 Fetch the remote's branch and tags
- T-3266 Evaluate the release policy, stopping unless it is met
- T-3267 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published 2026-11-13 by T-3267, pushed:

| Component | From | To | Tag or version | Items |
|-----------|------|----|----------------|-------|
| flai | 1.86.22 | 1.86.23 | `flai/v1.86.23` | S-0683, S-0723, S-0726, S-0732 |
| template | 1.1.34 | 1.1.35 | `template/template.yaml`, template repository v1.1.35 at 2c78d45 | S-0726 |

### Accepted by the orchestrator

- Verified: a84b3003d0222e4ff2544abf68e57d2237cc3ef9
- At: 2026-11-13T07:29:50Z

Verdict: met at a84b3003d0222e4ff2544abf68e57d2237cc3ef9
- 1: T-3265 done
- 2: T-3266 done
- 3: T-3267 done
