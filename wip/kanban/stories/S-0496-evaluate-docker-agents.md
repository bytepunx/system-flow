---
id: S-0496
type: story
nature: research
title: Evaluate Docker Agents
status: review
owner: alex
created: 2026-10-25T21:01:31Z
updated: 2026-10-25T21:12:03Z
transitions:
  - to: ready
    at: 2026-10-25T21:01:37Z
    by: alex
  - to: in-progress
    at: 2026-10-25T21:01:48Z
    by: agent-S-0496
  - to: review
    at: 2026-10-25T21:12:03Z
    by: agent-S-0496
tags: [cli]
touches: [design/adrs/README.md, design/system/agent-adapters.md, design/adrs/0202-flai-gains-a-docker-agent-harness-adapter-that-runs-a-story-s-agent-as-a-docker.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: medium
usage:
  source: log
  seconds: 633
  priced_by: harness
  turns:
    - day: 2026-10-25
      ceremony: 1
      hand_edits: 6
      work: 20
  models:
    - model: claude-fable-5-1
      input: 866
      output: 38449
      cache_read: 4856659
      cache_write: 244155
      cost: 8.0284
    - model: claude-haiku-5-5
      input: 264021
      output: 49236
      cache_read: 0
      cache_write: 0
      cost: 0.151
  strategic:
    - kind: planner
      seconds: 256
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 354
          output: 16780
          cache_read: 1211473
          cache_write: 129436
          cost: 3.7341
    - kind: orchestrator
      seconds: 774
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 94
          output: 1622
          cache_read: 9446423
          cache_write: 47624
          cost: 2.0817
        - model: claude-sonnet-5-5
          input: 8
          output: 53
          cache_read: 82546
          cache_write: 41380
          cost: 0.0848
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: alex
    at: 2026-10-25T21:01:31Z
  value: 50
  by: planner-S-0496
  at: 2026-10-25T21:04:35Z
forecast:
  duration: 35m
  delivery: 2026-10-25T21:49:00Z
  basis: "flai forecast gave 26m (median 218 s per unit over 4 done research stories in the medium band, size 7); raised to 35m, the median agent time of the three done research stories that wrote one finding document and one ADR (S-0339 2102 s, S-0380 2466 s, S-0413 1372 s), since the survey reads several Docker products; delivery is the start at 2026-10-25T21:02Z plus 35m times the cycle factor 1.34."
  by: planner-S-0496
  at: 2026-10-25T21:04:35Z
---
# S-0496 Evaluate Docker Agents

## Goal

Determine whether or not flai could leverage Docker Agents as a means to handle a mixture of local and hosted models (i.e. Anthropic Claude via subscription and local Ollama hosted models) as a way to make use of locally available models alongside hosted models.

## Acceptance criteria
- [x] Write a draft ADR proposing support for Docker Agents
- [x] Share trade-offs and caveats with this approach
- [x] Highlight any notable limitations or security concerns

## Tasks
- T-2113 The finding document says what Docker Agents is, how it runs local and hosted models side by side, and how it would start a story's agent for flai
- T-2114 The finding weighs Docker Agents against flai's contracts, the guard, the permission prompt, usage and cost, and the provider map, and records the trade-offs, limitations, and security concerns
- T-2115 A draft ADR proposes support for Docker Agents, states its trade-offs and caveats, and highlights its limitations and security concerns, citing the finding
- T-2116 The agent adapters design records the Docker Agent survey: what the harness, its providers, its hooks, its Claude Code harness, and its usage reporting give against flai's concerns
- T-2117 A draft ADR proposes a docker-agent harness adapter that mixes a Claude subscription and local models, with its trade-offs, caveats, limitations, and security concerns

## Notes

### Planning

Planned by planner-S-0496 on 2026-10-25.

Touches. `flai touches suggest S-0496` refused: the story declared no touches, so there were no seeds for co-change. Every touch comes from the layout of past research stories and from the design:

- `design/system/docker-agents.md` (layout): the finding of a research story is one document under `design/system`, as S-0339 wrote `agent-adapters.md`, S-0413 `agent-prose.md`, and S-0426 `agent-memory.md`. The ADR cites it rather than carrying the survey.
- `design/system/README.md` (layout): the index row of the finding document.
- `design/adrs/README.md` (layout): `flai adr new` writes the draft ADR's index row.
- `design/adrs` (design, folder kept): the draft ADR's number and slug are given by `flai adr new` when it is written, so no task can name the file. S-0495, in review, carries `0201`, so the number depends on which story lands first. `flai task done` widens the touches with the committed path (ADR-0144). The folder is the only one kept, and only for that file.

Forecast. See the forecast's basis: it is `flai forecast`'s figure after the touches above were written, reviewed against the agent seconds of the done research stories that wrote one finding document and one ADR (S-0339 2102 s, S-0380 2466 s, S-0413 1372 s).

Cost of delay. 50 USD a week, `flai cod`'s value from the operator's input `penalty_per_week: 50` and nothing else. It stands: the story has no epic and no revenue or time-lost input.
