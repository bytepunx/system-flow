---
id: TH-0855
title: "S-0585 stays in review: flai blocks acceptance while criterion 3 is unticked"
anchor:
  path: wip/kanban/stories/S-0585-publish-what-is-accepted-and-not-yet-released.md
  item: S-0585
status: open
participants: [orchestrator]
created: 2026-11-07T08:18:05Z
updated: 2026-11-07T08:18:05Z
---

# TH-0855 S-0585 stays in review: flai blocks acceptance while criterion 3 is unticked

On wip/kanban/stories/S-0585-publish-what-is-accepted-and-not-yet-released.md.

## Entries

### 2026-11-07T08:18:05Z orchestrator
**Recommended:** reword or remove criterion 3 of S-0585 so the stop at step 2 counts as the outcome. For example, change it to "Publish what is accepted and not yet released, or record the evaluation that stopped the run". Then tick it with `flai criteria tick S-0585 3`, and I will accept S-0585 on my next run. You can also accept it yourself with `flai accept S-0585 --by alex`.

**Why I have not accepted it:** on TH-0854 you confirmed that S-0585 should be accepted with the stop recorded. But `flai accept S-0585 --by orchestrator --verified ff4e47115391 --dry-run` refuses it:

```text
blocked: S-0585 has unchecked acceptance criteria
blocked: acceptance criterion 3 of S-0585 is unticked: Publish what is accepted and not yet released
```

Everything else is clean. The branch head is ff4e47115391, the same as main. The close-out passed at that commit. T-2546 and T-2548 are done, and T-2549 is cancelled. The story's agent left criterion 3 unticked on purpose, because no publish ran. I do not tick criteria myself, and changing a criterion's words changes scope, so this is yours to decide.

**Alternative:** leave it as it is until S-0562's T-2506 settles how an instance stopped at step 2 ends. That may reword the formula's criterion for later instances.
