---
id: S-0820
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-17T02:48:43Z
updated: 2026-11-17T03:29:08Z
transitions:
  - to: ready
    at: 2026-11-17T02:48:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T03:23:30Z
    by: agent-S-0820
  - to: review
    at: 2026-11-17T03:27:43Z
    by: agent-S-0820
  - to: done
    at: 2026-11-17T03:29:08Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0820-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 260
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 1
      work: 27
  models:
    - model: claude-opus-5-5
      input: 60
      output: 7699
      cache_read: 4183720
      cache_write: 163519
      cost: 2.2991
  strategic:
    - kind: orchestrator
      seconds: 107
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 191
          cache_read: 939086
          cache_write: 115807
          cost: 0.3251
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 7329
          cache_write: 5732
          cost: 0.0105
formula: publish
---
# S-0820 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3500 Fetch the remote's branch and tags
- T-3501 Evaluate the release policy, stopping unless it is met
- T-3502 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: f703b97112becbf3f896577f5d2c95c603d4faf3
- At: 2026-11-17T03:29:08Z

Verdict: S-0820 meets all three criteria; the close-out passed at f703b97112becbf3f896577f5d2c95c603d4faf3.
- 1: T-3500 done
- 2: T-3501 done
- 3: T-3502 done
