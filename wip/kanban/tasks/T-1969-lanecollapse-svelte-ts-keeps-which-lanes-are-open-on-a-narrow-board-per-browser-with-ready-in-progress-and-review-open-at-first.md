---
id: T-1969
type: task
nature: improvement
title: lanecollapse.svelte.ts keeps which lanes are open on a narrow board per browser, with Ready, In progress, and Review open at first
status: done
parent: S-0462
owner: alex
created: 2026-10-21T07:04:42Z
updated: 2026-10-21T09:06:41Z
transitions:
  - to: ready
    at: 2026-10-21T09:05:37Z
    by: agent-S-0462
  - to: in-progress
    at: 2026-10-21T09:05:38Z
    by: agent-S-0462
  - to: done
    at: 2026-10-21T09:06:41Z
    by: agent-S-0462
stream: S-0462
tags: [dashboard]
touches: [flaiover/src/lib/lanecollapse.svelte.ts, flaiover/src/lib/lanecollapse.svelte.test.ts]
usage:
  source: log
  seconds: 63
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 3815
      cache_read: 524257
      cache_write: 19954
      cost: 0.3133
---
# T-1969 lanecollapse.svelte.ts keeps which lanes are open on a narrow board per browser, with Ready, In progress, and Review open at first

## Work

- Write `flaiover/src/lib/lanecollapse.svelte.ts` after `$lib/boardtypes.svelte.ts`. It holds the set of open lanes as rune state, kept in `localStorage` under `flaiover-board-lanes`, with `ready`, `in-progress`, and `review` open when nothing is kept. It exposes `isOpen`, `toggle`, and the set.
- Ignore a kept value that is not a list of the board's states.
- Test the default, the toggle, the persistence across a reload, and a corrupt stored value.
- Waits for nothing in this story, which waits for S-0458.

## Done when

`lanecollapse.svelte.test.ts` passes.

## Notes
