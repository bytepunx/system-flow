---
id: D-0660
title: "Accept ADR-0305: the metrics are readings of the workflow's timers: Recommended: accept ADR-0305 as written. agent-S-0750 then accepts it with `flai adr accept` and writes `design/system/metrics.md`…"
kind: answer
sources: [TH-1223]
topics: [cli, analysis, template, flai, go, code, all]
paths: [design/adrs/README.md, design/system/metrics.md, flai/internal/workflow/timers.go, flai/internal/workflow/timers_test.go, flai/internal/workflow/baseline/board.md, flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go, flai/internal/metrics/claims.go, flai/internal/metrics/waiting.go, flai/internal/metrics/coordination.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/costorder.go, flai/internal/metrics/forecast.go, flai/internal/metrics/digest.go, flai/internal/metrics/usage.go, flai/internal/metrics/strategic.go, flai/internal/metrics/spend.go, design/adrs/0305-lead-cycle-queue-and-review-time-and-time-in-state-are-readings-of-the-workflow.md, flai/internal/metrics/claims_test.go, flai/internal/metrics/coordination_test.go, flai/internal/metrics/costofdelay_test.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/concurrency.go, flai/internal/metrics/reopen_test.go, flai/internal/metrics/stats_golden_test.go, flai/internal/metrics/testdata/good.stats.golden.json]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0660 Accept ADR-0305: the metrics are readings of the workflow's timers: Recommended: accept ADR-0305 as written. agent-S-0750 then accepts it with `flai adr accept` and writes `design/system/metrics.md`…

The operator, alex, answered TH-1223 at 2026-11-17T07:46:10Z.

What was asked:

```text
**Recommendation: accept ADR-0305 as written** (`design/adrs/0305-lead-cycle-queue-and-review-time-and-time-in-state-are-readings-of-the-workflow.md`, on `story/S-0750`). Criterion 1 needs an accepted ADR, and `metrics.md` changes only with one.
```

The answer:

```text
**Recommended:** accept ADR-0305 as written. agent-S-0750 then accepts it with `flai adr accept` and writes `design/system/metrics.md` (T-3242).

**Why:**

- ADR-0301 §5 asks for this ADR. It says S-0750 "writes the ADR that changes the metrics contract". It also fixes the timers' forms, points, and the baseline's `lead`, `cycle`, `queue`, and `review`. ADR-0305 §2, §3, and §5 restate these unchanged.
- The five rules agent-S-0750 lists are new, but each one keeps a figure `flai stats` gives today. Criterion 4 asks for exactly that: the same bytes on `testdata/good` and `windows.golden.json`.
- Rule 1 (no last interval for a closed item) is how time in state works today. Without it, time in `done` grows without end.
- Rule 2 (a bounded sum has no value when the item entered its states only after the point) keeps ADR-0190's queue time. That rule gives no value to an item that started before it was ready.
- Rules 3 and 4 restate ADR-0055 and ADR-0301 §5's `counts: accepted` in terms of kinds.
- Rule 5 keeps the JSON field names, so flaiover needs no change for the baseline.

**Why you:** accepting an ADR records a decision, and `metrics.md` changes only with one. That call is yours.

**Alternative:** accept with changes, naming the rule. A change to rule 1 or 2 changes past figures, which criterion 4 forbids on the fixtures.

Source: design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md § 5. Timers
```

Source: TH-1223.
