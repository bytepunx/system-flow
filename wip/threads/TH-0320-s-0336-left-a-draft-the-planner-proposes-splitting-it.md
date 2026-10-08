---
id: TH-0320
title: "S-0336 left a draft: the planner proposes splitting it"
anchor:
  path: wip/kanban/stories/S-0336-the-dashboard-shows-the-messages-between-stories-agents-on-each-story-s-page-and-in-a-messages-view.md
  item: S-0336
status: open
participants: [orchestrator, alex]
created: 2026-10-07T20:24:26Z
updated: 2026-10-08T04:29:27Z
---

# TH-0320 S-0336 left a draft: the planner proposes splitting it

On wip/kanban/stories/S-0336-the-dashboard-shows-the-messages-between-stories-agents-on-each-story-s-page-and-in-a-messages-view.md.

## Entries

### 2026-10-07T20:24:26Z orchestrator
Recommendation: split S-0336 as the planner proposes on TH-0318. One story would hold the reads, the Messages view, and the story page, after S-0330 only. The other would hold the held-card part, after S-0334. You would see messages on the dashboard three stories sooner. Then finalize both.

`flai promote --drafts` lists S-0336 as complete: touches, a 48m forecast, and a 124.03 USD/week value. I am leaving it a draft because the open split proposal would rewrite its words, and a finalized story's words are yours. It cannot be pulled before S-0334 in any case.

If you would rather keep it whole, finalize it as it is.

### 2026-10-08T04:29:10Z alex
take the recommendation for the split

### 2026-10-08T04:29:27Z orchestrator
Follow-up: I cannot make the split myself. It needs a new story, which only the planner or you can create, and `flai plan --candidates` no longer lists S-0336 or E-0018, so I cannot start the planner on them. To split it, run `flai plan E-0018` (or Plan on the dashboard) and ask for the split TH-0318 proposed: the reads, the Messages view, and the story page after S-0330 only; the held-card part after S-0334. When the planner's thread comes back I will approve the plan and finalize both drafts. Until then S-0336 stays a draft.
