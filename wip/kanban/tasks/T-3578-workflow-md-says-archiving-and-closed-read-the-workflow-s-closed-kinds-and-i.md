---
id: T-3578
type: task
nature: remediation
title: workflow.md says archiving and Closed read the workflow's closed kinds, and I-0377 is closed
status: done
parent: S-0842
owner: alex
created: 2026-11-17T12:29:51Z
updated: 2026-11-18T08:08:17Z
transitions:
  - to: ready
    at: 2026-11-18T08:08:01Z
    by: agent-S-0842
  - to: in-progress
    at: 2026-11-18T08:08:01Z
    by: agent-S-0842
  - to: done
    at: 2026-11-18T08:08:17Z
    by: agent-S-0842
stream: S-0842
tags: [docs, workflow]
touches: [design/system/workflow.md]
after: [T-3572, T-3576]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 28
      cache_read: 527817
      cache_write: 1682
      cost: 0.2026
---
# T-3578 workflow.md says archiving and Closed read the workflow's closed kinds, and I-0377 is closed

## Work

In `design/system/workflow.md`, in the paragraph of § Transitions and who makes them that says the command line reads its moves from the project's workflow (S-0744), add that `flai accept` and `flai archive` archive an item in a state of the accepted or cancelled kind, whatever its name, and that every reader of an item's closed state reads those kinds (S-0842, I-0377). Then close the issue with `flai issue close I-0377 --reason "<what fixed it>"`, naming T-3572's change and T-3576's tests. Waits for T-3572, whose behaviour it records, and T-3576, whose tests show the issue no longer occurs before it is closed.

## Done when

- `design/system/workflow.md` says what reads the accepted and cancelled kinds, and its `updated` is bumped.
- I-0377 is closed with a reason naming the fix and its tests.
- `flai check --strict` and the markdown lint pass on the changed paths.

## Notes
