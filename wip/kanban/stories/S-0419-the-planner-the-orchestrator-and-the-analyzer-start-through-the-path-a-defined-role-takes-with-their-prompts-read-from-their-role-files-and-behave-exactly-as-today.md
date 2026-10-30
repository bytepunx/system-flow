---
id: S-0419
type: story
nature: improvement
title: The planner, the orchestrator, and the analyzer start through the path a defined role takes, with their prompts read from their role files, and behave exactly as today
status: review
parent: E-0024
owner: alex
created: 2026-10-16T05:40:19Z
updated: 2026-10-30T14:06:45Z
transitions:
  - to: ready
    at: 2026-10-28T13:05:18Z
    by: alex
  - to: backlog
    at: 2026-10-28T22:20:32Z
    by: alex
  - to: ready
    at: 2026-10-30T00:10:23Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-30T00:10:58Z
    by: agent-S-0419
  - to: review
    at: 2026-10-30T14:06:45Z
    by: agent-S-0419
blocked:
  - from: 2026-10-30T00:56:44Z
    until: 2026-10-30T01:57:41Z
    reason: awaiting the operator on TH-0760
  - from: 2026-10-30T13:34:19Z
    until: 2026-10-30T13:34:19Z
    reason: awaiting the operator on TH-0760
  - from: 2026-10-30T13:34:42Z
    until: 2026-10-30T14:05:54Z
    reason: awaiting the operator on TH-0760
tags: [cli, template]
topics: [cli, planning, orchestration, template]
touches: [design/conventions/strategic-agents.md, design/system/strategic-agents.md, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/serve/agents.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, flai/internal/check/activity_test.go, flai/internal/check/check_test.go, flai/internal/check/epic_lags_test.go, flai/internal/metrics/testdata/good/design/workflows/roles/analyzer.md, flai/internal/metrics/testdata/good/design/workflows/roles/curator.md, flai/internal/metrics/testdata/good/design/workflows/roles/orchestrator.md, flai/internal/metrics/testdata/good/design/workflows/roles/planner.md, flai/internal/harness/adapters.go, design/system/flai-cli.md, design/system/repository-layout.md, design/workflows/roles/README.md, docs/operators/index.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/check.go, template/root/design/workflows/roles/README.md, design/system/workflow.md, flai/cmd/check_stats_test.go]
after: [S-0418]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5900
  priced_by: harness
  turns:
    - day: 2026-10-30
      ceremony: 9
      test_runs: 2
      hand_edits: 3
      work: 127
  models:
    - model: claude-haiku-5-5
      input: 74
      output: 31375
      cache_read: 4568031
      cache_write: 212627
      cost: 0.3577
    - model: claude-opus-5-5
      input: 438
      output: 119148
      cache_read: 40340969
      cache_write: 1069270
      cost: 17.9763
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
      seconds: 1457
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 232
          output: 2514
          cache_read: 30771326
          cache_write: 422270
          cost: 6.887
        - model: claude-sonnet-5-5
          input: 10
          output: 70
          cache_read: 144835
          cache_write: 54514
          cost: 0.1377
cost_of_delay:
  value: 48.45
  by: planner-S-0419
  at: 2026-10-21T12:15:14Z
forecast:
  duration: 1h
  delivery: 2026-10-30T11:45:00Z
  basis: "Its own forecast of 1h; 25th in the pull order with an in-progress limit of 5, behind S-0424, S-0425, S-0527, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0506, S-0508, S-0509, S-0510, S-0511, S-0512, S-0516, S-0517, S-0518, S-0520, S-0523, S-0524 and S-0529."
  by: flai
  at: 2026-10-30T00:12:18Z
finalized:
  by: alex
  at: 2026-10-28T13:05:08Z
---
# S-0419 The planner, the orchestrator, and the analyzer start through the path a defined role takes, with their prompts read from their role files, and behave exactly as today

## Goal

One path serves every role. `serve.Plan`, the orchestrator's launcher, and `serve.Analyze` become the defined-role path over `design/workflows/roles/planner.md`, `orchestrator.md`, and `analyzer.md`: their prompts are their files' bodies, their host actions and triggers their files' fields, and the text a project writes below the marker reaches the prompt. Their commands, records, activity documents, and tests stay as they are, and a project with no role files still runs them from flai's embedded baseline.

## Acceptance criteria

- [x] The three roles' definition prompts carry their role file's body, its project additions included, from the project's file or the embedded baseline, and `flai/internal/harness/harness_test.go` (`TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, `TheAnalyzer`) passes with the same start prompt text.
- [x] `flai plan`, `plan.run`, the MCP tool `plan`, the replanner, the orchestrator's start, stop, and hold, `flai analyze`, `analyze.run`, and the analysis schedule start their role through the same start path a defined role takes; `serve/agents.json` keeps `plans`, `orchestrator`, `analyzer`, `past_orchestrators`, and `past_analyzers` as they are beside `roles`; and the tests under `flai/internal/serve` and `flai/cmd` for the three pass unchanged.
- [x] Text a project writes below the marker of `planner.md` reaches the planner's prompt, shown by a test.
- [x] A project with no `design/workflows/roles/` runs the three roles from the baseline files embedded in flai, and `flai check` warns `roles.baseline`, naming `flai upgrade`.
- [x] `design/system/strategic-agents.md` (How it runs, What it is told, for each role), the convention `strategic-agents.md`, and `template/CHANGELOG.md` say that the three roles are role files run through the same path.

## Tasks
- T-1753 The planner's, orchestrator's, and analyzer's definitions carry their role file's body, its additions included, when the harness is given the file
- T-1754 A project with no design/workflows/roles runs the three roles from baseline files embedded in flai, and flai check warns roles.baseline naming flai upgrade
- T-1756 Every start of a built-in role goes through startOf with its role file, the project's or the baseline, so the text below planner.md's marker reaches the planner's prompt
- T-1757 serve.Analyze and the analysis schedule start the analyzer through the defined-role path, keeping analyzer and past_analyzers in agents.json
- T-1759 The orchestrator's launcher, its stop, start, and hold, runs the orchestrator as the continuous role orchestrator.md defines, keeping orchestrator and past_orchestrators in agents.json
- T-1761 The design, the convention, and the changelog say the three roles are role files run through the path every role takes

## Notes

- Drafted by the planner for E-0024. Waits for the trigger story, which gives the path the `continuous` and schedule triggers the orchestrator and the analyzer need.
- The orchestrator's hold, the replanner's queue, and the analysis scheduler keep their behaviour; this story changes where they start a run, not when.

- **Criteria 1 and 2 reworded (2026-10-30).** The operator confirmed on TH-0760 the orchestrator's recommendation, option A. A role file's body is the definition's prompt, and the start prompts stay in Go (ADR-0207 § 4). Criterion 1 says so, and criterion 2 names the same start path a defined role takes, `startOf`, in place of "the defined-role path". Reworded by agent-S-0419.

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
