---
id: S-0367
type: story
nature: improvement
title: Perform lint fix on every markdown file a story touches as part of validation
status: in-progress
owner: alex
created: 2026-10-12T16:48:56Z
updated: 2026-10-12T16:49:57Z
transitions:
  - to: ready
    at: 2026-10-12T16:48:58Z
    by: alex
  - to: in-progress
    at: 2026-10-12T16:49:57Z
    by: agent-S-0367
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 443
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 110
          cache_read: 2015371
          cache_write: 4401
          cost: 0.5092
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: alex
    at: 2026-10-12T16:48:56Z
---
# S-0367 Perform lint fix on every markdown file a story touches as part of validation

## Goal

Run `scripts/lint-md.sh --fix [markdown file]` for every markdown file that gets touched as part of a story to prevent bad markdown from slipping into main and breaking build checks.

## Acceptance criteria
- [ ] All markdown files get a lint fix call made on them if they were modified as part of a story
- [ ] Any markdown files outside of a story that get added also get the lint fix call
- [ ] Lint fix should get called right after a rebase to ensure that old files with bad markdown don't continue polluting recently changed markdown files

## Tasks

## Notes
