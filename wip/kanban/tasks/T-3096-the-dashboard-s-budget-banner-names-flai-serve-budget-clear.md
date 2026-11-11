---
id: T-3096
type: task
nature: remediation
title: The dashboard's budget banner names flai serve budget clear
status: backlog
parent: S-0720
owner: alex
created: 2026-11-11T08:54:39Z
updated: 2026-11-11T08:54:39Z
transitions: []
stream: S-0720
tags: [flaiover, budget]
touches: [flaiover/src/lib/components/BudgetBanner.svelte, flaiover/src/lib/components/BudgetBanner.svelte.test.ts]
---
# T-3096 The dashboard's budget banner names flai serve budget clear

## Work

- In `flaiover/src/lib/components/BudgetBanner.svelte`, after the provider, the limit, and until, name the command that clears the hold by hand: `flai serve budget clear <provider>`, with the held entry's `provider` as the key.
- Extend `flaiover/src/lib/components/BudgetBanner.svelte.test.ts` to expect the command in the banner's text.
- Waits for nothing: the command's name and argument are fixed by the story's third criterion, and the banner reads `provider` from `agent.status` as it does today.

## Done when

- The banner shows `flai serve budget clear <provider>` for a held state, in its test.
- `flai test` on the changed paths passes.

## Notes
