---
id: S-0731
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-12T03:44:58Z
updated: 2026-11-12T03:49:52Z
transitions:
  - to: ready
    at: 2026-11-12T03:44:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-12T03:45:07Z
    by: agent-S-0731
  - to: review
    at: 2026-11-12T03:48:11Z
    by: agent-S-0731
  - to: done
    at: 2026-11-12T03:49:52Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0731-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 197
  priced_by: harness
  turns:
    - day: 2026-11-12
      ceremony: 2
      hand_edits: 5
      work: 20
  models:
    - model: claude-opus-5-5
      input: 58
      output: 7962
      cache_read: 3858368
      cache_write: 156300
      cost: 2.1815
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 22
          cache_read: 75717
          cache_write: 15259
          cost: 0.0258
formula: publish
---
# S-0731 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3135 Fetch the remote's branch and tags
- T-3136 Evaluate the release policy, stopping unless it is met
- T-3137 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on 2026-11-12: `flai/v1.86.22` (patch, from 1.86.21), for S-0717, S-0701, and S-0681; the branch and the tag were pushed. No other component had anything pending.

### Accepted by the orchestrator

- Verified: 894c56c011542677cbfb69acfcd3d5c333bb94dc
- At: 2026-11-12T03:49:52Z

Verdict: pass, the publish formula's three steps are done, and the close-out passed at 894c56c
- 1: T-3135 done
- 2: T-3136 done
- 3: T-3137 done
