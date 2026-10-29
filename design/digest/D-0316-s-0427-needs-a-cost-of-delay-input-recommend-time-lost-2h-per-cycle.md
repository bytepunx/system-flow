---
id: D-0316
title: "S-0427 needs a cost of delay input: recommend time lost 2h per cycle: Recommendation: set `time_lost_per_cycle: 2h` on S-0427, as the planner recommends: `flai edit S-0427 --time-lost-per-cycle 2h`"
kind: answer
sources: [TH-0474]
topics: [server-side, all]
paths: [design/experiments/S-0427-sqlite-and-postgres-read-models.md, design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md, design/issues/summary.md]
date: 2026-10-16
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0756]
---

# D-0316 S-0427 needs a cost of delay input: recommend time lost 2h per cycle: Recommendation: set `time_lost_per_cycle: 2h` on S-0427, as the planner recommends: `flai edit S-0427 --time-lost-per-cycle 2h`

The operator, alex, answered TH-0474 at 2026-10-16T06:49:19Z.

What was asked:

```text
S-0427 has no cost of delay inputs and no epic to take a share from, so `flai cod S-0427` cannot work out a value.
```

The answer:

```text
**Recommendation:** set `time_lost_per_cycle: 2h` on S-0427, as the planner recommends: `flai edit S-0427 --time-lost-per-cycle 2h`.

## Why 2h

- **The goal asks for it.** It measures the "operator wait time and agent overhead" that file reads cost on each board, chart, and prime call. That is time lost, not revenue or a penalty.
- **Nothing in the story points to an alternative.** The story does not say the board feels fast enough, which would argue for 1h. It does not price agent time as a penalty either, which would argue for 250 USD.

## Why this needs you

I tried to set it, and flai refused: "S-0427 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0427 to ready, so its input is yours to set.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
```

Source: TH-0474.
