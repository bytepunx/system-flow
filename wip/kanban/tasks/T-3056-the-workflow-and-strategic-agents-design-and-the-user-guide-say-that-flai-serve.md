---
id: T-3056
type: task
nature: improvement
title: The workflow and strategic-agents design and the user guide say that flai serve asks the operator on one thread when nothing can be made ready
status: backlog
parent: S-0691
owner: alex
created: 2026-11-10T10:35:17Z
updated: 2026-11-10T10:35:17Z
transitions: []
stream: S-0691
tags: [docs, design]
touches: [design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md]
after: [T-3050]
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
