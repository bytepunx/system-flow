---
id: ADR-0244
title: "Over HTTP, one wait_for_events or wait_for_work call holds at most the 45 seconds permission_prompt holds there, and then answers timed_out for the agent to call again"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0215]
topics: [mcp, harness]
---

# ADR-0244 Over HTTP, one wait_for_events or wait_for_work call holds at most the 45 seconds permission_prompt holds there, and then answers timed_out for the agent to call again

## Context

[ADR-0215](0215-permission-prompt-holds-a-write-at-most-four-minutes-over-stdio-and-at-most-45.md) bounds `permission_prompt`'s hold over HTTP at 45 seconds, 15 below the 60 after which Claude Code drops a call to flai's HTTP MCP server. That server answers in JSON response mode, so `keepAlive`'s progress notifications reach the client only when the call ends, and nothing a held call sends keeps it alive. ADR-0215 leaves every other long call alone, and [I-0228](../issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md) records them: `wait_for_events` (60 seconds by default, up to 30 minutes), `wait_for_work` (five minutes by default), `verify`, and a `task_done` or `story_start` slower than 60 seconds. Over HTTP each fails with "The operation timed out." while the server goes on with it. So an idle agent in a container cannot hold either wait, and S-0506's first criterion asks `wait_for_events` to answer over HTTP as it does on stdio.

Both waits already take a bound, `Options.MaxWait`, which is `LongestWait`, 30 minutes, unless set ([I-0059](../issues/I-0059-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md)). A wait whose time passes answers `timed_out`, and both tools' descriptions tell the agent to call again then.

Since S-0546, `verify` asks `flai serve` for the close-out as a job and answers at once ([ADR-0229](0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)); the agent waits for the job's `verified` event with `wait_for_events`.

## Decision

**Over HTTP, one `wait_for_events` or `wait_for_work` call holds at most the 45 seconds `permission_prompt` holds there, and then answers `timed_out` for the agent to call again.** This refines ADR-0215 by giving these two calls its HTTP bound.

1. The bound is `mcpserver.HTTPLongestWait`, defined as `HTTPPermissionHold`: one figure for every hold over HTTP, measured against Claude Code's 60 seconds once. The server `flai mcp http` runs sets it as its `MaxWait`. The server `flai mcp` runs over stdio keeps `LongestWait`.
2. A `timeout_seconds` above the bound, or none, is held for the bound, and one below it as asked. A change within the bound answers within a second, as on stdio.
3. Each tool's description states its default and the most it holds on the server that offers it, as `permission_prompt`'s states its hold. An agent over HTTP reads 45 seconds, and one over stdio 30 minutes. The `timeout_seconds` field defers to the description.
4. What stays open: `verify` needs no bound, since it answers at once and its result comes as an event within this bound. A `task_done` or `story_start` slower than 60 seconds over HTTP still meets Claude Code's limit while the server finishes it. Neither waits by design, and cutting them short would leave their work half done. A per-server `timeout` in the MCP file the container adapter writes (S-0508) may lift Claude Code's limit for those; it does not replace this bound, which holds whatever the client's configuration says.

## Consequences

- An idle agent in a container that holds `wait_for_work` or `wait_for_events` gets `timed_out` every 45 seconds and calls again, instead of an error.
- Each return is a model turn that reads the agent's context again (I-0059), so an agent idle over HTTP takes up to 40 turns in the half hour one over stdio holds in one. That is the price of an answer at all, until an answer can hold longer over HTTP.
- I-0228 is fixed for the two waits and for `verify`, and stays open for a slow `task_done` or `story_start`.
- A newer Claude Code that changes the 60 seconds needs the figure measured again, once, for every hold over HTTP.

## Alternatives considered

- **SSE answers that carry `keepAlive`'s progress.** They might let a wait hold longer, but they change the HTTP server's response mode, and how Claude Code counts progress over SSE was not measured. They can follow, and raise the bound.
- **A per-server `timeout` in the MCP file.** It rests on each client's configuration rather than flai, and a hand-written configuration stays at 60 seconds, as ADR-0215 found.
- **A figure of its own for the waits, such as 55 seconds.** It saves a few turns but narrows the margin, and two figures need two measurements.
- **Refusing the waits over HTTP.** An agent in a container could then never wait on the host's events.
