---
id: S-0788
type: story
nature: improvement
title: Stories in review that cannot be merged need an action
status: in-progress
owner: alex
created: 2026-11-14T18:02:48Z
updated: 2026-11-14T18:03:27Z
transitions:
  - to: ready
    at: 2026-11-14T18:02:49Z
    by: alex
  - to: in-progress
    at: 2026-11-14T18:03:27Z
    by: agent-S-0788
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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
- [ ] Stories in review that cannot be merged should not just report this to the operator without a mechanism to have the agent resolve it
- [ ] A rebase and resolve button should be visible that will take the recommended actions shown to the operator

## Tasks

## Notes
