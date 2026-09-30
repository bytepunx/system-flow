---
id: S-0170
type: story
nature: improvement
title: The Operator can stop agents
status: done
parent: E-0013
owner: alex
created: 2026-09-29T23:56:11Z
updated: 2026-09-30T01:11:30Z
transitions:
  - to: ready
    at: 2026-09-30T00:14:24Z
    by: alex
  - to: in-progress
    at: 2026-09-30T00:46:33Z
    by: agent-S-0170
  - to: review
    at: 2026-09-30T01:05:50Z
    by: agent-S-0170
  - to: done
    at: 2026-09-30T01:11:30Z
    by: alex
tags: [dashboard, cli]
topics: [server-side, client-side]
touches: [flaiover/src, flai/cmd, flai/internal/serve, flai/internal/hostapi, docs, design]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1185
  models:
    - model: claude-opus-5-5
      input: 240
      output: 87338
      cache_read: 21515523
      cache_write: 262893
      cost: 8.154
---
# S-0170 The Operator can stop agents

## Goal

Right now, a stuck agent will remain in the activity page forever since it is technically just going to sit there and persist between restarts and even machine restarts.

The operator should be able to stop an agent from the activity page (with confirmation to protect against accidental clicks).

## Acceptance criteria
- [x] The operator can click stop on any agent running in the activity page
- [x] When clicking stop, the operator must confirm their intention (should include a warning explaining the impact)
- [x] When the operator confirms their intent is to stop the agent, the agent should be terminated

## Tasks
- T-0602 flai serve agent stop ends a story's running agent, and agent.stop asks for it
- T-0603 The activity page stops an agent after the operator confirms it
- T-0604 Record the stop in an ADR, the design, and the docs
- T-0605 Tell an agent's process by its start in clock ticks, and settle a run whose PID another process has

## Notes

- Criterion 1: each activity card whose agent runs, or ended waiting for an answer, has **Stop** while the `agent` host action is on (`ActivityView.svelte`, `stoppable` in `flaiover/src/lib/activity.ts`). Verified by the `ActivityView stopping an agent (S-0170)` component tests and the `stopping an agent` tests in `activity.test.ts`. It was not clicked in a running dashboard.
- Criterion 2: Stop opens `AgentStopConfirm.svelte`. Its warning says the process and everything it started end now and cannot be resumed, that the worktree keeps what the agent left, and that the story gets no agent until Retry or a move back to ready. Nothing is posted until **Stop the agent** is pressed. Verified by the same component tests.
- Criterion 3: confirming posts `{ action: 'stop' }` to `/api/items/:id/agent`, flai's `agent.stop`, which runs `flai serve agent stop`. That command sends SIGTERM to the agent's process group, then SIGKILL after 10 s, and records the run `stopped` (`serve.Stop`, ADR-0058). Verified against real processes by `TestTheOperatorStopsAStorysAgent` (an agent ended, one killed after the grace period, a reused PID left alone, an asking agent not resumed, the refusals), by `TestServeAgentStopEndsAStorysAgent`, and by the route and host-method tests.
