---
id: T-3043
type: task
nature: improvement
title: An ADR records what I-0288's instance shows and decides that flai serve asks the operator on one thread when nothing can be made ready
status: done
parent: S-0691
owner: alex
created: 2026-11-10T10:34:39Z
updated: 2026-11-11T05:10:50Z
transitions:
  - to: ready
    at: 2026-11-11T05:09:36Z
    by: agent-S-0691
  - to: in-progress
    at: 2026-11-11T05:09:37Z
    by: agent-S-0691
  - to: done
    at: 2026-11-11T05:10:50Z
    by: agent-S-0691
stream: S-0691
tags: [flai, adr]
touches: [design/adrs/0285-when-nothing-can-be-made-ready-flai-serve-opens-one-thread-to-the-operator-on.md, design/adrs/README.md]
usage:
  source: log
  seconds: 73
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 10759
      cache_read: 1600174
      cache_write: 68099
      cost: 1.0006
---
# T-3043 An ADR records what I-0288's instance shows and decides that flai serve asks the operator on one thread when nothing can be made ready

## Work

Write the ADR with `flai adr new`, then `flai adr topics` and `flai adr accept`. Confirm or correct the findings under `### Planning` in S-0691's Notes, and record them in its Context: the empty board of 2026-11-03 to 2026-11-05, the drafts behind TH-0532 and TH-0533, the `ready` reason pack that names no thread, and the reminder that had no channel.

The decision covers these points:

- When flai serve opens the thread: the ready column holds no story that can start under ADR-0279's count, the in-progress limit has room, and `flai promote --candidates` lists no candidate.
- Where it goes: one thread on `wip/kanban/board.md`, by flai serve's author.
- What it lists: each backlog story that cannot be promoted and the decision it waits on. That is a draft to finalize with its open threads that await the operator, a story held by `after` with the root of its chain, or a blocked story with its reason. It also says when the host has no reminder channel.
- When it adds an entry: only when that list changes.
- When it resolves: once a story can start again.
- Why the orchestrator leaves it, and how the reminder of ADR-0205 covers it.

Under Alternatives, name the ones you weigh: the orchestrator asking on its `ready-empty` reason, a reminder with no thread, and `flai promote --candidates` naming threads.

This task waits on none: it is the first layer. The new ADR's file is under the story's `design/adrs` folder touch, since `flai adr new` names it. Add that file to this task's touches with `flai touches` once it exists.

## Done when

- The ADR is accepted, and `design/adrs/README.md` lists it.
- `flai check --strict` passes.

## Notes
