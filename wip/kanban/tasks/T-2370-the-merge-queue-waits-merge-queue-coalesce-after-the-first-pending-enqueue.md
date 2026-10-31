---
id: T-2370
type: task
nature: improvement
title: The merge queue waits merge_queue.coalesce after the first pending enqueue before it builds a batch, and shows the window
status: in-progress
parent: S-0552
owner: alex
created: 2026-10-31T07:33:01Z
updated: 2026-10-31T07:34:54Z
transitions:
  - to: ready
    at: 2026-10-31T07:34:53Z
    by: system-flow
  - to: in-progress
    at: 2026-10-31T07:34:54Z
    by: system-flow
stream: S-0552
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, docs/operators/settings.md, design/system/project-manifest.md]
---
# T-2370 The merge queue waits merge_queue.coalesce after the first pending enqueue before it builds a batch, and shows the window

## Work

Add `merge_queue.coalesce` to the manifest: a Go duration, default 3m, 0 turning the window off, refused when it does not parse or is negative, with its settings catalog row, its row in `docs/operators/settings.md`, and its rows in `design/system/project-manifest.md`. `mergequeue.RunBatch` builds no batch while the stories it would pick are fewer than `merge_queue.batch_max` and the window, from the earliest enqueue among them plus `coalesce`, has not ended; a story enqueued while a batch runs joins the next batch, whose window has then usually ended. The queue's state (`Queue`, read by the board, `flai queue`, and the dashboard) carries the window: when it ends and the stories filling it. `flai queue` prints it. Waits for nothing: it shares no file with the other layer-1 tasks.

## Done when

- [ ] A test shows an enqueue alone builds no batch before the window ends, a second enqueue inside it joins the same batch, and `batch_max` stories start at once.
- [ ] The key's default, check, and docs rows exist, and the settings tests pass.

## Notes

Layer 1.
