---
id: T-1910
type: task
nature: feature
title: The orchestrator's agent definitions name queue_add
status: done
parent: S-0403
owner: alex
created: 2026-10-20T00:24:40Z
updated: 2026-10-20T04:29:32Z
transitions:
  - to: ready
    at: 2026-10-20T00:25:06Z
    by: agent-S-0403
  - to: in-progress
    at: 2026-10-20T00:25:07Z
    by: agent-S-0403
  - to: done
    at: 2026-10-20T04:29:32Z
    by: agent-S-0403
stream: S-0403
tags: []
touches: [".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-1663]
usage:
  source: log
  seconds: 681
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 116
      output: 40517
      cache_read: 8500193
      cache_write: 308383
      cost: 4.4241
---
# T-1910 The orchestrator's agent definitions name queue_add

## Work

- Add `mcp__flai__queue_add` to the tools of `template/root/.claude/agents/orchestrator.md` and its copy `.claude/agents/orchestrator.md`, and one sentence in step 11: while the merge queue runs, as the prompt says, enqueue each verified story with `queue_add` instead of accepting it.
- Both are protected paths: the story's agent writes them through the operator's permission prompt (TH-0522).
- Waits for T-1663, whose convention and prompt the sentence follows. Split out of T-1663 so its code and convention could close while the write waits on the operator.

## Done when

Both files name `queue_add`, stay identical, and the harness test that reads them passes.

## Notes
