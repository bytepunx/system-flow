---
id: T-0923
type: task
nature: feature
title: flai stats reports the longest waits of the window with their story, thread, and who was awaited
status: backlog
parent: S-0215
owner: alex
created: 2026-10-05T05:44:57Z
updated: 2026-10-05T05:44:57Z
transitions: []
stream: S-0215
tags: [flai]
touches: [flai/internal/metrics]
after: [T-0919]
---
# T-0923 flai stats reports the longest waits of the window with their story, thread, and who was awaited

## Work

In `flai/internal/metrics/waiting.go`, `threadWait` keeps only a span. Make it keep the thread's ID and the author of the entry that ended the wait. Then build `Waiting.Longest` as `metrics.md` defines it: the thread waits clipped to each story's `in-progress` intervals, and the review intervals with the `by` of the move out of `review`, all clipped to the window, longest first, at most 10. The totals and means already reported stay as they are.

It waits for T-0919, whose `metrics.md` section and ADR fix the fields and their rules.

## Done when

- `flai stats --json` carries `waiting.longest[]` with every field `metrics.md` defines.
- Tests in `flai/internal/metrics/waiting_test.go` pin it on fixtures, to the second: a thread answered by the operator, one still open, a review wait, a wait cut by the window's start, and the cap of 10.
- `scripts/flai-test.sh` passes.

## Notes
