---
id: S-0746
type: story
nature: feature
title: The MCP tools, story start, task done, the sync, item edits, and the inbox act on a state's kind and who may make a transition instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:16Z
updated: 2026-11-13T09:21:34Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/internal/storystart/start.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/storygit/sync.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/covers.go, flai/internal/inbox/inbox.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 108313
          cache_write: 18019
          cost: 0.0361
draft: true
forecast:
  duration: 42m
  delivery: 2026-11-13T17:44:00Z
  basis: "Its own forecast of 42m; 20th in the pull order with an in-progress limit of 3, behind S-0761, S-0763, S-0760, S-0745, S-0759, S-0751, S-0757, S-0740, S-0661, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743 and S-0744."
  by: flai
  at: 2026-11-13T09:21:34Z
---
# S-0746 The MCP tools, story start, task done, the sync, item edits, and the inbox act on a state's kind and who may make a transition instead of a state's name

## Goal

An agent's tools follow the workflow. `item_move` lets the agent make only the transitions the workflow gives it, `story_start` and `flai story start` make the pull transition, `task_done` closes a task along the task's transition to its accepted state, the sync, item edits, and claims tell an open item from a closed one by kind, and the inbox reports moves in the workflow's names.

## Acceptance criteria

- [ ] `flai/internal/mcpserver`, `storystart`, `taskdone`, `storygit/sync.go`, `itemedit`, and `inbox` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [ ] `item_move` refuses a transition its caller may not make, naming the transition and who may make it; the move to done stays refused as today, and the move to the wait state is still refused while the worktree has uncommitted changes
- [ ] `story_start` makes the transition with the pull effect from the queue state the workflow names, and refuses a story elsewhere with today's words and that state's name
- [ ] `task_done` closes a task along the task's transition to an accepted state, and a story's claim counts its open tasks by kind
- [ ] Tests run `item_move`, `story_start`, `task_done`, and the inbox's moved changes against a workflow with a state renamed and one added between `in-progress` and `review`
- [ ] With the baseline, every existing test of these packages passes unchanged

## Tasks

- T-3224 item_move refuses a transition its caller may not make, naming it and who may, and the MCP server's other tools read a state's kind
- T-3225 story_start and flai story start make the pull transition from the queue state the workflow names
- T-3226 task done, the sync, item edits and claims, and the inbox tell open from closed by kind

## Notes

- `mcpserver` holds 14 references (`server.go` 6, `work.go` 4); `itemedit` 3; `storystart`, `taskdone`, `storygit`, and `inbox` one each.

### Planning

- Touches, all files:
  - Layout: the files of these packages that name a state, with the tests of the tools whose behaviour changes.
  - Co-change: `flai/internal/taskdone/taskdone_test.go` (8% in `flai touches suggest`) and `mcpserver/story_start_test.go`.
  - Left out from `flai touches suggest`: `mcpserver/folder_test.go` (19%), since no tool is added; the user guide and reference (17% to 18%), since no tool's description changes; `cmd/mcp.go` and `mcp_http.go` (8% to 9%).
- `after`: S-0737. `item_move` and `story_start` read who may make each transition and the pull effect from the workflow S-0737 gives, through its mover check, and the tests use its test hook. The second criterion's uncommitted-worktree refusal stays as today whether S-0738 has landed or not, so this story does not wait for it.
- Tasks: T-3224 (`item_move` and the server), T-3225 (story start), and T-3226 (task done, sync, edits, inbox), one layer: no two share a path.
- Forecast 42m, `flai forecast`'s rate for size 23, up from its 39m for the first draft's 21.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
