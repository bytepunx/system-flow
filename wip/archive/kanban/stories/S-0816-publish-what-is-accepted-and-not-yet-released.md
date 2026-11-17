---
id: S-0816
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-17T01:33:37Z
updated: 2026-11-17T01:38:59Z
transitions:
  - to: ready
    at: 2026-11-17T01:33:39Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T01:33:58Z
    by: agent-S-0816
  - to: review
    at: 2026-11-17T01:37:45Z
    by: agent-S-0816
  - to: done
    at: 2026-11-17T01:38:59Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0816-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 240
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 1
      work: 30
  models:
    - model: claude-opus-5-5
      input: 66
      output: 8386
      cache_read: 4386146
      cache_write: 164746
      cost: 2.3632
  strategic:
    - kind: orchestrator
      seconds: 104
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 137
          cache_read: 856468
          cache_write: 114735
          cost: 0.2982
formula: publish
---
# S-0816 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3488 Fetch the remote's branch and tags
- T-3489 Evaluate the release policy, stopping unless it is met
- T-3490 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Step 3 published S-0804 and pushed the branch and tags:

| Component | Version | Tag |
|-----------|---------|-----|
| flai | 1.101.2 → 1.101.3 | `flai/v1.101.3` |
| flaiover | 0.61.2 → 0.61.3 | `flaiover/v0.61.3` |

### Accepted by the orchestrator

- Verified: c2a68420ae5cdf6cdedcdde9486bfb0de503d9a2
- At: 2026-11-17T01:38:59Z

Verdict: passed. Close-out passed at c2a68420 (rebase, sync, narrative, check); formula instance, every step's task done.
- 1: T-3488 done
- 2: T-3489 done
- 3: T-3490 done
