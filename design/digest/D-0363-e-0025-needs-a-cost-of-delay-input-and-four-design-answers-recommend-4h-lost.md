---
id: D-0363
title: "E-0025 needs a cost of delay input and four design answers: recommend 4h lost per cycle, one file per entry under design/digest, a curator role of its own, a 30-day compaction window: Recommended…"
kind: answer
sources: [TH-0476]
topics: [cli, planning, analysis, conventions, flai, go, template, code, all]
paths: [flai/internal/context, flai/internal/serve, flai/cmd, design/system/agent-context.md, design/system/strategic-agents.md, design/system/metrics.md, design/conventions/agent-narrative.md, docs/users/flai.md, docs/operators/index.md, template]
date: 2026-10-21
status: absorbed
absorbed_by: "design/adrs/0191-the-digest-is-one-markdown-file-per-entry-under-the-design-folder-s-digest.md#decision"
asked_again: []
cited: [S-0431, S-0432]
---

# D-0363 E-0025 needs a cost of delay input and four design answers: recommend 4h lost per cycle, one file per entry under design/digest, a curator role of its own, a 30-day compaction window: Recommended…

The operator, alex, answered TH-0476 at 2026-10-21T12:59:17Z.

What was asked:

```text
I am drafting E-0025's stories now. Five inputs are yours. Each has my recommendation first. I draft on the recommended answers and change the stories if you answer otherwise.
```

The answer:

```text
Recommended: "yes to all five", as planner-E-0025's entry of 2026-10-21T05:57Z puts it.

1. Set `time_lost_per_cycle: 4h` on E-0025. At 150 USD an hour over a 168h cycle, that is 600 USD a week, shared across S-0428 to S-0435 by forecast. This is the same figure you chose for E-0024 on TH-0471. Without inputs on the epic, `flai cod` refuses all eight stories, and none of them can be promoted.
2. Keep one file per entry under `design/digest`. Your TH-0471 answer moves roles and formulas to `design/workflows`, but it does not touch the digest.
3. Build the curator now as a role of its own on the analyzer machinery. E-0024 is back in the backlog, so waiting for S-0419 would hold S-0431 to S-0435. If E-0024 lands first, S-0419 moves the curator onto the role path (TH-0473).
4. and 5. Take the planner's recommendations as written into the drafts.

Confirming this sets no input by itself: the epic's input is yours to set, with `flai edit E-0025` and its cost of delay flags.

Source: design/system/strategic-agents.md § Cost of delay
```

Source: TH-0476.
