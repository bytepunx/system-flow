---
id: T-3456
type: task
nature: feature
title: flai queue names the suites the batch runs, in its listing and its --json
status: done
parent: S-0809
owner: alex
created: 2026-11-16T23:49:32Z
updated: 2026-11-18T09:07:47Z
transitions:
  - to: ready
    at: 2026-11-18T09:03:18Z
    by: agent-S-0809
  - to: in-progress
    at: 2026-11-18T09:03:18Z
    by: agent-S-0809
  - to: done
    at: 2026-11-18T09:07:47Z
    by: agent-S-0809
stream: S-0809
tags: [flai]
touches: [flai/cmd/queue.go, flai/cmd/queue_test.go, docs/users/flai-reference.md]
after: [T-3453]
usage:
  source: log
  seconds: 269
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 3323
      cache_read: 913750
      cache_write: 31332
      cost: 0.4453
---
# T-3456 flai queue names the suites the batch runs, in its listing and its --json

## Work

ADR-0312 §6, criterion 2. In `printBatch` (`flai/cmd/queue.go`), print the suites the batch runs beside its stories, as `suites: behavior, smoke` or `suites: all`. `--json` gives them already once `Batch` carries `Suites`, through `mergequeue.Queue`; check that, and that the MCP `board` tool's `queue` carries them too. Say in the command's long help that the batch line names its suites, and regenerate the reference with `make flai-reference`, which writes `docs/users/flai-reference.md`. Waits for T-3453, which adds the field this prints.

## Done when

- [ ] A test in `flai/cmd/queue_test.go`, beside `TestQueueShowsTheBatchsResultAndDrop`, shows `flai queue` naming a batch's suites, for a list and for `all`, in text and in `--json`.
- [ ] `TestReferenceIsCurrent` passes.
- [ ] `flai test` on the changed paths passes.

## Notes
