---
id: T-2689
type: task
nature: improvement
title: The request answered event gives the send's time in its own field, send_ms, logged after the write
status: done
parent: S-0611
owner: alex
created: 2026-11-07T23:55:15Z
updated: 2026-11-09T17:09:13Z
transitions:
  - to: ready
    at: 2026-11-09T17:00:55Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:00:55Z
    by: system-flow
  - to: done
    at: 2026-11-09T17:09:13Z
    by: agent-S-0611
stream: S-0611
tags: [cli]
touches: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flai/internal/perf/perf.go, flai/internal/perf/perf_test.go, docs/users/flai.md]
usage:
  source: log
  seconds: 498
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 55
      output: 355
      cache_read: 7694512
      cache_write: 40548
      cost: 2.8233
---
# T-2689 The request answered event gives the send's time in its own field, send_ms, logged after the write

## Work

Criterion 3. Today the channel's request goroutine logs the event with `rec.Log` (`flai/internal/channel/channel.go`, about line 827) before `send(out)` (line 828), so `duration_ms` leaves out the envelope's marshal, the wait for the connection's write lock `wmu`, and `conn.Write`.

- Add a field for the send to `perf.Answered` (`flai/internal/perf/perf.go`), written by `Recorder.Log` as `send_ms` when it is set.
- In the channel, time `send(out)` and log the event after it returns, with `send_ms` from the start of the send to the end of the write, the lock's wait included. Keep `duration_ms` as it is: the time to the answer ready to write. A failed send still logs, with its `err`.
- Decide whether the slow threshold compares `duration_ms` alone or with `send_ms` added, and say which in the docs. Recommended: the sum, so a large answer stuck behind the lock is logged at `INFO`.
- MCP answers go through the SDK, which flai does not time; leave them without `send_ms` and say so.
- Update `docs/users/flai.md` § Request timing (S-0152): the sentence on what is measured, and a `send_ms` row in the field table.

Waits for no task: no other task touches these files.

## Done when

- [ ] `channel_test.go`'s test of the `request answered` event checks a `send_ms` field, and a test shows it counts a wait on the write lock.
- [ ] `perf_test.go` checks the field is written when set and left out when not.
- [ ] `docs/users/flai.md` § Request timing names `send_ms` and what it covers.
- [ ] `flai test` passes on the changed paths.

## Notes

- Planned by planner-S-0611.
