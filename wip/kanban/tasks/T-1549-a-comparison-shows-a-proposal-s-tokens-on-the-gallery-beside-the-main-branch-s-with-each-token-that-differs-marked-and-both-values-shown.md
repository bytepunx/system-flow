---
id: T-1549
type: task
nature: feature
title: A comparison shows a proposal's tokens on the gallery beside the main branch's, with each token that differs marked and both values shown
status: backlog
parent: S-0384
owner: alex
created: 2026-10-13T17:48:03Z
updated: 2026-10-13T17:48:03Z
transitions: []
stream: S-0384
tags: [dashboard]
touches: [flaiover/src/lib/design/diff.ts, flaiover/src/lib/design/diff.test.ts, flaiover/src/lib/components/design/ProposalCompare.svelte, flaiover/src/lib/components/design/ProposalCompare.svelte.test.ts]
---
# T-1549 A comparison shows a proposal's tokens on the gallery beside the main branch's, with each token that differs marked and both values shown

## Work

`flaiover/src/lib/design/diff.ts` lists the tokens two sets differ in, added, removed, or changed, with both values, with tests. `ProposalCompare.svelte` shows two `PreviewFrame`s, the main branch's and the proposal's, side by side, and the list of differences beside them, each with a swatch or sample of both values. Component tests cover a proposal with changes and one with none. It waits for nothing.

## Done when

- `diff.test.ts` and `ProposalCompare.svelte.test.ts` pass.

## Notes

- Layer 1.
