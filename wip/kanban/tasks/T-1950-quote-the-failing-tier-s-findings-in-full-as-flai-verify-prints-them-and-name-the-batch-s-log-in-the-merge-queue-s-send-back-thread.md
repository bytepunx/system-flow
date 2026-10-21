---
id: T-1950
type: task
nature: improvement
title: Quote the failing tier's findings in full, as flai verify prints them, and name the batch's log in the merge queue's send-back thread
status: done
parent: S-0457
owner: alex
created: 2026-10-21T06:10:12Z
updated: 2026-10-21T06:54:50Z
transitions:
  - to: ready
    at: 2026-10-21T06:48:35Z
    by: agent-S-0457
  - to: in-progress
    at: 2026-10-21T06:48:36Z
    by: agent-S-0457
  - to: done
    at: 2026-10-21T06:54:50Z
    by: agent-S-0457
stream: S-0457
tags: [flai]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/bisect_test.go]
after: [T-1947, T-1949]
usage:
  source: log
  seconds: 374
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 47
      output: 15867
      cache_read: 2067499
      cache_write: 104238
      cost: 1.3447
---
# T-1950 Quote the failing tier's findings in full, as flai verify prints them, and name the batch's log in the merge queue's send-back thread

## Work

- Keep on `Result` in `flai/internal/mergequeue/state.go` the failing step's findings as T-1947's formatter prints them, such as `Findings string`: a plain tier's failure lines and the tail of its output, every line, not the first alone. Cap it, such as at 8,000 characters, marked where it is cut, so `state.json` and the thread stay readable. `result` in `flai/internal/mergequeue/batch.go` fills it from the failing step; `Summary` stays the short one-line form it is.
- In `failingText` in `flai/internal/mergequeue/bisect.go`, quote those findings in the fenced block in place of the summary, falling back to the summary when a result has none (one kept by an older flai), and name the log folder T-1949 records, so the reader can open the whole output. The move's reason stays one line.
- Reproduce I-0172 in `flai/internal/mergequeue/bisect_test.go`: a plain tier whose first line is a heading such as `smoke: template render and check` and whose failure is lines further down fails a batch; the send-back thread quotes those lines and the tail and names the log. Update `TestTheResultSummarisesTheVerification` in `flai/internal/mergequeue/batch_test.go` only where the result's new field changes what it asserts.

It waits for T-1947, whose formatter it calls, and for T-1949, which changes the same files and records the log folder it names.

## Done when

- The test reproducing I-0172 passes, and fails against the first-line summary.
- The send-back thread of a red batch quotes the failing tier's findings in full, within the cap, and names its log folder.
- `flai test` on the changed paths passes.

## Notes
