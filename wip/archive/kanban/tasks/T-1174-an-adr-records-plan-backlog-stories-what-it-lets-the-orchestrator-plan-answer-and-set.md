---
id: T-1174
type: task
nature: feature
title: "An ADR records plan_backlog_stories: what it lets the orchestrator plan, answer, and set"
status: done
parent: S-0328
owner: alex
created: 2026-10-07T19:40:08Z
updated: 2026-10-07T19:41:16Z
transitions:
  - to: ready
    at: 2026-10-07T19:40:47Z
    by: agent-S-0328
  - to: in-progress
    at: 2026-10-07T19:40:48Z
    by: agent-S-0328
  - to: done
    at: 2026-10-07T19:41:16Z
    by: agent-S-0328
stream: S-0328
tags: []
touches: [design/adrs]
usage:
  source: log
  seconds: 28
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 5
      output: 2375
      cache_read: 402903
      cache_write: 10128
      cost: 0.1914
---
# T-1174 An ADR records plan_backlog_stories: what it lets the orchestrator plan, answer, and set

## Work

Record with `flai adr new --commit`, refining ADR-0087, the decision S-0328's other tasks implement:

- `orchestration.permissions.plan_backlog_stories`, a boolean, off by default, listed after `plan_backlog_epics`.
- `flai plan --candidates` also lists stories: a story in the backlog that lacks a plan, meaning no touches, no forecast duration, no cost of delay value, or no open or done task. Left out: a planner runs for it or its epic now; its newest planner run asked on a thread still awaiting an answer; its newest planner run ended and the story has not changed since.
- With the permission on, the orchestrator starts the planner on a story the candidates list, with `plan` or `flai plan`; flai and the guard refuse any other story.
- On a thread a story's planner opened (opener `planner-S-nnnn`), it may reply as an answer, with no source, or as a recommendation, and resolve the thread. That covers approving a plan, answering a question, and choosing a cost of delay input. Confirming a recommendation stays the operator's.
- It may set a backlog story's cost of delay, inputs and value, and nothing else of the story, with `item_edit` `cost_of_delay` or the `flai edit` cost flags.

No task waits on another to start this one; it runs in layer 1.

## Done when

- [ ] The ADR is in `design/adrs` and its README, committed on `story/S-0328`.

## Notes
