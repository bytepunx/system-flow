---
id: S-0272
type: story
nature: improvement
title: "An agent with an open question ends instead of waiting: flai serve restarts it on the answer, and wait_for_events keeps a timeout only for an agent with work in hand"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:31Z
updated: 2026-10-06T21:00:33Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, conventions]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/usage/log.go, flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/metrics.md, design/adrs, design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 96
  by: planner-E-0017
  at: 2026-10-06T11:36:16Z
forecast:
  duration: 50m
  delivery: 2026-10-07T06:55:00Z
  basis: "Its own forecast of 50m; 29th in the pull order with an in-progress limit of 3, behind S-0224, S-0227, S-0229, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0261, S-0264, S-0265, S-0269, S-0270 and S-0271."
  by: flai
  at: 2026-10-06T21:00:33Z
---
# S-0272 An agent with an open question ends instead of waiting: flai serve restarts it on the answer, and wait_for_events keeps a timeout only for an agent with work in hand

## Goal

650 turns across 38 story runs woke from `wait_for_events` to find nothing, each a full-context model call. flai serve already restarts a story agent when its thread is answered, so an agent whose only pending work is a question has nothing to wait for. The convention and the harness prompt say to end the session on an open question after writing the narrative's state; `wait_for_events` answers at once with `end: true` when the caller's story has an open thread and no task in progress, and keeps its timeout for an agent that has work to go on with. The MCP server's instructions say the same.

## Acceptance criteria
- [ ] `wait_for_events` answers `end: true` and why when the calling agent's story has an open thread to the designer and no task in progress, and the harness prompt and conventions tell the agent to end then
- [ ] flai serve's restart on the answer is tested end to end: an agent that ended on a question is started again when the thread is answered, with the answer in its first inbox
- [ ] `design/system/flai-serve.md` (or where serve is described), the conventions, the template's copies, and the user guide describe the loop
- [ ] `flai stats` counts the empty wakes so that the saving is measured

## Tasks

## Notes

From the epic's log classification: MCP wait_for_events 650 turns, 1,235 minutes, 38 stories.

### Planning

Touches, none declared before. `flai touches suggest S-0272` was seeded with `flai/internal/serve` and `flai/internal/mcpserver`, which 145 commits changed:

- `flai/internal/mcpserver/server.go`, `server_test.go`: layout. `wait_for_events` and the server's instructions are defined here; S-0285 last changed its description.
- `flai/internal/serve/restart.go`, `flai/internal/serve/restart_test.go`: layout (criterion 2). The restart on an answered thread lives in `restart.go`, and it has no test file of its own yet.
- `flai/internal/harness/harness.go`, `harness_test.go`: layout. The story prompt says today to hold `wait_for_events` until the thread is answered, and also that flai starts the agent again if it ends.
- `flai/internal/usage/log.go`, `flai/internal/metrics/waiting.go`, `waiting_test.go`, `flai/cmd/stats.go`, `flai/cmd/check_stats_test.go`, `design/system/metrics.md`: design (criterion 4). Empty wakes come from the run logs `usage` reads, and waits are computed in `waiting.go`.
- `design/adrs`: design. `metrics.md` is the contract with the dashboard and changes only with an ADR.
- `design/conventions/work-management.md`, `delegation.md`, their `template/root` copies, and `template/CHANGELOG.md`: layout. These two conventions name `wait_for_events` for a story's agent.
- `design/system/flai-cli.md`, `docs/users/flai.md`, `flai-reference.md`: co-change (34, 35, and 23 of 145). There is no `design/system/flai-serve.md`; serve is described in `flai-cli.md`'s `flai serve` section.
- Left out: `flai/internal/guard`. S-0285's refusal of `wait_for_events` while a sub-agent runs is untouched by an `end` answer. Also left out: `flai/cmd/serve_actions.go` and `flai/internal/hostapi`, co-changed but not reached.

Forecast: flai gave 38m (89 s per unit over 21 done large improvement stories, times size 25). I raised it to 50m for an end-to-end restart test through serve and a metrics change that needs an ADR. S-0285, on the same tool and prompt, took 42m without either. The delivery is flai's, 2026-10-07T01:18Z, moved by the added 12m.

Cost of delay: 96 USD a week, against flai's 131.58. This is E-0017's 900 USD a week shared by the turns each story removes. This one removes the 650 empty wakes.
