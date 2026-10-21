---
id: T-1947
type: task
nature: improvement
title: Move flai verify's step findings formatter into internal/verify so the merge queue can quote a failing tier as flai verify prints it
status: done
parent: S-0457
owner: alex
created: 2026-10-21T06:09:33Z
updated: 2026-10-21T06:27:02Z
transitions:
  - to: ready
    at: 2026-10-21T06:13:35Z
    by: agent-S-0457
  - to: in-progress
    at: 2026-10-21T06:13:36Z
    by: agent-S-0457
  - to: done
    at: 2026-10-21T06:27:02Z
    by: agent-S-0457
stream: S-0457
tags: [flai]
touches: [flai/internal/verify/text.go, flai/internal/verify/text_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go]
usage:
  source: log
  seconds: 806
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 4013
      cache_read: 522939
      cache_write: 26365
      cost: 0.3401
---
# T-1947 Move flai verify's step findings formatter into internal/verify so the merge queue can quote a failing tier as flai verify prints it

## Work

- Move `stepFindings` from `flai/cmd/verify.go` into `flai/internal/verify/text.go` as an exported function, such as `StepFindings(s Step, indent string) string`: every line of each finding's `Text()`, the plain tier's failure lines and tail included, and the `… N more findings left out` line. Package `mergequeue` cannot import `cmd`, so the queue's thread needs it here.
- Make `verifyText` in `flai/cmd/verify.go` call it, so `flai verify` and `flai test` print exactly what they print now.
- Test it in `flai/internal/verify/text_test.go` with a plain finding of several lines and with findings left out.

It waits for nothing: it changes no path the other first-layer task changes, so the two run together.

## Done when

- `flai verify` prints a failing step's findings unchanged; `flai/cmd/verify_test.go` passes as it is.
- The exported formatter has its own test in `flai/internal/verify/text_test.go`, and `flai test` on the changed paths passes.

## Notes
