---
id: TH-0018
title: Which coordination design should flai build?
anchor:
  path: wip/kanban/stories/S-0124-survey-state-of-the-art-coordination-mechanisms.md
  item: S-0124
status: open
participants: [agent-S-0124]
created: 2026-09-26T07:24:35Z
updated: 2026-09-26T07:24:35Z
---

# TH-0018 Which coordination design should flai build?

On wip/kanban/stories/S-0124-survey-state-of-the-art-coordination-mechanisms.md.

## Entries

### 2026-09-26T07:24:35Z agent-S-0124
The survey is in design/system/agent-coordination.md, on story/S-0124. Today nothing holds a ready story back except the in-progress limit. `touches` only produces a `wip.overlap` warning, and a conflict between two open story branches shows up only after one of them is accepted.

My recommendation is design 1, built in this order:
1. **Hold at pull.** A ready story whose `touches` overlap an in-progress story's is held. Overlap uses the same `/`-prefix rule `flai check` already uses: `flai/cmd` overlaps `flai/cmd/serve`, but `flai/cmd` does not overlap `flaiover`. While a story is held, flai serve does not start its agent and wait_for_work does not offer it. The card turns yellow with the reason and what clears it, for example: "held: touches flai/cmd/serve, inside flai/cmd which S-0128 (in progress) touches; starts when S-0128 is accepted or cancelled".
2. **Explicit dependency.** An `after: [S-nnnn]` field holds a story until those stories are done. This covers dependencies that are not about files.
3. **Safety net at sync.** `flai stream sync` trial-merges the branch with every other open story branch using `git merge-tree`, which touches no worktree. It reports a conflict to both stories before either is accepted. It also says when the paths a branch actually changed have drifted outside its `touches`.
4. **Notice at accept.** Accepting a story tells every open story that overlaps it what changed.

The alternatives, briefly:
- **Design 2** is design 1 plus widening from git co-change history and component edges in the manifest. These show amber: the story starts, but it is flagged. I would do this later.
- **Design 3** is optimistic only: never hold, just trial-merge and report. It is cheapest, but it wastes agent work, and it gives you nothing to see before work starts.
- **Design 4** is design 1 plus a guard that refuses commits outside a story's claim. It is the strongest guarantee, but agents get stuck mid-story. The research (Cursor, MCP Agent Mail) backed away from hard locks for this reason.

This refines ADR-0019, which rejected locks because "a lock nobody releases blocks work". A hold here is released by acceptance, cancellation, or a return to backlog. It is visible with its reason, and you can override it with Start agent.

Do you want design 1 in that order? Or a different design, or a different order?
