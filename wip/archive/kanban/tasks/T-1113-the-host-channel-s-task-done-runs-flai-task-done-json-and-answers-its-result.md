---
id: T-1113
type: task
nature: improvement
title: The host channel's task.done runs flai task done --json and answers its result
status: done
parent: S-0269
owner: alex
created: 2026-10-06T22:53:42Z
updated: 2026-10-07T03:02:55Z
transitions:
  - to: ready
    at: 2026-10-07T02:58:34Z
    by: agent-S-0269
  - to: in-progress
    at: 2026-10-07T02:58:34Z
    by: agent-S-0269
  - to: done
    at: 2026-10-07T03:02:55Z
    by: agent-S-0269
stream: S-0269
tags: [hostapi]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts]
after: [T-1102]
usage:
  source: log
  seconds: 261
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 13480
      cache_read: 1976391
      cache_write: 73801
      cost: 1.114
---
# T-1113 The host channel's task.done runs flai task done --json and answers its result

## Work

Add the operation on the host channel. Host writes there are specs that build `flai` arguments; `item.move` and `item.criteria` are two of them.

- In `flai/internal/hostapi/writes.go`, add the spec `task.done`. Its params are `{id, message, log}`. It checks that `id` is a task and that `message` is not empty, then builds `task done <id> -m <message> [--log <entry>] --by=<owner> --json`; leave `--by` out if the command takes none.
- Map the command's exit codes for a refused sync and a failed check (T-1102) in `exits`, so the dashboard gets the result together with the step that stopped it.
- In `flai/internal/hostapi/writes_test.go`, cover the arguments built, a missing message, a story ID refused, and each mapped exit code.

This task waits for T-1102, whose flags and exit codes it builds on. It is layer 4, beside the conventions task, with which it shares no file.

## Done when

- [ ] `task.done` on the host channel runs `flai task done --json` and answers the same result.
- [ ] Tests cover its arguments, its refusals, and its exit codes, and `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner. No dashboard page calls `task.done` yet; this story only makes it reachable.
