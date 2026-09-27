---
id: T-0485
type: task
nature: remediation
title: A story does not go to review with uncommitted changes in its worktree
status: done
parent: S-0140
owner: alex
created: 2026-09-26T21:11:25Z
updated: 2026-09-26T21:15:15Z
transitions:
  - to: ready
    at: 2026-09-26T21:11:31Z
    by: agent-S-0140
  - to: in-progress
    at: 2026-09-26T21:11:32Z
    by: agent-S-0140
  - to: done
    at: 2026-09-26T21:15:15Z
    by: agent-S-0140
stream: S-0140
tags: []
touches: [flai/internal/workitem, flai/internal/mcpserver, flai/cmd, flai/internal/harness]
---
# T-0485 A story does not go to review with uncommitted changes in its worktree

## Work

- `flai move <story> review` and the MCP `item_move` refuse a story whose worktree under `.flai-cache/worktrees/` has uncommitted changes, naming the paths and saying to commit them on the story branch.
- The acceptance preview (`flai accept --dry-run`, the dashboard's confirmation) reports the story worktree's uncommitted paths as their own field and as a blocker, before anything is merged.
- The prompt flai serve gives an agent says to commit everything in the worktree before moving the story to review.

## Done when

- [x] Behaviour tests show the refusal from the CLI and from MCP, and a clean worktree moving as before.
- [x] The dry run lists the worktree's uncommitted paths.
- [x] The harness prompt test pins the commit instruction.

## Notes
