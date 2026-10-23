---
id: D-0351
title: "S-0419 plan: six tasks in five layers, moved to design/workflows/roles, with the embedded baseline first: Recommended for the operator"
kind: answer
sources: [TH-0579]
topics: [cli, planning, orchestration, template, conventions, flai, go, code, all]
paths: [design/conventions/strategic-agents.md, design/workflows/roles/analyzer.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/system/strategic-agents.md, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/serve/agents.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analyze.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/plan.go, flai/internal/serve/replan.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0351 S-0419 plan: six tasks in five layers, moved to design/workflows/roles, with the embedded baseline first: Recommended for the operator

The operator, alex, answered TH-0579 at 2026-10-21T12:33:31Z.

What was asked:

```text
S-0419 is replanned on your TH-0471 answers. It keeps its six tasks. The embedded baseline now comes before the prompt move. The story has a forecast and a cost of delay value.
```

The answer:

```text
Recommended for the operator:

- **Question 2:** take the planner's answer. Generate the built-in roles' definitions in a new story after S-0419, not inside it. S-0419's criterion 1 keeps the three claude-code start tests unchanged, and those tests write fixtures under `.claude/agents`, so generating definitions there would break that promise. The new story would generate the three definitions from their role files and pass them in `--agents`, as S-0417 now does for defined roles. It would check them against today's files and settle the symlink for sessions started by hand. This replaces my recommendation on TH-0567, question 2, which put generation in S-0419. A new story is a scope addition, so it is yours; the epic's planner can draft it.
- **Question 3:** as on TH-0567, run `flai edit E-0024 --time-lost-per-cycle 4h`. Then `flai cod` values all of E-0024's stories in place of the planners' hand-worked figures. S-0419's would be 48.45 USD a week.

Source: design/conventions/strategic-agents.md § As the planner
```

Source: TH-0579.
