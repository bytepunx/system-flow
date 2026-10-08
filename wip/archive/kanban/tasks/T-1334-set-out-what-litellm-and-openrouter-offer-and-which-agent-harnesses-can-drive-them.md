---
id: T-1334
type: task
nature: research
title: Set out what LiteLLM and OpenRouter offer and which agent harnesses can drive them
status: done
parent: S-0339
owner: alex
created: 2026-10-08T07:20:39Z
updated: 2026-10-08T07:57:29Z
transitions:
  - to: ready
    at: 2026-10-08T07:48:44Z
    by: agent-S-0339
  - to: in-progress
    at: 2026-10-08T07:48:44Z
    by: agent-S-0339
  - to: done
    at: 2026-10-08T07:57:29Z
    by: agent-S-0339
stream: S-0339
tags: [research, agents]
touches: [design/system/agent-adapters.md]
after: [T-1333]
usage:
  source: log
  seconds: 525
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 296
      output: 16276
      cache_read: 1568563
      cache_write: 88232
      cost: 2.5867
---
# T-1334 Set out what LiteLLM and OpenRouter offer and which agent harnesses can drive them

## Work

Add a `## LiteLLM and OpenRouter` section to `design/system/agent-adapters.md`. Base it on each vendor's current documentation, and cite each page by URL. For each of the two, state:

- whether it is self-hosted, like the LiteLLM proxy, or hosted, like OpenRouter
- the APIs it serves: OpenAI chat completions, the OpenAI Responses API, an Anthropic-compatible messages endpoint, and streaming
- how it names models and routes between providers
- how tool calls and MCP reach a model through it
- how it authenticates, and which keys it needs
- what usage and cost each response reports, and where a caller reads them.

Then set out the agent harnesses that could run a story's agent over each. Include Claude Code pointed at the gateway by its base URL and model settings, OpenAI's Codex CLI, and other headless coding agents that accept an OpenAI-compatible endpoint, such as OpenCode, Goose, and aider. Include an agent loop of flai's own. For each harness, state:

- whether it runs headless and can resume a session
- whether it loads MCP servers, including flai's
- whether it can run sub-agents with their own model
- whether it has a hook to refuse a tool call, which `flai guard` needs, and a permission handler, which `permission_prompt` needs
- what usage it logs.

State plainly what the documentation does not settle. Make no paid call to either service without the operator's word on the story's thread. Waits for T-1333 because it edits the same file and compares against `## Today`.

## Done when

- The section has one subsection per vendor and one per harness.
- It has a table of harnesses against the concerns in `## Today`.
- Every claim about a vendor or harness cites its source.

## Notes
