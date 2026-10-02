---
id: TH-0069
title: "Forks are still not offered headless: the repeat measurement cannot run in S-0230"
anchor:
  path: wip/kanban/stories/S-0230-tasks-carry-after-the-plan-shows-on-the-board-and-a-story-s-agent-hands-tasks-to-sub-agents-without-parallel-layers.md
  item: S-0230
status: resolved
participants: [agent-S-0230, alex]
created: 2026-10-02T17:11:51Z
updated: 2026-10-02T17:16:46Z
---

# TH-0069 Forks are still not offered headless: the repeat measurement cannot run in S-0230

On wip/kanban/stories/S-0230-tasks-carry-after-the-plan-shows-on-the-board-and-a-story-s-agent-hands-tasks-to-sub-agents-without-parallel-layers.md.

## Entries

### 2026-10-02T17:11:51Z agent-S-0230
The fourth criterion asks that the measurement be repeated once forked sub-agents are offered to a headless session. They still are not: in this `claude -p` session (Claude Code 2.1.286), started by flai serve at 17:1xZ today, the Agent tool refused `subagent_type: fork` with "Agent type 'fork' not found. Available agents: claude, Explore, explorer, general-purpose, Plan, statusline-setup, verifier", the same as S-0176 found, and `claude --help` has no flag that offers them.

So the criterion's condition has not happened, and the repeat cannot be run in this story.

Recommended: take the fourth criterion out of S-0230 and file it as its own `experiment` story in the backlog under the same epic ("Measure planned stories again with forked task sub-agents once a headless session offers forks"), with a note that it waits for Claude Code to offer `fork` headless; S-0230's agent-context.md section records that forks were probed today and still refused. I can file that story and edit the criterion myself if you agree.

Alternatives: (b) keep the criterion and I block S-0230 after the other four are done, until forks arrive; (c) reword the criterion to "probe whether forks are offered and record the result", which I can meet now.

Meanwhile I am working the other four criteria (the host's flai 1.28.0 already carries `after` and `flai.minimum` is 1.28.0; I-0054; the conventions; agent-context.md).

### 2026-10-02T17:13:52Z alex
use rec

### 2026-10-02T17:16:46Z agent-S-0230
Resolved: Done as recommended: S-0241 filed in the backlog for the fork measurement, and S-0230's fourth criterion now asks agent-context.md to record the probe and the split.
