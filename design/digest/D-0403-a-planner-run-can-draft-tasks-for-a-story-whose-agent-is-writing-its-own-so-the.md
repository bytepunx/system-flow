---
id: D-0403
title: A planner run can draft tasks for a story whose agent is writing its own, so the story's agent cancels the duplicates on the plan thread before its layers run
kind: decision
sources: [S-0413, S-0473, S-0479, S-0496, TH-0472, TH-0574, TH-0577, S-0548, TH-0770]
topics: [planning, orchestration]
paths: [design/system/strategic-agents.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0403 A planner run can draft tasks for a story whose agent is writing its own, so the story's agent cancels the duplicates on the plan thread before its layers run

Five stories met planner drafts that duplicated tasks their own agent had written or finished:

- S-0413 cancelled five (TH-0472).
- S-0473 cancelled T-1986 and T-1989, which duplicated T-1992, already done (TH-0574).
- S-0479 cancelled four drafts that its T-2007 covered (TH-0577).
- S-0496 cancelled T-2113 to T-2115, which duplicated its T-2116 and T-2117.
- S-0548 was rewritten to nine criteria seconds after its pull (TH-0770). Its agent wrote T-2361 to T-2364 against the old four while planner-S-0548 drafted T-2365 to T-2369. This time the planner folded the agent's four into its plan and rewrote its own duplicate, T-2366. The agent kept the planner's plan and adjusted the tasks' scope (wip/archive/agents/S-0548.md, Decisions).

The story's agent reviews the planner's tasks when it pulls the story, and cancels a duplicate on the plan thread, as the planner's rules ask for any task it did not write. When the planner has already merged the two sets, as on S-0548, the agent reviews the merged plan instead and cancels nothing.

Sources: S-0413, S-0473, S-0479, S-0496, S-0548.
