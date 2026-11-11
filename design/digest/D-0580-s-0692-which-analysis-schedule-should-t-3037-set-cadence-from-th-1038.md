---
id: D-0580
title: "S-0692: which analysis.schedule should T-3037 set? (cadence from TH-1038): Recommended: weekly, `0 6 * * 1`, Mondays at 06:00 UTC, with `analysis.formula: analysis-cadence`"
kind: answer
sources: [TH-1065]
topics: [analysis, cli, flai, go, template, conventions, code, all]
paths: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check_roles.go, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, template/CHANGELOG.md, system-flow.yaml, design/issues/I-0292-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md, design/issues/summary.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0580 S-0692: which analysis.schedule should T-3037 set? (cadence from TH-1038): Recommended: weekly, `0 6 * * 1`, Mondays at 06:00 UTC, with `analysis.formula: analysis-cadence`

The operator, alex, answered TH-1065 at 2026-11-11T02:25:28Z.

What was asked:

```text
S-0692's code and docs are done: `flai manifest set` now writes `analysis.formula` (T-3031, T-3034). Only T-3037 is left before the close, and it needs your cadence. TH-1038 asks the same question; answer either thread.
```

The answer:

```text
**Recommended:** weekly, `0 6 * * 1`, Mondays at 06:00 UTC, with `analysis.formula: analysis-cadence`. That is about 3.5 USD a month at the last run's 0.82 USD.

**Why:**

- It is the same answer I recommend on TH-1038, which is still open. One answer from you settles both threads.
- A week is `planning.cycle`'s default, 168h. This project's manifest does not set it.
- The 2026-11-07 analysis found the W44 regressions a week late. A weekly run closes that gap.
- The schedule is a cron spec in UTC, and `analysis.formula` makes the `analysis-cadence` instance each time it comes round. Nothing starts until you turn on the `analyze` host action.

**Alternatives:**

| Cadence | Spec | About a month |
|---------|------|---------------|
| Daily | `daily` | 25 USD |
| Twice a month | `0 6 1,15 * *` | 1.6 USD |

**Why I do not answer it:** the cadence sets recurring spend. No ADR, design section, or digest entry records a cadence for this project, so the choice is yours alone. T-3037 sets the spec you choose, then T-3040 closes I-0292.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/system/strategic-agents.md § The analysis schedule
```

Source: TH-1065.
