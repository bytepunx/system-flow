---
id: S-0343
type: story
nature: improvement
title: Add an archive action to context menus
status: ready
owner: alex
created: 2026-10-08T08:07:29Z
updated: 2026-10-08T08:07:30Z
transitions:
  - to: ready
    at: 2026-10-08T08:07:30Z
    by: alex
tags: [dashboard]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: medium
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 2
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 44
          cache_read: 363301
          cache_write: 3577
          cost: 0.0904
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: alex
    at: 2026-10-08T08:07:29Z
---
# S-0343 Add an archive action to context menus

## Goal

Cancelled cards are remaining in the cancelled column and causing flai check to fail. Make it simple to archive specific cancelled cards or archive all cards in the cancelled lane with an `Archive` and `Archive All` action in the respective context menus.

## Acceptance criteria
- [ ] When cards are in the cancelled lane, right clicking in the lane offers an `Archive All` action that will call archive on all the cancelled cards
- [ ] When a card in cancelled is right clicked, an `Archive` action is offered that allows archiving that specific card.

## Tasks

## Notes
