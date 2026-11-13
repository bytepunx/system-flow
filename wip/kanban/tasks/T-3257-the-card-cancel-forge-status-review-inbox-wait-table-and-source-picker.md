---
id: T-3257
type: task
nature: feature
title: The card, cancel, forge status, review, inbox, wait table, and source picker components read a state's kind, tested on a state of the workflow's own
status: backlog
parent: S-0758
owner: alex
created: 2026-11-12T16:10:05Z
updated: 2026-11-12T16:10:05Z
transitions: []
stream: S-0758
tags: [dashboard]
touches: [flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/CancelConfirm.svelte, flaiover/src/lib/components/ForgeStatus.svelte, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/InboxView.svelte, flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/explorer/SourcePicker.svelte]
after: [T-3252]
---
# T-3257 The card, cancel, forge status, review, inbox, wait table, and source picker components read a state's kind, tested on a state of the workflow's own

## Work

Waits for T-3252: these components read the module that task adds. Replace each comparison with a state's name in these components with the module's kind, colour, or moves. Add component tests that render a card and the review page for a story in a state added between `in-progress` and `review`, and keep the components' examples under `design/components/` in step with any prop that changes.

## Done when

None of these components compares a status with a state's name, and the new and existing tests pass.

## Notes
