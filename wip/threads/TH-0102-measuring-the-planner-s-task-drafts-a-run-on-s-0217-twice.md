---
id: TH-0102
title: "Measuring the planner's task drafts: a run on S-0217, twice"
anchor:
  path: wip/kanban/stories/S-0255-the-planner-drafts-a-story-s-tasks-into-the-backlog-and-revisits-the-children-it-already-has.md
  item: S-0255
status: open
participants: [agent-S-0255]
created: 2026-10-04T04:12:53Z
updated: 2026-10-04T04:21:03Z
---

# TH-0102 Measuring the planner's task drafts: a run on S-0217, twice

On wip/kanban/stories/S-0255-the-planner-drafts-a-story-s-tasks-into-the-backlog-and-revisits-the-children-it-already-has.md.

## Entries

### 2026-10-04T04:12:53Z agent-S-0255
Criterion 5 needs a planner run on a story with the S-0255 build, judged by you here. The code and docs are done on story/S-0255; S-0255 goes to review after S-0209 is accepted (TH-0099, TH-0100), because the item_new body check and the activity entry's item list it relies on come from S-0209.

Recommended plan, once S-0209 is accepted and I have synced:
1. I build bin/flai in the S-0255 worktree and tell you here.
2. You turn on the plan action (`flai serve enable plan`) and, from the main checkout, run `.flai-cache/worktrees/S-0255/bin/flai plan S-0217`. S-0217 (the orchestrator's deterministic operations as commands) is a backlog story of E-0016 with no tasks and plenty of code to plan, so it tests drafting. Note that the planner session's guard is main's (`scripts/flai.sh guard`), which lacks `task new` until S-0255 is accepted, so the planner may write its tasks with item_new instead.
3. When it ends, run the same command again: the second run finds the tasks the first drafted and tests revisiting.
4. You judge the drafts and the plan thread here: are the tasks the ones you would want, are touches and after right, does the revisit propose sensibly. I record the runs, their cost from wip/agents/planner.md, and your judgement in design/system/strategic-agents.md.

I cannot start the planner myself: the MCP tool plan refuses an agent flai serve started, and planning is yours to ask for. If you would rather use another story, or run it from the dashboard after acceptance instead (then I record it in a follow-up), say so.

### 2026-10-04T04:21:03Z agent-S-0255
agent-S-0255: correction from what TH-0103 found. The host runs the installed release (~/.flai/bin/flai 1.29.1), which predates the plan action, so the run needs the accepted work published and the host upgraded first, as TH-0103 recommends. Once S-0209 is accepted and I have synced and built, the command is `.flai-cache/worktrees/S-0255/bin/flai --config ~/.flai/config.json plan S-0217`, run twice. I will post the exact command here when the build is ready. Nothing needed from you until then.
