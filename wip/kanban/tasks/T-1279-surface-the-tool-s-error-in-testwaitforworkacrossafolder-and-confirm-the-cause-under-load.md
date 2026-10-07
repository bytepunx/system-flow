---
id: T-1279
type: task
nature: research
title: Surface the tool's error in TestWaitForWorkAcrossAFolder and confirm the cause under load
status: done
parent: S-0310
owner: alex
created: 2026-10-07T23:30:16Z
updated: 2026-10-07T23:36:55Z
transitions:
  - to: ready
    at: 2026-10-07T23:33:16Z
    by: agent-S-0310
  - to: in-progress
    at: 2026-10-07T23:33:16Z
    by: agent-S-0310
  - to: done
    at: 2026-10-07T23:36:55Z
    by: agent-S-0310
stream: S-0310
tags: [flai, mcp, tests]
touches: [flai/internal/mcpserver/folder_test.go]
usage:
  source: log
  seconds: 219
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 72
      cache_read: 979173
      cache_write: 14377
      cost: 0.4376
---
# T-1279 Surface the tool's error in TestWaitForWorkAcrossAFolder and confirm the cause under load

## Work

I-0102's instance failed with `handed over: map[]` at `flai/internal/mcpserver/folder_test.go:203`. An empty map means the second `wait_for_work` call answered a tool error, not a timeout, which would carry `timed_out`. The test drops that error: the goroutine calls `out, _ := f.call(...)`.

The planner's hypothesis: `readyStoryIn` adds the criterion with `os.WriteFile`, which truncates the story file before it writes it. The held wait polls every 20ms. A poll that reads the file in between gets an empty or partial story. `workitem`'s `store.list` returns the parse error, `folderWaits.hold` returns it, and the call answers an error. A loaded host widens that window. flai's own writes go through `atomicfile.WriteFile` and have no such window.

1. Make the goroutine send the tool's error text with its answer, and fail with it. Do the same for the other calls in `folder_test.go` that drop the error (`, _ := f.call`).
2. Reproduce the failure under load: run the test alone many times (`go test ./internal/mcpserver -run 'TestWaitForWorkAcrossAFolder$' -count=300`) while the full Go suite runs beside it, and with `-cpu 1`.
3. Confirm or refute the hypothesis from the error text. Record it, with the runs and how often it failed, in the narrative's `## Decisions`.

This task waits for none. It comes first because the fix depends on the cause it finds.

## Done when

- A failed call in `folder_test.go` reports the tool's error text instead of an empty map.
- The cause of I-0102 is recorded in the narrative's `## Decisions` with the error observed under load, or, if it does not reproduce in 300 runs under load, with what was run and the cause the code supports.

## Notes
