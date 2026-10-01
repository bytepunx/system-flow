---
id: TH-0056
title: Should wip/ warnings outside a story stop its close-out?
anchor:
  path: design/conventions/work-management.md
status: open
participants: [agent-S-0187]
created: 2026-10-01T10:01:41Z
updated: 2026-10-01T10:01:41Z
---

# TH-0056 Should wip/ warnings outside a story stop its close-out?

On design/conventions/work-management.md.

## Entries

### 2026-10-01T10:01:41Z agent-S-0187
From S-0187. `scripts/close-out.sh` runs `flai check --strict`, which reads the main checkout's `wip/` and fails on warnings the story's agent may not fix: `threads.archived` (TH-0032 today, which holds an open request of yours), `item.archive` for done epics, another story's narrative lint. The new rule in `work-management.md` says to move to review only when the close-out ends clean, so the agent either waits on you or finishes the steps by hand. S-0187 went to review with only TH-0032's warning left, and its narrative says why.

Recommended: a story that gives `flai check` a story scope for the close-out (`flai check --strict --story S-nnnn`, or a flag that reports `wip/` findings outside the story's items as notes rather than failures), with `close-out.sh` using it. Alternative: the work-management rule names the exception: a close-out that stops only on findings outside the story may go to review with them recorded in the narrative. Want the story written?
