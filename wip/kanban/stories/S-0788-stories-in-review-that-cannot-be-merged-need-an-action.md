---
id: S-0788
type: story
nature: improvement
title: Stories in review that cannot be merged need an action
status: review
owner: alex
created: 2026-11-14T18:02:48Z
updated: 2026-11-14T23:14:06Z
transitions:
  - to: ready
    at: 2026-11-14T18:02:49Z
    by: alex
  - to: in-progress
    at: 2026-11-14T18:03:27Z
    by: agent-S-0788
  - to: review
    at: 2026-11-14T23:14:06Z
    by: agent-S-0788
tags: [dashboard]
touches: [flaiover/src, design/adrs/0302-the-operator-sends-a-review-story-that-cannot-be-accepted-back-to-its-agent-to.md, design/adrs/README.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/accept.go, flai/cmd/accept_resolve.go, flai/cmd/accept_resolve_test.go, flai/internal/hostapi/guardreads_test.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_items_test.go, flai/internal/hostapi/writes_test.go, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3263
  priced_by: harness
  turns:
    - day: 2026-11-14
      ceremony: 1
      hand_edits: 2
      work: 37
  models:
    - model: claude-haiku-5-5
      input: 64
      output: 29178
      cache_read: 3134985
      cache_write: 156284
      cost: 0.2478
    - model: claude-opus-5-5
      input: 256
      output: 86247
      cache_read: 15418990
      cache_write: 795102
      cost: 9.5684
  strategic:
    - kind: orchestrator
      seconds: 34
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 70
          cache_read: 299793
          cache_write: 74572
          cost: 0.1117
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: alex
    at: 2026-11-14T18:02:48Z
---
# S-0788 Stories in review that cannot be merged need an action

## Goal

I keep getting stories in the review queue that cannot be merged - often due to missing commits that landed on main.

Provide a rebase and reverify for this scenario.

## Acceptance criteria
- [x] Stories in review that cannot be merged should not just report this to the operator without a mechanism to have the agent resolve it
- [x] A rebase and resolve button should be visible that will take the recommended actions shown to the operator

## Tasks
- T-3364 flai accept --resolve sends a story in review back to its agent to rebase and resolve what holds its acceptance, and the host method accept.resolve runs it
- T-3365 The review page's Rebase and resolve button asks the story's agent, through accept.resolve, to take the recommended actions that hold its acceptance

## Notes
