---
id: S-0343
type: story
nature: improvement
title: Add an archive action to context menus
status: review
owner: alex
created: 2026-10-08T08:07:29Z
updated: 2026-10-09T16:52:24Z
transitions:
  - to: ready
    at: 2026-10-08T08:07:30Z
    by: alex
  - to: in-progress
    at: 2026-10-09T16:37:26Z
    by: agent-S-0343
  - to: review
    at: 2026-10-09T16:52:20Z
    by: agent-S-0343
tags: [dashboard]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, design/system/flai-cli.md, flaiover/src/routes/api/board/archive/+server.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/writes.test.ts, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/lib/cardmenu.ts, flaiover/src/lib/cardmenu.test.ts, flaiover/src/lib/components/LaneMenu.svelte, flaiover/src/lib/components/CardMenu.svelte, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/issues/I-0132-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md, design/issues/I-0133-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: medium
usage:
  source: log
  seconds: 907
  turns:
    - day: 2026-10-09
      ceremony: 2
      hand_edits: 5
      work: 45
  models:
    - model: claude-opus-5-5
      input: 198
      output: 55829
      cache_read: 10855240
      cache_write: 348945
      cost: 5.6485
  strategic:
    - kind: planner
      seconds: 314
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 322
          output: 7465
          cache_read: 1772501
          cache_write: 84159
          cost: 0.3201
        - model: claude-opus-5-5
          input: 54
          output: 17666
          cache_read: 2893663
          cache_write: 143605
          cost: 2.0811
    - kind: orchestrator
      seconds: 4956
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 76
          output: 1190
          cache_read: 6362813
          cache_write: 67975
          cost: 1.6187
        - model: claude-sonnet-5-5
          input: 13
          output: 82
          cache_read: 201897
          cache_write: 29560
          cost: 0.1871
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: alex
    at: 2026-10-08T08:07:29Z
  value: 50
  by: planner-S-0343
  at: 2026-10-09T16:38:04Z
forecast:
  duration: 25m
  delivery: 2026-10-09T17:04:00Z
  basis: "flai forecast: median 78 s per unit of size over 54 done large-band improvement stories on claude-opus-5-5, times size 19 (2 criteria, 17 touches), counted from its start at 2026-10-09T16:38Z."
  by: planner-S-0343
  at: 2026-10-09T16:38:04Z
---
# S-0343 Add an archive action to context menus

## Goal

Cancelled cards are remaining in the cancelled column and causing flai check to fail. Make it simple to archive specific cancelled cards or archive all cards in the cancelled lane with an `Archive` and `Archive All` action in the respective context menus.

## Acceptance criteria
- [x] When cards are in the cancelled lane, right clicking in the lane offers an `Archive All` action that will call archive on all the cancelled cards
- [x] When a card in cancelled is right clicked, an `Archive` action is offered that allows archiving that specific card.

## Tasks
- T-1430 flai's host channel offers board.archive, which archives the done or cancelled items it names with flai archive
- T-1431 The dashboard's POST /api/board/archive archives the items it names through board.archive
- T-1432 The cancelled lane's menu offers Archive All and a cancelled card's menu offers Archive

## Notes

### Planning

Planned by planner-S-0343 on 2026-10-09. The story's agent pulled S-0343 at 16:37Z, during the run, so its tasks are that agent's to write. The planner's proposed tasks are on the plan's thread on S-0343.

Touches, all files, no folder touch:

- **Declared:** none before this plan.
- **Layout:**
  - `flai/internal/hostapi/writes.go` and `writes_test.go`: no host API method archives yet. A new `item.archive` runs `flai archive <id>...`, which already takes IDs.
  - `flaiover/src/routes/api/board/archive/+server.ts`: a new route beside `api/board/limit`, taking the IDs for both actions.
  - `flaiover/src/lib/server/agent.ts`: `REQUIRED_METHODS` lists each write the dashboard calls. `flaiover/src/lib/server/writes.test.ts` covers the writes.
  - `flaiover/src/lib/lanes.ts`, `flaiover/src/lib/cardmenu.ts`, and their tests: `laneEntries` and `cardMenu` build the two menus.
  - `flaiover/src/lib/components/LaneMenu.svelte`, `CardMenu.svelte`, and `flaiover/src/routes/board/+page.svelte` (`pickLane`, `pickCard`), with `lanemenu.svelte.test.ts` and `cardmenu.svelte.test.ts`.
- **Design:**
  - `design/system/flaiover-dashboard.md` § Views, the `/board` row, which lists both menus' entries.
  - `design/system/flai-cli.md`, the `flai hostapi` row, which lists every method.
  - `docs/users/flaiover.md` § Board, the card menu and lane menu tables.
- **Co-change:** `flai touches suggest` gave `docs/users/flai.md` (46%), `docs/operators/index.md` (21%), and `flaiover/src/lib/server/repo.ts` (5%). None was taken: none of them names a host API method, and `repo.ts`'s `write` takes any method name.
- **Overlap:** S-0337 (in progress) also touches `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`, in other sections.

Forecast 25m, from `flai forecast`: median 78 s per unit over 54 done large-band improvement stories, times size 19. It stands: one Go write method, one route, two menu entries, and docs, each beside a pattern it copies.

Cost of delay 50 USD a week: `flai cod` from the operator's 50 USD a week penalty. It stands.
