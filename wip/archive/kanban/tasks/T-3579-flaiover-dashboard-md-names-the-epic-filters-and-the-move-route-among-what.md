---
id: T-3579
type: task
nature: remediation
title: flaiover-dashboard.md names the epic filters and the move route among what reads a state's kind, and I-0399 is closed
status: done
parent: S-0838
owner: alex
created: 2026-11-17T12:29:52Z
updated: 2026-11-17T12:50:42Z
transitions:
  - to: ready
    at: 2026-11-17T12:50:01Z
    by: agent-S-0838
  - to: in-progress
    at: 2026-11-17T12:50:02Z
    by: agent-S-0838
  - to: done
    at: 2026-11-17T12:50:42Z
    by: agent-S-0838
stream: S-0838
tags: [flaiover, workflow, docs]
touches: [design/system/flaiover-dashboard.md, design/issues/I-0399-flaiover-still-compares-state-names-in-the-epic-filters-and-the-move-route-so.md, design/issues/summary.md]
after: [T-3571, T-3575]
usage:
  source: log
  seconds: 40
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 8436
      cache_read: 1828425
      cache_write: 68648
      cost: 0.9941
---
# T-3579 flaiover-dashboard.md names the epic filters and the move route among what reads a state's kind, and I-0399 is closed

## Work

In `design/system/flaiover-dashboard.md` § A state's kind (S-0758), add the design page and `ProposeChange`, `NewItemForm`, and `DeploySettings`, which read the flow from `/api/board` or take it as a prop, and the move route, which reads it from the repo's board view on the server, to the list of what reads a state's kind. Say that the move route judges the cancellation's preview and the acceptance's timeout by the target lane's kind (S-0838). Bump `updated`.

Search `flaiover/src` outside the tests and fixtures for a status or a move's target compared with a state's name, such as `grep -rnE "(status|to|state) *[!=]== *'(done|cancelled|ready|backlog|review|in-progress)'"`. Confirm that only what the section already keeps remains: the task plan's states in `$lib/taskplan.ts` and the charts' baseline states in `$lib/viz/charts.ts`, which S-0754 draws from the workflow. Name any other in the story's Notes.

Close the issue with `flai issue close I-0399 --reason "<what fixed it>"`, naming the tasks of S-0838 and the views each changed; it writes the issue's file and `design/issues/summary.md`. It waits for T-3571 and T-3575, whose changes the section and the reason describe.

## Done when

- The section lists every view and the route that S-0838 changed
- The search finds only the comparisons the section keeps
- `flai issue show I-0399` reads closed, with the reason
- `flai test` on the task's paths passes, the markdown lint included

## Notes
