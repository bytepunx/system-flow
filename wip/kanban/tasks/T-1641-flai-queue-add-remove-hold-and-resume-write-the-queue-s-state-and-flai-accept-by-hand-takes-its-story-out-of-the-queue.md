---
id: T-1641
type: task
nature: feature
title: flai queue add, remove, hold, and resume write the queue's state, and flai accept by hand takes its story out of the queue
status: backlog
parent: S-0399
owner: alex
created: 2026-10-16T05:00:51Z
updated: 2026-10-16T05:00:51Z
transitions: []
stream: S-0399
tags: [cli]
touches: [flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go]
---
# T-1641 flai queue add, remove, hold, and resume write the queue's state, and flai accept by hand takes its story out of the queue

## Work

- Give `state.go` a writer, under a file lock, so that `flai serve` and a shell never write the state at once.
- `flai queue add S-nnnn` records the story with `by` (`FLAI_AGENT`, else the config's author) and `at`. It refuses a story that is not eligible, naming each reason from `eligible.go`.
- `flai queue remove S-nnnn` dequeues a story. `flai queue hold` and `resume` set the held flag.
- `flai accept`, once it has merged a story, removes the story from the queue and marks a batch that holds it as dropped.
- Test each command, its `--json`, its refusals, and the acceptance's removal.
- First in its layer: it waits for no task.

## Done when

The four commands and the acceptance change the state as above, and their tests pass.

## Notes
