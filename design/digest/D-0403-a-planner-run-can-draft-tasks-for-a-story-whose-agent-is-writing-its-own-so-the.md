---
id: D-0403
title: A planner run can draft tasks for a story whose agent is writing its own, so the story's agent cancels the duplicates on the plan thread before its layers run
kind: decision
sources: [S-0413, S-0473, S-0479, S-0496, TH-0472, TH-0574, TH-0577]
topics: [planning, orchestration]
paths: [design/system/strategic-agents.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0403 A planner run can draft tasks for a story whose agent is writing its own, so the story's agent cancels the duplicates on the plan thread before its layers run

Four stories met planner drafts that duplicated tasks their own agent had written or finished:

- S-0413 cancelled five (TH-0472).
- S-0473 cancelled T-1986 and T-1989, which duplicated T-1992, already done (TH-0574).
- S-0479 cancelled four drafts that its T-2007 covered (TH-0577).
- S-0496 cancelled T-2113 to T-2115, which duplicated its T-2116 and T-2117.

The story's agent reviews the planner's tasks when it pulls the story, and cancels a duplicate on the plan thread, as the planner's rules ask for any task it did not write.

Sources: S-0413, S-0473, S-0479, S-0496.
