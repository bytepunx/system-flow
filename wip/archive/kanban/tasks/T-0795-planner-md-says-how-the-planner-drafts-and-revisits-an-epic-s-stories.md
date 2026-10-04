---
id: T-0795
type: task
nature: feature
title: planner.md says how the planner drafts and revisits an epic's stories
status: done
parent: S-0209
owner: alex
created: 2026-10-04T04:03:02Z
updated: 2026-10-04T04:18:17Z
transitions:
  - to: ready
    at: 2026-10-04T04:03:32Z
    by: agent-S-0209
  - to: in-progress
    at: 2026-10-04T04:10:01Z
    by: agent-S-0209
  - to: done
    at: 2026-10-04T04:18:17Z
    by: agent-S-0209
stream: S-0209
tags: []
touches: [".claude/agents/planner.md", template/root/.claude/agents/planner.md]
after: [T-0791]
usage:
  source: log
  seconds: 496
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 10425
      cache_read: 3982571
      cache_write: 45254
      cost: 1.3046
---
# T-0795 planner.md says how the planner drafts and revisits an epic's stories

## Work

- `template/root/.claude/agents/planner.md` and its copy `.claude/agents/planner.md`: step 3 and step 6 say what the prompt says after T-0791: stories as drafts, the plan thread with stories, order, and assumptions, proposals for splits, merges, additions, and drops in it, and a final line naming the stories created and revisited.
- Agents cannot write under `.claude/`: the operator pastes the contents the story's agent sends on a thread, and the agent checks them.
- Waits for T-0791: it says what the prompt says.

## Done when

- [ ] The two files are the same and say what `planPrompt` says.

## Notes

Shares the second layer with the documentation task.
