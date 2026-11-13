---
id: S-0771
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T13:41:40Z
updated: 2026-11-13T14:38:52Z
transitions:
  - to: ready
    at: 2026-11-13T13:41:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T14:35:00Z
    by: agent-S-0771
  - to: review
    at: 2026-11-13T14:37:57Z
    by: agent-S-0771
  - to: done
    at: 2026-11-13T14:38:52Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0771-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 190
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 3
      work: 24
  models:
    - model: claude-opus-5-5
      input: 60
      output: 9279
      cache_read: 4043033
      cache_write: 175952
      cost: 2.402
  strategic:
    - kind: orchestrator
      seconds: 75
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 216
          cache_read: 838443
          cache_write: 71554
          cost: 0.2645
formula: publish
---
# S-0771 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3300 Fetch the remote's branch and tags
- T-3301 Evaluate the release policy, stopping unless it is met
- T-3302 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-13 by step 3, `flai release --pending --require-pending --json`, which exited 0 and pushed:

| Component | Version | Tag | Items |
|-----------|---------|-----|-------|
| flai | 1.91.0 to 1.92.0 (minor) | `flai/v1.92.0` | S-0661, S-0745 |
| template | 1.1.36 to 1.1.37 (patch) | `v1.1.37` (8693665) on `bytepunx/system-flow-template` | S-0745 |

### Accepted by the orchestrator

- Verified: 9194f2cd6ed28144fdd16a182f9436debaa0e171
- At: 2026-11-13T14:38:52Z

Verdict: pass — formula instance of publish; every step's task is done and the close-out passed at 9194f2cd.
- 1: T-3300 done
- 2: T-3301 done
- 3: T-3302 done
