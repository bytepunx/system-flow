---
id: TH-0790
title: "S-0241 is back in flai promote --candidates (rank 2): make it a draft so I do not promote it again before forks are offered headless"
anchor:
  path: wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md
  item: S-0241
status: resolved
participants: [orchestrator]
created: 2026-11-01T23:46:04Z
updated: 2026-11-10T09:21:21Z
---

# TH-0790 S-0241 is back in flai promote --candidates (rank 2): make it a draft so I do not promote it again before forks are offered headless

On wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md.

## Entries

### 2026-11-01T23:46:04Z orchestrator
**Recommended:** make S-0241 a draft again, so that `flai promote --candidates` stops listing it until a headless session offers `subagent_type: fork`:

```bash
flai edit S-0241 --draft
```

**Why:**

- You moved S-0241 from ready to backlog at 2026-11-01T23:44:38Z. I promoted it to ready on 2026-10-30, because `flai promote --candidates` listed it.
- Its Notes say: "Do not move it to ready until a headless session offers forks: it cannot be worked before then." It has also been blocked since 2026-10-08, with the reason "No candidate solutions yet".
- `flai promote --candidates` lists it again now, at position 2 (2h, throughput). It does not weigh the Notes or the open block. Under `promote_to_ready`, the next run of mine started for the ready column would promote it again.
- A draft is never a candidate. Finalizing a draft is yours, so you decide when forks are offered and it can go to ready.

**Alternative:** leave it as it is, and I will skip S-0241 in every promotion while its block stays open. The rule would live only in my activity log, so a later run could miss it.

I have not promoted it, and I have not changed it.

### 2026-11-10T09:21:21Z orchestrator
Resolved: The operator made S-0241 a draft at 2026-11-10T09:19:43Z, as recommended. I leave it a draft: flai promote --drafts lists it complete, but its Notes say it cannot be worked until a headless session offers subagent_type fork, and it is blocked since 2026-10-08. Finalizing it is the operator's once forks are offered.
