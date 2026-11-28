---
id: T-3645
type: task
nature: improvement
title: "flai serve starts the orchestrator with the reason review:S-nnnn on a passed verified event of a story in review"
status: backlog
parent: S-0854
owner: alex
created: 2026-11-27T23:59:39Z
updated: 2026-11-27T23:59:39Z
transitions: []
stream: S-0854
tags: [flai, serve]
touches: [flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate_test.go]
after: [T-3642]
---
# T-3645 flai serve starts the orchestrator with the reason review:S-nnnn on a passed verified event of a story in review

## Work

This task waits for T-3642. It matches events against the `verified` event and the `passed-in-review` condition that T-3642 adds, and `Trigger.Holds` reads the `Occurrence` fields that task gives it.

In `flai/internal/serve/orchestrate_reasons.go`:

- Add `{event: verified, when: passed-in-review}` to `builtinTriggers`, the triggers used when the project has no role file.
- Give `eventNames` a case for `inbox.KindVerified` that names it `verified`.
- In `occurrenceOf`, fill the verdict from the event's `to`, the record flag from its `recorded` field, and the story's state from the item it reads.
- Map a matched verified reason to `review:S-nnnn` in `prime`, so that the run is primed with the review's story and its last verify result. Then `primeReasons` no longer falls back to the general pack. Keep `line()`'s text for it plain: the story passed its close-out and is in review.

Write a test in `flai/internal/serve/orchestrate_reasons_test.go` that reproduces I-0413. A story sits in review, and its resolve job's `verified` event with verdict `passed` comes after the orchestrator's last run. Before the fix nothing starts the orchestrator. With the fix, one reason `review:S-nnnn` does. Cover the cases that start nothing: the same event with `recorded` true, a `failed` verdict, and a story no longer in review. Cover the coalescing case too: the story's move to review and its verified event in one window give one reason. Extend `flai/internal/serve/orchestrate_test.go` where it asserts the built-in triggers or the reasons a look gathers.

## Done when

- A passed, non-recorded `verified` event on a story in review starts the orchestrator with the reason `review:S-nnnn`, primed as a review.
- The test reproducing I-0413 passes, and so do the cases that start nothing.
- `flai test` on the changed paths passes.

## Notes
