---
id: T-1598
type: task
nature: improvement
title: Close I-0147 with the remedy that fixed it
status: done
parent: S-0390
owner: alex
created: 2026-10-15T02:59:52Z
updated: 2026-10-15T03:20:17Z
transitions:
  - to: ready
    at: 2026-10-15T03:20:03Z
    by: agent-S-0390
  - to: in-progress
    at: 2026-10-15T03:20:04Z
    by: agent-S-0390
  - to: done
    at: 2026-10-15T03:20:17Z
    by: agent-S-0390
stream: S-0390
tags: [issues]
touches: [design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md, design/issues/summary.md]
after: [T-1596, T-1597]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2813
      cache_read: 572995
      cache_write: 17311
      cost: 0.2849
---
# T-1598 Close I-0147 with the remedy that fixed it

## Work

Run `flai issue close I-0147 --reason` from the story's worktree. The reason names the new ADR and says that an outdated session's `task_done` now closes the task with the installed flai once the host has caught up. It waits for T-1597 so the issue closes only once the remedy is built and documented.

## Done when

- I-0147 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- The story's second acceptance criterion can be ticked.

## Notes

Run it in the worktree, not the main checkout, so the change lands on the story branch.
