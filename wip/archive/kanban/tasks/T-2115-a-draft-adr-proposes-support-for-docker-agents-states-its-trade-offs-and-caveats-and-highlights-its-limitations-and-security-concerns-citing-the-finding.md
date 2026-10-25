---
id: T-2115
type: task
nature: research
title: A draft ADR proposes support for Docker Agents, states its trade-offs and caveats, and highlights its limitations and security concerns, citing the finding
status: cancelled
parent: S-0496
owner: alex
created: 2026-10-25T21:05:33Z
updated: 2026-10-25T21:09:50Z
transitions:
  - to: cancelled
    at: 2026-10-25T21:09:50Z
    by: system-flow
stream: S-0496
tags: [cli]
touches: [design/adrs/README.md]
after: [T-2114]
---
# T-2115 A draft ADR proposes support for Docker Agents, states its trade-offs and caveats, and highlights its limitations and security concerns, citing the finding

## Work

Write the draft ADR with `flai adr new` in the story's worktree, with the topic `cli`, and leave it `proposed`: the story asks for a draft, so do not run `flai adr accept`. Its Context says what the operator asked, the mixture of a hosted Anthropic model and local Ollama models, and cites `design/system/docker-agents.md` for every fact. Its Decision proposes, or declines to propose, a `docker-agents` harness beside `claude-code` under the contracts of [ADR-0130](../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md): what a story's `agent` and the host's `providers` would carry, how flai would start and read it, and under which conditions, as the finding's recommendation line says. Its Consequences carry three sections the criteria name, in these words so each criterion can be ticked against one: `Trade-offs and caveats`, `Limitations`, and `Security concerns`, each a list drawn from the finding's closing sections. Its Alternatives considered name at least the LiteLLM and OpenRouter path already decided in `agent-adapters.md` § Decision and a plain OpenAI-compatible provider entry pointing at Ollama with no Docker runtime.

Then link the ADR from the finding document's end, as `agent-adapters.md` § Decision links its decision, so the two read as one.

Waits for T-2114: the ADR cites the finding's sections and takes its recommendation. The ADR's file is created by `flai adr new`, so its path is not in the touches; `flai task done` adds it when the task closes (ADR-0144). S-0495, in review, carries `0201`: if it reaches main before this task runs, the new ADR is `0202`; otherwise run `flai stream sync` after S-0495 is accepted and renumber with `flai adr` before review, so the two do not collide.

## Done when

The ADR exists under `design/adrs/` with status `proposed` and topic `cli`, has its row in `design/adrs/README.md`, carries the sections `Trade-offs and caveats`, `Limitations`, and `Security concerns`, cites `design/system/docker-agents.md`, is linked from that document, and `flai check --strict --story S-0496` reports nothing on it.

## Notes

The operator decides on the ADR after the story: it stays proposed, and the story's close-out ticks the three criteria against it.
- 2026-10-25T21:09:50Z: moved to cancelled: Planned in parallel by planner-S-0496 while agent-S-0496 wrote T-2116 and T-2117, which cover the same work: the survey goes into the existing harness survey design/system/agent-adapters.md rather than a new finding document, and the draft ADR is T-2117
