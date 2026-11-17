---
id: T-3255
type: task
nature: feature
title: The item page and the Deployments page read a state's kind and offer only the moves the operator may make
status: done
parent: S-0758
owner: alex
created: 2026-11-12T16:09:57Z
updated: 2026-11-17T08:12:13Z
transitions:
  - to: ready
    at: 2026-11-17T08:01:27Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:01:28Z
    by: system-flow
  - to: done
    at: 2026-11-17T08:12:13Z
    by: agent-S-0758
stream: S-0758
tags: [dashboard]
touches: ["flaiover/src/routes/items/[id]/+page.svelte", flaiover/src/routes/ops/deployments/+page.svelte, "flaiover/src/routes/items/[id]/item.svelte.test.ts", flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, design/components/flaiover/src/lib/components/BoardCard.examples.json, design/components/flaiover/src/lib/components/CancelConfirm.examples.json, flaiover/src/lib/components/explorer/SourcePicker.svelte.test.ts]
after: [T-3252]
usage:
  source: log
  seconds: 644
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 92
      output: 457
      cache_read: 4678806
      cache_write: 137440
      cost: 1.8342
---
# T-3255 The item page and the Deployments page read a state's kind and offer only the moves the operator may make

## Work

Waits for T-3252: both pages read the module that task adds. Replace the state names in the item page with the module's kind and moves: offer only the moves its transitions allow the operator, and ask for a reason where one is needed. Make the Deployments page tell a story's namespace by the kind of its state.

## Done when

Neither page compares a status with a state's name, and the flaiover tests pass.

## Notes
