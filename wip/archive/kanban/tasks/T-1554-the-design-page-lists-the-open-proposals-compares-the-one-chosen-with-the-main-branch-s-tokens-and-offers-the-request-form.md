---
id: T-1554
type: task
nature: feature
title: The Design page lists the open proposals, compares the one chosen with the main branch's tokens, and offers the request form
status: done
parent: S-0384
owner: alex
created: 2026-10-13T17:48:27Z
updated: 2026-10-15T03:56:37Z
transitions:
  - to: ready
    at: 2026-10-15T03:17:16Z
    by: system-flow
  - to: in-progress
    at: 2026-10-15T03:46:56Z
    by: system-flow
  - to: done
    at: 2026-10-15T03:56:37Z
    by: system-flow
stream: S-0384
tags: [dashboard]
touches: [flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, flaiover/src/lib/components/design/ProposalCompare.svelte, flaiover/src/lib/components/design/ProposalCompare.svelte.test.ts]
after: [T-1549, T-1550, T-1553]
usage:
  source: log
  seconds: 581
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 118
      output: 42612
      cache_read: 6114969
      cache_write: 185572
      cost: 3.1447
---
# T-1554 The Design page lists the open proposals, compares the one chosen with the main branch's tokens, and offers the request form

## Work

On `flaiover/src/routes/design/+page.svelte`, list the stories `/api/design/proposals` answers, each linked to its item page; choosing one loads `/api/design?story=` and shows `ProposalCompare` with the main branch's tokens; and a Propose a change button opens `ProposeChange`. The list follows story and document changes, so an accepted proposal leaves it and its tokens become the page's own. Page tests cover no proposals, a chosen proposal, and an accepted one. It waits for T-1549, T-1550, and T-1553, whose components and endpoints it places.

## Done when

- `design.svelte.test.ts` passes with the new cases.

## Notes

- Layer 4.
