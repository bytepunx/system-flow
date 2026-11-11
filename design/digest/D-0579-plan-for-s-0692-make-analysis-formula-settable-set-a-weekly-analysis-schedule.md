---
id: D-0579
title: "Plan for S-0692: make analysis.formula settable, set a weekly analysis schedule, four tasks in three layers; two questions: The revisited plan is approved, and S-0692 is finalized"
kind: answer
sources: [TH-1038]
topics: [analysis, cli, flai, go, template, conventions, code, all]
paths: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check_roles.go, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, template/CHANGELOG.md, system-flow.yaml, design/issues/I-0292-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md, design/issues/summary.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0579 Plan for S-0692: make analysis.formula settable, set a weekly analysis schedule, four tasks in three layers; two questions: The revisited plan is approved, and S-0692 is finalized

The operator, alex, answered TH-1038 at 2026-11-11T02:24:25Z.

What was asked:

```text
The plan for S-0692 is written. It needs two answers from you before it is finalized: the analysis cadence and the cost of delay inputs.
```

The answer:

```text
**The revisited plan is approved, and S-0692 is finalized.**

- The value is 300 USD a week, from `flai cod` on the 2h input. The forecast stays at 38m, and the 20 file touches stay as they are.
- The four tasks in three layers cover both criteria.
- The `analyze` action's wording in `flai/internal/hostapi/writes_actions.go` stays out of scope, as the planner left it. It is not a criterion of S-0692.

**Question 1, still recommended:** weekly, `0 6 * * 1`. That is about 3.5 USD a month at the last run's 0.82 USD. The alternatives are daily, about 25 USD a month, or twice a month, about 1.6 USD a month. Only T-3037 waits on your answer. T-3031 and T-3034 can run before it.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/conventions/strategic-agents.md § As the orchestrator
```

Source: TH-1038.
