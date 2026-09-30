---
id: T-0605
type: task
nature: feature
title: Tell an agent's process by its start in clock ticks, and settle a run whose PID another process has
status: done
parent: S-0170
owner: alex
created: 2026-09-30T01:05:21Z
updated: 2026-09-30T01:05:46Z
transitions:
  - to: ready
    at: 2026-09-30T01:05:27Z
    by: agent-S-0170
  - to: in-progress
    at: 2026-09-30T01:05:27Z
    by: agent-S-0170
  - to: done
    at: 2026-09-30T01:05:46Z
    by: agent-S-0170
stream: S-0170
tags: []
touches: [flai/internal/serve]
usage:
  source: log
  seconds: 19
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 9
      output: 3263
      cache_read: 803888
      cache_write: 9823
      cost: 0.3047
---

# T-0605 Tell an agent's process by its start in clock ticks, and settle a run whose PID another process has

## Work

- `Owns` compared the process's start, worked out from `/proc/stat`'s `btime`, with the run's wall-clock start, give or take ten minutes. On WSL2 `btime` moves when the host sleeps, so a real agent could fail the test and not be stopped. Record the process's start in clock ticks since boot in the run (`AgentRun.Start`) when it starts, and compare that exactly.
- Use the same test (`AgentRun.running`) in `settleOrphans` and in the refusals of start, restart, and commit. A run whose PID another process has after a reboot then settles at the next look, rather than reading as working until someone stops it: the story's goal names machine restarts.
- Correct ADR-0058 (not yet on main), the design, the operator guide, and the command help.

## Done when

- The flai behavior and integration tiers and golangci-lint pass, with a test that a run whose PID another session leader has is settled without signalling it.

## Notes
