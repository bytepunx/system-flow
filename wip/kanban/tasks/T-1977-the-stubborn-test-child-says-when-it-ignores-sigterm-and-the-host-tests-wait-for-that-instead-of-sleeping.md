---
id: T-1977
type: task
nature: remediation
title: The stubborn test child says when it ignores SIGTERM, and the host tests wait for that instead of sleeping
status: done
parent: S-0463
owner: alex
created: 2026-10-21T07:32:07Z
updated: 2026-10-21T07:36:28Z
transitions:
  - to: ready
    at: 2026-10-21T07:34:15Z
    by: agent-S-0463
  - to: in-progress
    at: 2026-10-21T07:34:16Z
    by: agent-S-0463
  - to: done
    at: 2026-10-21T07:36:28Z
    by: agent-S-0463
stream: S-0463
tags: [flai, tests]
touches: [flai/internal/host/host_test.go]
usage:
  source: log
  seconds: 132
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 4126
      cache_read: 1516569
      cache_write: 24120
      cost: 0.5789
---
# T-1977 The stubborn test child says when it ignores SIGTERM, and the host tests wait for that instead of sleeping

## Work

- In `flai/internal/host/host_test.go`, have the `stubborn` child of `TestMain` write a file named by an environment variable (as `run` writes `FLAI_HOST_TEST_ENVFILE`) once `signal.Ignore(syscall.SIGTERM)` has returned.
- Give `TestAChildThatWillNotStopIsKilled` and `TestAChildIsStoppedOnlyOnceItsProcessHasEnded` that variable through `spec("stubborn", ...)`, with a path under `t.TempDir()`.
- Replace each `time.Sleep(100 * time.Millisecond) // it has ignored SIGTERM by now` with a poll for that file, failing the test after a deadline of a few seconds.
- Reproduce the cause I-0173 names: add a delay knob to the `stubborn` child, an environment variable that makes it sleep before it ignores SIGTERM. Show in the commit message that with a delay over 100 ms the old sleep fails and the wait for the file passes. Keep the knob only if a test uses it.
- First task of the story, waiting for none: the issue can only be closed once this fix is in.

## Done when

- No `time.Sleep` stands in for the stubborn child's readiness in `flai/internal/host/host_test.go`.
- `go test -race -count=20 -run 'TestAChild' ./internal/host` passes in `flai/`, alongside a parallel load such as a second `go test ./...`.
- `flai test flai/internal/host/host_test.go` passes.

## Notes

Drafted by the planner from I-0173's instance, which proposes this fix.
