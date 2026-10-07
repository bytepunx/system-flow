---
id: S-0334
type: story
nature: feature
title: A story held on overlap is asked about by message, and starts when the holding story's agent shares the paths or narrows its claim
status: backlog
parent: E-0018
owner: alex
created: 2026-10-07T20:11:04Z
updated: 2026-10-07T23:47:55Z
transitions: []
tags: [flai, template]
topics: [cli, conventions, template]
touches: [design/adrs, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/boardview.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/workflow.md, design/system/agent-coordination.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0332]
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
      seconds: 300
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 45
          output: 819
          cache_read: 3661515
          cache_write: 10356
          cost: 0.9048
draft: true
cost_of_delay:
  value: 149.87
  by: planner-E-0018
  at: 2026-10-07T20:22:03Z
forecast:
  duration: 58m
  delivery: 2026-10-08T07:59:00Z
  basis: "Its own forecast of 58m; 34th in the pull order with an in-progress limit of 3, behind S-0232, S-0310, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0287, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0309, S-0312, S-0313, S-0315, S-0316, S-0317, S-0318, S-0319, S-0320, S-0321, S-0322, S-0323, S-0324, S-0325, S-0326 and S-0327."
  by: flai
  at: 2026-10-07T23:47:55Z
---
# S-0334 A story held on overlap is asked about by message, and starts when the holding story's agent shares the paths or narrows its claim

## Goal

Raise throughput by letting the agent that holds a story decide whether the hold is needed. Today a ready story whose claim overlaps a story in progress waits until that story moves to review, even when the holding agent will not change the paths, or will change them only in a section the held story does not touch. When flai serve finds a ready story held on overlap alone, it messages the holding story's agent with the paths and the held story's goal. That agent answers by narrowing its touches, which clears the hold as it does today, or by sharing the paths with a stated split of who changes what. A shared overlap does not hold, and the held story's agent is told the split when it starts.

## Acceptance criteria

- [ ] When a ready story is held on overlap alone, flai serve messages each story in progress that holds it, once per pair while the hold lasts, `about` the overlapping paths, with the held story's goal and what the agent can answer.
- [ ] `flai message share <conversation> --paths <path>… "<split>"`, and the MCP tool `message_share`, record that the holding story shares those paths with the held story and how the work is split; only the holding story's agent, or the operator, may share.
- [ ] An overlap on shared paths does not hold: `flai board`, `inbox`, `wait_for_work`, and the launcher take the story, and its `held` reason names any paths still held; a share ends when either story leaves the open columns, or the held story leaves ready before it starts.
- [ ] The held story's agent, once started, finds the conversation and the split in its first `inbox`, and the prompt flai serve gives it names them.
- [ ] Two stories that share paths are still trial-merged at `flai stream sync`, and a conflict between them opens their conversation as S-0332 has it.
- [ ] An ADR refining ADR-0046 and ADR-0096 records the share, and `design/system/workflow.md` § Branches and collisions, `agent-coordination.md`, and `work-management.md` in both copies describe it.

## Tasks

- T-1202 An ADR refining ADR-0046 and ADR-0096 lets the holding story's agent share overlapping paths, so the overlap no longer holds
- T-1203 flai message share and the MCP tool message_share record a share of paths with a held story, and flai can ask about a hold
- T-1204 An overlap on shared paths does not hold, and the held reason names only the paths still held
- T-1205 flai serve asks a holding story's agent about each story held on overlap alone, and tells a story started on a share its split
- T-1206 The workflow, the coordination design, the convention, the template, and the guides describe asking about a hold and sharing paths

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07. It waits for S-0332 (`after`): its ADR refines S-0332's, and a share leans on the conflict conversation to catch a split that did not hold.

This is the throughput lever of E-0018: a hold the holding agent judges needless ends while that agent is still working, not when it reaches review. It refines ADR-0046 and ADR-0096, so T-1202 records it first.

Layers:

1. T-1202, the ADR.
2. T-1203, the share and the request in the messages package and its front ends.
3. T-1204, the hold honours shares; T-1205, flai serve asks and passes the split on. They share no path and run together.
4. T-1206, the docs.

Touches:

- **Declared:** none before planning.
- **Layout:**
  - The messages package, `flai/cmd/message.go`, `flai/internal/mcpserver/messages.go`, and their tests: the share; `flai/internal/guard/guard.go` and its test, so a sub-agent cannot share.
  - `flai/internal/workitem/hold.go` and its test: `NewHolds` and `WithShared` judge overlaps; `boardview.go` builds the holds the board, `inbox`, and `wait_for_work` read.
  - `flai/internal/serve/agents.go` and its test: the launcher's look; `flai/internal/harness/harness.go` and its test: the prompt a started agent gets.
- **Co-change:** `flai touches suggest` from `hold.go` and `agents.go` gave `agents_test.go` (48%), `flai/internal/serve/start.go` (17%), `harness.go` and its test (14 to 17%), and `hold_test.go` (12%).
- **Design:** `design/system/workflow.md` § Branches and collisions and `agent-coordination.md` § Decision; both copies of `work-management.md` and `template/CHANGELOG.md`; `flai-cli.md`, `docs/users/flai.md`, and the generated reference.
- **Folder touch kept:** `design/adrs`, for T-1202's ADR. Inside `claims.shared`.
- **Not taken:** `flai/cmd/serve_actions.go` (50%), `flai/internal/hostapi/writes.go` (24%), and `flai/internal/serve/restart.go` (24%): no host action, host write, or restart changes. `flai/internal/serve/start.go` and `flai/internal/mcpserver/work.go` may need the shares passed in; T-1204 widens its touches to each caller it changes.

Forecast 58m, delivery 2026-10-08T10:47Z: `flai forecast` gave it, 114 s per unit over 29 done large-band feature stories, times size 30 (6 criteria, 24 touches). It stands.

Cost of delay 149.87 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 58m of 6h27m. It stands; this story is where the epic's throughput comes from, but no input prices that apart.
