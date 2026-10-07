---
id: T-1191
type: task
nature: feature
title: flai guard refuses a sub-agent's message_send and message_reply
status: done
parent: S-0331
owner: alex
created: 2026-10-07T20:14:37Z
updated: 2026-10-07T21:12:42Z
transitions:
  - to: ready
    at: 2026-10-07T21:01:58Z
    by: agent-S-0331
  - to: in-progress
    at: 2026-10-07T21:01:58Z
    by: agent-S-0331
  - to: done
    at: 2026-10-07T21:12:42Z
    by: agent-S-0331
stream: S-0331
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
usage:
  source: log
  seconds: 644
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 11800
      cache_read: 2011737
      cache_write: 71901
      cost: 1.0733
---
# T-1191 flai guard refuses a sub-agent's message_send and message_reply

## Work

Hold sub-agents to reading, as ADR-0060 does for thread writes. It waits for nothing: it needs only the tool names, which S-0331's criteria fix. It shares no path with T-1190 and runs beside it.

- Add `message_send`, `message_reply`, and the `flai message send` and `reply` commands to what `flai guard` refuses a sub-agent; `message_get` and `flai message list` and `show` pass.
- The planner, the orchestrator, and the analyzer have no story and are refused by the tools themselves; the guard needs no rule for them.

## Done when

- Tests in `guard_test.go` cover a sub-agent refused each write and allowed each read, and the story's agent allowed both.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
