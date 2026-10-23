---
id: D-0214
title: "Commit or drop the uncommitted I-0078 bump in S-0212's worktree: Recommended: (a), drop the two edits"
kind: answer
sources: [TH-0249]
topics: [planning, orchestration, analysis, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/routes/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, design/system/metrics.md, flaiover/src/lib/components/ForecastTable.svelte, flaiover/src/lib/components/ForecastTable.svelte.test.ts, design/adrs/0111-flai-stats-gives-each-item-the-model-its-forecasts-are-grouped-under.md, design/adrs/README.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/summary.md, flai/internal/metrics/forecast.go, flai/internal/metrics/forecast_test.go, flai/internal/metrics/metrics.go, design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0214 Commit or drop the uncommitted I-0078 bump in S-0212's worktree: Recommended: (a), drop the two edits

The operator, alex, answered TH-0249 at 2026-10-07T08:07:07Z.

What was asked:

```text
The only uncommitted change in S-0212's worktree is a bump of I-0078 (count 37 to 38, a 2026-10-07T07:33:15Z occurrence for S-0212, S-0250 is cancelled; run flai archive) and the matching row in design/issues/summary.md. A close-out run made it after commit 94d8b7bb, which took S-0212's earlier bump of the same occurrence back out so S-0212 and S-0271 merge cleanly: story/S-0271 already moves I-0078 to 38 with its own 07:26:28Z entry for the same warning. If I commit this bump, the two branches edit the same lines of I-0078 and summary.md again and will conflict at acceptance. If I drop it…
```

The answer:

```text
Recommended: (a), drop the two edits.

- They came from my verifier's close-out run at 94d8b7bb, not from yours. Your branch had already chosen to leave I-0078 as main has it.
- ADR-0085 has a close-out commit the issue files it writes. Read alone, that favours (b). But main already records this same S-0250 archive warning in I-0078, through S-0271, now accepted. Committing it again on story/S-0212 would recreate the conflict TH-0247 found, for an occurrence main already counts.
- After dropping them, run `flai stream sync S-0212` to take S-0271's merge in.

This is a recommendation, because no recorded decision settles a duplicate occurrence. It awaits the operator's confirmation.

Source: design/adrs/0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md § Decision
```

Source: TH-0249.
