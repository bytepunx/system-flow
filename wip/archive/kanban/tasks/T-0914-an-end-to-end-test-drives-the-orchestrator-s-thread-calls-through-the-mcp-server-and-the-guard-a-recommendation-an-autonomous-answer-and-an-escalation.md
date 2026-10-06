---
id: T-0914
type: task
nature: feature
title: "An end-to-end test drives the orchestrator's thread calls through the MCP server and the guard: a recommendation, an autonomous answer, and an escalation"
status: done
parent: S-0220
owner: alex
created: 2026-10-05T04:48:41Z
updated: 2026-10-06T06:47:44Z
transitions:
  - to: ready
    at: 2026-10-06T06:41:14Z
    by: agent-S-0220
  - to: in-progress
    at: 2026-10-06T06:41:14Z
    by: agent-S-0220
  - to: done
    at: 2026-10-06T06:47:44Z
    by: agent-S-0220
stream: S-0220
tags: [flai]
touches: [flai/internal/mcpserver/orchestrator_threads_test.go]
after: [T-0906, T-0908, T-0910, T-0912]
usage:
  source: log
  seconds: 390
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 85
      output: 27935
      cache_read: 3431287
      cache_write: 174842
      cost: 2.2121
---
# T-0914 An end-to-end test drives the orchestrator's thread calls through the MCP server and the guard: a recommendation, an autonomous answer, and an escalation

## Work

The story's criteria ask for tests of a recommendation, an autonomous answer, and an escalation. Each earlier task tests its own part; this one drives the whole path against a fixture repository with a story's agent waiting on a thread, calling as the orchestrator does, in role `orchestrate`, through the guard's decision and the MCP server's tools:

- Recommendation, with `answer_threads: recommend`: the orchestrator posts a recommendation with a source. The thread stays awaiting the operator, the asking agent is not woken, the decision log names the source, and an answer without `recommendation` is refused. The operator's `flai thread confirm` then sets `answered`, and `flai stats` counts the wait as confirmed.
- Autonomous answer, with `autonomous`: an answer with a source sets `answered`, wakes the asking agent, is logged with its source, and is counted under `orchestrator` in `flai stats`.
- Escalation, with `autonomous`: an answer without a source is refused, and the same reply posted as a recommendation leaves the thread awaiting the operator. A resolve of the thread is refused, and so is an answer on a thread the orchestrator opened.

Put the test in `flai/internal/mcpserver/orchestrator_threads_test.go`, a file of its own, so that no other task's paths are shared.

It waits for T-0906, T-0908, T-0910, and T-0912, whose parts it drives together. It runs with T-0913, whose paths it does not share.

## Done when

- the three scenarios pass, each asserting the thread's status, who is awaited, the decision log entry, and the wait's count in `flai stats`
- `go test ./internal/mcpserver/` passes, and `scripts/flai-test.sh` passes

## Notes
