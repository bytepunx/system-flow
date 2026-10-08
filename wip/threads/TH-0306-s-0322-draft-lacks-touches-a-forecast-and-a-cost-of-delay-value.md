---
id: TH-0306
title: S-0322 draft lacks touches, a forecast, and a cost of delay value
anchor:
  path: wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md
  item: S-0322
status: resolved
participants: [orchestrator, alex, planner-S-0322]
created: 2026-10-07T19:08:05Z
updated: 2026-10-08T04:36:39Z
---

# TH-0306 S-0322 draft lacks touches, a forecast, and a cost of delay value

On wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md.

## Entries

### 2026-10-07T19:08:05Z orchestrator
Recommendation: cancel S-0322 as a duplicate of S-0312, or merge the two, then plan the one you keep.

S-0322 ("flai task done commits every uncommitted file in the story worktree…") describes the same fault as S-0312 ("flai task done commits everything in the worktree…"). TH-0288 is already open on S-0312.

I am leaving S-0322 as a draft. `flai promote --drafts` lists it as incomplete: no touches, no forecast duration, no forecast delivery, no cost of delay value.

### 2026-10-08T04:34:41Z alex
Resolved.

### 2026-10-08T04:36:28Z planner-S-0322
Settled: the operator cancelled S-0312 as the duplicate and kept S-0322, adding I-0104's close as criterion 2. S-0322 is replanned with S-0312's fix: 17 file touches, a 50m forecast, a 12.50 USD a week value, and tasks T-1330, T-1328, T-1331, and T-1329. The plan is on its own thread on S-0322.

### 2026-10-08T04:36:39Z orchestrator
Resolved: Settled: S-0312 cancelled as the duplicate; S-0322 replanned with the fix (TH-0358) and finalized by the orchestrator
