---
id: T-1554
type: task
nature: feature
title: The Design page lists the open proposals, compares the one chosen with the main branch's tokens, and offers the request form
status: backlog
parent: S-0384
owner: alex
created: 2026-10-13T17:48:27Z
updated: 2026-10-13T17:48:27Z
transitions: []
stream: S-0384
tags: [dashboard]
touches: [flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts]
after: [T-1549, T-1550, T-1553]
---
# T-1554 The Design page lists the open proposals, compares the one chosen with the main branch's tokens, and offers the request form

## Work

On `flaiover/src/routes/design/+page.svelte`, list the stories `/api/design/proposals` answers, each linked to its item page; choosing one loads `/api/design?story=` and shows `ProposalCompare` with the main branch's tokens; and a Propose a change button opens `ProposeChange`. The list follows story and document changes, so an accepted proposal leaves it and its tokens become the page's own. Page tests cover no proposals, a chosen proposal, and an accepted one. It waits for T-1549, T-1550, and T-1553, whose components and endpoints it places.

## Done when

- `design.svelte.test.ts` passes with the new cases.

## Notes

- Layer 4.
