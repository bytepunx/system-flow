---
id: T-2541
type: task
nature: research
title: Measure how long a per-server timeout in the MCP file holds a call in the agent container, and record the figure in an ADR refining ADR-0244
status: done
parent: S-0583
owner: alex
created: 2026-11-07T07:06:10Z
updated: 2026-11-07T07:26:32Z
transitions:
  - to: ready
    at: 2026-11-07T07:12:28Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T07:12:29Z
    by: system-flow
  - to: done
    at: 2026-11-07T07:26:32Z
    by: agent-S-0583
stream: S-0583
tags: [cli, docker]
touches: [design/adrs/README.md, design/adrs/0258-the-mcp-file-flai-writes-for-its-http-server-names-a-per-server-timeout-of.md]
usage:
  source: log
  seconds: 843
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 60
      output: 15016
      cache_read: 7751405
      cache_write: 115093
      cost: 2.7281
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 18267
          cache_write: 781
          cost: 0.0047
---
# T-2541 Measure how long a per-server timeout in the MCP file holds a call in the agent container, and record the figure in an ADR refining ADR-0244

## Work

S-0519 T-2227 measured a per-server `"timeout": 600000` holding 400 s on the host and left the container untested (`wip/archive/agents/S-0519.md`). `MCP_TOOL_TIMEOUT` meets a second idle limit at 300 s, so this task settles whether the per-server `timeout` does too.

Measure as T-2227 did, in `flai-agent-claude-code:local`, with a scratch MCP server under the temporary folder that answers in JSON response mode, as `flai/internal/mcphttp` answers. Do not use flai's own `wait_for_events` or `permission_prompt`: the first returns on any board event and the second opens real threads. Keep the scratch server and its logs out of the repository.

1. Control: with no `timeout` in the MCP file, a call to a tool that answers after 90 s fails at 60 s.
2. With the `timeout` this task proposes, a tool that answers after 660 s returns its answer. That passes both the 300 s idle limit and the 600 s of the earlier figure.
3. Record the Claude Code version the image pins, and each measured duration with its error text.

Choose the figure. Start from 1800000 ms, 30 minutes, the HTTP server's default `--idle`, and lower it only if the measurement or the slowest real call gives a reason. `flai test` on a Go package with many importers took 3m59.7s in I-0228's third instance, and the figure must clear `flai test --all` too.

Write one ADR with `flai adr new`, refining ADR-0244 and citing ADR-0215 and D-0417. Decide that the MCP file flai writes for an HTTP server names this per-server `timeout`, as a constant in flai and not a setting, and that it lifts the limit for calls that do work: `task_done`, `story_start`, `test`, and `prime`. The 45-second holds of `permission_prompt`, `wait_for_events`, and `wait_for_work` stay, and they hold whatever a client's configuration says. A hand-written MCP file without the `timeout` stays at 60 s. Say that a newer Claude Code needs the figure measured again.

If the per-server `timeout` does not hold in the container, stop: write no ADR, record the figures on a thread on the story, and ask which way to go, since SSE answers carrying `keepAlive`'s progress are the alternative ADR-0244 left open.

## Done when

- The control, the held call, the Claude Code version, and the error text are recorded in the ADR's Context.
- The ADR is accepted, names the figure in milliseconds, and refines ADR-0244.
- `design/adrs/README.md` lists it, and `flai check --strict` is clean for the story.

## Notes
