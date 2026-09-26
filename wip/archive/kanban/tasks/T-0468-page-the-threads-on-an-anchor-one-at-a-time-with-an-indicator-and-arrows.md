---
id: T-0468
type: task
nature: feature
title: Page the threads on an anchor one at a time with an indicator and arrows
status: done
parent: S-0133
owner: alex
created: 2026-09-26T17:49:20Z
updated: 2026-09-26T17:51:55Z
transitions:
  - to: ready
    at: 2026-09-26T17:49:31Z
    by: agent-S-0133
  - to: in-progress
    at: 2026-09-26T17:49:31Z
    by: agent-S-0133
  - to: done
    at: 2026-09-26T17:51:55Z
    by: agent-S-0133
stream: S-0133
tags: []
touches: [flaiover/src/lib/components]
---

# T-0468 Page the threads on an anchor one at a time with an indicator and arrows

## Work

In `flaiover/src/lib/components/Threads.svelte`, show one thread of the list at a time. When there is more than one, put an indicator ("2 of 5") above and below the thread with previous and next arrows to its right. Keep the thread being read across reloads by its ID; a thread just opened becomes the one shown. Behaviour tests in `Threads.svelte.test.ts`.

## Done when

- With several threads only one `article` renders, and the indicator above and below names its place.
- The arrows move between threads and are disabled at the ends.
- A single thread shows no indicator; `make test` for flaiover and its lint and type check pass.

## Notes
