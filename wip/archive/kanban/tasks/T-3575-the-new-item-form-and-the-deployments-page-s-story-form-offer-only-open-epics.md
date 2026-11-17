---
id: T-3575
type: task
nature: remediation
title: The new-item form and the Deployments page's story form offer only open epics by their lane's kind
status: done
parent: S-0838
owner: alex
created: 2026-11-17T12:29:43Z
updated: 2026-11-17T12:49:19Z
transitions:
  - to: ready
    at: 2026-11-17T12:42:15Z
    by: agent-S-0838
  - to: in-progress
    at: 2026-11-17T12:42:15Z
    by: agent-S-0838
  - to: done
    at: 2026-11-17T12:49:19Z
    by: agent-S-0838
stream: S-0838
tags: [flaiover, workflow]
touches: [flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/NewItemForm.svelte.test.ts, flaiover/src/lib/components/DeploySettings.svelte, flaiover/src/lib/components/DeploySettings.svelte.test.ts]
after: [T-3568]
usage:
  source: log
  seconds: 424
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 7463
      cache_read: 1617658
      cache_write: 60735
      cost: 0.8795
---
# T-3575 The new-item form and the Deployments page's story form offer only open epics by their lane's kind

## Work

In `flaiover/src/lib/components/NewItemForm.svelte`, filter the epics with `isOpen(flow, e.status)` from `$lib/workflow` over the flow the form already reads from `/api/board`, in place of `e.status !== 'done' && e.status !== 'cancelled'`. Derive the list from the epics read and the flow, so that the reset of a parent that has closed, and the choice of a sole epic, run once both are read, whichever answers first.

In `flaiover/src/lib/components/DeploySettings.svelte`, read the flow from `/api/board` beside the epics in `onMount`, and filter them with `isOpen` the same way, so that the story for the scripts goes under an open epic only.

Tests: in `NewItemForm.svelte.test.ts` and `DeploySettings.svelte.test.ts`, with the renamed flow T-3568 adds to `lanes.fixture.ts`, an epic in its accepted or cancelled lane is not offered and one in an open lane is; under the baseline flow they behave as before. It waits for T-3568, whose fixture its tests read, and runs beside the design page's task, with which it shares no file.

## Done when

- Neither file compares a state name: `grep -nE "'(done|cancelled)'"` on the two `.svelte` files finds nothing
- Both test files pass, and `flai test` on the task's paths passes, svelte-check included

## Notes
