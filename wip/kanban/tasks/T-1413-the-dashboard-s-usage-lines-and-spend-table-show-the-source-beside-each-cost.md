---
id: T-1413
type: task
nature: feature
title: The dashboard's usage lines and spend table show the source beside each cost
status: backlog
parent: S-0358
owner: alex
created: 2026-10-08T08:51:49Z
updated: 2026-10-08T08:51:49Z
transitions: []
stream: S-0358
tags: [dashboard]
touches: [flaiover/src/lib/usage.ts, flaiover/src/lib/usage.test.ts, flaiover/src/lib/components/SpendTable.svelte, flaiover/src/lib/components/SpendTable.svelte.test.ts]
after: [T-1412]
---
# T-1413 The dashboard's usage lines and spend table show the source beside each cost

## Work

- `usage.ts`: the types gain `priced_by` and `cost_by_source`; the usage line prints `(gateway)`, `(harness)`, or `(estimated)` after the cost, falling back to today's `(estimated)` rule for an item without `priced_by`.
- `SpendTable.svelte`: the mark after a cost names its source rather than `*` alone, with a title saying what each means.
- Tests in `usage.test.ts` and `SpendTable.svelte.test.ts` for each source and an older item.

It waits for T-1412, whose JSON it reads.

## Done when

- `flai test flaiover/src/lib/usage.ts flaiover/src/lib/components/SpendTable.svelte` passes.

## Notes

Layer 3 of S-0358.
