---
id: S-0718
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-11T01:44:41Z
updated: 2026-11-11T22:17:45Z
transitions:
  - to: ready
    at: 2026-11-11T01:44:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T22:13:46Z
    by: agent-S-0718
  - to: review
    at: 2026-11-11T22:16:27Z
    by: agent-S-0718
  - to: done
    at: 2026-11-11T22:17:45Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0718-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 180
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 2
      work: 25
  models:
    - model: claude-opus-5-5
      input: 60
      output: 8195
      cache_read: 4071380
      cache_write: 163760
      cost: 2.2885
  strategic:
    - kind: orchestrator
      seconds: 312
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 95
          output: 560
          cache_read: 3520714
          cache_write: 468475
          cost: 1.115
formula: publish
---
# S-0718 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3080 Fetch the remote's branch and tags
- T-3081 Evaluate the release policy, stopping unless it is met
- T-3082 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published 2026-11-11: flai 1.86.17 → 1.86.18 (patch), tag `flai/v1.86.18`, branch and tag pushed. It covers S-0689, S-0707, S-0714, and S-0715. No other component was released.

### Accepted by the orchestrator

- Verified: 4528074d00ce2d5a0b60698a6eaedb68c99426a3
- At: 2026-11-11T22:17:45Z

Verdict: passed; close-out passed at 4528074d00ce, every step of the publish formula's instance done.
- 1: T-3080 done
- 2: T-3081 done
- 3: T-3082 done
