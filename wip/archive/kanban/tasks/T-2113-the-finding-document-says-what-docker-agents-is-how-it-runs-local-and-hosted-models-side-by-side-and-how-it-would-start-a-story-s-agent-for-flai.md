---
id: T-2113
type: task
nature: research
title: The finding document says what Docker Agents is, how it runs local and hosted models side by side, and how it would start a story's agent for flai
status: cancelled
parent: S-0496
owner: alex
created: 2026-10-25T21:04:50Z
updated: 2026-10-25T21:09:48Z
transitions:
  - to: cancelled
    at: 2026-10-25T21:09:48Z
    by: system-flow
stream: S-0496
tags: [cli]
touches: [design/system/docker-agents.md, design/system/README.md]
---
# T-2113 The finding document says what Docker Agents is, how it runs local and hosted models side by side, and how it would start a story's agent for flai

## Work

Write `design/system/docker-agents.md`, the finding of S-0496, in the shape of `agent-adapters.md` (S-0339): front matter with `title`, `updated`, `status`, and `topics: [cli]`, a first paragraph saying what was asked and that nothing is built, and every statement about Docker dated and linked to its source. Add its row to `design/system/README.md` beside `agent-adapters.md`.

Settle first what "Docker Agents" names, from Docker's documentation: Docker's agent runtime (`cagent`, surfaced as `docker agent`), Docker Model Runner (`docker model`, the OpenAI-compatible endpoint for local models), Docker Sandboxes for coding agents, and the Docker MCP Gateway, and say which of them the goal's mixture needs. If the documentation leaves it open which product the operator means, ask on a thread on S-0496 and go on with the rest.

Then answer, with sources:

- How an agent definition names its model and provider, and whether one definition, or one run, can route to a hosted Anthropic model and to a local model (Ollama, or Docker Model Runner) by role or by step, which is the mixture the goal asks for.
- How a hosted Anthropic model is reached: an API key billed by usage, or a Claude subscription. Say plainly whether a subscription is usable outside Claude Code, since the goal names one.
- How Ollama is reached: as an OpenAI-compatible provider by base URL, through Docker Model Runner, or not at all.
- How flai would start it on a story: the command, the working directory, how the story's prompt and `flai mcp` reach it (MCP over stdio from inside a container, or over HTTP with the token, [ADR-0030](../../design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)), how it ends, and whether a session can be resumed, read against the harness contract of [ADR-0130](../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md) and the provider map of [ADR-0129](../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md).

Leave the guard, the permission prompt, usage and cost, and the security concerns to the next task, which writes the following sections of the same document. Waits for nothing: it is the first layer.

## Done when

`design/system/docker-agents.md` exists with the sections above, each statement about Docker dated and sourced, `design/system/README.md` lists it, and `flai check --strict --story S-0496` reports nothing on either file.

## Notes

The operator's goal names Ollama and a Claude subscription; the document says what each of them can and cannot be under Docker Agents before the next task weighs it.
- 2026-10-25T21:09:48Z: moved to cancelled: Planned in parallel by planner-S-0496 while agent-S-0496 wrote T-2116 and T-2117, which cover the same work: the survey goes into the existing harness survey design/system/agent-adapters.md rather than a new finding document, and the draft ADR is T-2117
