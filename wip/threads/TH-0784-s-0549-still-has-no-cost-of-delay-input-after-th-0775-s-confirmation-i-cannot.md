---
id: TH-0784
title: "S-0549 still has no cost of delay input after TH-0775's confirmation; I cannot set it while set_cost_of_delay is off"
anchor:
  path: wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md
  item: S-0549
status: open
participants: [orchestrator]
created: 2026-10-31T09:51:16Z
updated: 2026-10-31T09:51:16Z
---

# TH-0784 S-0549 still has no cost of delay input after TH-0775's confirmation; I cannot set it while set_cost_of_delay is off

On wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md.

## Entries

### 2026-10-31T09:51:16Z orchestrator
**Recommended:** run this in the main checkout, or from the dashboard:

```bash
flai edit S-0549 --time-lost-per-cycle 20m
```

Or turn on `orchestration.permissions.set_cost_of_delay`, and I set it the next time I am started.

**Why:**

- You confirmed the 20m input on TH-0775 and resolved the thread at 2026-10-31T09:50:28Z.
- A confirmation does not write the input. S-0549 still has no cost of delay.
- `set_cost_of_delay` is not in `system-flow.yaml`, so it is off. I am not allowed to write the input you decided.
- Until the input is set, `flai cod S-0549` refuses, and S-0549 has no value.
