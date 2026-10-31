---
id: S-0554
type: story
nature: feature
title: A role file's start prompt and its triggers replace the planner's, the orchestrator's, the analyzer's, and the curator's start prompts and schedules in Go, so the built-in roles are defined by their files alone
status: backlog
parent: E-0024
owner: alex
created: 2026-10-30T20:45:40Z
updated: 2026-10-31T07:45:23Z
transitions: []
tags: [cli, template]
topics: [cli, orchestration, planning, template, conventions]
touches: [design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, design/workflows/roles/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/analyzer.md, design/workflows/roles/curator.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/curator.md, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/roles/baseline/curator.md, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/start.go, flai/internal/roles/start_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/replan.go, flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/serve.go, docs/operators/index.md, template/CHANGELOG.md]
after: [S-0553]
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
      seconds: 91
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 33
          output: 159
          cache_read: 1449087
          cache_write: 91757
          cost: 0.3471
draft: true
cost_of_delay:
  value: 71.06
  by: planner-E-0024
  at: 2026-10-30T20:48:20Z
forecast:
  duration: 1h19m
  delivery: 2026-11-01T00:16:00Z
  basis: "Its own forecast of 1h19m; 32nd in the pull order with an in-progress limit of 5, behind S-0548, S-0552, S-0550, S-0551, S-0546, S-0549, S-0547, S-0425, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0506, S-0508, S-0509, S-0510, S-0516, S-0517, S-0518, S-0520, S-0420, S-0421, S-0422, S-0423, S-0424, S-0540, S-0543, S-0545 and S-0553."
  by: flai
  at: 2026-10-31T07:45:23Z
---
# S-0554 A role file's start prompt and its triggers replace the planner's, the orchestrator's, the analyzer's, and the curator's start prompts and schedules in Go, so the built-in roles are defined by their files alone

## Goal

After S-0419 and S-0553, each built-in strategic role's definition comes from its role file. Two things still live in Go. The start prompts, `planPrompt`, `orchestratePrompt`, `analyzePrompt`, `curatePrompt`, and `definedPrompt` in `flai/internal/harness/harness.go`, say what each run is for. The rules for when the planner, the analyzer, and the curator run are in `flai/internal/serve/replan.go`, `analysis_schedule.go`, and `curate_schedule.go`. Their role files' `trigger` entries are validated but never read, and the conditions `stale-figures`, `replan-agent`, and `each-ready-story` are named but not implemented. This story moves both into the role files: a start-prompt section templated over the run's fields, under a new ADR refining ADR-0207 § 4, and the triggers read through the event and schedule launchers S-0418 built. Then the four built-in roles are defined by their files alone, as the operator asked on TH-0760, and a project changes what a built-in role is told, or when it runs, with no flai release.

## Acceptance criteria

- [ ] A new ADR refining ADR-0207 § 4 settles the start prompt as data. It says where it lives in a role file, apart from the body, which stays the definition's prompt. It says which run fields it may use: the item and its kind, who asked, the reasons and their prime forms, the merge queue, the focus, the window since the last run, and a defined role's inputs. It also says how a project's additions reach it and how `flai upgrade` merges it.
- [ ] `planPrompt`, `orchestratePrompt`, `analyzePrompt`, `curatePrompt`, and `definedPrompt` give way to rendering the role file's start prompt. The four baseline role files, here, in the template, and embedded in flai, carry today's text. The harness tests, `TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, `TheAnalyzer`, `TheCurator`, and `TestClaudeCodeStartsADefinedRole`, pass with the same start prompt text.
- [ ] `stale-figures`, `replan-agent`, and `each-ready-story` are implemented in `flai/internal/roles/conditions.go`. The planner starts from its role file's `trigger` entries through the event and schedule launchers, which take an item-taking role. `flai/internal/serve/replan_test.go` passes unchanged.
- [ ] The analyzer's `analysis.schedule`, and the curator's `accepted` and `digest.schedule`, start them from their role files' triggers through the same launchers. `analysis_schedule_test.go` and `curate_schedule_test.go` pass unchanged, and an analysis schedule that names a formula, if S-0424 has landed, still starts it.
- [ ] A test shows that a project's role file changes a built-in role's start prompt, by its project additions, and its triggers, with no change to flai.
- [ ] `design/system/strategic-agents.md`, the strategic-agents convention and `design/workflows/roles/README.md` here and in the template, `docs/operators/index.md`, and `template/CHANGELOG.md` say that a built-in role's start prompt and triggers are read from its file.

## Tasks
- T-2354 An ADR refining ADR-0207 § 4 settles a role file's start prompt as data and the built-in roles' triggers read from their files
- T-2355 flai reads a role file's start prompt and renders it over a run's fields, and flai check reports one that does not parse or names an unknown field
- T-2356 The trigger conditions stale-figures, replan-agent, and each-ready-story are implemented in flai/internal/roles/conditions.go
- T-2357 The harness renders each role's start prompt from its role file, and the four baseline role files carry today's start prompts
- T-2358 The planner starts from its role file's triggers through the event and schedule launchers, which take an item-taking built-in role, in place of the replanner's own triggers
- T-2359 The analyzer and the curator start from their role files' triggers through the same launchers, in place of the analysis and digest schedulers' own rules
- T-2360 The strategic-agents design and convention, the roles README, the operator docs, flai-cli.md, and the changelog say a built-in role's start prompt and triggers are read from its file

## Notes

- Drafted by the planner of E-0024 on 2026-10-30, on the operator's question on TH-0760. It is the follow-up the orchestrator's confirmed recommendation of 2026-10-30T12:57:49Z names there: a templated start-prompt section under a new ADR refining ADR-0207 § 4, and the replan rules and the analysis schedule as `event` and `schedule` triggers, after the definition-generation story (D-0351), which is S-0553. The curator's triggers are included, since its file declares them the same way and leaving it alone would keep one role in Go.
- The orchestrator's triggers are already read from its file (`orchestrate.go`, `orchestrate_reasons.go`), and its reasons-and-end run (ADR-0209) stays. ADR-0207's `continuous` trigger is not built here: no role needs it.
- The records stay keyed as they are (`plans`, `orchestrator`, `analyzer`, `curator`). Keying them by the role's name is bookkeeping a later story can do once nothing else tells the roles apart.

### Planning

- **After.** S-0553, as the confirmed recommendation on TH-0760 orders it: the start-prompt section is a new part of the role file that S-0553's generation must leave out of the definition, so this story builds on the role-file shape S-0553 lands. No other story's change is needed on main. S-0424 changes `analysis_schedule.go` too; the touch overlap holds whichever is pulled second, and criterion 4 keeps S-0424's formula schedule working if it lands first.
- **Touches.** 38, all files but one.
  - **From layout** (S-0419's worktree, read 2026-10-30): the five start prompts in `flai/internal/harness/harness.go`, dispatched by `Prompt`; the replanner's triggers in `serve/replan.go`; `serve/analysis_schedule.go` and `serve/curate_schedule.go`; the launchers `serve/role_events.go` and `serve/role_schedule.go`, which skip built-in roles and an item-taking schedule today; `serve/serve.go`, which wires the schedulers; `flai/internal/roles/conditions.go`, where the three trigger conditions are named but `Matches` reports them not implemented; `roles.go` for the parsed start prompt. `roles/start.go` and its test are new, for rendering.
  - **From design:** ADR-0207 § 3 and § 4; the four role files in all three copies, which gain the start prompt; `strategic-agents.md`, the convention pair, and the roles README pair, which describe the roles; `docs/operators/index.md`; the changelog.
  - **Co-change, added:** `design/system/flai-cli.md` (36% of the seed commits), which describes when `flai serve` starts each role.
  - **Co-change, not added:** `docs/users/flai.md` and `flai-reference.md`: no command or flag changes. `docs/operators/settings.md`: no setting changes. `flai/internal/serve/agents.go`: `startOf` already loads the role file; `flai task done` widens the touches if the request's fields move.
  - **Folder touch kept:** `design/adrs`, since the new ADR's file name is not known until it is written.
- **Forecast.** flai's 1h19m kept: no live run, and the 37 touches count the port of five prompts and three trigger paths. The delivery is flai's play-out after S-0553's adjusted delivery, 34th in the pull order.
- **Cost of delay.** 71.06 USD a week, worked out by hand as the planners of S-0414 to S-0425 did, because `flai cod` refuses while E-0024 has no inputs (TH-0567, question 3, confirmed; TH-0587 asks the operator to write 4h on the epic). 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week. S-0554's share is its 1h19m of the 11h7m forecast for E-0024's nine open stories, S-0419 to S-0425, S-0553, and S-0554: 79 / 667 × 600.
- **Size.** Six criteria and two halves that share only the role files: the start prompts (harness) and the triggers (serve). The plan thread on E-0024 proposes splitting it there if the operator prefers two smaller stories; the tasks are layered so that the split needs no rewrite.
