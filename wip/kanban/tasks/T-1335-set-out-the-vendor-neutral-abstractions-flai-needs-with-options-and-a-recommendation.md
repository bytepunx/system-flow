---
id: T-1335
type: task
nature: research
title: Set out the vendor-neutral abstractions flai needs, with options and a recommendation
status: backlog
parent: S-0339
owner: alex
created: 2026-10-08T07:20:49Z
updated: 2026-10-08T07:20:49Z
transitions: []
stream: S-0339
tags: [research, agents]
touches: [design/system/agent-adapters.md]
after: [T-1334]
---
# T-1335 Set out the vendor-neutral abstractions flai needs, with options and a recommendation

## Work

Add an `## Abstractions` section to `design/system/agent-adapters.md`. It states what flai needs from any agent, apart from Anthropic's SDK API and OpenAI's. First, separate three things a story's `agent` now runs together:

- the harness: the program that runs the agent loop, its tools, and its sub-agents
- the provider or gateway the model calls go to: Anthropic, a LiteLLM proxy, or OpenRouter
- the model's name on that provider.

Then give one subsection per concern in `## Today`. In each, state the contract flai needs, the options for meeting it across the harnesses in `## LiteLLM and OpenRouter`, and what a harness that cannot meet it loses. The concerns are:

- starting and resuming an agent
- refusing a sub-agent's write, now `flai guard`
- asking the owner before a protected write, now `permission_prompt`
- a model per sub-agent role
- the sub-agent and strategic agent definitions, now `.claude/agents/`, and the instructions file, now `CLAUDE.md`
- usage and cost: a reader per harness or per gateway, and where a price comes from when the log gives none.

Also cover the manifest, story, and host settings each option adds, such as a base URL, a provider, and the name of a key's environment variable. Never store the key itself. Then compare the options in a table, and end with a recommendation: which adapters the epic builds, in what order, and which current ADRs each option would supersede or refine. Waits for T-1334 because it builds on what the vendors and harnesses offer.

## Done when

- The section has the harness, provider, and model split.
- It has one subsection per concern, a comparison table, and a recommendation with its reasons.
- It names each ADR the recommendation would supersede or refine.

## Notes
