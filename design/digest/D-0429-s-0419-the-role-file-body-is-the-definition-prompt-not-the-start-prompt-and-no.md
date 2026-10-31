---
id: D-0429
title: "S-0419: the role file body is the definition prompt, not the start prompt, and no continuous launcher exists; how should criteria 1, 2, and T-1759 read: if we go with A now, does that prevent us from…"
kind: recommendation
sources: [TH-0760]
topics: [cli, planning, orchestration, template, conventions, flai, go, code, all]
paths: [design/conventions/strategic-agents.md, design/workflows/roles/analyzer.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/system/strategic-agents.md, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/serve/agents.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analyze.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/plan.go, flai/internal/serve/replan.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, flai/internal/check/activity_test.go, flai/internal/check/check_test.go, flai/internal/check/epic_lags_test.go, flai/internal/metrics/testdata/good/design/workflows/roles/analyzer.md, flai/internal/metrics/testdata/good/design/workflows/roles/curator.md, flai/internal/metrics/testdata/good/design/workflows/roles/orchestrator.md, flai/internal/metrics/testdata/good/design/workflows/roles/planner.md]
date: 2026-10-30
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0429 S-0419: the role file body is the definition prompt, not the start prompt, and no continuous launcher exists; how should criteria 1, 2, and T-1759 read: if we go with A now, does that prevent us from…

The operator, alex, settled TH-0760 on S-0419 at 2026-10-30T14:05:53Z by confirming the orchestrator's recommendation of 12:57:49Z: option A.

What was asked: S-0419's plan assumed a built-in role file's body could become the text `planPrompt`, `orchestratePrompt`, and `analyzePrompt` produce, and that a continuous-role launcher existed. Neither held. The body of `design/workflows/roles/<name>.md` is the agent's definition prompt (the system prompt), not the start prompt (the per-run user message), as ADR-0207 § 4 says. No continuous launcher exists: ADR-0209 makes the orchestrator run for a reason and end.

The decision, option A:

- The planner, the orchestrator, the analyzer, and the curator load their role file, the project's or flai's embedded baseline, and start through the same start path a defined role takes (`startOf`). The harness passes the file's body, with the project's additions below the marker, as the definition's prompt in `--agents`.
- The start prompts, the run records (`plans`, `orchestrator`, `analyzer`), and when each role runs stay in Go. `role_continuous.go` was dropped (T-1759 cancelled).
- Option B, start-prompt text moved into the bodies as templates, was turned down because it contradicts ADR-0207 § 4.

The operator asked first whether A prevents the built-in roles from one day being defined by data files alone. The answer, which the confirmation took: no. Moving the start prompts needs a templated start-prompt section in the role file under a new ADR refining ADR-0207 § 4, leaving the body as the definition. The replan rules and the analysis schedule would move to `event` and `schedule` triggers. That follow-up is a backlog story under E-0024, after the definition-generation story (D-0351). Until that generation story lands, D-0411 still holds.

The operator, alex, answered TH-0760 again at 2026-10-30T20:38:30Z.

The answer:

```text
Are you able to write the story?
```

Source: TH-0760.
