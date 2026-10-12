---
id: T-1442
type: task
nature: remediation
title: Close I-0127 as fixed by task_done's refusal while the running flai is older than the project's history
status: backlog
parent: S-0363
owner: alex
created: 2026-10-12T14:36:44Z
updated: 2026-10-12T14:36:44Z
transitions: []
stream: S-0363
tags: [cli]
touches: [design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md, design/issues/summary.md]
after: [T-1440, T-1441]
---
# T-1442 Close I-0127 as fixed by task_done's refusal while the running flai is older than the project's history

## Work

Close the issue once the fix and its docs are in. Run, in the story's worktree:

```bash
flai issue close I-0127 --reason "S-0363: flai task done and the MCP tool task_done refuse, changing nothing, while the running flai is older than the newest flai release in the project's history (ADR-0136), so an outdated flai no longer commits another open task's paths"
```

It writes the issue file and regenerates `design/issues/summary.md`.

Waits for T-1440 and T-1441: the reason names a fix that must be built and documented first.

## Done when

- I-0127's status is `closed`, with the reason above.
- `design/issues/summary.md` no longer lists I-0127.
- `flai check --strict --story S-0363` passes.

## Notes

Drafted by the planner for S-0363.
