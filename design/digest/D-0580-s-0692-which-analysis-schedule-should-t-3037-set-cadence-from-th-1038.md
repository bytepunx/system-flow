---
id: D-0580
title: "S-0692: which analysis.schedule should T-3037 set? (cadence from TH-1038): Recommended: weekly, `0 6 * * 1`, Mondays at 06:00 UTC, with `analysis.formula: analysis-cadence`"
kind: answer
sources: [TH-1065, TH-1038, S-0692]
topics: [analysis, cli, flai, go, template, conventions, code, all]
paths: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check_roles.go, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, template/CHANGELOG.md, system-flow.yaml, design/issues/I-0292-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md, design/issues/summary.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: [S-0693, TH-1038, TH-1066]
---

# D-0580 S-0692: which analysis.schedule should T-3037 set? (cadence from TH-1038): Recommended: weekly, `0 6 * * 1`, Mondays at 06:00 UTC, with `analysis.formula: analysis-cadence`

This project's analyzer runs weekly. The schedule is `analysis.schedule: 0 6 * * 1`, Mondays at 06:00 UTC, and `analysis.formula: analysis-cadence` makes that formula's instance each time the schedule comes round. The operator, alex, chose this cadence. It is recurring spend, so it is the operator's choice alone.

S-0692's agent asked on TH-1065 which `analysis.schedule` T-3037 should set. The orchestrator recommended weekly, for the operator alone:

```text
**Recommended:** weekly, `0 6 * * 1`, Mondays at 06:00 UTC, with `analysis.formula: analysis-cadence`. That is about 3.5 USD a month at the last run's 0.82 USD.
```

The alternatives were `daily`, about 25 USD a month, and `0 6 1,15 * *`, about 1.6 USD a month. A week matches `planning.cycle`'s default, 168h. The 2026-11-07 analysis had found the W44 regressions a week late. alex confirmed the recommendation at 2026-11-11T02:25:28Z. A minute earlier, alex had also confirmed the orchestrator's entry on TH-1038, which recommended the same weekly cadence beside the approval of S-0692's plan (D-0579).

S-0692's narrative records the decision in its own words: "The operator confirmed on TH-1065 (2026-11-11T02:25:28Z) the weekly cadence, `analysis.schedule: 0 6 * * 1` with `analysis.formula: analysis-cadence`; T-3037 sets those." T-3037 set both keys with `flai manifest set` (4d052e4), and T-3040 closed I-0292. The orchestrator accepted S-0692 at 2026-11-11T05:07:30Z. The schedule starts nothing while the `analyze` host action is off for the project. The planner also assumed that the template's manifest gets no default schedule, since a default would make every new project spend on the analyzer without its operator choosing it.

Sources: TH-1065, TH-1038, S-0692 (`wip/archive/agents/S-0692.md`).
