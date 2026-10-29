---
id: T-2228
type: task
nature: remediation
title: "An ADR refining ADR-0124 gives permission_prompt a bound per transport: four minutes over stdio and, over HTTP, the bound the measurement settles below Claude Code's tool-call timeout"
status: done
parent: S-0519
owner: alex
created: 2026-10-28T13:44:51Z
updated: 2026-10-29T02:00:22Z
transitions:
  - to: ready
    at: 2026-10-29T01:59:46Z
    by: agent-S-0519
  - to: in-progress
    at: 2026-10-29T01:59:48Z
    by: agent-S-0519
  - to: done
    at: 2026-10-29T02:00:22Z
    by: agent-S-0519
stream: S-0519
tags: [cli, mcp]
touches: [design/adrs/README.md, design/adrs/0215-permission-prompt-holds-a-write-at-most-four-minutes-over-stdio-and-at-most-45.md]
after: [T-2227]
usage:
  source: log
  seconds: 34
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 3115
      cache_read: 682165
      cache_write: 27782
      cost: 0.3631
---
# T-2228 An ADR refining ADR-0124 gives permission_prompt a bound per transport: four minutes over stdio and, over HTTP, the bound the measurement settles below Claude Code's tool-call timeout

## Work

Write with `flai adr new` an ADR that refines [ADR-0124](../../../design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md), whose decision 1 says the four minutes sit "below Claude Code's idle timeout on both transports", and [ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md), which says each harness has its own bound. Context: I-0217, the figures the measurement task logged (the timeout over HTTP, on the host and in the container, with `MCP_TOOL_TIMEOUT` unset and set, and the Claude Code versions), and that the HTTP server answers as JSON so a progress notification reaches the client only when the call ends. Decision: the bound is per transport, four minutes over stdio as today and over HTTP the figure the operator confirmed on the thread, zero meaning the refusal at once naming the thread; the thread stays open and the same write made again takes the answer, as ADR-0124 has it; the refusal over HTTP says the same. Alternatives: `MCP_TOOL_TIMEOUT` in the container's environment, and SSE answers with progress. Accept it with `flai adr accept`, since the operator confirmed the remedy on the thread, and set its topics with `flai adr topics` (`mcp`, `harness`). It waits for the measurement task, since it records what that found and what the operator chose.

## Done when

The ADR is accepted, carries `refines: [ADR-0124, ADR-0130]`, has its topics, is indexed in `design/adrs/README.md`, and `flai check --strict` passes.

## Notes

The ADR's file is named by `flai adr new`; the story's touches keep the `design/adrs` folder for it.
