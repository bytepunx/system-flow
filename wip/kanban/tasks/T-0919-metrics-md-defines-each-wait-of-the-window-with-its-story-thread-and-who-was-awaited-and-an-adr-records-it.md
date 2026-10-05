---
id: T-0919
type: task
nature: feature
title: metrics.md defines each wait of the window, with its story, thread, and who was awaited, and an ADR records it
status: backlog
parent: S-0215
owner: alex
created: 2026-10-05T05:44:51Z
updated: 2026-10-05T05:44:51Z
transitions: []
stream: S-0215
tags: [flai]
touches: [design/system/metrics.md, design/adrs]
---
# T-0919 metrics.md defines each wait of the window, with its story, thread, and who was awaited, and an ADR records it

## Work

`flai stats --json` gives a story's thread and review waits only as totals (`items[].wait_threads_seconds`, `wait_review_seconds`, `waiting.weeks[]`). The table of longest waits needs each wait on its own. `metrics.md` is the contract between `flai stats` and the dashboard, so a change to it needs an ADR.

- Under `## Planning, waiting, and claims (S-0205)` › `### Waiting` in `design/system/metrics.md`, define `waiting.longest[]`: the longest waits that overlap the window, longest first, at most 10. Give each wait `item`, `kind` (`thread` or `review`), `thread` (the thread ID, absent for review), `started`, `ended` (absent while it is still open), `seconds` (the part inside the window, and inside the story's `in-progress` intervals for a thread wait, as `wait_threads_seconds` counts it), and `awaited`. For a thread, `awaited` is the author of the entry that ended the wait. For review, it is the `by` of the move out of `review`. It is absent while the wait is open.
- Write the precision rule for it beside the others.
- Write an ADR that records the addition and links ADR-0081, and add it to `design/adrs/README.md`.

This task waits for nothing: the contract comes first, and the flai and dashboard tasks build to it.

## Done when

- `metrics.md` defines `waiting.longest[]` and every field of it, with its precision rule.
- A new ADR records the addition and is listed in `design/adrs/README.md`.
- The markdown lint and `flai check --strict` are clean.

## Notes
