---
id: T-0737
type: task
nature: improvement
title: The pull respects the review limit
status: done
parent: S-0243
owner: arobson
created: 2026-10-03T02:58:39Z
updated: 2026-10-03T03:12:43Z
transitions:
  - to: ready
    at: 2026-10-03T02:59:05Z
    by: agent-S-0243
  - to: in-progress
    at: 2026-10-03T02:59:05Z
    by: agent-S-0243
  - to: done
    at: 2026-10-03T03:12:43Z
    by: agent-S-0243
stream: S-0243
tags: []
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/changes_test.go, flai/internal/serve/agents.go, flai/internal/serve/start.go, flai/internal/serve/restart.go, flai/internal/serve/agents_test.go, flai/internal/serve/start_test.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/agents.go, flai/internal/mcpserver/work_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 818
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 89
      output: 26490
      cache_read: 3380325
      cache_write: 96441
      cost: 1.8036
---
# T-0737 The pull respects the review limit

## Work

While review is at or over its limit, no story is pulled (TH-0078, option A). `flai serve` starts no new agent for a ready story, and says why: "review is full (5 of 5): accept or send back a story". The MCP `inbox` and `wait_for_work` report `can_pull` false with that reason, and `wait_for_work` offers no story. Stories already in progress finish. The operator's start-now still goes past it and names it.

## Done when

- Tests show serve, `inbox`, and `wait_for_work` holding a ready story while review is full, and pulling again once review has room.
- The in-progress limit's hold is unchanged.

## Notes
