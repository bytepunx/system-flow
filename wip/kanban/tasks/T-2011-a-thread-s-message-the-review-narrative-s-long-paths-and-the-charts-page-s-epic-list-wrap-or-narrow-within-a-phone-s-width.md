---
id: T-2011
type: task
nature: feature
title: A thread's message, the review narrative's long paths, and the charts page's epic list wrap or narrow within a phone's width
status: in-progress
parent: S-0461
owner: alex
created: 2026-10-21T12:49:03Z
updated: 2026-10-21T12:49:26Z
transitions:
  - to: ready
    at: 2026-10-21T12:49:25Z
    by: agent-S-0461
  - to: in-progress
    at: 2026-10-21T12:49:26Z
    by: agent-S-0461
stream: S-0461
tags: []
touches: [flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/Review.svelte, "flaiover/src/routes/charts/[kind]/+page.svelte"]
---

# T-2011 A thread's message, the review narrative's long paths, and the charts page's epic list wrap or narrow within a phone's width

## Work

- T-2010's end-to-end test found three things that make a page scroll sideways at 375px: a thread message's long unbroken path (313px past the window on the review and documents pages), an inline code path in the review narrative's state (57px, and 4px at 1280px in the two-column grid), and the charts page's epic list, as wide as its longest epic title (78px).
- Let a thread message and the narrative's state and next steps break a long word anywhere (`wrap-anywhere`), and cap the epic list's width so a long title is cut in the closed list.
- Waits for nothing. T-2010 waits for it, as its test fails until these are fixed.

## Done when

`scripts/flaiover-e2e.sh e2e/page-sidebar.e2e.ts` passes with T-2010's realistic data, and `flai test` passes on the paths changed.

## Notes
