---
id: S-0720
type: story
nature: remediation
title: A usage limit on one Claude login holds every start on claude-code, and changing the harness's program or login does not lift it
status: backlog
owner: alex
created: 2026-11-11T05:05:59Z
updated: 2026-11-11T08:58:19Z
transitions: []
tags: [flai, flaiover, serve, budget]
touches: [flai/internal/serve/budget.go, flai/cmd/serve_agent.go, design/system/flai-cli.md, flai/internal/serve/budget_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_start.go, flai/internal/serve/notify.go, flai/internal/serve/notify_test.go, flai/cmd/serve.go, flai/cmd/serve_budget.go, flai/cmd/serve_budget_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_actions_test.go, flaiover/src/lib/components/BudgetBanner.svelte, flaiover/src/lib/components/BudgetBanner.svelte.test.ts, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
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
    - kind: planner
      seconds: 1
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 676
          cache_read: 205603
          cache_write: 4504
          cost: 0.0911
    - kind: orchestrator
      seconds: 25
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 114
          cache_read: 345129
          cache_write: 52377
          cost: 0.1109
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h30m
    by: orchestrator
    at: 2026-11-11T08:52:16Z
  value: 375
  by: planner-S-0720
  at: 2026-11-11T08:56:05Z
forecast:
  duration: 1h30m
  delivery: 2026-11-12T00:32:00Z
  basis: "Its own forecast of 1h30m; 27th in the pull order with an in-progress limit of 3, behind S-0694, S-0698, S-0702, S-0706, S-0719, S-0724, S-0716, S-0689, S-0705, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0699, S-0701, S-0703, S-0707, S-0709, S-0710, S-0711, S-0712, S-0714, S-0715 and S-0717."
  by: flai
  at: 2026-11-11T08:58:19Z
---
# S-0720 A usage limit on one Claude login holds every start on claude-code, and changing the harness's program or login does not lift it

## Goal

When the operator changes what a harness runs as (its program, its arguments, or the container's login), a usage limit recorded for the old login stops holding starts, and a limit on one login never holds a project that runs on another. The operator can also clear a budget state by hand.

## Acceptance criteria
- [ ] `flai serve agent harness <name>` (the dashboard's harness setting included), when it changes the program, arguments, or login a project's or the host's harness runs with, ends the hold on that harness's budget state for what it changed, so the next start goes as the probe; the change is journaled under `serve.budget` with why, and sent as a notice like any other clear.
- [ ] A budget state is kept under what tells one login from another (the provider, else the harness with its resolved program and login), so a run of `rmhclod` ending on the limit does not hold a project whose `claude-code` program is `meclod`; a `budget.json` in the old shape is still read.
- [ ] `flai serve budget clear [<provider>]` clears a state by hand and `flai serve budget` shows it; the hold's message (`budgetHeldWhy`) and the dashboard name that command.
- [ ] Tests cover: a harness change lifting the hold, two programs on one harness held apart, and the manual clear; `design/system/flai-cli.md` and the operator docs describe all three.

## Tasks
- T-3095 Keep a budget state under the harness's resolved program and login, and still read the old shape
- T-3096 The dashboard's budget banner names flai serve budget clear
- T-3098 Clear a budget state with a reason: journaled under serve.budget and sent as a notice, and budgetHeldWhy names flai serve budget clear
- T-3104 flai serve budget shows the budget states and flai serve budget clear clears one by hand
- T-3105 flai serve agent harness clears the budget state of the program, arguments, or login it changes
- T-3106 Describe the per-login budget state, the harness change's clear, and flai serve budget in the design and the docs

## Notes

What happened, 2026-11-11 (UTC): S-0692's run on `rmhclod` (`CLAUDE_CONFIG_DIR=~/.claude-rmh`) ended on `usage_limit_reached` (429) at 02:33:30Z, and flai serve set `~/.flai/serve/budget.json` key `claude-code`, held until the API's reset, 2026-11-13T22:00:00Z. At 02:35:06Z the operator set system-flow's `claude-code` program to `meclod` (another account) from the dashboard (`settings.harness`). The hold stayed: S-0690, S-0692, and S-0693 were logged "agents not started: the budget state holds them" through a flai serve restart at 05:01Z. Removing `budget.json` by hand at 05:04Z let all three start on `meclod` at the next pass.

Why: `budgetKey` (`flai/internal/serve/budget.go`) is the run's provider, or else the harness's name, `claude-code`. The program and login are not part of it, so every project on the host shares one state whichever account its program logs in with (signet was on `meclod` already and was held too). The state clears only when a probe ends clean, and a probe only starts after `until`; nothing else clears it, and there is no command for it. The error text in `Dir.budgets` tells the operator to remove the file.

Still true after the manual clear: the host default and authstar run `rmhclod`, so their next run ends on the limit and sets the host-wide hold again, system-flow and signet included. That is the second criterion.

### Planning

Touches, by where each came from:

| Source | Paths |
|--------|-------|
| Declared | `flai/internal/serve/budget.go`, `flai/cmd/serve_agent.go`, `design/system/flai-cli.md` |
| Layout: the key and the run's identity | `flai/internal/serve/budget_test.go`, `flai/internal/serve/agents.go` (`AgentRun`), `flai/internal/serve/agents_start.go` (`startOf` and the gate) |
| Layout: the clear's notice | `flai/internal/serve/notify.go`, `flai/internal/serve/notify_test.go` |
| Layout: the new commands | `flai/cmd/serve.go` (registration), `flai/cmd/serve_budget.go`, `flai/cmd/serve_budget_test.go` (new) |
| Layout: the harness change | `flai/cmd/serve_harness.go` (`RunE`), `flai/cmd/serve_actions_test.go` (its tests) |
| Layout: the banner | `flaiover/src/lib/components/BudgetBanner.svelte` and its test |
| Co-change and design | `docs/users/flai.md` (76% co-change), `docs/operators/index.md` (30%), `docs/users/flai-reference.md` (25%, generated), `docs/operators/settings.md` (17%, generated flag table), `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` (the banner) |

No folder touches. The dashboard's `settings.harness` runs `flai serve agent harness` as a subprocess (`flai/internal/hostapi/settings.go`), so it needs no change and is not touched. The six end callers of `budgetEnded` are not touched: T-3095 keeps the identity on `AgentRun`. `design/system/project-manifest.md` does not describe the key, so it is not touched. The tags gained `flai` and `flaiover`, which `story.component-tag` asks for now the story touches both. No `after`: the story reads no other story's output.

Tasks, in four layers:

1. T-3095 (the key) and T-3096 (the banner), with no path in common.
2. T-3098 (the clear, its journal, notice, and `budgetHeldWhy`), after T-3095: it clears under T-3095's key, in the same files.
3. T-3104 (`flai serve budget`) and T-3105 (the harness change's clear), both after T-3098, whose clear they call; no path in common.
4. T-3106 (design and docs), after T-3104, T-3105, and T-3096, to describe them as built.

Figures:

- Forecast 1h30m, against flai's 46m (109 s per unit of size over 66 done remediation stories, size 25). Raised for four sequential layers over the budget file's lock and probe logic, and a close-out that runs the integration and smoke tiers for a change under `flai/`. Delivery is flai's, moved by the same 44m.
- Cost of delay 375 USD a week, as `flai cod` works it out from `time_lost_per_cycle` 2h30m (TH-1070, set by the orchestrator): 2h30m a week at 150 USD an hour. It stands: the input measures the 2026-11-11 hold, and the host default still runs `rmhclod`.

Assumption: whether an old-shape entry keyed `claude-code` holds every program on the harness until it clears is left to T-3095, which records the choice.
