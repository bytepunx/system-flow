---
id: T-1877
type: task
nature: feature
title: The Components page sends the tried token values to the frame as they change, says it shows tried tokens, and drops them on request
status: backlog
parent: S-0441
owner: alex
created: 2026-10-16T21:36:37Z
updated: 2026-10-16T21:36:37Z
transitions: []
stream: S-0441
tags: [dashboard]
touches: [flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/lib/components/explorer/ComponentFrame.svelte]
after: [T-1874, T-1875]
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
