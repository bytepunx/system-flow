---
id: S-0415
type: story
nature: feature
title: flai prime --role gives a role a project defines its pack within the budget, from the topics and the item kind its role file names, and the three built-in roles prime from their files unchanged
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:39:36Z
updated: 2026-10-16T07:30:22Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, conventions]
touches: [design/system/agent-context.md, design/system/conventions.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/conventions/conventions.go, flai/internal/conventions/conventions_test.go, flai/internal/mcpserver/prime_test.go, flai/internal/mcpserver/server.go, design/system/flai-cli.md]
after: [S-0414]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 29m
  delivery: 2026-10-16T18:14:00Z
  basis: "Its own forecast of 29m; 24th in the pull order with an in-progress limit of 5, behind S-0426, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412 and S-0414."
  by: flai
  at: 2026-10-16T07:30:22Z
---
# S-0415 flai prime --role gives a role a project defines its pack within the budget, from the topics and the item kind its role file names, and the three built-in roles prime from their files unchanged

## Goal

A role a project defines primes as the planner, the orchestrator, and the analyzer do. `flai prime --role <role>` and the MCP tool `prime` read the role's file for the topics its pack selects and the item it takes, and give it the conventions whose `roles` are empty or list it, briefs of the design its topics select, what its item names, and a catalog, within the budget (ADR-0049) and in parts (ADR-0104). The built-in roles' topics and item kinds move from Go into their role files, and their packs stay exactly as they are.

## Acceptance criteria

- [ ] `flai prime --role <role>` and the MCP tool `prime` with that role give a defined role its pack: its topics from its role file, the conventions whose `roles` are empty or list it, what its item names when its file takes one (`--epic`, `--story`, or none), the ranked sections and briefs, and a catalog, within `--budget` and in parts.
- [ ] `plan`, `orchestrate`, and `analyze` take their topics and item kinds from their role files, and their packs are byte for byte what they were: the existing tests in `flai/internal/context/strategic_test.go`, `flai/cmd/prime_test.go`, and the MCP prime tests pass unchanged.
- [ ] A role no file defines, an item kind its file does not take, and a project whose role file does not parse are refused, naming the defined roles or the file's finding.
- [ ] `flai check` accepts in a convention's `roles` every role a project defines beside `story`, `explore`, and `verify`, and warns `conventions.roles` for any other, as ADR-0075 does today.
- [ ] `design/system/agent-context.md` § Strategic packs, `design/system/conventions.md` § Roles, `docs/users/flai.md`, and the regenerated `flai-reference.md` say how a defined role primes.

## Tasks
- T-1725 The strategic pack takes a role's topics and item kind from its role file in place of the Go map, with the three built-in packs unchanged
- T-1726 flai check accepts in a convention's roles every role a project defines and warns conventions.roles for any other
- T-1727 flai prime --role and the MCP tool prime take any role a project defines, with its item flag, and refuse a role or item kind its file does not allow
- T-1728 The design and the user guide say how a defined role primes and which conventions it reads

## Notes

- Drafted by the planner for E-0024. Waits for the role file story, whose loader it reads.
- Runs beside the guard story: the two share no file.

### Planning

- **Touches.** The story declared none. Every touch comes from its four tasks.
  - **From layout:** `flai/internal/context/strategic.go` and its test, `flai/internal/conventions/conventions.go` and its test, `flai/cmd/prime.go` and its test, `flai/internal/mcpserver/server.go` (the `prime` tool) and `prime_test.go`, and `docs/users/flai-reference.md`, regenerated.
  - **From design:** `agent-context.md` § Strategic packs, `conventions.md` § Roles, and `docs/users/flai.md` for criterion 5.
  - **Co-change, added:** `design/system/flai-cli.md` (53%), which describes `flai prime`; shared.
  - **Co-change, not added:** the operator docs, the dashboard's design, `design/adrs/README.md` (no ADR here), and `mcpserver/folder.go`, since the `prime` tool is registered already.
  - No folder touch.
- **Forecast.** Kept at flai's 29m. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
