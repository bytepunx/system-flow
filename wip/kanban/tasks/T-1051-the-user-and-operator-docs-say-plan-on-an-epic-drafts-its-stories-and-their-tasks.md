---
id: T-1051
type: task
nature: improvement
title: The user and operator docs say Plan on an epic drafts its stories and their tasks
status: done
parent: S-0300
owner: alex
created: 2026-10-06T21:47:28Z
updated: 2026-10-06T23:08:40Z
transitions:
  - to: ready
    at: 2026-10-06T23:07:14Z
    by: agent-S-0300
  - to: in-progress
    at: 2026-10-06T23:07:15Z
    by: agent-S-0300
  - to: done
    at: 2026-10-06T23:08:40Z
    by: agent-S-0300
stream: S-0300
tags: [planner, docs]
touches: [docs/users/flaiover.md, docs/users/flai.md, docs/operators/index.md]
after: [T-1047]
usage:
  source: log
  seconds: 85
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 19
      output: 82
      cache_read: 456811
      cache_write: 49491
      cost: 0.2231
---
# T-1051 The user and operator docs say Plan on an epic drafts its stories and their tasks

## Work

User-facing behaviour changes, so the docs change with it (documentation.md). Each change follows the design and convention wording from T-1047, so this task waits for T-1047 and is in layer 2.

- `docs/users/flaiover.md`, The item page: the sentence on **Plan** says "for an epic it drafts the stories, or revisits those it has". Make it say that it drafts the stories and their tasks. The stories appear under the epic as drafts, each with its tasks, for you to read and finalize. Check the card menu table row for **Plan** too.
- `docs/users/flai.md`, Running the planner: say the same for `flai plan E-nnnn`.
- `docs/operators/index.md`: in the planner section, where it says what a run on an epic writes, add the tasks. One run on an epic now writes more, and costs more.

## Done when

- [ ] The three docs say that planning an epic drafts its stories and their tasks, in the same words as the design.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner on the recommendation in TH-0201.
