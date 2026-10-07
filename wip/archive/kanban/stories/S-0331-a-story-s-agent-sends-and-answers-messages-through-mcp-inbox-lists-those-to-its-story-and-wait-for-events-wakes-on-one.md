---
id: S-0331
type: story
nature: feature
title: A story's agent sends and answers messages through MCP, inbox lists those to its story, and wait_for_events wakes on one
status: done
parent: E-0018
owner: alex
created: 2026-10-07T20:10:32Z
updated: 2026-10-07T21:42:43Z
transitions:
  - to: ready
    at: 2026-10-07T21:00:31Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T21:00:43Z
    by: agent-S-0331
  - to: review
    at: 2026-10-07T21:41:42Z
    by: agent-S-0331
  - to: done
    at: 2026-10-07T21:42:43Z
    by: orchestrator
tags: [flai, template]
topics: [cli, conventions, template]
touches: [flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/inbox_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, flai/cmd/touches.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/items_write_test.go, flai/internal/storystart/start.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, design/conventions/delegation.md, design/system/workflow.md, template/root/design/conventions/delegation.md, flai/internal/check/check.go]
after: [S-0330]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2471
  turns:
    - day: 2026-10-07
      ceremony: 2
      test_runs: 1
      hand_edits: 3
      work: 56
  models:
    - model: claude-opus-5-5
      input: 336
      output: 120741
      cache_read: 20584074
      cache_write: 735692
      cost: 10.9822
  strategic:
    - kind: orchestrator
      seconds: 2859
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 112
          output: 2096
          cache_read: 13661361
          cache_write: 34123
          cost: 3.3746
        - model: claude-sonnet-5-5
          input: 20
          output: 133
          cache_read: 340219
          cache_write: 58276
          cost: 0.3435
cost_of_delay:
  value: 113.7
  by: planner-E-0018
  at: 2026-10-07T20:21:59Z
forecast:
  duration: 44m
  delivery: 2026-10-08T07:58:00Z
  basis: "Its own forecast of 44m; 38th in the pull order with an in-progress limit of 3, behind S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0280, S-0287, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0309, S-0310, S-0312, S-0313, S-0314, S-0315, S-0316, S-0317, S-0318, S-0319, S-0320, S-0321, S-0322, S-0323, S-0324, S-0325, S-0326 and S-0327."
  by: flai
  at: 2026-10-07T21:00:01Z
finalized:
  by: orchestrator
  at: 2026-10-07T20:24:15Z
---
# S-0331 A story's agent sends and answers messages through MCP, inbox lists those to its story, and wait_for_events wakes on one

## Goal

Put S-0330's messages in front of the agents that work stories. A story's agent sends and answers messages through MCP tools, sees the conversations awaiting its story in `inbox`, and is woken by `wait_for_events` within a second of one arriving, so two agents can settle who changes a shared path while both are working. The convention tells them when to message rather than open a thread.

## Acceptance criteria

- [x] The MCP tools `message_send`, `message_reply`, and `message_get` send from the calling agent's story (`FLAI_STORY`, or the story named by its agent name), answer, and read a conversation; a session with no story of its own is refused.
- [x] `inbox` returns `messages`: the open conversations of the agent's story, each with the other story, its `about` paths, its last entry, and whether it awaits this story; the operator's `awaiting_you` does not count them.
- [x] `wait_for_events` returns at once when a message to the agent's story arrives, as an event of kind `message` naming the conversation and the sender's story.
- [x] `flai guard` refuses a sub-agent's `message_send` and `message_reply`, as it refuses its thread writes (ADR-0060).
- [x] The MCP server's instructions and the `overlapped` guidance say to message the other story's agent, not to open a thread, and to ask the operator on a thread only when the two do not agree.
- [x] `work-management.md`, in `design/conventions` and in the template, says when a story's agent messages another and how it answers one, and `template/CHANGELOG.md` records it.
- [x] `design/system/agent-narrative.md` § What an agent is told through MCP and `design/system/flai-cli.md` describe the tools and the inbox field.

## Tasks

- T-1190 The MCP tools message_send, message_reply, and message_get work a conversation from the calling agent's story
- T-1191 flai guard refuses a sub-agent's message_send and message_reply
- T-1192 inbox lists the conversations of the agent's story under messages, and wait_for_events wakes on a message to it
- T-1193 The convention, the template, and the design say when a story's agent messages another and how it answers one

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07. It waits for S-0330 (`after`), whose package its tools call.

Layers:

1. T-1190, the tools; T-1191, the guard. They share no path and run together.
2. T-1192, the inbox and the wait. It waits for T-1190 because both change `server.go`.
3. T-1193, the convention and the docs.

Touches:

- **Declared:** none before planning.
- **Layout:**
  - `flai/internal/mcpserver/messages.go` and its test: a new tool file, as `task.go` and `issues.go` are.
  - `flai/internal/mcpserver/server.go`: tool registration, the instructions, and `watched()`.
  - `flai/internal/inbox/inbox.go` and its test: the inbox and its changes are built there.
  - `flai/internal/guard/guard.go` and its test.
- **Co-change:** `flai touches suggest` from `server.go`, `inbox.go`, and `guard.go` gave `flai/internal/mcpserver/folder.go` (49%), where tools are routed for a folder that is not a project, `server_test.go` (33%), `guard_test.go` (25%), `docs/users/flai.md` (23%), and `design/system/flai-cli.md` (20%). `inbox_test.go` in `mcpserver` is where the inbox tool is tested.
- **Design:** both copies of `work-management.md` and `template/CHANGELOG.md`, whose rules say to coordinate on a thread; `design/system/agent-narrative.md` § What an agent is told through MCP.
- **Not taken:** `flai/cmd/guard.go`, `flai/cmd/mcp.go`, `docs/operators/settings.md`, and `flai/internal/hostapi/writes.go`: `suggest` listed them at 10 to 19%, and nothing here adds a setting or a host write.

Forecast 44m, delivery 2026-10-08T08:31Z: `flai forecast` gave it, 114 s per unit over 29 done large-band feature stories, times size 23. It stands: four tasks over known files.

Cost of delay 113.70 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 44m of 6h27m. It stands.

### Accepted by the orchestrator

- Verified: 49bd891e30599a7590c4f2dbaa5cbee68b70ee94
- At: 2026-10-07T21:42:43Z

Verdict: all seven criteria are met at head 49bd891e; flai verify passed every step there, and the branch changes no .claude/ path.
- 1: flai/internal/mcpserver/messages.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/messages_test.go
- 2: flai/internal/inbox/inbox.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/inbox_test.go
- 3: flai/internal/mcpserver/server.go, flai/internal/mcpserver/cursor.go, flai/internal/inbox/inbox.go, flai/internal/mcpserver/messages_test.go
- 4: flai/internal/guard/guard.go, flai/internal/guard/guard_test.go
- 5: flai/internal/mcpserver/messages.go, flai/cmd/touches.go, flai/internal/mcpserver/items_write.go
- 6: design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md
- 7: design/system/agent-narrative.md, design/system/flai-cli.md
