---
id: S-0746
type: story
nature: feature
title: The MCP tools, story start, task done, the sync, item edits, and the inbox act on a state's kind and who may make a transition instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:16Z
updated: 2026-11-12T16:00:16Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/installed.go, flai/internal/storystart/start.go, flai/internal/taskdone/taskdone.go, flai/internal/storygit/sync.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/covers.go, flai/internal/inbox/inbox.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0746 The MCP tools, story start, task done, the sync, item edits, and the inbox act on a state's kind and who may make a transition instead of a state's name

## Goal

An agent's tools follow the workflow. `item_move` lets the agent make only the transitions the workflow gives it, `story_start` and `flai story start` make the pull transition, `task_done` closes a task along the task's transition to its accepted state, the sync, item edits, and claims tell an open item from a closed one by kind, and the inbox reports moves in the workflow's names.

## Acceptance criteria

- [ ] `flai/internal/mcpserver`, `storystart`, `taskdone`, `storygit/sync.go`, `itemedit`, and `inbox` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [ ] `item_move` refuses a transition its caller may not make, naming the transition and who may make it; the move to done stays refused as today, and the move to the wait state is refused while the worktree has uncommitted changes, by the rule its transition lists
- [ ] `story_start` makes the transition with the pull effect from the queue state the workflow names, and refuses a story elsewhere with today's words and that state's name
- [ ] `task_done` closes a task along the task's transition to an accepted state, and a story's claim counts its open tasks by kind
- [ ] Tests run `item_move`, `story_start`, `task_done`, and the inbox's moved changes against a workflow with a state renamed and one added between `in-progress` and `review`
- [ ] With the baseline, every existing test of these packages passes unchanged

## Tasks

## Notes

- `mcpserver` holds 14 references (`server.go` 6, `work.go` 4); `itemedit` 3; `storystart`, `taskdone`, `storygit`, and `inbox` one each.
