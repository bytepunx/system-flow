---
id: T-2255
type: task
nature: remediation
title: "An ADR refining ADR-0054 and ADR-0057 settles that metrics.md § Charts says what the dashboard draws, with a row for every chart it draws"
status: done
parent: S-0527
owner: alex
created: 2026-10-28T15:00:28Z
updated: 2026-10-30T00:12:45Z
transitions:
  - to: ready
    at: 2026-10-30T00:12:06Z
    by: agent-S-0527
  - to: in-progress
    at: 2026-10-30T00:12:07Z
    by: agent-S-0527
  - to: done
    at: 2026-10-30T00:12:45Z
    by: agent-S-0527
stream: S-0527
tags: [design, adr]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 38
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 2798
      cache_read: 1005846
      cache_write: 30903
      cost: 0.5044
---
# T-2255 An ADR refining ADR-0054 and ADR-0057 settles that metrics.md § Charts says what the dashboard draws, with a row for every chart it draws

## Work

Write the ADR with `flai adr new`, with the topics `dashboard` and `analysis`, and accept it with `flai adr accept`. It waits for no task, and the later tasks cite its number.

Context: I-0193's instance, and the planner's `### Planning` note on S-0527. ADR-0054 says `design/system/metrics.md` defines each chart's data and changes only with an ADR; ADR-0057 decided the dashed all-models line for Avg. Time / Model and Avg. Cost / Model alone. The charts in `flaiover/src/lib/viz/charts.ts`, their explainers in `flaiover/src/lib/viz/explain.ts` (S-0485), `docs/users/flaiover.md`, and `design/system/flaiover-dashboard.md` agree with each other; only § Charts differs. Read `cycleTime`, `burnUp`, `perModel`, `strategicRows`, and `strategicUseByDay` in `charts.ts` before writing, and the `controls` function, which offers Cycle Time an epic and no nature.

Decision to record: § Charts says what the dashboard draws, and the dashboard follows it from then on.

- Cycle Time: one point per item completed in the window, chosen per epic or whole, coloured by nature, with the p50 and p85 lines; no filter by nature.
- Burn-up: the scope line and the done line per day of the window, per epic or whole; no forecast line. A forecast line from throughput is not built; if wanted, it is a feature story of its own.
- Every chart per model, Tokens / Min and Tokens / $ included, adds a dashed line for all models with two models or more, as ADR-0057 decided for Avg. Time / Model and Avg. Cost / Model.
- The table gains a row for each chart it lacks: Agent Waiting, Strategic Cost, and Strategic Use. Strategic Use's mean waiting is the dashboard's: the mean over the items of the report's type completed that day of `wait_threads_seconds` plus `wait_review_seconds`, a wait an item lacks counted as zero, worked out from `items[]` so that `flai stats` and the contract stay unchanged (TH-0140); say so in § Strategic use per day too.
- Every chart kind the dashboard draws has a row in § Charts under its title, and a test in `flaiover/src/lib/viz/explain.test.ts` holds it, so a chart added without its row fails the flaiover tier.

Record the consequences: the explainers stay as S-0485 wrote them; no `flai stats` field changes; a later change to what a chart draws changes § Charts in the same story, by an ADR as ADR-0054 says.

## Done when

- The ADR exists under `design/adrs/`, status accepted, topics `dashboard` and `analysis`, refining ADR-0054 and ADR-0057, and `design/adrs/README.md` indexes it.
- Its decision names the five points above and says a Burn-up forecast line is a feature story of its own, not this remediation's.
- `flai check --strict` and the markdown lint are clean on the ADR and the index.

## Notes

Drafted by planner-S-0527. The ADR's file cannot be named in the touches until `flai adr new` allots its number; `flai task done` widens them with the committed path (ADR-0144).
