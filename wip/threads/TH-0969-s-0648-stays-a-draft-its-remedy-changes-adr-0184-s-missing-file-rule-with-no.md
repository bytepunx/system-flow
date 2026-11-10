---
id: TH-0969
title: "S-0648 stays a draft: its remedy changes ADR-0184's missing-file rule with no ADR in its touches; confirm the remedy and add the ADR"
anchor:
  path: wip/kanban/stories/S-0648-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md
  item: S-0648
status: open
participants: [orchestrator]
created: 2026-11-10T06:38:17Z
updated: 2026-11-10T06:38:17Z
---

# TH-0969 S-0648 stays a draft: its remedy changes ADR-0184's missing-file rule with no ADR in its touches; confirm the remedy and add the ADR

On wip/kanban/stories/S-0648-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md.

## Entries

### 2026-11-10T06:38:17Z orchestrator
**Recommended:** confirm the planner's remedy, and have the planner add an ADR refining ADR-0184 to S-0648's plan: `design/adrs/<new>.md` and `design/adrs/README.md` in T-2827's touches and the story's. Then I finalize S-0648.

**Why it is not finalized.** `flai promote --drafts` lists S-0648 as complete (position 1, 1h, throughput), and its four tasks, touches, 1h forecast, and 5 USD a week value fit each other. Two things stop me:

1. **The remedy reverses a recorded rule with no ADR.** ADR-0184 § Decision says: "A missing file, a file that does not parse, or an item made in memory is written as before." T-2827 makes `Repo.Save` refuse to recreate an item file that was there when it was read, and makes `Archive` take the items lock. That changes a tool-enforced rule. `design/conventions/decisions.md` asks for an ADR for that. S-0648's touches name `design/system/work-hierarchy.md` but no ADR.
2. **The goal asks for the remedy to be proposed before it is built.** The planner proposed it in § Proposed fix, and nobody has confirmed it. That is a decision not yet recorded, as for S-0637 on TH-0960, and as you answered for S-0468 on TH-0599 (D-0375).

**The remedy reads as sound.** Every I-0323 instance is the stray open copy of T-2760. It differs from the archived copy only in `usage`, which fits a `Measure` save that lands after `Archive`'s unlocked rename. T-2828 keeps the harness's figure by writing it to the archived copy.

**Alternative:** finalize as it stands, and let the story's agent write the ADR when T-2827 meets ADR-0184. I don't recommend it: the plan does not name the ADR, so nothing makes the agent write it.

I have not changed S-0648.
