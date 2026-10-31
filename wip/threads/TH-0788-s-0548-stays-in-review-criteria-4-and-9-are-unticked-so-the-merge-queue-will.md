---
id: TH-0788
title: "S-0548 stays in review: criteria 4 and 9 are unticked, so the merge queue will not take it"
anchor:
  path: wip/kanban/stories/S-0548-a-queue-verification-takes-under-five-minutes-the-slow-and-network-bound-tests.md
  item: S-0548
status: open
participants: [orchestrator]
created: 2026-10-31T10:00:25Z
updated: 2026-10-31T10:00:25Z
---

# TH-0788 S-0548 stays in review: criteria 4 and 9 are unticked, so the merge queue will not take it

On wip/kanban/stories/S-0548-a-queue-verification-takes-under-five-minutes-the-slow-and-network-bound-tests.md.

## Entries

### 2026-10-31T10:00:25Z orchestrator
**Recommended:** reword criteria 4 and 9 to what S-0548 delivered, and move what is left to follow-up work. Then tick them, and I enqueue S-0548 on my next run. The alternative is to send it back to in-progress until both are met.

**Why I did not enqueue it:**

- `flai verify` passed at the branch head, 797b165, in 4m53s. Integration took 3m33s and smoke 1m08s.
- The queue lists S-0548 as ineligible: `criteria_unticked`, criteria 4 and 9. `queue_add` refuses any story with a criterion unticked.
- **Criterion 4 is not met for `mcpserver`.** It took 136 to 162 s under race, and the target is under two minutes. `cmd` meets it at 75 to 87 s. The agent recorded the remainder as I-0234: about 70 s of serial tests that set environment variables the server reads when it starts.
- **Criterion 9 cannot be checked before landing.** The agent measured 4m40s for integration and smoke together, but at load 64 to 74, not "when nothing else runs". The fortnight's p50 in `flai stats` only exists after landing.

**What would let it land:**

1. Criterion 4: limit the two-minute bound to `cmd`, and leave `mcpserver` to a story on I-0234. Or accept the 136 to 162 s measured under load.
2. Criterion 9: keep the close-out's 4m40s as the measurement, and move the fortnight p50 check to a later review.
3. Tick both with `flai criteria tick S-0548 4 9`.

Rewording the criteria changes the story's scope, and that decision is yours.
