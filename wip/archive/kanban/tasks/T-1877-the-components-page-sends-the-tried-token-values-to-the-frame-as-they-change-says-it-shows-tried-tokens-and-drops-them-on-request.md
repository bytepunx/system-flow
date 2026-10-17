---
id: T-1877
type: task
nature: feature
title: The Components page sends the tried token values to the frame as they change, says it shows tried tokens, and drops them on request
status: done
parent: S-0441
owner: alex
created: 2026-10-16T21:36:37Z
updated: 2026-10-17T03:07:16Z
transitions:
  - to: ready
    at: 2026-10-17T02:59:46Z
    by: agent-S-0441
  - to: in-progress
    at: 2026-10-17T02:59:47Z
    by: agent-S-0441
  - to: done
    at: 2026-10-17T03:07:15Z
    by: agent-S-0441
stream: S-0441
tags: [dashboard]
touches: [flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts]
after: [T-1874, T-1875]
usage:
  source: log
  seconds: 448
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 28418
      cache_read: 3321412
      cache_write: 149494
      cost: 2.0752
---
# T-1877 The Components page sends the tried token values to the frame as they change, says it shows tried tokens, and drops them on request

## Work

- `+page.svelte`: read the draft with T-1874's `loadDraft` and follow it with `onDraft`. On each change, send T-1875's declarations as `tokens` through the frame. Show "showing tried tokens" with a button that clears the draft and sends an empty `tokens`.
- `ComponentFrame.svelte`: resend the last `tokens` after each `ready`, so a reloaded frame keeps them.

Waits for T-1874 for the draft and T-1875 for the declarations.

## Done when

Component tests cover a draft sent on load, a change from another tab, the notice, the drop, and the resend after `ready`. `flai test flaiover/src/routes/design/components flaiover/src/lib/components/explorer` passes.

## Notes

None.
