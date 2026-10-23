---
id: D-0317
title: "S-0426 needs a cost of delay input: recommend time lost 1h per cycle (150 USD a week): Recommendation: set `time_lost_per_cycle: 1h` on S-0426 (150 USD a week), as the planner recommends: `flai edit…"
kind: answer
sources: [TH-0475]
topics: [cli, all]
paths: [design/system/agent-memory.md, design/system/agent-context.md, design/system/README.md]
date: 2026-10-16
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0317 S-0426 needs a cost of delay input: recommend time lost 1h per cycle (150 USD a week): Recommendation: set `time_lost_per_cycle: 1h` on S-0426 (150 USD a week), as the planner recommends: `flai edit…

The operator, alex, answered TH-0475 at 2026-10-16T19:31:40Z.

What was asked:

```text
S-0426 has no cost of delay inputs and no epic to take a share from, so `flai cod S-0426` refuses. The inputs are yours.
```

The answer:

```text
**Recommendation:** set `time_lost_per_cycle: 1h` on S-0426 (150 USD a week), as the planner recommends: `flai edit S-0426 --time-lost-per-cycle 1h`. The planner then writes the value on its next run.

**Why you:** S-0426 is in progress. flai lets the orchestrator set cost of delay inputs only on a backlog story, as it refused for S-0427 on TH-0474, so this input is yours.

**The alternative:** take `10m` (25 USD a week) instead if you weigh this research as you weighed S-0413. Both are research stories that recommend without building anything.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
```

Source: TH-0475.
