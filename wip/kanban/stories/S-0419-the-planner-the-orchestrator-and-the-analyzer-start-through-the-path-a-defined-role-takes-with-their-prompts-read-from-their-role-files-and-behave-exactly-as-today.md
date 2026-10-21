---
id: S-0419
type: story
nature: improvement
title: The planner, the orchestrator, and the analyzer start through the path a defined role takes, with their prompts read from their role files, and behave exactly as today
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:19Z
updated: 2026-10-21T07:19:47Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, template]
touches: [design/conventions/strategic-agents.md, design/roles/analyzer.md, design/roles/orchestrator.md, design/roles/planner.md, design/system/strategic-agents.md, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/harness/harness.go, flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analyze.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/plan.go, flai/internal/serve/replan.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/roles/analyzer.md, template/root/design/roles/orchestrator.md, template/root/design/roles/planner.md, flai/internal/harness/harness_test.go]
after: [S-0418]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h
  delivery: 2026-10-21T16:06:00Z
  basis: "Its own forecast of 1h; 14th in the pull order with an in-progress limit of 5, behind S-0453, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417 and S-0418."
  by: flai
  at: 2026-10-21T07:19:47Z
---
# S-0419 The planner, the orchestrator, and the analyzer start through the path a defined role takes, with their prompts read from their role files, and behave exactly as today

## Goal

One path serves every role. `serve.Plan`, the orchestrator's launcher, and `serve.Analyze` become the defined-role path over `design/roles/planner.md`, `orchestrator.md`, and `analyzer.md`: their prompts are their files' bodies, their host actions and triggers their files' fields, and the text a project writes below the marker reaches the prompt. Their commands, records, activity documents, and tests stay as they are, and a project with no role files still runs them from flai's embedded baseline.

## Acceptance criteria

- [ ] `planPrompt`, `orchestratePrompt`, and `analyzePrompt` are read from the bodies of the three role files, in the template and here, and `flai/internal/harness/harness_test.go` (`TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, `TheAnalyzer`) passes with the same prompt text.
- [ ] `flai plan`, `plan.run`, the MCP tool `plan`, the replanner, the orchestrator's start, stop, and hold, `flai analyze`, `analyze.run`, and the analysis schedule start their role through the defined-role path; `serve/agents.json` keeps `plans`, `orchestrator`, `analyzer`, `past_orchestrators`, and `past_analyzers` as they are beside `roles`; and the tests under `flai/internal/serve` and `flai/cmd` for the three pass unchanged.
- [ ] Text a project writes below the marker of `planner.md` reaches the planner's prompt, shown by a test.
- [ ] A project with no `design/roles/` runs the three roles from the baseline files embedded in flai, and `flai check` warns `roles.baseline`, naming `flai upgrade`.
- [ ] `design/system/strategic-agents.md` (How it runs, What it is told, for each role), the convention `strategic-agents.md`, and `template/CHANGELOG.md` say that the three roles are role files run through the same path.

## Tasks
- T-1753 planPrompt, orchestratePrompt, and analyzePrompt are read from the three role files' bodies, with the harness tests' prompt text unchanged
- T-1754 A project with no design/roles runs the three roles from baseline files embedded in flai, and flai check warns roles.baseline naming flai upgrade
- T-1756 serve.Plan and the replanner start the planner through the defined-role path, keeping plans in agents.json, and the text below planner.md's marker reaches its prompt
- T-1757 serve.Analyze and the analysis schedule start the analyzer through the defined-role path, keeping analyzer and past_analyzers in agents.json
- T-1759 The orchestrator's launcher, its stop, start, and hold, runs the orchestrator as the continuous role orchestrator.md defines, keeping orchestrator and past_orchestrators in agents.json
- T-1761 The design, the convention, and the changelog say the three roles are role files run through the path every role takes

## Notes

- Drafted by the planner for E-0024. Waits for the trigger story, which gives the path the `continuous` and schedule triggers the orchestrator and the analyzer need.
- The orchestrator's hold, the replanner's queue, and the analysis scheduler keep their behaviour; this story changes where they start a run, not when.

### Planning

- **Touches.** The story declared none. Every touch comes from its five tasks.
  - **From layout:** `flai/internal/harness/harness.go`; a new `flai/internal/roles/baseline.go` with its test; the check; serve's `plan.go`, `replan.go`, `analyze.go`, `analysis_schedule.go`, `orchestrate.go`, `orchestrate_hold.go`, `role.go`, and `role_continuous.go`; and the role files in both copies.
  - **From design:** `strategic-agents.md`, the convention in both copies, and the changelog.
  - **Co-change, added:** `flai/internal/harness/harness_test.go` (28%), on the prompt task. Its prompt assertions must not change, but its fixtures may need the role files in place.
  - **Co-change, not added:** `docs/users/flai.md` and `design/system/flai-cli.md`. No command or flag changes here.
  - No folder touch.
- **Forecast.** Raised from flai's 44m to 1h. Rewiring three launchers with their tests unchanged is the epic's main regression risk. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
