---
id: T-0193
type: task
nature: feature
title: Activity and inbox pages, and the inbox count in the navigation
status: done
parent: S-0042
owner: alex
created: 2026-09-19T04:46:17Z
updated: 2026-09-19T04:50:43Z
transitions:
  - to: ready
    at: 2026-09-19T04:48:45Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:48:45Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:50:43Z
    by: system-flow
stream: S-0042
tags: []
touches: [flaiover/src/routes, flaiover/src/lib/components]
---

# T-0193 Activity and inbox pages, and the inbox count in the navigation

## Work
`/activity`: a table or cards of active streams with the fields above, ages shown with the existing `age` helper, blocked flagged, links to the narrative and the item. `/inbox`: entries grouped by kind, each linking to its page (the thread's anchor, the narrative, the review page for a story in review, the item for a blocked one), with an empty state. Navigation gains Activity and Inbox, the latter with a count badge that refreshes on the existing SSE change events. Both pages refresh on change events too. Component tests for the two views and the badge.

## Done when
- The tests pass; lint and svelte-check are clean
- A production build succeeds

## Notes
