---
id: T-3002
type: task
nature: improvement
title: Decide in an ADR that a rebase which wrote a generated file again at a stop writes it again at the branch's head and commits it
status: done
parent: S-0702
owner: alex
created: 2026-11-10T10:31:10Z
updated: 2026-11-11T08:49:35Z
transitions:
  - to: ready
    at: 2026-11-11T08:49:08Z
    by: agent-S-0702
  - to: in-progress
    at: 2026-11-11T08:49:09Z
    by: agent-S-0702
  - to: done
    at: 2026-11-11T08:49:35Z
    by: agent-S-0702
stream: S-0702
tags: [adr, git]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 26
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 32
      cache_read: 374517
      cache_write: 2565
      cost: 0.1395
---
# T-3002 Decide in an ADR that a rebase which wrote a generated file again at a stop writes it again at the branch's head and commits it

## Work

Write an ADR with `flai adr new` that refines [ADR-0098](../../../design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md) and [ADR-0228](../../../design/adrs/0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md). Its context is I-0266's cause, as the story's `### Planning` note gives it. A rebase stopped on `design/issues/summary.md` at one commit, and flai wrote the summary from the issue files at that stop. A later commit of the branch then added an issue file without touching the summary, so the branch's head carried a summary its issue files do not render.

The decision to weigh, the planner's proposal:

- After a rebase in which `ContinueOverResolved` wrote a generated file again, flai writes each generated file again at the branch's head: the issues summary and the ADR index.
- When what one renders differs from the head's, it commits the files on the branch. The summary's `updated` line is left out of that comparison, as the `issues.summary` check leaves it out.
- This holds in `flai stream sync`, and so in `flai task done`, the close-out, and `flai accept`, and in the merge queue's batch rebase.
- Consider the alternative of writing them again at the head after every clean rebase, not only after one that wrote one at a stop. Say why it was or was not chosen.

Accept it with `flai adr accept` once it stands. The predicted path in this task's touches is a guess at the number and slug. `flai task done` widens the touches with the file written.

This task waits for none: it is the first layer.

## Done when

- An accepted ADR records the decision, refines ADR-0098, and names I-0266.
- `design/adrs/README.md` lists it, as `flai adr` writes it.

## Notes
