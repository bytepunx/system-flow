---
id: T-1587
type: task
nature: remediation
title: The channel client logs an attempt's end before it publishes the state, with a test that makes the race certain
status: done
parent: S-0389
owner: alex
created: 2026-10-14T23:10:59Z
updated: 2026-10-14T23:14:33Z
transitions:
  - to: ready
    at: 2026-10-14T23:12:48Z
    by: agent-S-0389
  - to: in-progress
    at: 2026-10-14T23:12:48Z
    by: agent-S-0389
  - to: done
    at: 2026-10-14T23:14:33Z
    by: agent-S-0389
stream: S-0389
tags: [cli]
touches: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go]
usage:
  source: log
  seconds: 105
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 4979
      cache_read: 1918355
      cache_write: 35841
      cost: 0.7701
---
# T-1587 The channel client logs an attempt's end before it publishes the state, with a test that makes the race certain

## Work

`Client.Run` in `flai/internal/channel/channel.go` publishes an attempt's end with `setState` and only then logs it. A reader that waits for the state to read refused can read the log before the line is written. `TestTheGateIsAskedOnEveryAttempt` does that at `channel_test.go:883`, which is the race I-0146 records twice.

- Reorder `Run` so that every attempt's log line, held, unsigned, gate, or ended, is written before `setState` publishes that attempt's state. `proven` decides whether `lastLogged` is reset, and is read inside `setState` today: take it from the state before the log, since `Since` is the attempt's own and nothing else changes it once `attempt` has returned.
- Add a test that makes the race certain: a gated client whose logger writes through a writer that sleeps before it appends, such as 100ms. It waits for the state to read refused, then expects the refusal in the log at once. On the code before the reorder it fails; after it, it passes.
- Leave `TestTheGateIsAskedOnEveryAttempt` as it is. The reorder is what makes its read of the log safe; a polled wait in the test would hide the order instead of fixing it.

Waits for nothing: it is the first layer.

## Done when

- The new test fails on the code before the reorder and passes after it.
- `go test -race -count=50 -run 'TestTheGateIsAskedOnEveryAttempt|<the new test>' ./internal/channel/` in `flai/` passes.
- `flai test flai/internal/channel/channel.go flai/internal/channel/channel_test.go` passes.

## Notes
