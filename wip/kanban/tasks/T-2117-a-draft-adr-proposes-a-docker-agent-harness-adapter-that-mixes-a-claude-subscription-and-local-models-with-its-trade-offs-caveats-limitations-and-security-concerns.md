---
id: T-2117
type: task
nature: research
title: A draft ADR proposes a docker-agent harness adapter that mixes a Claude subscription and local models, with its trade-offs, caveats, limitations, and security concerns
status: done
parent: S-0496
owner: alex
created: 2026-10-25T21:07:54Z
updated: 2026-10-25T21:11:39Z
transitions:
  - to: ready
    at: 2026-10-25T21:09:34Z
    by: system-flow
  - to: in-progress
    at: 2026-10-25T21:09:35Z
    by: system-flow
  - to: done
    at: 2026-10-25T21:11:28Z
    by: agent-S-0496
stream: S-0496
tags: []
touches: [design/adrs/README.md, design/adrs/0202-flai-gains-a-docker-agent-harness-adapter-that-runs-a-story-s-agent-as-a-docker.md]
after: [T-2116]
usage:
  source: log
  seconds: 113
  estimated: true
  priced_by: harness
  models:
    - model: claude-fable-5-1
      input: 78
      output: 3484
      cache_read: 440076
      cache_write: 22124
      cost: 0.7275
---
# T-2117 A draft ADR proposes a docker-agent harness adapter that mixes a Claude subscription and local models, with its trade-offs, caveats, limitations, and security concerns

## Work

Record with `flai adr new --status proposed --commit` a draft ADR proposing that flai gain a `docker-agent` harness adapter, under ADR-0130's contracts, through which a story's agent runs as a Docker Agent team whose root agent is Docker Agent's `claude-code` harness on the operator's Claude subscription and whose sub-agents run on local models through Docker Model Runner or Ollama, or the reverse. Its Context states what Docker Agent gives against each of flai's concerns, from the subsection T-2116 writes; its Consequences carry the trade-offs and caveats (who holds the guard, two agent loops, cost figures, Anthropic's authentication terms) and the limitations and security concerns (`--dangerously-skip-permissions` in the harness, client-side permissions, telemetry on by default, plain-text key file, local models without tool calling, sandbox not isolating the harness). Waits for T-2116 so that it links the subsection by heading.

## Done when

- A proposed ADR exists under `design/adrs` with its index row, committed on the story's branch by `flai adr new --commit`
- Its Consequences name the trade-offs, caveats, limitations, and security concerns, each with its source
- The three acceptance criteria of S-0496 are met by the ADR and the design section together

## Notes

The ADR is a proposal: nothing is built by S-0496, a research story, and the operator decides at acceptance whether an epic follows.
