---
id: S-0338
type: story
nature: feature
title: A held card says when its holding story's agent was asked about the hold, and a story started on a share names the paths it shares
status: backlog
parent: E-0018
owner: alex
created: 2026-10-08T04:31:41Z
updated: 2026-10-08T04:38:47Z
transitions: []
tags: [flai, flaiover]
topics: [dashboard, cli]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/Messages.svelte, flaiover/src/lib/components/Messages.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, design/system/flai-cli.md]
after: [S-0334, S-0336]
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
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 48
          cache_read: 2064636
          cache_write: 4326
          cost: 0.5097
cost_of_delay:
  value: 217.39
  by: planner-E-0018
  at: 2026-10-08T04:33:24Z
forecast:
  duration: 40m
  delivery: 2026-10-08T13:55:00Z
  basis: "Its own forecast of 40m; 29th in the pull order with an in-progress limit of 3, behind S-0232, S-0316, S-0318, S-0320, S-0319, S-0309, S-0336, S-0326, S-0287, S-0322, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0313, S-0321, S-0323, S-0334 and S-0337."
  by: flai
  at: 2026-10-08T04:38:47Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:35:22Z
---
# S-0338 A held card says when its holding story's agent was asked about the hold, and a story started on a share names the paths it shares

## Goal

Show the operator, on the board, that a hold is being worked out between agents. S-0334 has flai serve ask the holding story's agent about each story held on overlap, and lets that agent share the paths with a stated split. S-0336 shows the conversations on a story's page and in a Messages view. Neither says so on the card, where the operator watches holds. Split from S-0336 at the operator's word on TH-0320.

## Acceptance criteria

- [ ] `flai board --json` gives a ready card held on overlap the conversation that asked its holding story's agent about the hold, and a card started on a share the shared paths, the split, and the conversation.
- [ ] A held card whose holding story's agent was asked says so and links the conversation.
- [ ] A card of a story started on a share names the shared paths and links the conversation.
- [ ] The Messages view and a story's page show a conversation's share and its split.
- [ ] `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, and `docs/users/flai.md` describe the card and the new board fields.

## Tasks

- T-1324 flai board gives a held card the conversation that asked about its hold, and a card started on a share its shared paths
- T-1325 The Messages component shows a conversation's share and its split
- T-1326 A held card says its holding story's agent was asked, and a card started on a share names the shared paths
- T-1327 The CLI design, the dashboard design, and the guides describe the asked and shared cards and the new board fields

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-08, split from S-0336 as the operator accepted on TH-0320. It waits for S-0334 (`after`), which makes the hold requests and the shares it shows, and for S-0336, which builds `Messages.svelte`.

Layers:

1. T-1324, the board's fields in flai; T-1325, the share in `Messages.svelte`. They share no path and run together.
2. T-1326, the card, which reads T-1324's fields.
3. T-1327, the docs.

Touches:

- **Declared:** the thirteen given at creation, all kept.
- **Layout:**
  - `flai/internal/workitem/hold.go`, `boardview.go`, and their tests: `Hold` is `{code, reason}` today, and `NewBoardView` puts it on the card.
  - `flaiover/src/lib/activity.ts` and its test: the dashboard's `Hold` type and `holdLine`; `BoardCard.svelte` and its test show it.
  - `flaiover/src/lib/components/Messages.svelte` and its test, which S-0336 creates.
- **Co-change:** `flai touches suggest` gave `design/system/flai-cli.md` (56%), which documents `held: {code, reason}` on the board's card; added.
- **Not taken:** `docs/operators/index.md` (23%) and `docs/users/flai-reference.md` (21%): no setting and no flag change. `design/system/workflow.md` (15%): S-0334 describes the share there. `flai/internal/serve/agents.go` and `flai/internal/mcpserver/work.go` may have to pass the asks and shares in; T-1324 widens its touches to each caller it changes.
- No folder touch.

Forecast 40m: `flai forecast` gave 33m, 104 s per unit over 33 done large-band feature stories, times size 19 (5 criteria, 14 touches). Raised to 40m, because `workitem` cannot import `messages` and T-1324 will pass the asks and shares in from callers no task can name yet, so the size will grow.

Cost of delay 217.39 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 40m of the 3h4m forecast over the four open stories. It stands.
