---
id: S-0588
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-07T09:24:05Z
updated: 2026-11-07T12:09:09Z
transitions:
  - to: ready
    at: 2026-11-07T09:24:27Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T12:05:12Z
    by: agent-S-0588
  - to: review
    at: 2026-11-07T12:07:48Z
    by: agent-S-0588
  - to: done
    at: 2026-11-07T12:09:09Z
    by: orchestrator
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 177
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      work: 23
  models:
    - model: claude-opus-5-5
      input: 54
      output: 6329
      cache_read: 3501450
      cache_write: 148393
      cost: 2.0142
  strategic:
    - kind: orchestrator
      seconds: 244
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 58
          output: 338
          cache_read: 2039868
          cache_write: 232191
          cost: 0.571
        - model: claude-sonnet-5-5
          input: 9
          output: 399
          cache_read: 74568
          cache_write: 31080
          cost: 0.0768
formula: publish
---
# S-0588 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2566 Fetch the remote's branch and tags
- T-2567 Evaluate the release policy, stopping unless it is met
- T-2568 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.

### Accepted by the orchestrator

- Verified: 0794a48d276966a39ad60ef490b8d171e434d1a2
- At: 2026-11-07T12:09:09Z

Verdict: pass
- 1: T-2566 done
- 2: T-2567 done
- 3: T-2568 done
