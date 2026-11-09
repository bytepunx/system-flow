---
id: D-0542
title: A 3h story that fixes one cause across flai and flaiover is kept whole, not split into a flai half and a dashboard half, when its task layers already order the halves
kind: decision
sources: [S-0611, TH-0903, TH-0904]
topics: [planning, orchestration, flai, flaiover]
paths: [design/conventions/work-management.md]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0542 A 3h story that fixes one cause across flai and flaiover is kept whole, not split into a flai half and a dashboard half, when its task layers already order the halves

planner-E-0033 proposed on TH-0903 (proposal 2) to split S-0611, a 3h story over flai and flaiover, into a flai half (T-2689, T-2690) and a dashboard half (T-2691 to T-2693), the second after the first. The orchestrator recommended against it, citing `design/conventions/work-management.md`:

```text
Split S-0611: no. `work-management.md` sizes a story as one demonstrable increment for one to a few sessions. 3h fits that. S-0611 fixes one cause (21), and criterion 4 measures both answers once. Its layers already run T-2690 before T-2691 reads the `finished` parameter. A split would add an `after` and a second close-out.
```

The operator did not reply on the thread. The operator finalized S-0611 whole and moved it to ready at 2026-11-09T16:59:36Z, and TH-0904 was resolved on that. agent-S-0611 recorded it in its narrative's `## Decisions` (`wip/archive/agents/S-0611.md`): "The operator finalized S-0611 whole, so the split planner-E-0033 proposed on TH-0903 is not made (TH-0904)." The story was then worked as one, with the flai tasks in the first layer and the dashboard task that reads the parameter after them.

A component split alone does not justify a second story. Every close-out runs the `integration` and `flaiover` tiers, which average about 238 s and 176 s, and the task layers already put the dependency in order.

Sources: TH-0903, TH-0904, S-0611.
