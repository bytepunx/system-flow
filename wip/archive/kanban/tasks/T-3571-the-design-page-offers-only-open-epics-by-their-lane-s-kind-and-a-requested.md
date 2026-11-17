---
id: T-3571
type: task
nature: remediation
title: The design page offers only open epics by their lane's kind, and a requested change moves to the workflow's pullable lane
status: done
parent: S-0838
owner: alex
created: 2026-11-17T12:29:38Z
updated: 2026-11-17T12:49:37Z
transitions:
  - to: ready
    at: 2026-11-17T12:42:15Z
    by: agent-S-0838
  - to: in-progress
    at: 2026-11-17T12:42:15Z
    by: agent-S-0838
  - to: done
    at: 2026-11-17T12:49:37Z
    by: agent-S-0838
stream: S-0838
tags: [flaiover, workflow]
touches: [flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, flaiover/src/lib/components/design/ProposeChange.svelte, flaiover/src/lib/components/design/ProposeChange.svelte.test.ts, design/components/flaiover/src/lib/components/design/ProposeChange.examples.json]
after: [T-3568]
usage:
  source: log
  seconds: 442
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 52
      output: 12657
      cache_read: 2743482
      cache_write: 103004
      cost: 1.4916
---
# T-3571 The design page offers only open epics by their lane's kind, and a requested change moves to the workflow's pullable lane

## Work

In `flaiover/src/routes/design/+page.svelte`, read the flow from `/api/board`, as `NewItemForm.svelte` does, and filter the epics with `isOpen(flow, e.status)` from `$lib/workflow` in place of `e.status !== 'done' && e.status !== 'cancelled'`. Derive the filtered list, so that epics read before the flow are filtered again when it arrives. Pass the flow to `ProposeChange` as a required prop.

In `flaiover/src/lib/components/design/ProposeChange.svelte`, replace the one move to `ready` with the moves that take the new story from the workflow's initial lane to its first start lane that is pullable: `startLanes(flow, 'story')` and `movesTo` from `$lib/lanes`, and `isPullable` from `$lib/workflow`, one move each, as `NewItemForm.svelte` makes them. Name that lane by its title in the message a refused move gives. With no pullable lane, leave the story in the initial lane and say so. Give `design/components/flaiover/src/lib/components/design/ProposeChange.examples.json` the new prop, from the baseline flow.

Tests: in `design.svelte.test.ts`, with the renamed flow T-3568 adds to `lanes.fixture.ts`, an epic in its accepted or cancelled lane is not offered, and one in an open lane is. In `ProposeChange.svelte.test.ts`, the move goes to `ready` under the baseline flow and to the renamed pullable lane under the renamed flow. It waits for T-3568, whose fixture its tests read.

## Done when

- Neither file compares a state name: `grep -nE "'(done|cancelled|ready)'"` on the two `.svelte` files finds nothing
- Both test files pass, and `flai test` on the task's paths passes, svelte-check and `components.examples` included

## Notes
