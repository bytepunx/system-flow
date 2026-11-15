---
id: S-0746
type: story
nature: feature
title: The MCP tools, story start, task done, the sync, item edits, and the inbox act on a state's kind and who may make a transition instead of a state's name
status: done
parent: E-0034
owner: alex
created: 2026-11-12T16:00:16Z
updated: 2026-11-15T01:56:41Z
transitions:
  - to: ready
    at: 2026-11-14T17:48:03Z
    by: alex
  - to: in-progress
    at: 2026-11-14T23:49:15Z
    by: agent-S-0746
  - to: review
    at: 2026-11-15T00:56:45Z
    by: agent-S-0746
  - to: done
    at: 2026-11-15T01:56:41Z
    by: alex
tags: [cli]
topics: [cli]
touches: [flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/internal/storystart/start.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/storygit/sync.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/covers.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/itemedit/claim_test.go, flai/internal/itemedit/covers_test.go, flai/internal/storygit/sync_test.go, flai/internal/storystart/start_test.go, flai/internal/inbox/starved.go, flai/internal/itemedit/widen.go, design/issues/I-0375-flai-check-finds-board-order-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4059
  priced_by: harness
  turns:
    - day: 2026-11-14
      work: 27
    - day: 2026-11-15
      ceremony: 3
      work: 66
  models:
    - model: claude-opus-5-5
      input: 390
      output: 98490
      cache_read: 17577372
      cache_write: 652606
      cost: 9.798
  strategic:
    - kind: orchestrator
      seconds: 48
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 159
          cache_read: 691699
          cache_write: 126843
          cost: 0.2428
cost_of_delay:
  value: 74.82
  by: system-flow
  at: 2026-11-14T17:44:07Z
forecast:
  duration: 42m
  delivery: 2026-11-15T00:37:00Z
  basis: "Its own forecast of 42m; 3rd in the pull order with an in-progress limit of 3, behind S-0738 and S-0741."
  by: flai
  at: 2026-11-14T23:48:54Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:54Z
---
# S-0746 The MCP tools, story start, task done, the sync, item edits, and the inbox act on a state's kind and who may make a transition instead of a state's name

## Goal

An agent's tools follow the workflow. `item_move` lets the agent make only the transitions the workflow gives it, `story_start` and `flai story start` make the pull transition, `task_done` closes a task along the task's transition to its accepted state, the sync, item edits, and claims tell an open item from a closed one by kind, and the inbox reports moves in the workflow's names.

## Acceptance criteria

- [x] `flai/internal/mcpserver`, `storystart`, `taskdone`, `storygit/sync.go`, `itemedit`, and `inbox` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [x] `item_move` refuses a transition its caller may not make, naming the transition and who may make it; the move to done stays refused as today, and the move to the wait state is still refused while the worktree has uncommitted changes
- [x] `story_start` makes the transition with the pull effect from the queue state the workflow names, and refuses a story elsewhere with today's words and that state's name
- [x] `task_done` closes a task along the task's transition to an accepted state, and a story's claim counts its open tasks by kind
- [x] Tests run `item_move`, `story_start`, `task_done`, and the inbox's moved changes against a workflow with a state renamed and one added between `in-progress` and `review`
- [x] With the baseline, every existing test of these packages passes unchanged

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
- Cost of delay: 74.82 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 42m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
