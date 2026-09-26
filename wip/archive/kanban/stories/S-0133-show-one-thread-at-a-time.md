---
id: S-0133
type: story
nature: improvement
title: Show one thread at a time
status: done
parent: E-0003
owner: alex
created: 2026-09-26T17:35:44Z
updated: 2026-09-26T17:53:35Z
transitions:
  - to: ready
    at: 2026-09-26T17:48:30Z
    by: alex
  - to: in-progress
    at: 2026-09-26T17:48:56Z
    by: agent-S-0133
  - to: review
    at: 2026-09-26T17:52:58Z
    by: agent-S-0133
  - to: done
    at: 2026-09-26T17:53:35Z
    by: alex
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0133 Show one thread at a time

## Goal

Instead of having to scroll vertically between threads, show only one thread at a time with an indicator above and below showing which thread the user is viewing (i.e. "4 of 5") with navigation arrows to the right of the indicator.

## Acceptance criteria
- [x] Stories with multiple threads only show one thread at a time
- [x] An indicator correctly shows which thread the operator is reading
- [x] Navigation arrows allow the operator to change which thread they are viewing

## Tasks
- T-0468 Page the threads on an anchor one at a time with an indicator and arrows
- T-0469 Document the thread pager for users and builders

## Notes

Verified by behaviour tests in `flaiover/src/lib/components/Threads.svelte.test.ts` (one `article` among several threads, "n of m" above and below, arrows move and stop at the ends, place kept across reloads), not in a browser.
