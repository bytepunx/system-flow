---
id: T-1589
type: task
nature: remediation
title: The channel client logs a proven connection before it publishes the state, with a test that makes the race certain
status: done
parent: S-0389
owner: alex
created: 2026-10-14T23:22:56Z
updated: 2026-10-14T23:23:35Z
transitions:
  - to: ready
    at: 2026-10-14T23:23:30Z
    by: agent-S-0389
  - to: in-progress
    at: 2026-10-14T23:23:31Z
    by: agent-S-0389
  - to: done
    at: 2026-10-14T23:23:34Z
    by: agent-S-0389
stream: S-0389
tags: [cli]
touches: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md]
after: [T-1587]
usage:
  source: log
  seconds: 3
  estimated: true
  priced_by: harness
  models: []
---
# T-1589 The channel client logs a proven connection before it publishes the state, with a test that makes the race certain

## Work

The close-out's integration tier failed `TestAConnectionRecordsAndWarnsOfTheUnsignedSides` with the same cause as I-0146 in another place: `serveOnce` in `flai/internal/channel/channel.go` publishes the connected state before it logs "connected to the dashboard" and "connection allowed unsigned", so a reader that sees the state connected can read the log before the warning is written.

- Move both log lines before the `setState` that publishes the connection.
- Add a test whose logger writes through the slow writer T-1587 added: it waits for the state to read connected, then expects the unsigned warning in the log at once.

Waits for T-1587, which adds the slow writer the test uses and changes the same two files.

## Done when

- The new test fails on the code before the move and passes after it.
- `go test -race -count=20 ./internal/channel/` in `flai/` passes.
- `flai test flai/internal/channel/channel.go flai/internal/channel/channel_test.go` passes.

## Notes
