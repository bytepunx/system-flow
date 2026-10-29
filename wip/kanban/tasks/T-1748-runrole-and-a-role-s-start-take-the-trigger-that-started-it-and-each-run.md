---
id: T-1748
type: task
nature: feature
title: RunRole and a role's start take the trigger that started it, and each run records it in agents.json and in its activity entry's Trigger line
status: done
parent: S-0418
owner: alex
created: 2026-10-16T05:45:51Z
updated: 2026-10-29T16:52:35Z
transitions:
  - to: ready
    at: 2026-10-29T16:49:49Z
    by: agent-S-0418
  - to: in-progress
    at: 2026-10-29T16:49:50Z
    by: agent-S-0418
  - to: done
    at: 2026-10-29T16:52:35Z
    by: agent-S-0418
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/cmd/role.go]
usage:
  source: log
  seconds: 165
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 14
      output: 83
      cache_read: 1398044
      cache_write: 13381
      cost: 0.5827
---
# T-1748 RunRole and a role's start take the trigger that started it, and each run records it in agents.json and in its activity entry's Trigger line

## Work

- `RunRole` and the launcher's `role` take the trigger that started the run, in place of the fixed `asked`. `flai role run`, MCP `role_run`, and the host API's `role.run` pass `asked`.
- Add `startRole`, the start every trigger takes: `roleCheck`, the budget gate, and `role` with the trigger, answering the run, the budget state that held it, or the refusal, without the hand start's refusal of a held budget (D-0421), so that a queued or scheduled start that the budget holds waits rather than drops.
- The run's record in `serve/agents.json` and its activity entry's `- Trigger:` line carry the trigger as given: `asked`, `<event> <ID>`, or `schedule <spec>`, joined by semicolons when coalesced. The journal entry names it in place of `asked`.
- Test a start with a trigger other than asked, its record, and its activity line.

First layer: T-1743 and T-1744 start roles through `startRole`.

## Done when

`flai test flai/internal/serve flai/cmd` passes.

## Notes
