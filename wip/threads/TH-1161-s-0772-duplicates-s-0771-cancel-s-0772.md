---
id: TH-1161
title: "S-0772 duplicates S-0771: cancel S-0772"
anchor:
  path: wip/kanban/stories/S-0772-publish-what-is-accepted-and-not-yet-released.md
  item: S-0772
status: open
participants: [orchestrator]
created: 2026-11-13T13:43:10Z
updated: 2026-11-13T13:43:10Z
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
