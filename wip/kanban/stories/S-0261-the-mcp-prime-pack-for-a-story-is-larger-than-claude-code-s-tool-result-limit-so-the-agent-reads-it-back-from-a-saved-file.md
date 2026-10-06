---
id: S-0261
type: story
nature: improvement
title: The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file
status: backlog
owner: alex
created: 2026-10-04T20:34:55Z
updated: 2026-10-06T21:31:43Z
transitions: []
tags: []
topics: [cli, conventions]
touches: [flai/internal/context, flai/internal/mcpserver, flai/cmd/prime.go, flai/cmd/prime_test.go, design/adrs, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/issues/I-0068-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-04T20:34:55Z
  value: 15
  by: planner-S-0261
  at: 2026-10-05T05:51:39Z
forecast:
  duration: 1h
  delivery: 2026-10-07T05:38:00Z
  basis: "Its own forecast of 1h; 20th in the pull order with an in-progress limit of 3, behind S-0227, S-0229, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251 and S-0254."
  by: flai
  at: 2026-10-06T21:31:43Z
---
# S-0261 The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file

## Goal

This story remediates [I-0068](../../../design/issues/I-0068-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md), "The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0068 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0068 is closed with `flai issue close I-0068 --reason` saying what fixed it

## Tasks
- T-0986 Measure I-0068's packs and propose the fix as an ADR refining ADR-0049
- T-0987 The context pack is fitted so that the prime result stays under the tool result limit, with a test that reproduces I-0068
- T-0988 MCP prime and flai prime return the fitted pack, and this repository's largest packs fit a tool result
- T-0989 The design and the guides say how the prime pack fits a tool result
- T-0990 Close I-0068 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0068. time_lost_per_cycle 6m: 3m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-03T21:21:24Z, 1 day before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0261 on 2026-10-05.

#### Touches

The story declared no touches, so `flai touches suggest S-0261` ran from `flai/internal/context` and `flai/internal/mcpserver`. Of 869 commits, 80 changed those two folders.

- `flai/internal/context`: layout. The pack is assembled and fitted to the budget here: `Pack.size`, `AddDesign`, `AddBriefs`, and `DefaultBudget` in `budget.go`.
- `flai/internal/mcpserver`: layout. The MCP `prime` tool is `prime` in `server.go`, its description is in `folder.go`, and its tests are in `prime_test.go`.
- `flai/cmd/prime.go` and `flai/cmd/prime_test.go`: co-change, 10% and 8%. `flai prime --json` prints the same pack.
- `design/adrs`: design. The fix refines ADR-0049, which set the budget and kept every brief.
- `design/system/agent-context.md`: design. It holds the section "Fitting the pack to a budget".
- `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, and `docs/operators/settings.md`: co-change, 29%, 35%, 15%, and 14%. Each describes `flai prime` or `prime.budget`.
- The I-0068 issue file and `design/issues/summary.md`: from the second criterion, since `flai issue close` writes both.
- Left out: `design/system/project-manifest.md`. Add it only if the fix adds a key to the config next to `prime.budget`.

#### Figures

- **Forecast: 1h, raised from flai's 21m.** flai's figure is 88 s a unit over 20 large-band improvement stories, times a size of 14. That size counts the criteria and the touches but not the measuring and the ADR the goal asks for before building. S-0146 took 39m from in-progress to review, and it built the budget from an ADR already decided. Delivery is flai's date for 31st in the pull order, 2026-10-05T23:59Z, moved by the extra 39m.
- **Cost of delay: 15 USD a week, as `flai cod` gives it.** That is 6m lost each 168h cycle at 150 USD an hour. The inputs are the operator's and stand as they are. They count 2 occurrences, but I-0068 now records 3, and this planner's own prime pack hit the limit as well. The plan's thread asks whether to raise them.
