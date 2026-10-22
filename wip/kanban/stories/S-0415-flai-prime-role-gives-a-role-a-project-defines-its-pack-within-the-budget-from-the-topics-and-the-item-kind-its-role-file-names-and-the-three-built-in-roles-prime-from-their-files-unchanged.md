---
id: S-0415
type: story
nature: feature
title: flai prime --role gives a role a project defines its pack within the budget, from the topics and the item kind its role file names, and the three built-in roles prime from their files unchanged
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:39:36Z
updated: 2026-10-22T21:15:43Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, conventions]
touches: [design/system/agent-context.md, design/system/conventions.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/context/story.go, flai/internal/conventions/conventions.go, flai/internal/conventions/conventions_test.go, flai/internal/check/check.go, flai/internal/mcpserver/prime_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, design/system/flai-cli.md]
after: [S-0414]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 330
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 122
          output: 1954
          cache_read: 14213780
          cache_write: 21995
          cost: 3.1436
draft: true
cost_of_delay:
  value: 32.52
  by: planner-S-0415
  at: 2026-10-21T11:17:24Z
forecast:
  duration: 38m
  delivery: 2026-10-23T04:28:00Z
  basis: "Its own forecast of 38m; 18th in the pull order with an in-progress limit of 5, behind S-0466, S-0477, S-0468, S-0467, S-0476, S-0469, S-0470, S-0472, S-0478, S-0481, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412 and S-0414."
  by: flai
  at: 2026-10-22T21:15:43Z
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

- **TH-0471's answers (2026-10-21).** None changes this story's words. Role files now live under `design/workflows/roles/` (answer 2), which `roles.Load` (S-0414) resolves, so no path here names the folder. Generated harness definitions (answer 6) and per-role agents (answer 4) are not prime's concern.
- **Touches.** 16, all files; no folder touch.
  - **Declared, kept:** the 13 earlier touches.
  - **From layout, added (code read 2026-10-21):**
    - `flai/internal/context/story.go`: `ForRole` refuses an unknown role naming `conventions.StrategicRoles` (T-1725).
    - `flai/internal/check/check.go`: `checker.conventions` is the one caller of `Set.Validate`, which needs the defined roles (T-1726).
    - `flai/internal/mcpserver/folder.go`: the `prime` tool's description, naming the three roles, is registered there, not in `server.go` (T-1727). The earlier note that left it out was wrong.
  - **Co-change, kept:** `design/system/flai-cli.md` (shared, holds nothing).
  - **Co-change, not added:** `docs/operators/index.md` and `settings.md` (no setting or flag is added; `--role` is in the flag index already), `flaiover-dashboard.md` and `docs/users/flaiover.md` (no page changes), `design/adrs/README.md` (no ADR), `workflow.md`, and `template/CHANGELOG.md` (no template file changes; S-0414 ships the role files and their topics).
  - **Left to later stories:** `internal/guard`, `internal/serve`, and `internal/harness` also read the fixed role list; S-0416 to S-0419 change them.
- **Forecast.** flai's 38m is kept: 107 s a unit times size 21, five criteria and 16 touches. The three touches added are small edits, which overstates the size a little. Keeping the three packs byte for byte through a new loader is the risk, and it offsets that. The delivery is flai's play-out, 2026-10-21T16:30Z.
- **Cost of delay.** 32.52 USD a week, worked out by hand from the operator's answer 1 on TH-0471, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week.
  - S-0415 takes the share its 38m is of the 11h41m forecast for E-0024's twelve open stories: 38 / 701 × 600.
  - Once `time_lost_per_cycle: 4h` is on E-0024, `flai cod S-0415` gives the same figure while the forecasts stand.
