---
id: T-3056
type: task
nature: improvement
title: The workflow and strategic-agents design and the user guide say that flai serve asks the operator on one thread when nothing can be made ready
status: done
parent: S-0691
owner: alex
created: 2026-11-10T10:35:17Z
updated: 2026-11-11T07:49:21Z
transitions:
  - to: ready
    at: 2026-11-11T07:48:27Z
    by: agent-S-0691
  - to: in-progress
    at: 2026-11-11T07:48:27Z
    by: agent-S-0691
  - to: done
    at: 2026-11-11T07:49:21Z
    by: agent-S-0691
stream: S-0691
tags: [docs, design]
touches: [design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md]
after: [T-3050]
usage:
  source: log
  seconds: 54
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 713
      cache_read: 206919
      cache_write: 6020
      cost: 0.1038
---
# T-3056 The workflow and strategic-agents design and the user guide say that flai serve asks the operator on one thread when nothing can be made ready

## Work

Describe what T-3050 built, linking T-3043's ADR:

- `design/system/workflow.md`: under Pull policy, say what happens when nothing can be made ready. Under Cadence, say that the board's thread is reminded like any other thread that awaits the operator.
- `design/system/strategic-agents.md`: in the orchestrator's `ready-empty` row and its `ready` reason, say that flai serve's thread on the board names the operator's decisions, and that the orchestrator leaves that thread alone.
- `docs/users/flai.md`: in the `flai serve` section, list the board's thread beside the other threads flai serve opens, with when it opens, updates, and resolves.

It waits on T-3050, so that it describes what was built.

## Done when

- All three documents say when the thread opens, what it lists, when it gains an entry, and when it resolves, and link the ADR.
- `flai test design/system/workflow.md design/system/strategic-agents.md docs/users/flai.md` passes.

## Notes
