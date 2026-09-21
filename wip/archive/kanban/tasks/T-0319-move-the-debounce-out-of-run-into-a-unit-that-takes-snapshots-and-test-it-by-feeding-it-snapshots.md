---
id: T-0319
type: task
nature: feature
title: Move the debounce out of Run into a unit that takes snapshots, and test it by feeding it snapshots
status: done
parent: S-0093
owner: alex
created: 2026-09-21T22:49:14Z
updated: 2026-09-21T22:50:18Z
transitions:
  - to: ready
    at: 2026-09-21T22:49:19Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T22:49:19Z
    by: system-flow
  - to: review
    at: 2026-09-21T22:50:18Z
    by: system-flow
  - to: done
    at: 2026-09-21T22:50:18Z
    by: system-flow
stream: S-0093
tags: []
---
# T-0319 Move the debounce out of Run into a unit that takes snapshots, and test it by feeding it snapshots

## Work
Extract the per-tick logic of `Run` in `flai/internal/watch/watch.go` into a small type holding the previous snapshot and the pending files, with one method that takes the next snapshot and returns the sorted, de-duplicated paths now ready to report. `Run` keeps the ticker, calls `snapshot`, and calls `emit`. Add tests that feed it snapshots: a file changing on consecutive ticks is reported once after it holds still for one tick; a file changed, held, then changed again is reported per settle; an added file; a removed file; a file removed while still pending is reported once; a file changing back to its earlier value while pending. Replace `TestAFileStillBeingWrittenIsReportedOnce` with the unit version and delete its sleeps.

## Done when
`Run` behaves exactly as before, the new tests do not import `time` or touch the filesystem, and `go test -race ./internal/watch` passes.

## Notes
A missed edge here would be a behavior change, so the existing real-ticker test stays as the end-to-end check.
