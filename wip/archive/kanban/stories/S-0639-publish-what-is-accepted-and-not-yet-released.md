---
id: S-0639
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-09T11:24:54Z
updated: 2026-11-09T11:27:49Z
transitions:
  - to: ready
    at: 2026-11-09T11:24:55Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:25:00Z
    by: agent-S-0639
  - to: review
    at: 2026-11-09T11:26:42Z
    by: system-flow
  - to: done
    at: 2026-11-09T11:27:49Z
    by: orchestrator
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 109
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 2
      work: 22
  models:
    - model: claude-opus-5-5
      input: 56
      output: 7567
      cache_read: 3603073
      cache_write: 149100
      cost: 2.065
  strategic:
    - kind: orchestrator
      seconds: 34
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 37
          cache_read: 334791
          cache_write: 2518
          cost: 0.0887
formula: publish
---
# S-0639 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2729 Fetch the remote's branch and tags
- T-2730 Evaluate the release policy, stopping unless it is met
- T-2731 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.

### Accepted by the orchestrator

- Verified: 19e81fe3b45a0d616c5dc2bd1d358bd98539047d
- At: 2026-11-09T11:27:49Z

Verdict: passed; formula instance publish, every step's task done, flai verify passed at 19e81fe3b45a
- 1: T-2729 done
- 2: T-2730 done
- 3: T-2731 done
