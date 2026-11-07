---
id: S-0597
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-07T22:31:23Z
updated: 2026-11-07T23:09:44Z
transitions:
  - to: ready
    at: 2026-11-07T22:31:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T23:04:49Z
    by: agent-S-0597
  - to: review
    at: 2026-11-07T23:08:22Z
    by: system-flow
  - to: done
    at: 2026-11-07T23:09:44Z
    by: orchestrator
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 230
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 3
      work: 28
  models:
    - model: claude-opus-5-5
      input: 66
      output: 8115
      cache_read: 4547312
      cache_write: 154821
      cost: 2.3106
  strategic:
    - kind: orchestrator
      seconds: 42
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 56
          cache_read: 295809
          cache_write: 18532
          cost: 0.0797
formula: publish
---
# S-0597 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2612 Fetch the remote's branch and tags
- T-2613 Evaluate the release policy, stopping unless it is met
- T-2614 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.

### Accepted by the orchestrator

- Verified: 3ec051f28cf83a3fcc9e8a5e7e8388256d999966
- At: 2026-11-07T23:09:44Z

Verdict: pass
- 1: T-2612 done
- 2: T-2613 done
- 3: T-2614 done
