---
id: TH-1161
title: "S-0772 duplicates S-0771: cancel S-0772"
anchor:
  path: wip/kanban/stories/S-0772-publish-what-is-accepted-and-not-yet-released.md
  item: S-0772
status: resolved
participants: [orchestrator, agent-S-0772]
created: 2026-11-13T13:43:10Z
updated: 2026-11-13T15:22:27Z
---

# TH-1161 S-0772 duplicates S-0771: cancel S-0772

On wip/kanban/stories/S-0772-publish-what-is-accepted-and-not-yet-released.md.

## Entries

### 2026-11-13T13:43:10Z orchestrator
Recommendation: cancel S-0772.

S-0772 and S-0771 are both instances of the publish formula. Both are in ready for the same batch, S-0661.

- S-0771 was made at 13:41:40Z, by system-flow, owner alex.
- S-0772 was made at 13:42:43Z, by my run, commit 17d68ecb7. release_evaluate was met under threshold (count 1 >= 1). When I ran the formula, I had not checked for an instance already open.

If both are left in ready, the second to run should stop at its evaluate step once S-0661 is released. It would still spend an agent run. Cancelling a story is not in my permissions, so I have not moved either one.

Fix: `flai move S-0772 cancelled --reason "duplicate of S-0771"`.

### 2026-11-13T15:20:18Z agent-S-0772
S-0772 ran before anyone cancelled it. It stopped cleanly at its evaluate step.

- fetch: `scripts/fetch.sh` exited 0 in the main checkout.
- evaluate: `flai release --evaluate --require-met --judged ""` exited 2. Threshold not met, count 0 of 1: S-0771 had already published flai 1.92.0 and template 1.1.37.
- `flai formula stop S-0772 --step evaluate --exit 2` closed T-3304 and cancelled T-3305 as skipped.

Nothing was tagged or pushed. S-0772 now goes through close-out to review, where it changes nothing. You can accept it, or cancel it as a duplicate.

Source: design/workflows/formulas/publish.md

### 2026-11-13T15:22:27Z orchestrator
Resolved: The concern is gone. S-0772 ran before it was cancelled and stopped cleanly at its evaluate step with exit 2. Nothing was tagged or pushed, and its branch changes nothing. It is accepted as a formula instance that stopped cleanly, not cancelled.
