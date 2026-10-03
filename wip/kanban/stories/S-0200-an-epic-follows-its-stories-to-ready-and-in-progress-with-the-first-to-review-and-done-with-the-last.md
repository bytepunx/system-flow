---
id: S-0200
type: story
nature: improvement
title: "An epic follows its stories: to ready and in-progress with the first, to review and done with the last"
status: in-progress
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:11Z
updated: 2026-10-03T07:38:05Z
transitions:
  - to: ready
    at: 2026-10-03T05:33:53Z
    by: alex
  - to: in-progress
    at: 2026-10-03T07:06:32Z
    by: agent-S-0200
tags: [flai, dashboard]
touches: [flai/internal/workitem, flai/internal/itemedit, flai/cmd/move.go, flai/cmd/move_cancel.go, flai/cmd/accept.go, flai/cmd/edit.go, flai/cmd/check.go, flai/cmd/accept_epic_test.go, flai/cmd/release_pending_test.go, flai/cmd/accept_research_test.go, flai/cmd/accept_test.go, flai/cmd/stream_diff_test.go, flai/cmd/workitems_test.go, flai/cmd/check_stats_test.go, flai/internal/mcpserver, design/system/workflow.md, flai/cmd/edit_test.go, flai/cmd/hostapi_reads_test.go, flai/cmd/move_cancel_test.go, flai/internal/preview, flai/internal/check, design/system/work-hierarchy.md, design/system/flai-cli.md, design/adrs, docs/users, design/issues/I-0062-two-stories-in-progress-at-once-take-the-same-adr-number-and-one-renumbers-by-hand.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3117
  models:
    - model: claude-haiku-4-5-20251001
      input: 218
      output: 9370
      cache_read: 1704955
      cache_write: 90168
      cost: 0.3303
    - model: claude-opus-5-5
      input: 512
      output: 181119
      cache_read: 31347201
      cache_write: 702418
      cost: 14.1863
    - model: claude-sonnet-5
      input: 54
      output: 16609
      cache_read: 1920718
      cache_write: 103164
      cost: 0.8083
---
# S-0200 An epic follows its stories: to ready and in-progress with the first, to review and done with the last

## Goal

An epic's status is moved by hand and drifts from its stories'. The designer decided on 2026-10-02 that an epic follows its children: it moves to `ready` when its first story is moved to ready or created in ready, to `in-progress` when its first story starts, to `review` when its last open story moves to review, and to `done` when its last story is accepted. This also gives the planner and orchestrator epics whose state means something.

## Acceptance criteria
- [x] When a story moves from backlog to ready, or is created in ready, and its epic is in backlog, the epic moves to ready in the same write, with a transition `by` the same actor and a reason naming the story
- [x] When a story moves to in-progress and its epic is in ready or backlog, the epic moves to in-progress
- [x] When the last story of an epic that is not done or cancelled moves to review, the epic moves to review; when the last one is accepted (done), the epic moves to done through the same acceptance flow (`flai accept`, `flai move done`, the dashboard), archiving it as an epic done by hand is archived today
- [x] A story moved back (review to in-progress, ready to backlog) moves the epic back only when no other story holds it where it is; cancelled stories do not count
- [x] The epic's move is reported in the command's output, the MCP result, and the inbox as a change on the epic
- [x] `flai check` warns about an epic whose status lags its stories under these rules, so epics from before this change are found
- [x] `design/system/workflow.md`, `work-hierarchy.md`, and the user guide describe it; an ADR records the rule, refining ADR-0004
- [x] Tests cover each transition, the back moves, and acceptance of the last story

## Tasks
- T-0745 An epic follows its story's move in the same write
- T-0748 Accepting an epic's last story accepts the epic, and flai move reports the epic's move
- T-0749 item_move and the inbox report an epic's move with its story
- T-0750 flai check warns about an epic that lags its stories
- T-0752 An ADR and the design and user guide describe how an epic follows its stories

## Notes
