---
id: S-0593
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-07T15:17:08Z
updated: 2026-11-07T15:21:07Z
transitions:
  - to: ready
    at: 2026-11-07T15:17:32Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T15:17:37Z
    by: agent-S-0593
  - to: review
    at: 2026-11-07T15:19:55Z
    by: system-flow
  - to: done
    at: 2026-11-07T15:21:07Z
    by: orchestrator
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 146
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      work: 24
  models:
    - model: claude-opus-5-5
      input: 56
      output: 6928
      cache_read: 3586495
      cache_write: 148754
      cost: 2.0461
  strategic:
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 27
          cache_read: 79654
          cache_write: 35995
          cost: 0.0292
formula: publish
---
# S-0593 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2586 Fetch the remote's branch and tags
- T-2587 Evaluate the release policy, stopping unless it is met
- T-2588 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.

### Accepted by the orchestrator

- Verified: 3ea7bb9259b448c2c3b2a9b4117d4eb475283abe
- At: 2026-11-07T15:21:07Z

Verdict: pass. Every step of the publish formula ran and exited 0; the last flai verify passed at 3ea7bb9.
- 1: T-2586 done
- 2: T-2587 done
- 3: T-2588 done
