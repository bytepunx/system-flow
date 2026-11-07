---
id: S-0585
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-07T08:03:26Z
updated: 2026-11-07T09:43:02Z
transitions:
  - to: ready
    at: 2026-11-07T08:03:48Z
    by: agent-S-0562
  - to: in-progress
    at: 2026-11-07T08:04:51Z
    by: agent-S-0585
  - to: review
    at: 2026-11-07T08:16:55Z
    by: system-flow
  - to: done
    at: 2026-11-07T09:43:02Z
    by: alex
blocked:
  - from: 2026-11-07T08:07:13Z
    until: 2026-11-07T08:15:35Z
    reason: awaiting the operator on TH-0854
tags: []
touches: [wip/kanban/stories/S-0585-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 229
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 3
      hand_edits: 1
      work: 30
  models:
    - model: claude-opus-5-5
      input: 76
      output: 11860
      cache_read: 5249456
      cache_write: 164411
      cost: 2.6027
  strategic:
    - kind: orchestrator
      seconds: 48
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 124
          cache_read: 531301
          cache_write: 83096
          cost: 0.1519
formula: publish
---
# S-0585 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2546 Fetch the remote's branch and tags
- T-2548 Evaluate the release policy, stopping unless it is met
- T-2549 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.

- 2026-11-07: step 2 stopped the run. `flai release --evaluate --require-met` exited 2 in the main checkout: the release policy is `judgement`, never met by itself, and nothing accepted is waiting for a release (count 0, value 0 USD/week). No release was made and step 3 was not run.
- 2026-11-07: on TH-0854 the operator confirmed agent-S-0562's recommendation. T-2548 is closed with its Notes as they stand, T-2549 is cancelled ("step 2 stopped the run: nothing accepted is waiting for a release"), and S-0585 goes to review for acceptance. Nothing was tagged or pushed.
- Criterion 3 is left unticked: `flai release --pending` never ran, and nothing was waiting to be published, so no publish could be observed.
