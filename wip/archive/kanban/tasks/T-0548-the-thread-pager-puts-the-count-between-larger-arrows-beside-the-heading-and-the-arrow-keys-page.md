---
id: T-0548
type: task
nature: feature
title: The thread pager puts the count between larger arrows beside the heading, and the arrow keys page
status: done
parent: S-0153
owner: alex
created: 2026-09-29T06:46:43Z
updated: 2026-09-29T06:50:45Z
transitions:
  - to: ready
    at: 2026-09-29T06:46:58Z
    by: agent-S-0153
  - to: in-progress
    at: 2026-09-29T06:46:58Z
    by: agent-S-0153
  - to: done
    at: 2026-09-29T06:50:45Z
    by: agent-S-0153
stream: S-0153
tags: []
touches: [flaiover/src/lib/components, design/system/flaiover-dashboard.md, docs/users]
usage:
  source: log
  seconds: 227
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 35
      output: 8355
      cache_read: 2176939
      cache_write: 35415
      cost: 0.8859
---

# T-0548 The thread pager puts the count between larger arrows beside the heading, and the arrow keys page

## Work

In `flaiover/src/lib/components/Threads.svelte`, lay the pager out as `← n of m →` with the count between the arrows. Show it once, in the Threads header beside the heading, with arrows the size of the header's other buttons, instead of small text above and below the thread. While focus is in the pager, Left and Right page. Update the component tests, `design/system/flaiover-dashboard.md`, and `docs/users`.

## Done when

- `Threads.svelte.test.ts` pins the order previous, count, next; the arrows paging and stopping at the ends; and the arrow keys paging.
- `make test` passes for flaiover, and lint and check pass.
- The design doc and the user docs describe the pager as built.

## Notes
