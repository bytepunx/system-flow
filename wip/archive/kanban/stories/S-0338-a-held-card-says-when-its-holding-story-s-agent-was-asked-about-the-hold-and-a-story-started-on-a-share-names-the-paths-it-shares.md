---
id: S-0338
type: story
nature: feature
title: A held card says when its holding story's agent was asked about the hold, and a story started on a share names the paths it shares
status: done
parent: E-0018
owner: alex
created: 2026-10-08T04:31:41Z
updated: 2026-10-08T20:55:08Z
transitions:
  - to: ready
    at: 2026-10-08T08:51:43Z
    by: alex
  - to: in-progress
    at: 2026-10-08T10:28:46Z
    by: agent-S-0338
  - to: review
    at: 2026-10-08T11:04:39Z
    by: agent-S-0338
  - to: done
    at: 2026-10-08T20:55:08Z
    by: orchestrator
tags: [flai, flaiover]
topics: [dashboard, cli]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/Messages.svelte, flaiover/src/lib/components/Messages.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, design/system/flai-cli.md, flai/internal/workitem/share.go, flai/internal/workitem/share_test.go, flaiover/src/lib/server/board.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md, design/issues/summary.md, design/issues/I-0128-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md, flai/cmd/story_start_test.go]
after: [S-0334, S-0336]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2175
  turns:
    - day: 2026-10-08
      ceremony: 1
      test_runs: 3
      hand_edits: 2
      work: 70
  models:
    - model: claude-opus-5-5
      input: 286
      output: 102155
      cache_read: 19194545
      cache_write: 566870
      cost: 9.5469
  strategic:
    - kind: orchestrator
      seconds: 1439
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 116
          output: 1537
          cache_read: 13316685
          cache_write: 112556
          cost: 3.3215
        - model: claude-sonnet-5-5
          input: 19
          output: 93
          cache_read: 304789
          cache_write: 66783
          cost: 0.3176
cost_of_delay:
  value: 217.39
  by: planner-E-0018
  at: 2026-10-08T04:33:24Z
forecast:
  duration: 40m
  delivery: 2026-10-08T11:16:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 5, behind S-0232."
  by: flai
  at: 2026-10-08T10:28:37Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:35:22Z
---
# S-0338 A held card says when its holding story's agent was asked about the hold, and a story started on a share names the paths it shares

## Goal

Show the operator, on the board, that a hold is being worked out between agents. S-0334 has flai serve ask the holding story's agent about each story held on overlap, and lets that agent share the paths with a stated split. S-0336 shows the conversations on a story's page and in a Messages view. Neither says so on the card, where the operator watches holds. Split from S-0336 at the operator's word on TH-0320.

## Acceptance criteria

- [x] `flai board --json` gives a ready card held on overlap the conversation that asked its holding story's agent about the hold, and a card started on a share the shared paths, the split, and the conversation.
- [x] A held card whose holding story's agent was asked says so and links the conversation.
- [x] A card of a story started on a share names the shared paths and links the conversation.
- [x] The Messages view and a story's page show a conversation's share and its split.
- [x] `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, and `docs/users/flai.md` describe the card and the new board fields.

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

### Accepted by the orchestrator

- Verified: 2d114ffc11fec9802810cbec19070f8857b0cc55
- At: 2026-10-08T20:55:08Z

Verdict: pass. All five criteria are met, every changed file is within the story's touches, and no convention break was found; flai verify passed every step at 2d114ffc.

- 1: flai/internal/workitem/hold.go, flai/internal/workitem/boardview.go, flai/internal/workitem/share.go, flai/internal/workitem/boardview_test.go, flai/internal/workitem/share_test.go, flai/cmd/story_start_test.go
- 2: flaiover/src/lib/activity.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/BoardCard.svelte.test.ts
- 3: flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/server/board.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/lib/components/BoardCard.svelte.test.ts
- 4: flaiover/src/lib/components/Messages.svelte, flaiover/src/lib/components/Messages.svelte.test.ts
- 5: design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md
