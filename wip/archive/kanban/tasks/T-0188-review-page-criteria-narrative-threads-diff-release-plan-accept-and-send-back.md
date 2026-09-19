---
id: T-0188
type: task
nature: feature
title: "Review page: criteria, narrative, threads, diff, release plan, accept, and send back"
status: done
parent: S-0041
owner: alex
created: 2026-09-19T04:08:29Z
updated: 2026-09-19T04:16:10Z
transitions:
  - to: ready
    at: 2026-09-19T04:12:36Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:12:36Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:16:10Z
    by: system-flow
stream: S-0041
tags: []
touches: [flaiover/src/routes, flaiover/src/lib/components]
---

# T-0188 Review page: criteria, narrative, threads, diff, release plan, accept, and send back

## Work
Add `/review/[id]`, linked from the item page and from cards in the review column. For a story in review it shows: the acceptance criteria with their ticks and how many are ticked; the narrative's Current state and Next steps; open threads on the story with the existing `Threads` component; the branch diff as a file list with additions and deletions and collapsible hunks, additions and deletions coloured and also marked with + and -; and the release plan and blockers from the existing acceptance preview, with the uncommitted-files choice from S-0051. Accept calls the streaming endpoint and shows each step as it completes, then the tags, or the error verbatim with the story still in review. Send back asks for a reason in the page, not with `prompt()`, and runs the existing move to in-progress. For a story not in review the page says what state it is in and links to the item. Component tests for the diff view and the accept progress; the board and item page keep their existing confirmation.

## Done when
- The tests pass; lint and svelte-check are clean
- A production build succeeds

## Notes
