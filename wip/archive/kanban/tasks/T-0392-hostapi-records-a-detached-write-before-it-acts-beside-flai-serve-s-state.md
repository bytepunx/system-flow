---
id: T-0392
type: task
nature: remediation
title: hostapi records a detached write before it acts, beside flai serve's state
status: done
parent: S-0109
owner: alex
created: 2026-09-24T07:37:30Z
updated: 2026-09-24T07:41:34Z
transitions:
  - to: ready
    at: 2026-09-24T07:37:34Z
    by: agent-S-0109
  - to: in-progress
    at: 2026-09-24T07:37:35Z
    by: agent-S-0109
  - to: done
    at: 2026-09-24T07:41:34Z
    by: agent-S-0109
stream: S-0109
tags: []
touches: [flai/internal/hostapi, flai/internal/serve, flai/cmd]
---
# T-0392 hostapi records a detached write before it acts, beside flai serve's state

## Work

- `hostapi.Host` gets `Requests`, a function that changes the records of detached writes under a lock. `cmd` wires it to `serve/requests.json` beside flai serve's state; with none, the records are kept in memory.
- A write with `detachTimeout` records its request ID as started, with serve's PID, before it acts, and its outcome when it ends. A repeat of that ID is answered from the record: the recorded outcome, or that it is under way, or that the serve that took it ended before recording the outcome. It is never done again.
- If the record cannot be written, the write is refused rather than done unrecorded.
- Behaviour tests with a fake runner and an in-memory store: a repeat after a finished write, a repeat while under way, a repeat after the serve that took it ended, and a refusal when the record cannot be written. A test in `serve` for the file.

## Done when

- `make test` and lint pass, and the new tests fail without the change.

## Notes
