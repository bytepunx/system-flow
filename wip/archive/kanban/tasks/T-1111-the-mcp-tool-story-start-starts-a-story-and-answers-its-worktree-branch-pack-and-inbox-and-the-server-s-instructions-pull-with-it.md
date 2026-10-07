---
id: T-1111
type: task
nature: improvement
title: The MCP tool story_start starts a story and answers its worktree, branch, pack, and inbox, and the server's instructions pull with it
status: done
parent: S-0274
owner: alex
created: 2026-10-06T22:53:37Z
updated: 2026-10-07T07:59:44Z
transitions:
  - to: ready
    at: 2026-10-07T07:40:25Z
    by: agent-S-0274
  - to: in-progress
    at: 2026-10-07T07:40:25Z
    by: agent-S-0274
  - to: done
    at: 2026-10-07T07:59:44Z
    by: agent-S-0274
stream: S-0274
tags: [mcp, go]
touches: [flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go]
after: [T-1099, T-1103]
usage:
  source: log
  seconds: 1159
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 157
      output: 55918
      cache_read: 8070827
      cache_write: 257907
      cost: 4.166
---
# T-1111 The MCP tool story_start starts a story and answers its worktree, branch, pack, and inbox, and the server's instructions pull with it

## Work

Add the tool `story_start`, taking `story`, an optional `budget`, and `project` where the server serves more than one. Write its handler in `flai/internal/mcpserver/story_start.go` and register it in `server.go` and, with `project`, in `folder.go`'s `addProjectTools`. It calls the storystart function under the server's agent name, then the exported inbox function, and answers `story`, `followed`, `worktree`, `branch`, `from`, `pack`, and `inbox`, the keys `flai story start --json` gives. It refuses a story that is not ready or is held, with the reason.

Rewrite where the server tells an agent to pull with two steps. Both servers' `Instructions`, and the `wait_for_work` descriptions in `server.go` and `folder.go`, say "item_move it to in-progress, then flai stream open on the host". They now say to call `story_start`, which does both and primes, and to call `wait_for_work` again if it says another agent took the story.

The guard needs no change: `story_start` is not in `MCPReads`, so a sub-agent stays refused it.

It waits for T-1099 and T-1103. It runs alongside the CLI task, which touches no file this one does.

## Done when

- `story_start_test.go` covers a ready story started with every key present, a held story refused with nothing changed, and the folder server's tool with `project`.
- The instructions and both `wait_for_work` descriptions name `story_start` and no longer name the two-step pull.
- `scripts/flai-test.sh` passes.

## Notes
