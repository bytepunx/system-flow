---
id: T-1626
type: task
nature: feature
title: The dashboard's acceptance confirmation lists the forge blockers by name
status: done
parent: S-0395
owner: alex
created: 2026-10-16T04:38:18Z
updated: 2026-10-17T17:57:11Z
transitions:
  - to: ready
    at: 2026-10-17T17:36:42Z
    by: agent-S-0395
  - to: in-progress
    at: 2026-10-17T17:46:34Z
    by: agent-S-0395
  - to: done
    at: 2026-10-17T17:57:11Z
    by: agent-S-0395
stream: S-0395
tags: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
after: [T-1623]
usage:
  source: log
  seconds: 637
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 44
      output: 13849
      cache_read: 2437193
      cache_write: 75149
      cost: 1.2136
---
# T-1626 The dashboard's acceptance confirmation lists the forge blockers by name

## Work

- Check whether the review page's acceptance confirmation lists the preview's blockers generically. Where it does not show T-1623's kinds by name, add them: each check or review linking its URL from the kept state.
- Add a test beside the page's other tests. If the page has none, test the component that renders the blockers.
- Waits for T-1623: it shows that task's blockers.

## Done when

The confirmation shows each forge blocker by name with its link, and the flaiover tests pass.

## Notes
