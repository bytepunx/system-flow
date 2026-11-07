---
id: T-1706
type: task
nature: feature
title: flai/internal/usage reads OpenCode's JSON events into neutral usage events, tested on recorded output from the known version
status: done
parent: S-0410
owner: alex
created: 2026-10-16T05:18:42Z
updated: 2026-11-07T21:02:38Z
transitions:
  - to: ready
    at: 2026-11-07T20:46:29Z
    by: agent-S-0410
  - to: in-progress
    at: 2026-11-07T20:46:30Z
    by: agent-S-0410
  - to: done
    at: 2026-11-07T21:02:38Z
    by: agent-S-0410
stream: S-0410
tags: [cli]
touches: [flai/internal/usage/opencode.go, flai/internal/usage/opencode_test.go, flai/internal/usage/testdata/opencode/run.jsonl, flai/internal/usage/testdata/opencode/resumed.jsonl]
usage:
  source: log
  seconds: 968
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 53
      output: 18313
      cache_read: 8197087
      cache_write: 143434
      cost: 2.9069
---
# T-1706 flai/internal/usage reads OpenCode's JSON events into neutral usage events, tested on recorded output from the known version

## Work

- Add `OpenCode()` in `flai/internal/usage/opencode.go`, a `Reader` like `ClaudeCode()`: `step_start` and the session ID to a session start; `step_finish` to a model call with its tokens by kind, its model, and `cost` when reported; a `tool_use` of the Task tool to the call that started a sub-agent; a `wait_for_events` that timed out to an empty wake, as `waitTool` does; and the end of the output to the run's end.
- Record a first run and a resumed one from the known `opencode` version, keys and paths scrubbed, as `testdata/opencode/run.jsonl` and `resumed.jsonl`.
- Test the totals, a sub-agent's calls kept apart, a run with no cost reported, and a malformed line, which is skipped and counted.
- Waits for nothing in this story.

## Done when

The reader gives the recorded runs' tokens, models, and costs, and its tests pass.

## Notes
