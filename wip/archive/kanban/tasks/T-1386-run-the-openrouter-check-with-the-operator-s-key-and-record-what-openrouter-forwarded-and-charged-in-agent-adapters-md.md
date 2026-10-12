---
id: T-1386
type: task
nature: research
title: Run the OpenRouter check with the operator's key and record what OpenRouter forwarded and charged in agent-adapters.md
status: done
parent: S-0351
owner: alex
created: 2026-10-08T08:45:16Z
updated: 2026-10-12T20:29:51Z
transitions:
  - to: ready
    at: 2026-10-12T20:27:20Z
    by: agent-S-0351
  - to: in-progress
    at: 2026-10-12T20:27:21Z
    by: agent-S-0351
  - to: done
    at: 2026-10-12T20:29:51Z
    by: agent-S-0351
stream: S-0351
tags: [cli]
touches: [design/system/agent-adapters.md]
after: [T-1384]
usage:
  source: log
  seconds: 150
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 28
      output: 13364
      cache_read: 2133543
      cache_write: 67828
      cost: 1.1484
---
# T-1386 Run the OpenRouter check with the operator's key and record what OpenRouter forwarded and charged in agent-adapters.md

## Work

- Run `TestGatewayOpenRouter` with the operator's `OPENROUTER_API_KEY` and a budget cap. Run it more than once only if the first run leaves a question open, and within the key's limit.
- In `design/system/agent-adapters.md`, add `### Checked against OpenRouter` under `## LiteLLM and OpenRouter`, with the date, the `claude` version, and the model slugs. Record:
  - whether the beta headers and body fields Claude Code sent reached the model, judged by the features that depend on them (tool use, MCP tools, prompt caching tokens in `usage`, the stream's events);
  - whether `haiku` resolved through `ANTHROPIC_DEFAULT_HAIKU_MODEL`;
  - whether the stream's message IDs are OpenRouter's generation IDs;
  - the gateway's cost for the run against `total_cost_usd`, and the cache tokens each reports.
- Change each `**to check**` the run settles to what was found, and say what stays open for LiteLLM.

It waits for T-1384, whose test it runs.

## Done when

- The section is in `agent-adapters.md` with figures that were observed, never estimated.
- `flai test design/system/agent-adapters.md` passes.

## Notes

Layer 2 of S-0351. Without the key in the agent's environment, ask the operator on a thread on S-0351 and wait.
