---
id: T-1050
type: task
nature: improvement
title: The planner agent definitions say an epic's planner drafts the tasks of the stories it drafts
status: done
parent: S-0300
owner: alex
created: 2026-10-06T21:47:23Z
updated: 2026-10-06T23:26:56Z
transitions:
  - to: ready
    at: 2026-10-06T23:07:14Z
    by: agent-S-0300
  - to: in-progress
    at: 2026-10-06T23:07:14Z
    by: agent-S-0300
  - to: done
    at: 2026-10-06T23:26:56Z
    by: agent-S-0300
stream: S-0300
tags: [planner, conventions]
touches: [".claude/agents/planner.md", template/root/.claude/agents/planner.md]
after: [T-1047]
usage:
  source: log
  seconds: 1182
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 70
      cache_read: 1283344
      cache_write: 30287
      cost: 0.5787
---
# T-1050 The planner agent definitions say an epic's planner drafts the tasks of the stories it drafts

## Work

The planner's sub-agent definition still describes the epic work without tasks. Bring it in line with the convention.

- In `.claude/agents/planner.md` and its template copy `template/root/.claude/agents/planner.md`, change step 3. For an epic, the planner drafts its stories, enriches each one, and drafts its tasks with `## Work`, `## Done when`, touches, and `after`.
- Change step 6 so the epic summary names the tasks created.
- While editing, fix the description's slip "revisits an story's open tasks".
- The two files stay identical.

This task waits for T-1047, the design and convention task, so its words follow the convention's. It is in layer 2.

A write under `.claude/` needs the operator: do it yourself, not through a sub-agent, once the rest of the story is done. The call opens an "Allow ...?" thread that holds until alex answers it (TH-0194).

## Done when

- [ ] Both `planner.md` files say that an epic's planner drafts the tasks of the stories it drafts, and that its summary names them.
- [ ] Both files are identical.

## Notes

Drafted by the planner on the recommendation in TH-0201.
