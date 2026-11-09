---
id: S-0641
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: review
owner: alex
created: 2026-11-09T16:56:49Z
updated: 2026-11-09T16:59:05Z
transitions:
  - to: ready
    at: 2026-11-09T16:56:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T16:56:55Z
    by: agent-S-0641
  - to: review
    at: 2026-11-09T16:59:05Z
    by: system-flow
tags: []
touches: [design/issues/I-0320-a-publish-instance-s-evaluate-and-publish-steps-race-a-publish-made-meanwhile.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 140
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      work: 36
  models:
    - model: claude-opus-5-5
      input: 82
      output: 9763
      cache_read: 5478860
      cache_write: 156305
      cost: 2.5418
  strategic:
    - kind: orchestrator
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 34
          cache_read: 89751
          cache_write: 53176
          cost: 0.0377
formula: publish
---
# S-0641 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2732 Fetch the remote's branch and tags
- T-2733 Evaluate the release policy, stopping unless it is met
- T-2734 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
