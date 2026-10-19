---
id: T-1641
type: task
nature: feature
title: flai queue add, remove, hold, and resume write the queue's state, and flai accept by hand takes its story out of the queue
status: done
parent: S-0399
owner: alex
created: 2026-10-16T05:00:51Z
updated: 2026-10-19T19:21:02Z
transitions:
  - to: ready
    at: 2026-10-19T19:00:27Z
    by: agent-S-0399
  - to: in-progress
    at: 2026-10-19T19:00:28Z
    by: agent-S-0399
  - to: done
    at: 2026-10-19T19:21:02Z
    by: agent-S-0399
stream: S-0399
tags: [cli]
touches: [flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go, docs/users/flai-reference.md, flai/internal/mergequeue/eligible.go]
usage:
  source: log
  seconds: 1234
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 29001
      cache_read: 4379401
      cache_write: 160626
      cost: 2.3177
---
# T-1641 flai queue add, remove, hold, and resume write the queue's state, and flai accept by hand takes its story out of the queue

## Work

- Give `state.go` a writer, under a file lock, so that `flai serve` and a shell never write the state at once.
- `flai queue add S-nnnn` records the story with `by` (`FLAI_AGENT`, else the config's author) and `at`. It refuses a story that is not eligible, naming each reason from `eligible.go`, and refuses under pull-request acceptance, saying the queue runs under local acceptance only.
- `flai queue remove S-nnnn` dequeues a story. `flai queue hold` and `resume` set the held flag.
- `flai accept`, once it has merged a story, removes the story from the queue and marks a batch that holds it as dropped.
- Test each command, its `--json`, its refusals (the pull-request mode among them), and the acceptance's removal.
- First in its layer: it waits for no task.

## Done when

The four commands and the acceptance change the state as above, and their tests pass.

## Notes
