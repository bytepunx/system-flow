---
id: T-3407
type: task
nature: remediation
title: Count a line main added as undone only by the ADR's rule in FindUndone and RecordedUndone, with unit tests
status: done
parent: S-0800
owner: alex
created: 2026-11-16T01:44:51Z
updated: 2026-11-16T02:09:41Z
transitions:
  - to: ready
    at: 2026-11-16T01:47:31Z
    by: agent-S-0800
  - to: in-progress
    at: 2026-11-16T01:48:55Z
    by: agent-S-0800
  - to: done
    at: 2026-11-16T02:09:41Z
    by: agent-S-0800
stream: S-0800
tags: [storygit, go]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/reads_test.go, flai/internal/mergequeue/eligible_test.go]
after: [T-3405]
usage:
  source: log
  seconds: 1246
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 39
      output: 10002
      cache_read: 2792721
      cache_write: 103282
      cost: 1.4266
---
# T-3407 Count a line main added as undone only by the ADR's rule in FindUndone and RecordedUndone, with unit tests

## Work

Implement T-3405's rule in `flai/internal/storygit/undone.go`. Under the position rule, `diffLines` gains the line numbers from each `-U0` hunk header, and the compare takes main's added lines at their positions in main's file at `onto` against the lines the branch's net change removes there. `RecordedUndone` does the same against main now, from the recorded old base, and keeps only the recorded lines. ADR-0297's markdown words rule stays as it is, applied to what the new count leaves. Keep `FindUndone`'s signature, so `sync.go`, `flai/internal/preview/accept.go`, and `flai/internal/mergequeue/eligible.go` need no change; if one does, widen this task's touches.

Add unit tests to `flai/internal/storygit/undone_test.go`:

- I-0385's case: main adds a function ending in a tab and a closing brace, and the branch removes another closing brace in the same file. Nothing is undone.
- I-0124's case: the branch's revert removes the lines main added, where main added them. The path is undone.
- The existing tests, the I-0338 case and ADR-0297's markdown cases among them, still pass.

It waits for T-3405, whose ADR sets the rule it implements.

## Done when

- `flai test flai/internal/storygit/undone.go flai/internal/storygit/undone_test.go` passes, with the I-0385 and I-0124 cases above among the tests.
- The comment above `FindUndone` and `notKept` names the new ADR beside ADR-0155 and ADR-0297.

## Notes
