---
id: T-1652
type: task
nature: feature
title: The queue appends each event it acts on to a log under .flai-cache
status: backlog
parent: S-0404
owner: alex
created: 2026-10-16T05:01:32Z
updated: 2026-10-16T05:03:57Z
transitions: []
stream: S-0404
tags: [cli]
touches: [flai/internal/mergequeue/log.go, flai/internal/mergequeue/log_test.go, flai/internal/mergequeue/state.go, flai/internal/serve/mergequeue.go]
---
# T-1652 The queue appends each event it acts on to a log under .flai-cache

## Work

- `log.go`: append one JSON line per event to a log under `.flai-cache`, and read it back. Events are enqueued, batch built, verified, bisected, landed, sent back, and dequeued. Each line carries the time, the batch, and the stories.
- Keep the log as the conflict log of ADR-0138 is kept: append only, host-local, readable by `flai stats`.
- Call the writer at each event in `serve/mergequeue.go`, and in `flai queue add` and `remove` through the package.
- Test the round trip and a log cut off mid-line.
- First in its layer: it waits for no task.

## Done when

Each queue event is logged and read back, and the tests pass.

## Notes
