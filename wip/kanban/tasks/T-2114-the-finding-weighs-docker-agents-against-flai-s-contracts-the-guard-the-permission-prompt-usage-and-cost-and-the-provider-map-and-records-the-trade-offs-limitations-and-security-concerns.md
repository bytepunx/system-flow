---
id: T-2114
type: task
nature: research
title: The finding weighs Docker Agents against flai's contracts, the guard, the permission prompt, usage and cost, and the provider map, and records the trade-offs, limitations, and security concerns
status: cancelled
parent: S-0496
owner: alex
created: 2026-10-25T21:05:13Z
updated: 2026-10-25T21:09:49Z
transitions:
  - to: cancelled
    at: 2026-10-25T21:09:49Z
    by: system-flow
stream: S-0496
tags: [cli]
touches: [design/system/docker-agents.md]
after: [T-2113]
---
# T-2114 The finding weighs Docker Agents against flai's contracts, the guard, the permission prompt, usage and cost, and the provider map, and records the trade-offs, limitations, and security concerns

## Work

Add to `design/system/docker-agents.md`, after the sections T-2113 wrote, the sections that judge Docker Agents under the concerns `agent-adapters.md` § Harnesses against the concerns judges every harness under, each with sources and dated:

- The guard: whether the runtime gives a hook or plugin before a tool call that `flai guard` could answer ([ADR-0130](../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md)), and what [ADR-0131](../../design/adrs/0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md) means for it when it gives none: a story's agent but no roles.
- The permission prompt: whether a protected write can be held and asked of the story's owner, or only refused under auto-approve.
- Sub-agents and strategic agents: whether the explorer, the verifier, the planner, and the others can be defined, each on its own model, as [ADR-0065](../../design/adrs/0065-a-story-s-agent-carries-a-model-per-sub-agent-role-and-claude-code-runs-each.md) wants.
- Usage and cost: what the runtime reports per call and per model, and how a story's cost would be priced for a local model, which has no spend log ([ADR-0132](../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md)).
- The provider map: whether a `providers` entry on the host ([ADR-0129](../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)) can name an Ollama or Model Runner endpoint and an Anthropic one at once, and what a story's `agent` would carry.

Then write three closing sections the ADR will cite: the trade-offs and caveats of running a story's agent in a Docker-managed runtime against Claude Code on the host (the worktree reached from a container, git and `flai` inside it, the image to keep, Docker Desktop against Docker Engine on this WSL2 host); the limitations (what the local models can and cannot do on a story of this repository, the subscription question T-2113 settled, resumption, the roles); and the security concerns (the keys the container is given and how, what of the repository and the host it can reach, the MCP token, images pulled by tag, and what the sandbox does and does not isolate).

Waits for T-2113, which writes the first half of the same file, so that the two do not change it at once.

## Done when

`design/system/docker-agents.md` has the five concern sections and the three closing sections, each statement sourced and dated, with a recommendation line saying whether Docker Agents is worth an adapter and under what conditions, and `flai check --strict --story S-0496` reports nothing on it.

## Notes

Where a question can only be settled by running the runtime, say so and record what to run; this story builds nothing and runs nothing on a story.
- 2026-10-25T21:09:49Z: moved to cancelled: Planned in parallel by planner-S-0496 while agent-S-0496 wrote T-2116 and T-2117, which cover the same work: the survey goes into the existing harness survey design/system/agent-adapters.md rather than a new finding document, and the draft ADR is T-2117
