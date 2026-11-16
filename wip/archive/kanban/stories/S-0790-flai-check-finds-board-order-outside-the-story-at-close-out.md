---
id: S-0790
type: story
nature: improvement
title: "flai check finds `board.order` outside the story at close-out"
status: done
owner: alex
created: 2026-11-15T01:56:43Z
updated: 2026-11-16T00:36:17Z
transitions:
  - to: ready
    at: 2026-11-15T21:11:49Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-15T21:16:43Z
    by: agent-S-0790
  - to: review
    at: 2026-11-15T22:38:19Z
    by: agent-S-0790
  - to: done
    at: 2026-11-16T00:36:17Z
    by: orchestrator
tags: [cli]
touches: [flai/internal/workitem/board.go, flai/internal/workitem/transition.go, flai/internal/workitem/board_update_test.go, flai/cmd/order.go, flai/cmd/accept.go, flai/cmd/board.go, flai/cmd/order_test.go, design/system/workflow.md, design/system/flai-cli.md, design/adrs/README.md, design/issues/I-0375-flai-check-finds-board-order-outside-the-story-at-close-out.md, design/issues/summary.md, design/adrs/0306-every-write-of-board-md-takes-one-board-lock-reads-the-board-again-inside-it.md, design/issues/I-0369-testcontainerstrategicstartunderalayerwaitsfortheharnessimagealone-fails-under.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4919
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 3
      work: 57
  models:
    - model: claude-opus-5-5
      input: 250
      output: 70556
      cache_read: 18430296
      cache_write: 774631
      cost: 9.8326
  strategic:
    - kind: planner
      seconds: 568
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 102
          output: 56329
          cache_read: 6333321
          cache_write: 199793
          cost: 0.4792
        - model: claude-opus-5-5
          input: 102
          output: 24740
          cache_read: 7374912
          cache_write: 188685
          cost: 3.4797
    - kind: orchestrator
      seconds: 238
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 71
          output: 606
          cache_read: 2307367
          cache_write: 352858
          cost: 0.8046
        - model: claude-sonnet-5-5
          input: 8
          output: 718
          cache_read: 34414
          cache_write: 45012
          cost: 0.0633
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: orchestrator
    at: 2026-11-15T21:05:33Z
  value: 5
  by: planner-S-0790
  at: 2026-11-15T21:11:05Z
forecast:
  duration: 1h
  delivery: 2026-11-15T22:49:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, behind S-0748 and S-0750."
  by: flai
  at: 2026-11-15T21:16:20Z
finalized:
  by: orchestrator
  at: 2026-11-15T21:11:43Z
---
# S-0790 flai check finds `board.order` outside the story at close-out

## Goal

This story remediates [I-0375](../../../design/issues/I-0375-flai-check-finds-board-order-outside-the-story-at-close-out.md), "flai check finds `board.order` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0375 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0375 is closed with `flai issue close I-0375 --reason` saying what fixed it

## Tasks
- T-3380 Decide in an ADR that board.md is written under a lock, reloaded inside it, and normalised on every write
- T-3381 workitem writes board.md under a lock, reloaded and normalised inside it, with a test that reproduces I-0375's lost update
- T-3382 flai order, flai accept, and flai board limit write board.md through the locked update
- T-3383 Close I-0375 with what fixed it

## Notes

### Planning

Proposed cause, from I-0375's two instances, the history of `wip/kanban/board.md`, and the archived items. S-0738, S-0746, and S-0739 moved to in-progress in the same second, 2026-11-14T23:49:15Z, when `flai serve` started their agents after S-0737's acceptance. Each start's `TransitionAll` loads `board.md`, removes its story from `order`, and saves it with `Board.Save`, which takes no lock. One save wrote back a board read before another start removed S-0738. S-0738 then stayed in the order while it was in progress, until S-0746's acceptance rewrote the board at 01:56Z.

Proposed fix: every write of `board.md` takes a lock, reloads the board inside it, and drops order names that are no longer ready or backlog stories. That makes true what `design/system/workflow.md` § The pull order already claims. `board.order` stays a finding a close-out records, since it names a real defect.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/board.go` | layout: `LoadBoard`, `Board.Save`, `RemoveFromOrder` |
| `flai/internal/workitem/transition.go` | layout: `TransitionAll` loads and saves the board |
| `flai/internal/workitem/board_update_test.go` | layout: new test that reproduces the lost update |
| `flai/cmd/order.go`, `flai/cmd/order_test.go` | layout: the other two board saves, and their test |
| `flai/cmd/accept.go` | layout: the acceptance's board save |
| `flai/cmd/board.go` | co-change (`flai touches suggest` from `board.go`, 50%) and layout: `flai board limit` saves the board |
| `design/system/workflow.md` | design: § The pull order says every write normalises the list |
| `design/system/flai-cli.md` | co-change (75%) |
| `design/adrs/README.md` | design: the new ADR's index entry; `flai task done` adds the ADR's file |
| `design/issues/I-0375-...md`, `design/issues/summary.md` | criteria: the issue is closed |

The story declared no touches, so all of them are predicted. No folder touch. No `after`: the story reads no other story's output. S-0739, in progress, changes `flai/internal/workitem`'s order and promotion code; the trial merge judges that overlap.

Forecast: flai gave 26m (108 s per unit, size 14). Raised to 1h. Similar remediations, S-0767, S-0730, S-0781, and S-0761, took 48 to 67 minutes of agent time, and a close-out of `flai/` changes runs about 11 minutes. Delivery moved later by the same 34 minutes, to 2026-11-16T00:25Z.

Cost of delay: 5 USD a week, from `flai cod`. That is 2m of time lost per 168h cycle at 150 USD an hour, the input the orchestrator set on TH-1225. It stands as computed.

### Accepted by the orchestrator

- Verified: 72897e27f7f60572a7bc08a46c231f93f91b7850
- Batch: 20261116-002222
- Enqueued verified: 697d56d48cf525c6e9b0cdb4530ca16ef81aae0c
- At: 2026-11-16T00:36:17Z

Verdict: S-0790 meets both criteria at 697d56d; the close-out passed every step there and every changed file is inside its touches.
- 1: flai/internal/workitem/board.go, flai/internal/workitem/transition.go, flai/internal/workitem/board_update_test.go, flai/cmd/order.go, flai/cmd/order_test.go, flai/cmd/accept.go, flai/cmd/board.go, design/adrs/0306-every-write-of-board-md-takes-one-board-lock-reads-the-board-again-inside-it.md
- 2: design/issues/I-0375-flai-check-finds-board-order-outside-the-story-at-close-out.md, design/issues/summary.md
