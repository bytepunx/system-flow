---
id: T-0602
type: task
nature: feature
title: flai serve agent stop ends a story's running agent, and agent.stop asks for it
status: done
parent: S-0170
owner: alex
created: 2026-09-30T00:48:33Z
updated: 2026-09-30T00:55:02Z
transitions:
  - to: ready
    at: 2026-09-30T00:48:49Z
    by: agent-S-0170
  - to: in-progress
    at: 2026-09-30T00:48:49Z
    by: agent-S-0170
  - to: done
    at: 2026-09-30T00:55:02Z
    by: agent-S-0170
stream: S-0170
tags: []
touches: [flai]
usage:
  source: log
  seconds: 373
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 62
      output: 22634
      cache_read: 5575899
      cache_write: 68131
      cost: 2.1132
---

# T-0602 flai serve agent stop ends a story's running agent, and agent.stop asks for it

## Work

- `serve.Stop` in `flai/internal/serve`: for a story whose newest run is live, mark the run stopped in `serve/agents.json`, send SIGTERM to its process group, then SIGKILL after a grace period, and record it ended with the outcome `stopped`. Signal only while the PID is still the agent flai started (a session leader that started when the run did), never a PID reused since. A run whose process is gone is only recorded as stopped. A run that ended waiting for an answer is recorded as stopped, so the answer does not start it again. Refuse, saying why, a story with no run or whose agent is not running.
- The launcher keeps the stop when it settles the run itself (its wait on the process, or `settleOrphans`), rather than judging it failed over it.
- `flai serve agent stop <story-id>` and the host method `agent.stop`, gated by the `agent` action, journalled.
- Behavior tests for each case.

## Done when

- `go test -race ./...` in `flai/` and golangci-lint pass, and the tests cover a running agent stopped, a reused PID left alone, an asked run stopped, and the refusals.

## Notes
