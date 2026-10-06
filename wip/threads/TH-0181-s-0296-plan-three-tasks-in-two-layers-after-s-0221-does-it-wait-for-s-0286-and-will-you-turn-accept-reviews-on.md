---
id: TH-0181
title: "S-0296 plan: three tasks in two layers after S-0221; does it wait for S-0286, and will you turn accept_reviews on?"
anchor:
  path: wip/kanban/stories/S-0296-a-story-waits-in-review-for-the-operator-while-nothing-else-can-start-so-the-board-stands-idle-until-a-person-accepts.md
  item: S-0296
status: open
participants: [planner-S-0296]
created: 2026-10-06T11:51:00Z
updated: 2026-10-06T11:51:00Z
---

# TH-0181 S-0296 plan: three tasks in two layers after S-0221; does it wait for S-0286, and will you turn accept_reviews on?

On wip/kanban/stories/S-0296-a-story-waits-in-review-for-the-operator-while-nothing-else-can-start-so-the-board-stands-idle-until-a-person-accepts.md.

## Entries

### 2026-10-06T11:51:00Z planner-S-0296
The fix for I-0088 is S-0221's: the orchestrator accepts a story in review under `orchestration.permissions.accept_reviews`. I made S-0296 wait for S-0221 (`after`). It is the story that shows the fix works and closes the issue.

Tasks and layers:

- Layer 1, run together:
  - T-1020: a `flai/internal/serve` test reproduces the idle board. A story is in review and an overlapping ready story is held. With `accept_reviews` off, nothing starts. With it on, a stand-in orchestrator accepts and the launcher starts the held story. `orchestrate.go` changes only if the review is noticed late, on the one-minute look.
  - T-1021: `workflow.md` and the operators' settings and index pages say that `accept_reviews` ends the wait, and which stories still wait.
- Layer 2: T-1022 closes I-0088 with `flai issue close`. It waits for T-1020 and T-1021 because its reason cites them.

Forecast 45m (flai gave 18m; raised for the serve test under the race detector). Cost of delay 420 USD/week, from flai's input, kept.

Assumptions:

- No new command or flag. S-0221 documents `flai accept --by orchestrator`.
- "The ones the operator chooses to review" (from the issue) needs no new mechanism in this story. You keep a story by leaving `accept_reviews` off, or by an open thread on it, which S-0221 makes a blocker.

Questions for you, each with my recommended answer first:

1. Should S-0296 also wait for S-0286? Recommended: yes. Without it, turning `accept_reviews` on lets the orchestrator accept a story that changes `.claude/`, which S-0286 keeps for you. The cost is that S-0286 is still a draft, so S-0296 starts later. If you say no, I leave `after` at S-0221 alone and T-1022's reason says S-0286 is still to come.
2. Will you turn `orchestrate` and `orchestration.permissions.accept_reviews` on for this project? Recommended: yes, once S-0221, and S-0286 if you agree to (1), are accepted. It is a hand edit of `system-flow.yaml`'s `orchestration` block, which agents may not make, and there is none in main yet. T-1022 closes I-0088 either way, and says in its reason whether the permission is on here.
3. Or, instead of all this: should I propose cancelling S-0296 and having S-0221 close I-0088? Recommended: no. S-0221 is almost done and has no test of the idle board, so S-0296 is where that test belongs.
