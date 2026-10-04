---
id: T-0826
type: task
nature: improvement
title: The design and the user guides describe when the planner runs again
status: done
parent: S-0211
owner: alex
created: 2026-10-04T20:42:18Z
updated: 2026-10-04T21:17:11Z
transitions:
  - to: ready
    at: 2026-10-04T20:42:49Z
    by: agent-S-0211
  - to: in-progress
    at: 2026-10-04T21:11:22Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T21:17:11Z
    by: agent-S-0211
stream: S-0211
tags: [flai]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-0824, T-0825]
usage:
  source: log
  seconds: 349
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 30531
      cache_read: 4903373
      cache_write: 123925
      cost: 2.3348
---
# T-0826 The design and the user guides describe when the planner runs again

## Work

`design/system/strategic-agents.md` says the planner also starts on its own, under a section on planning again: the gate, each trigger, the replan policy, the schedule, the coalescing queue, and the Trigger line, linking ADR-0084, and its table and "Starting it" no longer say nothing starts it on a move or a timer. `design/system/flai-cli.md` and `docs/users/flai.md` say the same where they describe `flai serve` and the planner; `docs/users/flaiover.md` mentions the Settings page's Planning section. It waits for T-0824 and T-0825, whose behavior it describes.

## Done when

- [ ] Each document describes the behavior as built, with no claim that nothing starts the planner on its own
- [ ] The markdown lint passes on them

## Notes
