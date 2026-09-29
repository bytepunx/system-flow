---
id: T-0549
type: task
nature: feature
title: A long thread shows its last two entries, with the earlier ones behind a show-earlier control
status: done
parent: S-0153
owner: alex
created: 2026-09-29T06:57:18Z
updated: 2026-09-29T06:59:36Z
transitions:
  - to: ready
    at: 2026-09-29T06:57:27Z
    by: agent-S-0153
  - to: in-progress
    at: 2026-09-29T06:57:27Z
    by: agent-S-0153
  - to: done
    at: 2026-09-29T06:59:36Z
    by: agent-S-0153
stream: S-0153
tags: []
touches: [flaiover/src/lib/components, design/system/flaiover-dashboard.md, docs/users]
usage:
  source: log
  seconds: 129
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 32
      output: 7616
      cache_read: 1984352
      cache_write: 32282
      cost: 0.8076
---

# T-0549 A long thread shows its last two entries, with the earlier ones behind a show-earlier control

## Work

In `flaiover/src/lib/components/Threads.svelte`, a thread with more than two entries shows its last two. Above them, a button reads "show N earlier entries" and expands the rest; once expanded, "hide earlier entries" collapses them again. The choice is per thread, kept while paging, and reset when the page reloads. The operator chose this on TH-0035. Update the component tests, `design/system/flaiover-dashboard.md`, and `docs/users/flaiover.md`.

## Done when

- `Threads.svelte.test.ts` pins that a thread with three or more entries shows two, that the button names the hidden count and expands and collapses, and that threads of two or fewer entries have no button.
- flaiover lint, svelte-check, and unit tests pass.
- The design doc and the user docs describe it.

## Notes
