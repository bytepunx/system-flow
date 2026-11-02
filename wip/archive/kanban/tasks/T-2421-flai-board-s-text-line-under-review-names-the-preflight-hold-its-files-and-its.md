---
id: T-2421
type: task
nature: improvement
title: flai board's text line under review names the preflight hold, its files, and its thread
status: done
parent: S-0549
owner: alex
created: 2026-11-02T00:38:18Z
updated: 2026-11-02T00:44:52Z
transitions:
  - to: ready
    at: 2026-11-02T00:39:02Z
    by: agent-S-0549
  - to: in-progress
    at: 2026-11-02T00:39:09Z
    by: agent-S-0549
  - to: done
    at: 2026-11-02T00:44:52Z
    by: agent-S-0549
stream: S-0549
tags: []
touches: [flai/cmd/board.go, flai/cmd/board_queue_test.go]
usage:
  source: log
  seconds: 343
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 10665
      cache_read: 2122291
      cache_write: 45947
      cost: 0.9162
---
# T-2421 flai board's text line under review names the preflight hold, its files, and its thread

## Work

ADR-0234 says `flai board` shows the preflight hold. Its JSON does, under `queue.preflight`, but the text line under the review column (`queueSummary` in `flai/cmd/board.go`) shows neither hold. Add the preflight hold and its thread to that line, with a test. It waits for no task: the hold and `Queue.Preflight` are done.

## Done when

`flai board` text names the preflight hold, its files, and its thread, and a test pins it. Story criterion 4.

## Notes

Found by T-2387 against ADR-0234.
