---
id: S-0419
type: story
nature: improvement
title: The planner, the orchestrator, and the analyzer start through the path a defined role takes, with their prompts read from their role files, and behave exactly as today
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:19Z
updated: 2026-10-23T17:43:05Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, template]
touches: [design/conventions/strategic-agents.md, design/workflows/roles/analyzer.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/system/strategic-agents.md, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/serve/agents.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analyze.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/plan.go, flai/internal/serve/replan.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md]
after: [S-0418]
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
      seconds: 120
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 51
          output: 5737
          cache_read: 3726197
          cache_write: 294016
          cost: 1.282
    - kind: orchestrator
      seconds: 579
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 136
          output: 1628
          cache_read: 23734706
          cache_write: 175728
          cost: 5.2796
draft: true
cost_of_delay:
  value: 48.45
  by: planner-S-0419
  at: 2026-10-21T12:15:14Z
forecast:
  duration: 1h
  delivery: 2026-10-24T03:39:00Z
  basis: "Its own forecast of 1h; 14th in the pull order with an in-progress limit of 5, behind S-0429, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417 and S-0418."
  by: flai
  at: 2026-10-23T17:43:05Z
---
# S-0419 The planner, the orchestrator, and the analyzer start through the path a defined role takes, with their prompts read from their role files, and behave exactly as today

## Goal

One path serves every role. `serve.Plan`, the orchestrator's launcher, and `serve.Analyze` become the defined-role path over `design/workflows/roles/planner.md`, `orchestrator.md`, and `analyzer.md`: their prompts are their files' bodies, their host actions and triggers their files' fields, and the text a project writes below the marker reaches the prompt. Their commands, records, activity documents, and tests stay as they are, and a project with no role files still runs them from flai's embedded baseline.

## Acceptance criteria

- [ ] `planPrompt`, `orchestratePrompt`, and `analyzePrompt` are read from the bodies of the three role files, in the template and here, and `flai/internal/harness/harness_test.go` (`TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, `TheAnalyzer`) passes with the same prompt text.
- [ ] `flai plan`, `plan.run`, the MCP tool `plan`, the replanner, the orchestrator's start, stop, and hold, `flai analyze`, `analyze.run`, and the analysis schedule start their role through the defined-role path; `serve/agents.json` keeps `plans`, `orchestrator`, `analyzer`, `past_orchestrators`, and `past_analyzers` as they are beside `roles`; and the tests under `flai/internal/serve` and `flai/cmd` for the three pass unchanged.
- [ ] Text a project writes below the marker of `planner.md` reaches the planner's prompt, shown by a test.
- [ ] A project with no `design/workflows/roles/` runs the three roles from the baseline files embedded in flai, and `flai check` warns `roles.baseline`, naming `flai upgrade`.
- [ ] `design/system/strategic-agents.md` (How it runs, What it is told, for each role), the convention `strategic-agents.md`, and `template/CHANGELOG.md` say that the three roles are role files run through the same path.

## Tasks
- T-1753 planPrompt, orchestratePrompt, and analyzePrompt are read from the three role files' bodies, with the harness tests' prompt text unchanged
- T-1754 A project with no design/workflows/roles runs the three roles from baseline files embedded in flai, and flai check warns roles.baseline naming flai upgrade
- T-1756 serve.Plan and the replanner start the planner through the defined-role path, keeping plans in agents.json, and the text below planner.md's marker reaches its prompt
- T-1757 serve.Analyze and the analysis schedule start the analyzer through the defined-role path, keeping analyzer and past_analyzers in agents.json
- T-1759 The orchestrator's launcher, its stop, start, and hold, runs the orchestrator as the continuous role orchestrator.md defines, keeping orchestrator and past_orchestrators in agents.json
- T-1761 The design, the convention, and the changelog say the three roles are role files run through the path every role takes

## Notes

- Drafted by the planner for E-0024. Waits for the trigger story, which gives the path the `continuous` and schedule triggers the orchestrator and the analyzer need.
- The orchestrator's hold, the replanner's queue, and the analysis scheduler keep their behaviour; this story changes where they start a run, not when.

### Planning

- **TH-0471's answers (2026-10-21).**
  - **Answer 2** moves the folder from `design/roles/` to `design/workflows/roles/`. The goal, criterion 4, and the six role-file touches now say so. T-1754's title and Work, and the Work of T-1753 and T-1756, were rewritten to the new folder and order, as the orchestrator approved on TH-0579.
  - **Answer 6** has flai generate a role's harness definition. This story keeps the three built-in roles on their `.claude/agents` definitions. `TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, and `TheAnalyzer` write their fixtures under `.claude/agents` and assert the refusal that names `flai upgrade`, and criterion 1 keeps them unchanged. Generating the built-in definitions is proposed on the plan thread as a story of its own.
  - **Answer 1's** 4h is not yet on E-0024 (TH-0567, question 3).
- **Order of the first two tasks.** T-1754, the embedded baseline, now comes first and T-1753 waits for it. The harness's prompt tests start roles in a temporary root with no role files. They pass unchanged only when the prompt falls back to the embedded bodies.
- **Touches.** 30, all files; no folder touch.
  - **Declared, kept:** the 19 earlier touches that are not role files, `harness_test.go` among them. That test must not change. It stays declared in case a fixture needs a role file.
  - **Declared, moved:** the six role files, now under `design/workflows/roles/` (answer 2).
  - **From layout, added (code read 2026-10-21):**
    - `flai/internal/roles/baseline/planner.md`, `orchestrator.md`, and `analyzer.md`, on T-1754 and T-1753. `go:embed` cannot reach `template/` from the `flai` module, so flai embeds a copy, as `internal/designsystem` does with `starter.md`. A test keeps the copy in step with the template's.
    - `flai/internal/serve/agents.go` on T-1756. A defined role's run is recorded under `roles`, and the built-in roles must keep `plans`, `orchestrator`, and `analyzer`.
    - `flai/internal/serve/serve.go` on T-1759. Today the launcher loop calls `orch.look` beside the replanner and the analysis scheduler. Once S-0418's role loop runs continuous roles, the orchestrator must be started by one of them, not both.
  - **Co-change, not added:** `docs/users/flai.md`, `design/system/flai-cli.md`, `template/template.yaml`, and `docs/operators/*`. No command, flag, manifest key, or template file changes here. `flai/cmd/plan.go`, `analyze.go`, and `serve_actions.go` call `serve` functions whose signatures stay.
- **Forecast.** Raised from flai's 45m to 1h. flai's 45m is 89 s a unit over 85 done improvement stories on this model, times size 30 (5 criteria, 25 touches). The raise has two reasons. Rewiring three launchers with their tests unchanged is the epic's main regression risk. The plan is now five serial layers. The delivery is S-0419's start in flai's play-out, 21:27Z after S-0418, plus 1h times the cycle factor 1.15: 2026-10-21T22:37Z. flai's replay since moved it to 22:36Z and rewrote the basis, keeping the 1h.
- **Cost of delay.** 48.45 USD a week. I worked it out by hand from answer 1, as the planners of S-0414, S-0417, and S-0418 did, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle, at `planning.hour_rate` 150 over the default 168h cycle, gives E-0024 600 USD a week.
  - S-0419's share is its 1h over the 12h23m forecast for E-0024's twelve open stories: 60 / 743 × 600.
  - Once 4h is on E-0024, `flai cod S-0419` gives the same figure while the forecasts stand.
