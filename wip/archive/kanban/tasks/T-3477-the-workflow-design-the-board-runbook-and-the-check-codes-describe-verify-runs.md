---
id: T-3477
type: task
nature: feature
title: The workflow design, the board runbook, and the check codes describe verify, runs, verify_at, and the suite findings
status: done
parent: S-0806
owner: alex
created: 2026-11-16T23:50:55Z
updated: 2026-11-17T09:13:06Z
transitions:
  - to: ready
    at: 2026-11-17T09:07:24Z
    by: agent-S-0806
  - to: in-progress
    at: 2026-11-17T09:07:24Z
    by: agent-S-0806
  - to: done
    at: 2026-11-17T09:13:06Z
    by: agent-S-0806
stream: S-0806
tags: [flai, docs]
touches: [design/system/workflow.md, docs/operators/runbooks/change-the-board.md, docs/users/flai.md, design/system/flai-cli.md, design/system/work-hierarchy.md]
after: [T-3471, T-3473, T-3474]
usage:
  source: log
  seconds: 342
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 48
      output: 14510
      cache_read: 3847190
      cache_write: 161514
      cost: 2.1069
---
# T-3477 The workflow design, the board runbook, and the check codes describe verify, runs, verify_at, and the suite findings

## Work

Describe what the story built, linking ADR-0312 rather than copying it:

- `design/system/workflow.md`, under `## The workflow file`: the three keys, their forms, the baseline's values, and that nothing reads them to decide a run until the stories after this one.
- `docs/operators/runbooks/change-the-board.md`: how an operator names suites on a move, a state, and `verify_at`, what `all` means, the two suite findings in `## Before the story is accepted`, and that an older flai stops on `flai.minimum` once a publish raises it.
- `docs/users/flai.md` and `design/system/flai-cli.md`: the two new codes in their tables of `flai check`'s workflow codes.
- `design/system/work-hierarchy.md`, under `### Fields an older flai does not know`: the fields file lists the workflow file's keys too.

It waits for T-3471, T-3473, and T-3474, whose values, codes, and minimum it describes.

## Done when

- Each document says what this story built, and the markdown lint passes on them.
- `flai test` on the changed paths passes.

## Notes
