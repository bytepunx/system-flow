---
id: ADR-0258
title: "The MCP file flai writes for its HTTP server names a per-server timeout of 1800000 milliseconds, so a call that does work runs past Claude Code's 60 seconds, while the 45-second holds stay"
status: accepted
date: 2026-11-07
supersedes: []
superseded_by: []
refines: [ADR-0244]
topics: [mcp, harness]
---

# ADR-0258 The MCP file flai writes for its HTTP server names a per-server timeout of 1800000 milliseconds, so a call that does work runs past Claude Code's 60 seconds, while the 45-second holds stay

## Context

Claude Code drops a call to flai's HTTP MCP server that has not answered after 60 seconds. That server answers in JSON response mode, so nothing it sends while a call runs reaches the client ([ADR-0215](0215-permission-prompt-holds-a-write-at-most-four-minutes-over-stdio-and-at-most-45.md)). [ADR-0244](0244-over-http-one-wait-for-events-or-wait-for-work-call-holds-at-most-the-45.md) bounds the two waits at 45 seconds over HTTP and leaves open, in its point 4, a call that does work and takes longer: a `task_done` or `story_start` slower than 60 seconds. [I-0228](../issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md)'s third instance is the MCP tool `test` on a Go package with many importers, which took 3m59.7s.

S-0519 measured a per-server `"timeout": 600000` in the MCP file holding a call 400 s on the host, and did not try it in the container. It also found that `MCP_TOOL_TIMEOUT` alone meets a second limit, 300 s with no response or progress (D-0417, ADR-0215).

S-0583 measured the per-server `timeout` in `flai-agent-claude-code:local`, which pins Claude Code 2.1.312, with a scratch MCP server on the host. Like `flai/internal/mcphttp`, it answers in JSON response mode, and its one tool answers after the seconds it is given. Each run was one headless `claude -p` with `--strict-mcp-config`, calling the tool once. Each figure is the time from the call reaching the server to Claude Code's answer:

| Run | `timeout` in the MCP file | Tool answers after | Claude Code |
|-----|---------------------------|--------------------|-------------|
| Control | none | 90 s | failed after 60.0 s: `The operation timed out.` |
| Held | 1800000 | 660 s | returned the answer, `slow answered after 660 s`, after 660.0 s |

Both runs were on 2026-11-07, with the login the operator set for container runs mounted read-only and `--model haiku`. The control's call reached the server at 07:13:30.801Z and failed at 07:14:30.809Z. The held call reached it at 07:14:46.377Z; the server answered at 07:25:46.320Z, and Claude Code recorded the result at 07:25:46.332Z.

The held call passed both the 300 s idle limit `MCP_TOOL_TIMEOUT` met and the 400 s S-0519 measured on the host. `flai mcp http`'s own session idle timeout, `--idle`, is 30 minutes by default, and the server pauses it while a request runs.

## Decision

**The MCP file flai writes for its HTTP server names a per-server `timeout` of 1800000 milliseconds, 30 minutes, so that a call that does work runs past Claude Code's 60 seconds.** This refines ADR-0244's point 4.

1. The figure is one constant in flai, `mcphttp.ClientTimeout`, not a setting. The file `httpMCPConfig` writes for an agent in a container names it as the `timeout` of the server named `flai`, and so does the sample `flai mcp status` prints for a hand-written client.
2. It lifts the limit for the calls that do work: `task_done`, `story_start`, `test`, and `prime`. A call that runs longer than 30 minutes still fails, as Claude Code drops it.
3. The 45-second holds of `permission_prompt`, `wait_for_events`, and `wait_for_work` stay. They hold whatever a client's configuration says, so a hand-written MCP file without the `timeout` still gets an answer from each wait.
4. A hand-written MCP file without the `timeout` stays at 60 seconds for the calls that do work.
5. The figure was measured to 660 s. Nothing was measured between 660 s and 30 minutes. A newer Claude Code needs it measured again, once, as ADR-0244 says of its own figure.

## Consequences

- An agent in a container can close a slow task, start a story, and run a slow `test` over HTTP, as it can over stdio. I-0228 closes.
- pilot in the agent image writes the MCP file, so an image built before this change keeps the 60-second limit until `make agent-images` builds it again.
- A call cut short at 30 minutes still leaves the server finishing its work, as a call cut at 60 seconds did. No call flai serves is expected to run that long: a close-out, the longest, is a `flai serve` job since ADR-0229.
- The two waits still return every 45 seconds over HTTP, as ADR-0244 decides.

## Alternatives considered

- **`MCP_TOOL_TIMEOUT` in the container's environment.** It stops at the 300 s idle limit, which S-0519 measured; the per-server `timeout` passed it.
- **SSE answers carrying `keepAlive`'s progress.** They would also lift the waits' bound, but they change the HTTP server's response mode, and how Claude Code counts progress over SSE was not measured. A larger change for no call that needs it now.
- **A figure below 30 minutes, such as 10.** It would cover I-0228's four-minute `test` but not `flai test --all` or a slow repository's `story_start`, and the server's own idle timeout already allows 30.
- **A setting in the manifest or host configuration.** One measured figure needs no knob; a setting the operator raises past what Claude Code holds would only fail later.
