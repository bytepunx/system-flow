---
id: T-1458
type: task
nature: improvement
title: flai task done runs the fix over the paths it commits, before it commits them
status: done
parent: S-0367
owner: alex
created: 2026-10-12T16:53:08Z
updated: 2026-10-12T17:09:11Z
transitions:
  - to: ready
    at: 2026-10-12T17:08:14Z
    by: agent-S-0367
  - to: in-progress
    at: 2026-10-12T17:08:14Z
    by: agent-S-0367
  - to: done
    at: 2026-10-12T17:09:11Z
    by: agent-S-0367
stream: S-0367
tags: []
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
after: [T-1456]
usage:
  source: log
  seconds: 57
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 18
      output: 7036
      cache_read: 1111852
      cache_write: 30773
      cost: 0.5542
---
# T-1458 flai task done runs the fix over the paths it commits, before it commits them

## Work

- `taskdone.Options` gets the fix function; the commit step runs it over the paths it is about to commit (the task's own and those no other open task covers) before staging them, so the fixed content is what it commits.
- Its sync step passes the same function to `storygit.Sync`, as T-1456 made it take.
- `flai task done` and the MCP tool `task_done` pass `verify.Fix` with the worktree's tiers.
- Waits for T-1456, whose `SyncOptions.Fix` its sync step passes on.

## Done when

- [ ] A task whose markdown a fix command rewrites is committed fixed, under a test.

## Notes
