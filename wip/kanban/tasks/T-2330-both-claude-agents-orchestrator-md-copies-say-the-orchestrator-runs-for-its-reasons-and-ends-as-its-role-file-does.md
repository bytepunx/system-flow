---
id: T-2330
type: task
nature: feature
title: Both .claude/agents/orchestrator.md copies say the orchestrator runs for its reasons and ends, as its role file does
status: done
parent: S-0536
owner: alex
created: 2026-10-29T08:02:00Z
updated: 2026-10-29T13:54:48Z
transitions:
  - to: ready
    at: 2026-10-29T13:54:38Z
    by: system-flow
  - to: in-progress
    at: 2026-10-29T13:54:39Z
    by: system-flow
  - to: done
    at: 2026-10-29T13:54:48Z
    by: agent-S-0536
stream: S-0536
tags: []
touches: [".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-2303]
usage:
  source: log
  seconds: 4
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 2
      output: 24
      cache_read: 413066
      cache_write: 334
      cost: 0.1925
---
# T-2330 Both .claude/agents/orchestrator.md copies say the orchestrator runs for its reasons and ends, as its role file does

## Work

Write the new description, opening paragraph, and step 16 of the orchestrator role file (T-2303) into `.claude/agents/orchestrator.md` and its template copy, the definition Claude Code runs (D-0411). Both paths are protected, so the story's agent writes them through `permission_prompt`, asked on TH-0737. It waits for T-2303, whose role file it copies.

## Done when

- [ ] Both copies carry the role file's new description and body, and are identical.

## Notes
